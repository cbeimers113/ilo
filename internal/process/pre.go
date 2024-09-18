package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cbeimers113/ilo/internal/constant"
	"cbeimers113/ilo/internal/config"
)

// SetUp creates the intermediate and output directories
func SetUp() error {
	interDir := filepath.Join(".", ".ilo_c")
	return os.MkdirAll(interDir, os.ModePerm)
}

// ParseArgs parses the command line arguments and returns the source and target file names
func ParseArgs(cfg *config.Config, args []string) (string, string, error) {
	if err := checkArgs(cfg, args); err != nil {
		return "", "", err
	}

	source := args[0]
	target := strings.TrimSuffix(source, constant.SourceFileExtension)

	if len(args) > 1 {
		target = strings.ToLower(args[1])

		for _, token := range args[2:] {
			target += strings.ToUpper(string(token[0]))

			if len(token) > 1 {
				target += strings.ToLower(token[1:])
			}
		}
	}

	return source, target, nil
}

// checkArgs checks the command line arguments to ensure that a source file is passed in
func checkArgs(cfg *config.Config, args []string) error {
	// Check for invalid characters
	r := regexp.MustCompile(`^[a-zA-Z_./]+[a-zA-Z0-9_./]*$`)
	for _, arg := range args {
		if !r.MatchString(arg) {
			return fmt.Errorf("%s \"%s\"", cfg.Message(constant.ErrInvalidChars), arg)
		}
	}

	// Make sure arguments are supplied
	if len(args) == 0 {
		return errors.New(cfg.Message(constant.ErrNoArguments))
	}

	// The first argument must be a source file
	if !strings.HasSuffix(args[0], constant.SourceFileExtension) {
		return errors.New(cfg.Message(constant.ErrNoSourceFile))
	}

	// Make sure the source file exists
	if _, err := os.Stat(args[0]); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%s: \"%s\"", cfg.Message(constant.ErrSourceNotExist), args[0])
	}

	return nil
}
