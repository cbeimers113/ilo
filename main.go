package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"cbeimers113/ilo/internal/config"
	"cbeimers113/ilo/internal/constant"
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
	log.Info(fmt.Sprintf("ilo v%s", version))
	process.SetUp()

	// Read command line arguments and begin processing
	if src, target, err := process.ParseArgs(cfg, os.Args[1:]); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Info(fmt.Sprintf("%s %s -> %s", cfg.Message(constant.MsgCompiling), src, target))
		startTime := time.Now()
		// TODO: transpile
		deltaTime := time.Since(startTime).Milliseconds()
		log.Info(fmt.Sprintf("%s %d ms", cfg.Message(constant.MsgFinished), deltaTime))
	}
}
