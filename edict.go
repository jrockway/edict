// Package edict implements a parser for the EDICT2 Japanese/English dictionary.
package edict

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Gloss encodes an English definition for a Japanese word.
type Gloss struct {
	Definition  string   // English translation.
	Information []Detail // Information about this particular definition.
	Xref        []string // Xref to related entries (to the Kanji key), "see also".
}

// Entry encodes a line of edict2 input.
type Entry struct {
	Kanji              []string // Kanji key.
	Kana               []string // Kana transcription of keys.
	Information        []Detail // Information about the word; part of speech, conjugation type, etc.
	Gloss              []Gloss  // The "glosses", English definitions, ordered by frequency.
	Sequence           string   // The entry's unique identifier.
	RecordingAvailable bool     // True if an audio clip of the entry reading is available from the JapanesePod101.com site.
}

// String formats an Entry as a single line; not in the edict2 format, but familiar enough.
func (e Entry) String() string {
	recording := ""
	if e.RecordingAvailable {
		recording = "X"
	}

	return fmt.Sprintf("%v %v /%v %v/%s%s/", e.Kanji, e.Kana, e.Information, e.Gloss, e.Sequence, recording)
}

// These lines contain the record separator as part of the entry, making the whole thing
// signficantly more difficult to parse.  We'll skip these for now, and I'll patch the dictionary to
// not do this :)
var blacklist = []int{31179, 104168, 104171}

func Parse(in io.Reader) ([]Entry, error) {
	result := []Entry{}
	scanner := bufio.NewScanner(in)
	line := 0
lines:
	for scanner.Scan() {
		line++
		entry, err := parseLine(scanner.Text())
		if err != nil {
			for _, knownBadLine := range blacklist {
				if knownBadLine == line {
					continue lines
				}
			}
			return result, fmt.Errorf("parse: line %d: %s", line, err)
		}
		result = append(result, entry)
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("parse: past EOF (line %d): %s", line, err)
	}

	return result, nil
}

// Regular expressions for parsing entry lines.
var (
	// TODO(jrockway): kanji field is optional, according to the docs.  Key part looks like
	// "key1;key2;... [reading1;reading2;...] " (note the space at the end)
	parseKeys = regexp.MustCompile(`^([^[:space:]]+) \[([^\]]+)\] `)
)

type parseGlossState int

const (
	start parseGlossState = iota
	capture
	closed
	definition
)

func parseIdentifier(s string) (*Detail, *string, *string) {
	if _, err := strconv.Atoi(s); err == nil {
		return nil, nil, nil
	} else if strings.HasPrefix(s, "See ") {
		word := strings.TrimPrefix(s, "See ")
		return nil, &word, nil
	} else if detail, ok := DetailFor[s]; ok {
		return &detail, nil, nil
	} else {
		return nil, nil, &s
	}
}

func parseGloss(gloss string) (def string, details []Detail, xrefs []string, err error) {
	gloss = strings.TrimSpace(gloss)

	// This is the state machine for parsing the gloss.  We start in the start state, looking
	// for a ( starting an identifier, or the start of a definition (anything other than an
	// opening paren).  Upon seeing a ( we transition to capture, capturing everything that's
	// not a ).  Upon reaching the ), we then transition to closed.  In the closed state, we
	// look for a space, and finding it, transition to start.  At the end of the loop, we must
	// be in the definition-capture state.  If not, we raise an error.
	state := start
	captured := []rune{}
	defcapture := []rune{}

	for idx, c := range gloss {
		switch state {
		case start:
			if c == '(' {
				state = capture
			} else {
				state = definition
				defcapture = append(defcapture, c)
			}
		case definition:
			defcapture = append(defcapture, c)
		case capture:
			if c == ')' {
				state = closed
				d, x, u := parseIdentifier(string(captured))

				if d != nil {
					details = append(details, *d)
				} else if x != nil {
					xrefs = append(xrefs, *x)
				} else if u != nil {
					defcapture = append(defcapture, '(')
					for _, c := range captured {
						defcapture = append(defcapture, c)
					}
					defcapture = append(defcapture, ')')
					state = definition
				}
				captured = []rune{}
			} else {
				captured = append(captured, c)
			}
		case closed:
			if c == ' ' {
				state = start
			} else {
				err = fmt.Errorf("unexpected '%c' while in closed state (expecting space)", c)
			}
		default:
			err = fmt.Errorf("in unexpected state %v at byte %d '%c'", state, idx, c)
		}
	}

	if state != definition {
		err = fmt.Errorf("not in definition state after parsing:\ndetails=%v, xref=%v, def=%s", details, xrefs, def)
		return
	}

	def = string(defcapture)

	return
}

func parseLine(line string) (Entry, error) {
	result := Entry{}
	parts := strings.Split(line, "/")
	last := parts[len(parts)-1]
	if last != "" {
		return result, fmt.Errorf("parseLine: last component should be blank, but is %s", last)
	}

	// Parse the sequence number part, since having this in the result makes misparsing lines
	// easier to grep for.
	result.Sequence = parts[len(parts)-2]
	if strings.HasSuffix(result.Sequence, "X") {
		result.RecordingAvailable = true
		result.Sequence = strings.TrimSuffix(result.Sequence, "X")
	}
	result.Sequence = seqParts[1]

	// Parse the first part into Kanji/Kana fields.
	key := parts[0]
	keyParts := parseKeys.FindStringSubmatch(key)
	if len(keyParts) == 0 {
		result.Kanji = strings.Split(key, ";")
	} else if keyParts[0] != key || len(keyParts) != 3 {
		return result, fmt.Errorf("incomplete match on key '%s':\n got '%v'", key, keyParts)
	} else {
		result.Kanji = strings.Split(keyParts[1], ";")
		result.Kana = strings.Split(keyParts[2], ";")
	}

	// Next we get some details from the first gloss.
	glosses := []string{parts[1]}

	if len(parts) > 4 {
		// If there's more than one gloss, the entry-wide details come before the (1)
		// marker.
		firstGlossParts := strings.Split(parts[1], "(1)")
		if len(firstGlossParts) == 2 {
			_, detail, xref, err := parseGloss(firstGlossParts[0] + "fake definition")
			if err != nil {
				return result, fmt.Errorf("parsing entry details: %s", err)
			}
			if len(xref) != 0 {
				return result, fmt.Errorf("unexpected xref in global details section")
			}
			result.Information = detail
			glosses[0] = firstGlossParts[1]
		}
	}

	// We already have the first gloss in glosses, add the rest here.
	if len(parts) > 4 {
		for _, gloss := range parts[2 : len(parts)-2] {
			glosses = append(glosses, gloss)
		}
	}

	result.Gloss = []Gloss{}
	for _, gloss := range glosses {
		if gloss == "(P)" { // what a terrible file format
			result.Information = append(result.Information, Common)
			continue
		}

		def, detail, xref, err := parseGloss(gloss)
		if err != nil {
			return result, fmt.Errorf("parsing gloss %s got err %s", gloss, err)
		}
		result.Gloss = append(result.Gloss, Gloss{def, detail, xref})
	}

	// In the event that there's only one gloss, transfer the details to the entry.
	if len(parts) <= 4 {
		result.Information = result.Gloss[0].Information
		result.Gloss[0].Information = []Detail{}
	}

	return result, nil
}
