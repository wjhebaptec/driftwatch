package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/runner"
)

const version = "0.1.0"

func main() {
	var (
		configPath  = flag.String("config", "", "path to config file (default: driftwatch.yaml)")
		showVersion = flag.Bool("version", false, "print version and exit")
		format      = flag.String("format", "", "output format: text or json (overrides config)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("driftwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if *format != "" {
		cfg.Output.Format = *format
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	r := runner.New(cfg)
	drifted, err := r.Run()
	if err != nil {
		log.Fatalf("run failed: %v", err)
	}

	if drifted {
		os.Exit(1)
	}
}
