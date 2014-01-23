// Package edict implements a parser for the EDICT2 Japanese/English dictionary.
package edict

import (
	"bufio"
	"fmt"
	"io"
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
var blacklist = []int{5189, 31179, 104168, 104171, 148763}

func Parse(in io.Reader) ([]Entry, error) {
	var result []Entry

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

type identifierClass int

const (
	none identifierClass = iota
	xref
	detail
	text
)

func parseIdentifier(s string) (identifierClass, string) {
	if _, err := strconv.Atoi(s); err == nil {
		return none, ""
	} else if strings.HasPrefix(s, "See ") {
		return xref, strings.TrimPrefix(s, "See ")
	} else if _, ok := DetailFor[s]; ok {
		return detail, s
	} else {
		return text, s
	}
}

type parseGlossState int

const (
	startGS parseGlossState = iota
	captureGS
	closedGS
	definitionGS
)

func parseGloss(gloss string) (def string, details []Detail, xrefs []string, err error) {
	gloss = strings.TrimSpace(gloss)

	// This is the state machine for parsing the gloss.  We start in the start state, looking
	// for a ( starting an identifier, or the start of a definition (anything other than an
	// opening paren).  Upon seeing a ( we transition to capture, capturing everything that's
	// not a ).  Upon reaching the ), we then transition to closed.  In the closed state, we
	// look for a space, and finding it, transition to start.  At the end of the loop, we must
	// be in the definition-capture state.  If not, we raise an error.
	state := startGS
	captured := make([]rune, 0, len(gloss))
	defcapture := make([]rune, 0, len(gloss))

	for i, c := range gloss {
		switch state {
		case startGS:
			if c == '(' {
				state = captureGS
			} else {
				state = definitionGS
				defcapture = append(defcapture, c)
			}
		case definitionGS:
			defcapture = append(defcapture, c)
		case captureGS:
			if c == ')' {
				state = closedGS
				class, identifier := parseIdentifier(string(captured))

				switch class {
				case detail:
					details = append(details, DetailFor[identifier])
				case xref:
					xrefs = append(xrefs, identifier)
				case text:
					defcapture = append(defcapture, '(')
					for _, c := range captured {
						defcapture = append(defcapture, c)
					}
					defcapture = append(defcapture, ')')
					state = definitionGS
				}
				captured = captured[:0]
			} else if c == ',' {
				// Sometimes things are grouped together, like "(n,adj-no)" instead
				// of "(n) (adj-no)".  If we see a comma, we treat it like a ), but
				// don't transition into a different state.
				class, identifier := parseIdentifier(string(captured))
				if class == detail {
					details = append(details, DetailFor[identifier])
					captured = captured[:0]
				} else {
					// TODO(jrockway): We should blow up here if we get a
					// non-detail result from parseIdentifier, but
					// cross-references can also contain commas.  So if we get
					// no detail, we just continue accumulating as though the
					// comma means nothing.
					captured = append(captured, ',')
				}

			} else {
				captured = append(captured, c)
			}
		case closedGS:
			if c == ' ' {
				state = startGS
			} else {
				err = fmt.Errorf("unexpected '%c' while in closed state (expecting space)", c)
			}
		default:
			err = fmt.Errorf("in unexpected state %v at byte %d '%c'", state, i, c)
		}
	}

	if state != definitionGS {
		err = fmt.Errorf("not in definition state after parsing:\ndetails=%v, xref=%v, def=%s", details, xrefs, def)
		return
	}

	def = string(defcapture)

	return
}

type parseKeyState int

const (
	kanjiKS parseKeyState = iota
	spaceKS
	kanaKS
	doneKS
)

func parseKey(key string) (kanji []string, kana []string, err error) {
	key = strings.TrimSpace(key)

	kanji = make([]string, 0, 5)
	kana = make([]string, 0, 5)

	capture := make([]rune, 0, len(key))
	state := kanjiKS

	// This is a state machine to parse the key field.  Keys look like:
	// KANJI1;KANJI2;... [KANA1;KANA2;...]
	// KANJI1;KANJI2;...
	for _, c := range key {
		if c == ';' || state == kanaKS && c == ']' || state == kanjiKS && c == ' ' {
			// We've just seen a record terminator; ';' for the next element, ']' for
			// the last kana, or ' ' for the switch from kanji to kana.
			if state == kanjiKS {
				kanji = append(kanji, string(capture))
				if c == ' ' {
					state = spaceKS
				}
			} else if state == kanaKS {
				kana = append(kana, string(capture))
				if c == ']' {
					state = doneKS
				}
			}
			capture = capture[:0]
		} else if c == ' ' && state == spaceKS {
			// another space?  ignore.
		} else if c == '[' && state == spaceKS {
			// If we just saw a space and now see a [, we know it's time to start
			// accumulating kana.
			state = kanaKS
		} else {
			// By default, we capture the character.
			capture = append(capture, c)
		}
	}
	if state == kanjiKS {
		kanji = append(kanji, string(capture))
		capture = capture[:0]
		state = doneKS
	}

	if !(state == doneKS || state == spaceKS) {
		err = fmt.Errorf("not in done or space state (in %v) after parsing key %s", state, key)
		return
	}

	if len(capture) != 0 {
		err = fmt.Errorf("chars still in capture buffer after parsing %s! %v (state=%v)", key, capture, state)
	}

	return
}

func fixKey(key string) string {
	if strings.ContainsRune(key, '(') {
		parts := strings.Split(key, "(")
		return parts[0]
	}
	return key
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

	var err error
	result.Kanji, result.Kana, err = parseKey(parts[0])
	if err != nil {
		return result, err
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

	// TODO(jrockway): Kanji and Kana keys can also contain (information) identifiers like
	// entries and glosses.  For now, just remove those, though they are valuable.  (Showing the
	// most common reading, for example.)
	for i, kanji := range result.Kanji {
		result.Kanji[i] = fixKey(kanji)
	}
	for i, kana := range result.Kana {
		result.Kana[i] = fixKey(kana)
	}

	return result, nil
}
