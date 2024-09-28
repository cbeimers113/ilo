package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/locale"
)

// Flag parsing
type FlagId uint8
type FlagOpts map[FlagId]bool

const (
	FlagDebug FlagId = iota
)

// Command line flag options
var flagMap = map[string]FlagId{
	"d":      FlagDebug,
	"-debug": FlagDebug,
}

// SetUp creates the intermediate and output directories
func SetUp() error {
	interDir := filepath.Join(".", ".ilo_go")
	return os.MkdirAll(interDir, os.ModePerm)
}

// ParseArgs parses the command line arguments and returns the source
// and target file names as well as flag options
func ParseArgs(cfg *config.Config, args []string) (string, string, FlagOpts, error) {
	// Filter out flags
	pArgs := make([]string, 0)
	pFlags := make([]string, 0)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			pFlags = append(pFlags, strings.TrimPrefix(arg, "-"))
			continue
		}

		pArgs = append(pArgs, arg)
	}

	if err := checkArgs(cfg, pArgs); err != nil {
		return "", "", nil, err
	}

	source := pArgs[0]
	target := strings.TrimSuffix(source, locale.SourceFileExtension)

	if len(pArgs) > 1 {
		target = strings.ToLower(pArgs[1])

		for _, token := range pArgs[2:] {
			target += strings.ToUpper(string(token[0]))

			if len(token) > 1 {
				target += strings.ToLower(token[1:])
			}
		}
	}

	// Parse flags
	var (
		flags FlagOpts
		err   error
	)

	if flags, err = parseFlags(cfg, pFlags); err != nil {
		return "", "", nil, err
	}

	return source, target, flags, nil
}

// checkArgs checks the command line arguments to ensure that a source file is passed in
func checkArgs(cfg *config.Config, args []string) error {
	// Check for invalid characters
	r := regexp.MustCompile(`^[a-zA-Z_./]+[a-zA-Z0-9_./]*$`)
	for _, arg := range args {
		if !r.MatchString(arg) {
			return fmt.Errorf("%s \"%s\"", cfg.Message(locale.ErrInvalidChars), arg)
		}
	}

	// Make sure arguments are supplied
	if len(args) == 0 {
		return errors.New(cfg.Message(locale.ErrNoArguments))
	}

	// The first argument must be a source file
	if !strings.HasSuffix(args[0], locale.SourceFileExtension) {
		return errors.New(cfg.Message(locale.ErrNoSourceFile))
	}

	// Make sure the source file exists
	if _, err := os.Stat(args[0]); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s: \"%s\"", cfg.Message(locale.ErrSourceNotExist), args[0])
	}

	return nil
}

func parseFlags(cfg *config.Config, flags []string) (FlagOpts, error) {
	flagOpts := make(FlagOpts)
	for _, flag := range flags {
		if flagVal, ok := flagMap[flag]; ok {
			flagOpts[flagVal] = true
			continue
		}

		return nil, fmt.Errorf("%s: \"-%s\"", cfg.Message(locale.ErrInvalidOption), flag)
	}

	return flagOpts, nil
}
