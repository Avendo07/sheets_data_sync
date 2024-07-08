package main

import (
	"log"
	"strconv"
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

func CreateIndEqActivity(sheetRow []interface{}, accountId string) Activity {
	company, _ := sheetRow[1].(string)
	mkt, _ := sheetRow[2].(string)

	date, err := isoDate(sheetRow[0].(string))
	unitPrice, err := strconv.ParseFloat(sheetRow[4].(string), 64)
	buyQuantity, err := strconv.ParseFloat(sheetRow[5].(string), 64) //TODO: This is to maintain compaitability with US Eq
	sellQuantity, err := strconv.ParseFloat(sheetRow[6].(string), 64)
	ticker, currency := getTicker(company, mkt)
	action, qty := getAction(buyQuantity, sellQuantity)

	log.Printf("quant price date err: %d %s %s %s", qty, unitPrice, date, err)

	return Activity{
		Currency:   currency,
		DataSource: Yahoo,
		Date:       date,
		Fee:        0,
		Quantity:   qty,
		Symbol:     ticker,
		Type:       action,
		UnitPrice:  unitPrice,
		AccountID:  accountId,
		Tags:       []string{company},
		Comment:    nil,
	}
}

func getTicker(company string, mkt string) (string, Currency) {
	var suffixes string
	var currency Currency
	switch {
	case mkt == "NSE":
		suffixes = ".NS"
		currency = INR
	case mkt == "BSE":
		suffixes = ".BO"
		currency = INR
	default:
		currency = USD
	}
	ticker := company + suffixes
	return ticker, currency
}

func getAction(buyQty float64, sellQty float64) (ActivityType, float64) {
	if buyQty != 0 && sellQty == 0 {
		return Buy, buyQty
	} else {
		return Sell, sellQty
	}
}
