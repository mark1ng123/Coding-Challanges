package utilities

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
