package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
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

	if (len(p.currentKey) > 0 && len(p.currentValue) == 0 && len(p.result) > 0) || (len(p.currentKey) == 0 && len(p.currentValue) > 0 && len(p.result) > 0) {
		return nil, errors.New("incomplete JSON structure")
	}
	return p.result, p.scanner.Err()
}

func (p *JSONParser) resetForNextPair() {
	p.currentKey = nil
	p.currentValue = nil
	p.state = "end"
}

func isInteger(currValue string) bool {
	_, intErr := strconv.Atoi(currValue)
	return intErr == nil
}

func isFloat(currValue string) bool {
	_, floatErr := strconv.ParseFloat(currValue, 64)
	return floatErr == nil
}

func isNIL(currValue string) bool {
	return currValue == "null"
}

func setRightTypeForValue(currValue string) (interface{}, error) {
	var castedVal interface{}
	switch {
	case currValue == "true" || currValue == "false": // Not using ParseBool because it accepts 1,0,TRUE,False and its not allowed in json
		castedVal = currValue == "true"
		return castedVal, nil
	case isInteger(currValue):
		num, _ := strconv.Atoi(currValue)
		castedVal = num
		return castedVal, nil
	case isFloat(currValue):
		num, _ := strconv.ParseFloat(currValue, 64)
		castedVal = num
		return castedVal, nil
	case isNIL(currValue):
		castedVal = nil
		return castedVal, nil
	default:
		if len(currValue) > 1 && currValue[0] == '"' && currValue[len(currValue)-1] == '"' {
			castedVal = currValue[1 : len(currValue)-1]
			return castedVal, nil
		}
		return nil, fmt.Errorf("invalid type in values")
	}
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
		if currentChar == "," || currentChar == "}" {
			if len(p.currentValue) > 0 {
				value, err := setRightTypeForValue(string(p.currentValue))
				if err != nil {
					return err
				}
				p.result[string(p.currentKey)] = value
				p.resetForNextPair()
				if currentChar == "," {
					p.state = "key"
				} else if currentChar == "}" {
					p.state = "end"
				}
			}
		} else if currentChar != " " && currentChar != "\n" {
			p.currentValue = append(p.currentValue, rune(currentChar[0]))
		}
	case "end":
		if currentChar == "}" {
			return nil
		} else if currentChar != " " && currentChar != "\n" {
			return fmt.Errorf("unexpected character '%s' after value", currentChar)
		}
	}
	return nil
}
