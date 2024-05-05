package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	sheetId := os.Getenv("SHEET_ID")
	// creds := os.Getenv("SA_JSON")
	creds, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}
	sheetName := os.Getenv("DATASHEET_NAME")
	sheetRange := os.Getenv("DATASHEET_RANGE")
	var dataRange = sheetName + "!" + sheetRange
	fmt.Printf("Helo\n")
	fmt.Printf("Helllo %s %s\n", sheetId, dataRange)
	resp, err := readSheetData(sheetId, dataRange, creds)
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

func readSheetData(sheetId string, dataRange string, creds []byte) (any, error) {
	service, err := getSheetRef(creds)
	if err != nil {
		return nil, err
	}
	resp, err := service.Spreadsheets.Values.Get(sheetId, dataRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
		return nil, err
	}
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			fmt.Printf("%s\n", row)
			// fmt.Printf("%s %s", row[0], row[2])

			// response, err := createPostRequest(ghostfolioIp, row)
			// fmt.Printf("%s, %s\n", response.Status, err)
			// postReq()
		}
	}
	return resp, nil
}

func writeSheetData(sheetId string, dataRange string, creds []byte, rowData [][]interface{}) (string, error) {
	service, err := getSheetRef(creds)
	if err != nil {
		return "", err
	}
	values := sheets.ValueRange{Values: rowData}
	resp, err := service.Spreadsheets.Values.Update(sheetId, dataRange, &values).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to add data to google sheets service: %v", err)
		return "", err
	}
	fmt.Printf(resp.UpdatedRange)
	return resp.UpdatedRange, nil
}

func getSheetRef(creds []byte) (*sheets.Service, error) {
	var sheetsService *sheets.Service
	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to create JWT config: %v", err)
		return nil, err
	}
	client := config.Client(context.Background())
	sheetsService, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
		return nil, err
	}
	return sheetsService, err
}

func postCall(url string, payload []byte, headers map[string]string) (int, error) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}
