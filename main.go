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

const (
	USD Currency = "USD"
	INR Currency = "INR"
)
const (
	Yahoo DataSource = "YAHOO"
)

type Activity struct {
	Currency   Currency     `json:"currency"`
	DataSource DataSource   `json:"dataSource"`
	Date       string       `json:"date"`
	Fee        float64      `json:"fee"`
	Quantity   float64      `json:"quantity"`
	Symbol     string       `json:"symbol"`
	Type       ActivityType `json:"type"`
	UnitPrice  float64      `json:"unitPrice"`
	AccountID  string       `json:"accountId"`
	Comment    interface{}  `json:"comment,omitempty"` // Optional comment field
	Tags       []string     `json:"tags"`
}

type Payload struct {
	Activities []Activity `json:"activities"`
}

/* type EquityTransaction interface {                       //TODO: Implement interface or modules
    CreateActivity() Activity
} */

func main() {
	sheetId := os.Getenv("SHEET_ID")
	creds := os.Getenv("SA_JSON")
	accountId := os.Getenv("GH_ACC_ID")
	sheetName := os.Getenv("DATASHEET_NAME")
	equityType := os.Getenv("EQUITYTYPE") //TODO: Move out of here
	dataStoreSheetName := os.Getenv("DATASTORE")
	sheetRange, err := readProgressData()
	// creds, err := os.ReadFile("client_secret.json")                      //This is to emulate without env variables
	/*if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}*/
	// sheetRange := os.Getenv("DATASHEET_RANGE")                           //Emulate a dynamic sheet range
	fmt.Printf("Sheet Range %d %s\n", sheetRange, err)
	dataRange := sheetName + "!" + "A" + strconv.Itoa(sheetRange) + ":I"
	fmt.Printf("Data Range: %s\n", dataRange)
	fmt.Printf("Helo\n")
	fmt.Printf("Helllo %s %s\n", sheetId, dataRange)
	resp, err := readSheetData(sheetId, dataRange, []byte(creds))
	log.Printf("%s %s\n", resp, err)
	startPoint, err := readProgressData(dataStoreSheetName)

	for index, row := range resp.Values {
		// 	    transaction EquityTransaction;
		var activity Activity

		switch equityType {
		case "IND":
			activity = CreateIndEqActivity(row, accountId)
		case "US":
			activity = CreateUSEqActivity(row, accountId)
		case "MF", "ELSS", "US-MF": // Arrange them all in MF by passing accounts using maps, Is it a good pattern?
			activity = CreateMFActivity(row, accountId)
		default:
			activity = CreateIndEqActivity(row, accountId)
		}

		payload := Payload{
			Activities: []Activity{
				activity,
				// Add more activity objects here if needed
			},
		}
		fmt.Print(payload)
		status := createGhostfolioEntry(payload)
		if status != 201 {
			fmt.Printf("Status: %d", status)
			break
		}
		fmt.Printf("Status: %d", status)
		writeProgressData([]interface{}{startPoint + index + 1, status}, dataStoreSheetName)
	}
}

func writeProgressData(data []interface{}, dataStoreTarget string) (string, error) {
	sheetName := dataStoreTarget
	sheetRange := "A1:B2"
	dataRange := sheetName + "!" + sheetRange
	creds := os.Getenv("SA_JSON")
	sheetId := os.Getenv("SHEET_ID")
	sheetData := [][]interface{}{{"Entry No for Error", "Error"}, data}
	writeResp, err := writeSheetData(sheetId, dataRange, []byte(creds), sheetData)
	fmt.Printf("%s\n", writeResp)
	return writeResp, err
}

func readProgressData(dataStoreSource string) (int, error) {
	sheetName := dataStoreSource
	sheetRange := "A2:B2"
	dataRange := sheetName + "!" + sheetRange
	creds := os.Getenv("SA_JSON")
	sheetId := os.Getenv("SHEET_ID")
	readResp, err := readSheetData(sheetId, dataRange, []byte(creds))
	fmt.Printf("%s\n", readResp)
	readValue, _ := strconv.Atoi(readResp.Values[0][0].(string))
	return readValue, err
}

func createGhostfolioEntry(payload Payload) int {
	log.Printf("Payload : %s\n", payload)
	json, err := json.Marshal(payload)
	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + os.Getenv("API_JWT")}
	status, err := postCall("http://ghostfolio.ghostfolio.svc.cluster.local:3333/api/v1/import", []byte(json), headers)
	fmt.Printf("%d   %s\n", status, err)
	return status
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
