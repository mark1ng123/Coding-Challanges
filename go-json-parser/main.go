package main

import (
	"bufio"
	"fmt"
	"log"
	"utilities"
)

func main() {
	filepath := "step2/invalid2.json"
	isValid, err := utilities.DeclareValidJson(filepath)
	if err != nil {
		log.Fatalf("Error in json: %v", err)
	}
	log.Printf("Json validation func returned %v", isValid)
	loadedFile, err := utilities.LoadFile(filepath)
	if err != nil {
		log.Fatalf("Error in file loading: %v", err)
	}
	scanner := bufio.NewScanner(loadedFile)
	scanner.Split(bufio.ScanRunes)
	jsonParser := utilities.NewJSONParser(scanner)
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
