package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"c:\Users\iamma\Programs\trae\Go-Simple-CI-CD-Pipeline\src\pipeline"
)

func main() {
	// Define command line flags
	configFile := flag.String("config", "config/pipeline.yaml", "Path to pipeline configuration file")
	flag.Parse()

	// Initialize logger
	logger := log.New(os.Stdout, "[CI/CD] ", log.LstdFlags)
	logger.Println("Starting CI/CD pipeline")

	// Load pipeline configuration
	config, err := pipeline.LoadConfig(*configFile)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and run the pipeline
	p := pipeline.NewPipeline(config, logger)
	if err := p.Run(); err != nil {
		logger.Fatalf("Pipeline execution failed: %v", err)
	}

	fmt.Println("Pipeline completed successfully!")
}