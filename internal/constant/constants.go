package constant

const (
	// Locale indices
	ErrInvalidChars int = iota
	ErrNoArguments
	ErrNoSourceFile
	ErrInvalidFileName
	ErrSourceNotExist
	MsgCompiling
	MsgFinished

	SourceFileExtension = ".ilo"
)

var LocalizedStrings = map[string][]string{
	"en": {
		"invalid characters in text",
		"missing argument(s)",
		"source file not supplied",
		"invalid file name",
		"source file doesn't exist",
		"compiling",
		"finished in",
	},
	"eo": {
		"nevalidaj signoj en teksto",
		"mankas argumento(j)n",
		"fonta dosiero ne provizita",
		"nevalida nomo de dosiero",
		"fonta dosiero ne ekzistas",
		"kompilanta",
		"finita en",
	},
}
