package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adrianliechti/go-cli"
)

func main() {
	// Check for light theme flag
	for _, arg := range os.Args[1:] {
		if arg == "--light" || arg == "-l" {
			cli.SetLightTheme()
			break
		}
	}

	cli.Title("CLI Demo")
	fmt.Println()

	// Input demo
	name, err := cli.Input("What is your name?", "anonymous")
	if err != nil {
		cli.Fatal("Input error:", err)
	}

	// Select demo with filtering
	colors := []string{
		"Red",
		"Orange",
		"Yellow",
		"Green",
		"Cyan",
		"Blue",
		"Purple",
		"Pink",
		"Magenta",
		"White",
		"Black",
		"Gray",
	}

	_, color, err := cli.Select("Pick your favorite color (type to filter):", colors)
	if err != nil {
		cli.Fatal("Select error:", err)
	}

	// Confirm demo
	confirmed, err := cli.Confirm("Do you want to continue?", true)
	if err != nil {
		cli.Fatal("Confirm error:", err)
	}

	if !confirmed {
		cli.Warn("Cancelled by user")
		return
	}

	// Text demo
	description, err := cli.Text("Enter a description:", "")
	if err != nil {
		cli.Fatal("Text error:", err)
	}

	// File demo
	file, err := cli.File("Select a config file:", []string{".json", ".yaml", ".toml", ".yml"})
	if err != nil {
		if err == cli.ErrUserAborted {
			cli.Warn("File selection skipped")
			file = "(none)"
		} else {
			cli.Fatal("File error:", err)
		}
	}

	fmt.Println()

	// Spinner demo
	err = cli.Run("Processing your data...", func() error {
		time.Sleep(2 * time.Second)
		return nil
	})
	if err != nil {
		cli.Fatal("Run error:", err)
	}

	fmt.Println()
	cli.Info("✓ Processing complete!")
	fmt.Println()

	// Table demo - Summary
	cli.Title("Summary")
	fmt.Println()

	headers := []string{"Field", "Value"}
	rows := [][]string{
		{"Name", name},
		{"Favorite Color", color},
		{"Description", truncate(description, 40)},
		{"Config File", file},
	}

	cli.Table(headers, rows)
	fmt.Println()

	cli.Info("Thank you for trying the CLI demo!")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
