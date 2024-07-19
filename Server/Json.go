package main

import (
	"io"
	"os"
	"log"
	"encoding/json"
)

type JsonConstsFile struct {
	ServerIP string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	MaxBytesForReader int `json:"server_max_bytes"`
}

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
