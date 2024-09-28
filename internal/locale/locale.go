package locale

// Locale indices
const (
	// Error messages
	ErrInvalidChars int = iota
	ErrNoArguments
	ErrNoSourceFile
	ErrInvalidFileName
	ErrSourceNotExist
	ErrInvalidOption
	ErrInvalidChar
	ErrUnclosedQuote
	ErrInvalidToken

	// Debug flags
	DbgTokenizing
	DbgParsing

	// Info
	MsgCompiling
	MsgFinished

	SourceFileExtension = ".ilo"
)

var LocalizedStrings = map[string][]string{
	"en": {
		"invalid characters in text",
		"no arguments supplied, nothing to do",
		"source file not supplied",
		"invalid file name",
		"source file doesn't exist",
		"invalid option",
		"invalid character on line",
		"unclosed quote on line",
		"invalid token on line",

		"tokenizing",
		"parsing",

		"compiling",
		"finished in",
	},
	"eo": {
		"nevalidaj signoj en teksto",
		"neniuj argumentoj donitaj, nenio por fari",
		"fonta dosiero ne provizita",
		"nevalida nomo de dosiero",
		"fonta dosiero ne ekzistas",
		"nevalida opcio",
		"nevalida signo en linio",
		"nefermita citajxo en linio",
		"nevalida jxetono en linio",

		"jxetoniganta",
		"analizanta",

		"kompilanta",
		"finita en",
	},
}
