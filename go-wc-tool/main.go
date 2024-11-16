package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	file_utilities "coding-challenge/wc-tool-go/utilities"
)

func main() {
	args := os.Args[1:]
	log.Println("Received: ", args)

	// Validate and parse arguments
	command, file, err := parseArgs(args)
	if err != nil {
		log.Fatal(err)
	}

	// Handle stdin buffering if no file is provided
	var inputData []byte
	if file == "" {
		inputData, err = readStdin()
		if err != nil {
			log.Fatalf("Error reading stdin: %v", err)
		}
	}

	// Execute the command
	if err := executeCommand(command, file, inputData); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

// parseArgs validates and parses the command-line arguments.
func parseArgs(args []string) (command, file string, err error) {
	if len(args) > 2 {
		return "", "", fmt.Errorf("wc tool can only handle one or two arguments")
	}

	if len(args) == 1 {
		// Check if it's a command or a file
		if strings.ContainsAny(args[0], ".") {
			return "", args[0], nil // Default file processing
		}
		return extractCommand(args[0]), "", nil // Command with stdin
	} else if len(args) == 0 {
		return "", "", nil
	}

	// Two arguments: command and file
	return extractCommand(args[0]), args[1], nil
}

// extractCommand extracts the command from the argument.
func extractCommand(arg string) string {
	if strings.HasPrefix(arg, "-") && len(arg) > 1 {
		return arg[1:] // Remove the dash and return the command
	}
	return ""
}

// readStdin reads all input from stdin and returns it as a byte slice.
func readStdin() ([]byte, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %v", err)
	}
	return data, nil
}

// executeCommand runs the appropriate function based on the command and file.
func executeCommand(command, file string, inputData []byte) error {
	commandHandlers := map[string]func(string, []byte) (string, error){
		"c": countBytes,
		"l": countLines,
		"d": countWords,
		"m": countChars,
	}

	if command == "" {
		// Run all commands if no specific command is provided
		lineRes, err := countLines(file, inputData)
		if err != nil {
			return err
		}
		wordRes, err := countWords(file, inputData)
		if err != nil {
			return err
		}
		byteRes, err := countBytes(file, inputData)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", lineRes, wordRes, byteRes)
		return nil
	}

	// Run the specific command
	handler, exists := commandHandlers[command]
	if !exists {
		return fmt.Errorf("unsupported command: %s", command)
	}

	result, err := handler(file, inputData)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

// countBytes counts the bytes in the file or stdin.
func countBytes(file string, inputData []byte) (string, error) {
	var count int
	var err error

	if file == "" {
		count, err = file_utilities.NumberOfBytesCommand(bytes.NewReader(inputData))
	} else {
		count, err = file_utilities.NumberOfBytesCommand(file)
	}

	if err != nil {
		return "", fmt.Errorf("error counting bytes: %v", err)
	}
	return fmt.Sprintf("%d bytes", count), nil
}

// countLines counts the lines in the file or stdin.
func countLines(file string, inputData []byte) (string, error) {
	var count int
	var err error

	if file == "" {
		count, err = file_utilities.NumberOfLinesCommand(bytes.NewReader(inputData))
	} else {
		count, err = file_utilities.NumberOfLinesCommand(file)
	}

	if err != nil {
		return "", fmt.Errorf("error counting lines: %v", err)
	}
	return fmt.Sprintf("%d lines", count), nil
}

// countWords counts the words in the file or stdin.
func countWords(file string, inputData []byte) (string, error) {
	var count int
	var err error

	if file == "" {
		count, err = file_utilities.NumberOfWordsCommand(bytes.NewReader(inputData))
	} else {
		count, err = file_utilities.NumberOfWordsCommand(file)
	}

	if err != nil {
		return "", fmt.Errorf("error counting words: %v", err)
	}
	return fmt.Sprintf("%d words", count), nil
}

// countChars counts the characters in the file or stdin.
func countChars(file string, inputData []byte) (string, error) {
	var count int
	var err error

	if file == "" {
		count, err = file_utilities.NumberOfChars(bytes.NewReader(inputData))
	} else {
		count, err = file_utilities.NumberOfChars(file)
	}

	if err != nil {
		return "", fmt.Errorf("error counting characters: %v", err)
	}
	return fmt.Sprintf("%d characters", count), nil
}
