package file_utilities

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CheckForError(err error) bool {
	return err != nil
}

func GetFullFilePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func LoadFile(filePath string) (*os.File, error) {
	fullFilePath, err := GetFullFilePath("data/" + filePath)
	if CheckForError(err) {
		return nil, err
	}
	filePointer, err := os.OpenFile(fullFilePath, os.O_RDONLY, 0444)
	if CheckForError(err) {
		return nil, err
	}
	return filePointer, nil
}

func NumberOfBytesCommand(inputSource interface{}) (int, error) {
	var bytesCounter int
	var sourceReader io.Reader
	var filePointer *os.File

	switch v := inputSource.(type) {
	case string:
		filePointer, err := LoadFile(v)
		if CheckForError(err) {
			return 0, err
		}
		sourceReader = filePointer
	case io.Reader:
		sourceReader = v
	default:
		return 0, fmt.Errorf("invalid input type")
	}

	scanner := bufio.NewScanner(sourceReader)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		bytesCounter++
	}
	filePointer.Close()
	if err := scanner.Err(); CheckForError(err) {
		return bytesCounter, err
	}
	return bytesCounter, nil
}

func NumberOfLinesCommand(inputSource interface{}) (int, error) {
	var linesCounter int
	var sourceReader io.Reader
	var filePointer *os.File

	switch v := inputSource.(type) {
	case string:
		filePointer, err := LoadFile(v)
		if CheckForError(err) {
			return 0, err
		}
		sourceReader = filePointer
	case io.Reader:
		sourceReader = v
	default:
		return 0, fmt.Errorf("invalid input type")
	}

	scanner := bufio.NewScanner(sourceReader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		linesCounter++
	}
	filePointer.Close()
	if err := scanner.Err(); CheckForError(err) {
		return linesCounter, err
	}
	return linesCounter, nil

}

func NumberOfWordsCommand(inputSource interface{}) (int, error) {
	var wordCounter int
	var sourceReader io.Reader
	var filePointer *os.File

	switch v := inputSource.(type) {
	case string:
		filePointer, err := LoadFile(v)
		if CheckForError(err) {
			return 0, err
		}
		sourceReader = filePointer
	case io.Reader:
		sourceReader = v
	default:
		return 0, fmt.Errorf("invalid input type")
	}

	scanner := bufio.NewScanner(sourceReader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCounter++
	}
	filePointer.Close()

	if err := scanner.Err(); CheckForError(err) {
		return wordCounter, err
	}

	return wordCounter, nil
}

func NumberOfChars(inputSource interface{}) (int, error) {
	var charCounter int
	var sourceReader io.Reader
	var filePointer *os.File

	switch v := inputSource.(type) {
	case string:
		filePointer, err := LoadFile(v)
		if CheckForError(err) {
			return 0, err
		}
		sourceReader = filePointer
	case io.Reader:
		sourceReader = v
	default:
		return 0, fmt.Errorf("invalid input type")
	}

	scanner := bufio.NewScanner(sourceReader)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		charCounter++
	}
	filePointer.Close()

	if err := scanner.Err(); CheckForError(err) {
		return charCounter, err
	}

	return charCounter, nil
}

func DeclareValidJson(filePath string) (bool, error) {
	fileContent, err := LoadFile(filePath)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(fileContent)
	scanner.Split(bufio.ScanRunes)

	var CLOSING, OPENING byte = 125, 123 // Byte ref from ASCII
	parenthesesStack := NewStack[byte]()
	var iterationFlag bool = false
	for scanner.Scan() {
		iterationFlag = true
		currentByte := scanner.Bytes()[0]
		if parenthesesStack.isEmpty() && currentByte == CLOSING {
			return false, fmt.Errorf("invalid json")
		} else if topVal, _ := parenthesesStack.Top(); topVal == OPENING && currentByte == CLOSING {
			parenthesesStack.Pop()
		} else if currentByte == OPENING {
			parenthesesStack.Push(scanner.Bytes()[0])
		}
	}
	fileContent.Close()
	if iterationFlag {
		return parenthesesStack.isEmpty(), nil
	}
	return false, fmt.Errorf("invalid json")
}

func ParseJson(filePath string) (map[string]interface{}, error) {
	fileContent, err := LoadFile(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fileContent)
	scanner.Split(bufio.ScanRunes)
	var (
		state        string = "start" // start, key, value, colon, end
		currentKey   []rune
		currentValue []rune
		result       = make(map[string]interface{})
	)

	for scanner.Scan() {
		currentChar := scanner.Text()
		switch state {
		case "start":
			if currentChar == "{" {
				state = "key"
			}
		case "key":
			if currentChar == "\"" {
				if len(currentKey) > 0 {
					state = "colon"
				}
			} else if currentChar != " " && currentChar != "\n" {
				currentKey = append(currentKey, rune(currentChar[0]))
			}

		case "colon":
			if currentChar == ":" {
				state = "value"
			} else if currentChar != " " && currentChar != "\n" {
				return nil, fmt.Errorf("expected ':' after key, got '%s'", currentChar)
			}

		case "value":
			if currentChar == "\"" {
				if len(currentValue) > 0 {
					// End of value
					result[string(currentKey)] = string(currentValue)
					currentKey = nil
					currentValue = nil
					state = "end"
				}
			} else if currentChar != " " && currentChar != "\n" {
				currentValue = append(currentValue, rune(currentChar[0]))
			}

		case "end":
			if currentChar == "," {
				// Another key-value pair expected
				state = "key"
			} else if currentChar == "}" {
				// End of JSON
				fileContent.Close()
				return result, nil
			} else if currentChar != " " && currentChar != "\n" {
				return nil, fmt.Errorf("unexpected currentCharacter '%s' after value", currentChar)
			}
		}
	}

	fileContent.Close()
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return result, nil
}
