package main

import (
	"log"
	file_utilities "utilities"
)

func main() {
	filepath := "step1/valid.json"
	isValid, err := file_utilities.DeclareValidJson(filepath)
	if err != nil {
		log.Fatalf("Error in json: %v", err)
	}
	log.Printf("Json validation func returned %v", isValid)
}
