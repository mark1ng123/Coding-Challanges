package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func NewJSONParser(scanner *bufio.Scanner) *JSONParser {
	return &JSONParser{
		state:        "start",
		result:       make(map[string]interface{}),
		scanner:      scanner,
		nestingLevel: 0,
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

func isArray(currValue string) bool {
	return len(currValue) > 1 && currValue[0] == '[' && currValue[len(currValue)-1] == ']'
}

func isDict(currValue string) bool {
	return len(currValue) > 1 && currValue[0] == '{' && currValue[len(currValue)-1] == '}'
}

func parseArray(currValue string) ([]interface{}, error) {
	var resultArray []interface{}
	if len(currValue) == 0 {
		return resultArray, nil
	}
	arrayElements := strings.Split(currValue, ",")
	var firstType reflect.Type
	for elementIdx, element := range arrayElements {
		element = strings.TrimSpace(element)
		valueWithType, err := setRightTypeForValue(element)
		if err != nil {
			return nil, fmt.Errorf("error parsing array element at index %d: %v", elementIdx, err)
		}
		if elementIdx == 0 {
			firstType = reflect.TypeOf(valueWithType)
		} else {
			if reflect.TypeOf(valueWithType) != firstType {
				return nil, fmt.Errorf("type mismatch in array: expected %s, got %s", firstType, reflect.TypeOf(valueWithType))
			}
		}
		resultArray = append(resultArray, valueWithType)
	}
	return resultArray, nil

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
	case isArray(currValue):
		return parseArray(currValue[1 : len(currValue)-1])
	case isDict(currValue):
		// something for it like our own json parser
		nestedScanner := bufio.NewScanner(strings.NewReader(currValue))
		nestedParser := NewJSONParser(nestedScanner)
		parsedDict, err := nestedParser.JsonParse()
		if err != nil {
			return nil, fmt.Errorf("error parsing nested dictionary: %v", err)
		}
		castedVal = parsedDict
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
		if currentChar == "[" || currentChar == "{" {
			p.nestingLevel++
			p.currentValue = append(p.currentValue, rune(currentChar[0]))
		} else if currentChar == "]" || (currentChar == "}" && p.nestingLevel > 0) {
			if p.nestingLevel > 0 {
				p.nestingLevel--
				p.currentValue = append(p.currentValue, rune(currentChar[0]))
			}
		} else if (currentChar == "," || currentChar == "}") && p.nestingLevel == 0 {
			if len(p.currentValue) > 0 {
				value, err := setRightTypeForValue(string(p.currentValue))
				fmt.Println(value)
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
