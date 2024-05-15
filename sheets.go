package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func readSheetData(sheetId string, dataRange string, creds []byte) (*sheets.ValueRange, error) {
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
