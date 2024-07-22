package main

import (
	"io"
	"os"
	"log"
	"encoding/json"
)

func ReadJson() JsonConstsFile {
	var jsonConsts JsonConstsFile

	file, err := os.Open("JSON/JsonConsts.json")
	if err != nil {
		log.Fatalf("Something went wrong: \n\n%s", err)
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Something went wrong: \n\n%s", err)
	}

	json.Unmarshal(fileContent, &jsonConsts)

	return jsonConsts
}
