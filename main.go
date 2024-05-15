package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
	// entries, err := mapDataToPayload(resp.Values)
	quant, err := strconv.Atoi(resp.Values[0][5].(string))
	price, err := strconv.ParseFloat(resp.Values[0][4].(string), 64)
	log.Printf("quant price err: %s %s %s", quant, price, err)

	sheetName = "data-store"
	sheetRange = "A1:B2"
	dataRange = sheetName + "!" + sheetRange
	data := [][]interface{}{{"das", "update"}, {"2", "dsaasd"}}
	writeResp, err := writeSheetData(sheetId, dataRange, []byte(creds), data)
	fmt.Printf("%s\n", writeResp)

	payload := Payload{
		Activities: []Activity{
			{
				Currency:   INR,
				DataSource: Yahoo,
				Date:       "2023-09-17T00:00:00.000Z",
				Fee:        0,
				Quantity:   quant,
				Symbol:     resp.Values[0][0].(string) + ".NS",
				Type:       Buy,
				UnitPrice:  price,
				AccountID:  "4fe741a5-88e2-4c67-9431-8727274387c8",
				Comment:    nil,
			},
			// Add more activity objects here if needed
		},
	}
	json, err := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + os.Getenv("API_JWT")}
	status, err := postCall("http://ghostfolio.ghostfolio.svc.cluster.local:3333/api/v1/import", []byte(json), headers)
	fmt.Printf("%d\n", status)
}

/*func mapDataToPayload(data [][]interface{}) (Payload, error) {
	var payload Payload
	payload.Activities = make([]Activity, 0) // Initialize empty slice

	for _, row := range data {
		if len(row) != reflect.TypeOf(Activity{}).Elem().NumField() {
			return payload, fmt.Errorf("Invalid row length: expected %d, got %d", reflect.TypeOf(Activity{}).Elem().NumField(), len(row))
		}

		activity := Activity{}
		for i, element := range row {
			switch field := reflect.TypeOf(activity).Elem().Field(i); field.Type.Kind() {
			case reflect.String:
				reflect.ValueOf(&activity).Elem().Field(i).SetString(element.(string))
			case reflect.Int:
				reflect.ValueOf(&activity).Elem().Field(i).SetInt(element.(int64)) // Use int64 for wider range
			case reflect.Float64:
				reflect.ValueOf(&activity).Elem().Field(i).SetFloat(element.(float64))
			default:
				return payload, fmt.Errorf("Unexpected type for field %s: %v", field.Name, element)
			}
		}
		payload.Activities = append(payload.Activities, activity)
	}
	return payload, nil
}*/

/*func createGhostfolioEntry(ticker string, date string, transType string, quantity int, unitPrice float32) {
	baseUrl := "ghostfolio.ghostfolio.svc.cluster.local:3333/api/v1/import"

}*/
