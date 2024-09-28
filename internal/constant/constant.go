package constant

import "cbeimers113/ilo/internal/util"

// ANSI escape codes for colors
const (
	ColReset  = "\033[0m"
	ColRed    = "\033[31m"
	ColGreen  = "\033[32m"
	ColYellow = "\033[33m"
	ColBlue   = "\033[34m"
)

// Esperanto orthography mapping
var Orthography = map[rune]string {
	'ĉ': "cx",
	'ĝ': "gx",
	'ĥ': "hx",
	'ĵ': "jx",
	'ŝ': "sx",
	'ŭ': "ux",
}

// Ilo keywords
var Keywords = util.NewSet(
	"se",
	"alie",
	"dum",
	"por",
	"redonu",
	"aux",
	"kaj",
	"ne",
	"estas",
	"pli",
	"malpli",
	"ol",
	"jen",
	"uzante",
	"de",
	"kiel",
	"ago",
	"agoj",
	"tuto",
	"tutoj",
	"punkto",
	"punktoj",
	"vero",
	"veroj",
	"teksto",
	"tekstoj",
	"veras",
	"malveras",
)

// Mathematical, logical and syntactic operators
var Operators = util.NewSet(
	'+',
	'-',
	'*',
	'/',
	'^',
	'%',
	'=',
	',',
	'.',
	':',
	'!',
	'?',
	'#',
	'\\',
)
