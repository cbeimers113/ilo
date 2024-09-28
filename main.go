package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/locale"
	"cbeimers113/ilo/internal/log"
	"cbeimers113/ilo/internal/process"
)

//go:embed version
var version string

// loadConfig loads the local config file (must be called ilo.cfg)
func loadConfig() *config.Config {
	errMsg := "failed to load config | ne povis sxargxi agordon"
	cfgData, err := os.ReadFile("ilo.cfg")
	if err != nil {
		log.Fatal(errMsg)
	}

	cfg, err := config.New(cfgData)
	if err != nil {
		log.Fatal(errMsg)
	}

	return cfg
}

func main() {
	// Load config
	cfg := loadConfig()

	// Prepare environment
	log.Info(fmt.Sprintf("Ilo v%s", version))
	process.SetUp()

	// Read command line arguments and begin processing
	var (
		src    string
		target string
		flags  process.FlagOpts
		err    error
	)

	if src, target, flags, err = process.ParseArgs(cfg, os.Args[1:]); err != nil {
		log.Fatal(err.Error())
	}

	log.Info(fmt.Sprintf("%s %s -> %s", cfg.Message(locale.MsgCompiling), src, target))
	startTime := time.Now()

	// Read the source file
	data, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Preprocess the raw source code to convert Esperanto characters to their x-mode counterparts
	if flags[process.FlagDebug] {
		log.Debug(cfg.Message(locale.DbgPreprocessing))
	}
	sourceCode := process.Preprocess(string(data))

	// Create tokens from the raw source code
	if flags[process.FlagDebug] {
		log.Debug(cfg.Message(locale.DbgTokenizing))
	}
	tokens := process.Tokenize(cfg, sourceCode)
	if flags[process.FlagDebug] {
		for _, token := range tokens {
			fmt.Println(token)
		}
	}

	// Parse the tokens into an abstract syntax tree
	if flags[process.FlagDebug] {
		log.Debug(cfg.Message(locale.DbgParsing))
	}

	deltaTime := time.Since(startTime).Milliseconds()
	log.Info(fmt.Sprintf("%s %d ms", cfg.Message(locale.MsgFinished), deltaTime))
}
