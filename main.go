package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	sheetId := os.Getenv("SHEET_ID")
	creds := os.Getenv("SA_JSON")
	// creds, err := os.ReadFile("client_secret.json")
	/*if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}*/
	sheetName := os.Getenv("DATASHEET_NAME")
	sheetRange := os.Getenv("DATASHEET_RANGE")
	var dataRange = sheetName + "!" + sheetRange
	fmt.Printf("Helo\n")
	fmt.Printf("Helllo %s %s\n", sheetId, dataRange)
	resp, err := readSheetData(sheetId, dataRange, []byte(creds))
	log.Printf("%s %s", resp, err)

	sheetName = "data-store"
	sheetRange = "A1:B2"
	dataRange = sheetName + "!" + sheetRange
	data := [][]interface{}{{"das", "asd"}, {"2", "dsaasd"}}
	writeResp, err := writeSheetData(sheetId, dataRange, creds, data)
	fmt.Printf("%s\n", writeResp)

	payload := map[string]string{
		"username": "michael",
		"password": "success-password",
	}
	json, err := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json"}
	status, err := postCall("https://json-placeholder.mock.beeceptor.com/login", []byte(json), headers)
	fmt.Printf("%d\n", status)
}
