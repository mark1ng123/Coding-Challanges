package main

import (
	"fmt"
	"log"
	file_utilities "utilities"
)

func main() {
	filepath := "step2/valid2.json"
	isValid, err := file_utilities.DeclareValidJson(filepath)
	if err != nil {
		log.Fatalf("Error in json: %v", err)
	}
	log.Printf("Json validation func returned %v", isValid)

	resultMap, err := file_utilities.ParseJson(filepath)
	if err != nil {
		log.Fatalf("Error in json parsing: %v", err)
	}
	for key, value := range resultMap {
		fmt.Println(key, value)
	}
	fmt.Println(resultMap)
}
