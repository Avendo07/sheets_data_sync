/*
type TransactionType string

	const {
	    BUY TransactionType
	    SELL TransactionType
	}

	type struct USEquity{
	    ticker string
	    timeStamp string
	    transactionType TransactionType
	    quantity float64
	    unitPrice float64
	}

func (transaction USEquity) CreateActivity() Activity{

}
*/
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func CreateUSEqActivity(sheetRow []interface{}, accountId string) Activity {
	ticker := sheetRow[1].(string)
	date, err := isoTimeStamp(sheetRow[0].(string))
	unitPrice, err := strconv.ParseFloat(sheetRow[4].(string), 64)
	quantity, err := strconv.ParseFloat(sheetRow[5].(string), 64) //TODO: This is to maintain compaitability with US Eq
	action := getUSAction(sheetRow[6].(string))
	fee, err := strconv.ParseFloat(sheetRow[7].(string), 64)

	log.Printf("quant price date err: %f %f %f %s %s", quantity, unitPrice, fee, date, err)

	return Activity{
		Currency:   USD,
		DataSource: Yahoo,
		Date:       date,
		Fee:        fee,
		Quantity:   quantity,
		Symbol:     ticker,
		Type:       action,
		UnitPrice:  unitPrice,
		AccountID:  accountId,
		Tags:       []string{ticker},
		Comment:    nil,
	}
}

func isoTimeStamp(date string) (string, error) {
	layout := "02 JAN 2006, 03:04 AM" // YYYY-MM-DD format
	parsedDate, err := time.Parse(layout, date)

	// Handle potential parsing errors
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", err
	}

	// Format the parsed date into ISO8601 format
	isoFormattedDate := parsedDate.Format(time.RFC3339)
	return isoFormattedDate, nil
}

func getUSAction(action string) ActivityType {
	if strings.ToLower(action) != "buy" {
		return Sell
	} else {
		return Buy
	}
}
