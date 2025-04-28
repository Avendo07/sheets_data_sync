package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

/*
type struct Exchange

const {
    Exchange NSE
    Exchange BSE
}

type struct USEquity{
    ticker string
    date string
    exchange Exchange
    buyQuantity int64
    sellQuantity int64
    unitPrice float64
}
*/

type ActivityType string

const (
	Buy  ActivityType = "BUY"
	Sell ActivityType = "SELL"
)

func CreateMFActivity(sheetRow []interface{}, accountId string) Activity {
	mf_name, _ := sheetRow[1].(string)
// 	mkt, _ := sheetRow[2].(string)

	nav_date, err := isoDate(sheetRow[0].(string))
	nav, err := strconv.ParseFloat(sheetRow[2].(string), 64)
	quantity, err := strconv.ParseFloat(sheetRow[3].(string), 64) //TODO: This is to maintain compaitability with US Eq
    action := getUSAction(sheetRow[4].(string))
    fee, err := strconv.ParseFloat(sheetRow[5].(string), 64)

	log.Printf("quant price date err: %f %f %f %s %s\n", quantity, unitPrice, fee, date, err)

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

func getAction(buyQty float64, sellQty float64) (ActivityType, float64) {
	if buyQty != 0 && sellQty == 0 {
		return Buy, buyQty
	} else {
		return Sell, sellQty
	}
}

func isoDate(date string) (string, error) {
	layout := "02-01-2006"

	istLocation, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return "", err
	}

	parsedDate, err := time.ParseInLocation(layout, date, istLocation)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", err
	}
	// Format the parsed date into ISO8601 format
	isoFormattedDate := parsedDate.Format(time.RFC3339)
	return isoFormattedDate, nil
}
