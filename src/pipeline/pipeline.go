package pipeline

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	config *Config
	logger *log.Logger
	workDir string
}

// NewPipeline creates a new pipeline instance
func NewPipeline(config *Config, logger *log.Logger) *Pipeline {
	return &Pipeline{
		config: config,
		logger: logger,
		workDir: filepath.Join(os.TempDir(), fmt.Sprintf("pipeline-%d", time.Now().Unix())),
	}
}

// Run executes the pipeline
func (p *Pipeline) Run() error {
	p.logger.Printf("Starting pipeline: %s", p.config.Name)
	
	// Create working directory
	if err := os.MkdirAll(p.workDir, 0755); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}
	defer os.RemoveAll(p.workDir)
	
	// Execute pipeline stages
	if err := p.checkout(); err != nil {
		return err
	}
	
	if err := p.build(); err != nil {
		return err
	}
	
	if err := p.test(); err != nil {
		return err
	}
	
	if err := p.deploy(); err != nil {
		return err
	}
	
	return nil
}

// checkout clones the repository
func (p *Pipeline) checkout() error {
	p.logger.Println("Stage: Checkout")
	
	cmd := exec.Command("git", "clone", "--branch", p.config.Branch, p.config.Repository, p.workDir)
	cmd.Stdout = p.logger.Writer()
	cmd.Stderr = p.logger.Writer()
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("checkout failed: %w", err)
	}
	
	p.logger.Println("Checkout completed successfully")
	return nil
}

// build runs the build command
func (p *Pipeline) build() error {
	p.logger.Println("Stage: Build")
	
	if p.config.BuildCmd == "" {
		p.logger.Println("No build command specified, skipping build stage")
		return nil
	}
	
	return p.runCommand(p.config.BuildCmd, "build")
}

// test runs the test command
func (p *Pipeline) test() error {
	p.logger.Println("Stage: Test")
	
	if p.config.TestCmd == "" {
		p.logger.Println("No test command specified, skipping test stage")
		return nil
	}
	
	return p.runCommand(p.config.TestCmd, "test")
}

// deploy runs the deployment command
func (p *Pipeline) deploy() error {
	p.logger.Println("Stage: Deploy")
	
	if p.config.DeployCmd == "" {
		p.logger.Println("No deploy command specified, skipping deploy stage")
		return nil
	}
	
	return p.runCommand(p.config.DeployCmd, "deploy")
}

// runCommand executes a shell command
func (p *Pipeline) runCommand(cmdStr string, stage string) error {
	// Split the command string into command and arguments
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return fmt.Errorf("empty command for %s stage", stage)
	}
	
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = p.workDir
	cmd.Stdout = p.logger.Writer()
	cmd.Stderr = p.logger.Writer()
	
	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range p.config.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", stage, err)
	}
	
	p.logger.Printf("%s completed successfully", stage)
	return nil
}