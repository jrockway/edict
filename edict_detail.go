package edict

// A part of speech "detail" marking from http://www.edrdg.org/jmdict/edict_doc.html
type Detail int

// Exactly as listed by the documentation, with first letter capitalized and - replaced by _.
const (
	// Parts of speech
	AdjI   Detail = iota // adjective (keiyoushi)
	AdjNa                // adjectival nouns or quasi_adjectives (keiyodoshi)
	AdjNo                // nouns which may take the genitive case particle `no'
	AdjPn                // pre_noun adjectival (rentaishi)
	AdjT                 // `taru' adjective
	AdjF                 // noun or verb acting prenominally (other than the above)
	Adj                  // former adjective classification (being removed)
	Adv                  // adverb (fukushi)
	AdvN                 // adverbial noun
	AdvTo                // adverb taking the `to' particle
	Aux                  // auxiliary
	AuxV                 // auxiliary verb
	AuxAdj               // auxiliary adjective
	Conj                 // conjunction
	Ctr                  // counter
	Exp                  // Expressions (phrases, clauses, etc.)
	Int                  // interjection (kandoushi)
	Iv                   // irregular verb
	N                    // noun (common) (futsuumeishi)
	NAdv                 // adverbial noun (fukushitekimeishi)
	NPref                // noun, used as a prefix
	NSuf                 // noun, used as a suffix
	NT                   // noun (temporal) (jisoumeishi)
	Num                  // numeric
	Pn                   // pronoun
	Pref                 // prefix
	Prt                  // particle
	Suf                  // suffix
	V1                   // Ichidan verb
	V2aS                 // Nidan verb with 'u' ending (archaic)
	V4h                  // Yodan verb with `hu/fu' ending (archaic)
	V4r                  // Yodan verb with `ru' ending (archaic)
	V5                   // Godan verb (not completely classified)
	V5aru                // Godan verb _ _aru special class
	V5b                  // Godan verb with `bu' ending
	V5g                  // Godan verb with `gu' ending
	V5k                  // Godan verb with `ku' ending
	V5kS                 // Godan verb _ iku/yuku special class
	V5m                  // Godan verb with `mu' ending
	V5n                  // Godan verb with `nu' ending
	V5r                  // Godan verb with `ru' ending
	V5rI                 // Godan verb with `ru' ending (irregular verb)
	V5s                  // Godan verb with `su' ending
	V5t                  // Godan verb with `tsu' ending
	V5u                  // Godan verb with `u' ending
	V5uS                 // Godan verb with `u' ending (special class)
	V5uru                // Godan verb _ uru old class verb (old form of Eru)
	V5z                  // Godan verb with `zu' ending
	Vz                   // Ichidan verb _ zuru verb _ (alternative form of _jiru verbs)
	Vi                   // intransitive verb
	Vk                   // kuru verb _ special class
	Vn                   // irregular nu verb
	Vs                   // noun or participle which takes the aux. verb suru
	VsC                  // su verb _ precursor to the modern suru
	VsI                  // suru verb _ irregular
	VsS                  // suru verb _q special class
	Vt                   // transitive verb

	// Field of application
	Buddh   // Buddhist term
	MA      // martial arts term
	Comp    // computer terminology
	Food    // food term
	Geom    // geometry term
	Gram    // grammatical term
	Ling    // linguistics terminology
	Math    // mathematics
	Mil     // military
	Physics // physics terminology

	// Miscellaneous markings
	X      // rude or X-rated term
	Abbr   // abbreviation
	Arch   // archaism
	Ateji  // ateji (phonetic) reading
	Chn    // children's language
	Col    // colloquialism
	Derog  // derogatory term
	EK     // exclusively kanji
	Ek     // exclusively kana
	Fam    // familiar language
	Fem    // female term or language
	Gikun  // gikun (meaning) reading
	Hon    // honorific or respectful (sonkeigo) language
	Hum    // humble (kenjougo) language
	Ik     // word containing irregular kana usage
	IK     // word containing irregular kanji usage
	Id     // idiomatic expression
	Io     // irregular okurigana usage
	MSl    // manga slang
	Male   // male term or language
	MaleSl // male slang
	OK     // word containing out-dated kanji
	Obs    // obsolete term
	Obsc   // obscure term
	Ok     // out-dated or obsolete kana usage
	OnMim  // onomatopoeic or mimetic word
	Poet   // poetical term
	Pol    // polite (teineigo) language
	Rare   // rare (now replaced by "obsc")
	Sens   // sensitive word
	Sl     // slang
	UK     // word usually written using kanji alone
	Uk     // word usually written using kana alone
	Vulg   // vulgar expression or word

	// Indicators for common words
	Common
)

var DetailString = map[Detail]string{
	AdjI:    "adj-i",
	AdjNa:   "adj-na",
	AdjNo:   "adj-no",
	AdjPn:   "adj-pn",
	AdjT:    "adj-t",
	AdjF:    "adj-f",
	Adj:     "adj",
	Adv:     "adv",
	AdvN:    "adv-n",
	AdvTo:   "adv-to",
	Aux:     "aux",
	AuxV:    "aux-v",
	AuxAdj:  "aux-adj",
	Conj:    "conj",
	Ctr:     "ctr",
	Exp:     "exp",
	Int:     "int",
	Iv:      "iv",
	N:       "n",
	NAdv:    "n-adv",
	NPref:   "n-pref",
	NSuf:    "n-suf",
	NT:      "n-t",
	Num:     "num",
	Pn:      "pn",
	Pref:    "pref",
	Prt:     "prt",
	Suf:     "suf",
	V1:      "v1",
	V2aS:    "v2a-s",
	V4h:     "v4h",
	V4r:     "v4r",
	V5:      "v5",
	V5aru:   "v5aru",
	V5b:     "v5b",
	V5g:     "v5g",
	V5k:     "v5k",
	V5kS:    "v5k-s",
	V5m:     "v5m",
	V5n:     "v5n",
	V5r:     "v5r",
	V5rI:    "v5r-i",
	V5s:     "v5s",
	V5t:     "v5t",
	V5u:     "v5u",
	V5uS:    "v5u-s",
	V5uru:   "v5uru",
	V5z:     "v5z",
	Vz:      "vz",
	Vi:      "vi",
	Vk:      "vk",
	Vn:      "vn",
	Vs:      "vs",
	VsC:     "vs-c",
	VsI:     "vs-i",
	VsS:     "vs-s",
	Vt:      "vt",
	Buddh:   "buddh",
	MA:      "mA",
	Comp:    "comp",
	Food:    "food",
	Geom:    "geom",
	Gram:    "gram",
	Ling:    "ling",
	Math:    "math",
	Mil:     "mil",
	Physics: "physics",
	X:       "x",
	Abbr:    "abbr",
	Arch:    "arch",
	Ateji:   "ateji",
	Chn:     "chn",
	Col:     "col",
	Derog:   "derog",
	EK:      "eK",
	Ek:      "ek",
	Fam:     "fam",
	Fem:     "fem",
	Gikun:   "gikun",
	Hon:     "hon",
	Hum:     "hum",
	Ik:      "ik",
	IK:      "iK",
	Id:      "id",
	Io:      "io",
	MSl:     "m-sl",
	Male:    "male",
	MaleSl:  "male-sl",
	OK:      "oK",
	Obs:     "obs",
	Obsc:    "obsc",
	Ok:      "ok",
	OnMim:   "on-mim",
	Poet:    "poet",
	Pol:     "pol",
	Rare:    "rare",
	Sens:    "sens",
	Sl:      "sl",
	UK:      "uK",
	Uk:      "uk",
	Vulg:    "vulg",
	Common:  "P",
}

var DetailFor map[string]Detail

func init() {
	DetailFor = make(map[string]Detail, len(DetailString))
	for detail, str := range DetailString {
		DetailFor[str] = detail
	}
}

func (d Detail) String() string {
	return DetailString[d]
}
