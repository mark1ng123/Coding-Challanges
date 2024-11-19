package file_utilities

import (
	"bufio"
	"errors"
	"fmt"
)

func NewJSONParser(scanner *bufio.Scanner) *JSONParser {
	return &JSONParser{
		state:   "start",
		result:  make(map[string]interface{}),
		scanner: scanner,
	}
}
func (p *JSONParser) JsonParse() (map[string]interface{}, error) {
	for p.scanner.Scan() {
		currentChar := p.scanner.Text()
		if err := p.handleState(currentChar); err != nil {
			return nil, err
		}
	}

	if len(p.currentKey) > 0 || len(p.currentValue) > 0 {
		return nil, errors.New("incomplete JSON structure")
	}
	return p.result, p.scanner.Err()
}

func (p *JSONParser) resetForNextPair() {
	p.currentKey = nil
	p.currentValue = nil
	p.state = "end"
}

func (p *JSONParser) handleState(currentChar string) error {
	switch p.state {
	case "start":
		if currentChar == "{" {
			p.state = "key"
		}
	case "key":
		if currentChar == "\"" {
			p.numberOfDoubleQuotes++
			if len(p.currentKey) > 0 {
				p.state = "colon"
			}
		} else if currentChar != " " && currentChar != "\n" {
			p.currentKey = append(p.currentKey, rune(currentChar[0]))
		}
	case "colon":
		if p.numberOfDoubleQuotes != 2 {
			return errors.New("mismatched quotes in key")
		} else if currentChar == ":" {
			p.state = "value"
			p.numberOfDoubleQuotes = 0
		} else if currentChar != " " && currentChar != "\n" {
			return fmt.Errorf("expected ':' after key, got '%s'", currentChar)
		}
	case "value":
		if currentChar == "\"" {
			p.numberOfDoubleQuotes++
			if len(p.currentValue) > 0 {
				p.result[string(p.currentKey)] = string(p.currentValue)
				p.resetForNextPair()
			}
		} else if currentChar != " " && currentChar != "\n" {
			p.currentValue = append(p.currentValue, rune(currentChar[0]))
		}
	case "end":
		if currentChar == "," {
			p.state = "key"
			p.numberOfDoubleQuotes = 0
		} else if currentChar == "}" {
			return nil
		} else if currentChar != " " && currentChar != "\n" {
			return fmt.Errorf("unexpected character '%s' after value", currentChar)
		}
	}
	return nil
}
