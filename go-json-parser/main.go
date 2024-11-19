package main

import (
	"bufio"
	"fmt"
	"log"
	file_utilities "utilities"
)

func main() {
	filepath := "step2/invalid2.json"
	isValid, err := file_utilities.DeclareValidJson(filepath)
	if err != nil {
		log.Fatalf("Error in json: %v", err)
	}
	log.Printf("Json validation func returned %v", isValid)
	loadedFile, err := file_utilities.LoadFile(filepath)
	if err != nil {
		log.Fatalf("Error in file loading: %v", err)
	}
	scanner := bufio.NewScanner(loadedFile)
	scanner.Split(bufio.ScanRunes)
	jsonParser := file_utilities.NewJSONParser(scanner)
	resultMap, err := jsonParser.JsonParse()
	if err != nil {
		log.Fatalf("Error in JSON parsing: %v", err)
	}
	loadedFile.Close()
	// for key, value := range resultMap {
	// 	fmt.Println(key, value)
	// }
	fmt.Println(resultMap)
}
