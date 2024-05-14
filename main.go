package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Currency string
type DataSource string
type ActivityType string

const (
	USD Currency = "USD"
	INR Currency = "INR"
)
const (
	Yahoo DataSource = "YAHOO"
)
const (
	Buy  ActivityType = "BUY"
	SELL ActivityType = "SELL"
)

type Activity struct {
	Currency   Currency     `json:"currency"`
	DataSource DataSource   `json:"dataSource"`
	Date       string       `json:"date"`
	Fee        int          `json:"fee"`
	Quantity   int          `json:"quantity"`
	Symbol     string       `json:"symbol"`
	Type       ActivityType `json:"type"`
	UnitPrice  float64      `json:"unitPrice"`
	AccountID  string       `json:"accountId"`
	Comment    interface{}  `json:"comment,omitempty"` // Optional comment field
}

type Payload struct {
	Activities []Activity `json:"activities"`
}

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
	writeResp, err := writeSheetData(sheetId, dataRange, []byte(creds), data)
	fmt.Printf("%s\n", writeResp)

	payload := Payload{
		Activities: []Activity{
			{
				Currency:   USD,
				DataSource: Yahoo,
				Date:       "2023-09-17T00:00:00.000Z",
				Fee:        19,
				Quantity:   5,
				Symbol:     "MSFT",
				Type:       Buy,
				UnitPrice:  298.58,
				AccountID:  "4fe741a5-88e2-4c67-9431-8727274387c8",
				Comment:    nil,
			},
			// Add more activity objects here if needed
		},
	}
	json, err := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + os.Getenv("API_JWT")}
	status, err := postCall("ghostfolio.ghostfolio.svc.cluster.local:3333/api/v1/import", []byte(json), headers)
	fmt.Printf("%d\n", status)
}

/*func createGhostfolioEntry(ticker string, date string, transType string, quantity int, unitPrice float32) {
	baseUrl := "ghostfolio.ghostfolio.svc.cluster.local:3333/api/v1/import"

}*/
