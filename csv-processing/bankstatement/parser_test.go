package bankstatement_test

import (
	"strings"
	"testing"

	"github.com/fikryfahrezy/forward/csv-processing/bankstatement"
)

func TestParseCSV_AllLines_Success(t *testing.T) {
	text := "timestamp,name,type,amount,status,description"
	text += "\n1624507883,JOHN DOE,DEBIT,250000,SUCCESS,restaurant"
	text += "\n1624608050,E-COMMERCE A,DEBIT,150000,FAILED,clothes"
	text += "\n1624512883,COMPANY A,CREDIT,12000000,SUCCESS,salary"
	text += "\n1624615065,E-COMMERCE B,DEBIT,150000,PENDING,clothes"

	transactions, err := bankstatement.ParseCSV(strings.NewReader(text))

	if len(transactions) != 4 {
		t.Fatal("transaction length not match")
	}
	if len(err) != 0 {
		t.Fatal("error length not match")
	}

	actualFirstTransaction := transactions[0]
	if actualFirstTransaction.Timestamp != int64(1624507883) {
		t.Fatal("first transaction timestamp not match")
	}
	if actualFirstTransaction.Name != "JOHN DOE" {
		t.Fatal("first transaction name not match")
	}
	if actualFirstTransaction.Type != bankstatement.TransactionTypeDebit {
		t.Fatal("first transaction type not match")
	}
	if actualFirstTransaction.Amount != float64(250000) {
		t.Fatal("first transaction amount not match")
	}
	if actualFirstTransaction.Status != bankstatement.TransactionStatusSuccess {
		t.Fatal("first transaction status not match")
	}
	if actualFirstTransaction.Description != "restaurant" {
		t.Fatal("first transaction description not match")
	}

	actualSecondTransaction := transactions[1]
	if actualSecondTransaction.Timestamp != int64(1624608050) {
		t.Fatal("second transaction timestamp not match")
	}
	if actualSecondTransaction.Name != "E-COMMERCE A" {
		t.Fatal("second transaction name not match")
	}
	if actualSecondTransaction.Type != bankstatement.TransactionTypeDebit {
		t.Fatal("second transaction type not match")
	}
	if actualSecondTransaction.Amount != float64(150000) {
		t.Fatal("second transaction amount not match")
	}
	if actualSecondTransaction.Status != bankstatement.TransactionStatusFailed {
		t.Fatal("second transaction status not match")
	}
	if actualSecondTransaction.Description != "clothes" {
		t.Fatal("second transaction description not match")
	}

	actualThirdTransaction := transactions[2]
	if actualThirdTransaction.Timestamp != int64(1624512883) {
		t.Fatal("third transaction timestamp not match")
	}
	if actualThirdTransaction.Name != "COMPANY A" {
		t.Fatal("third transaction name not match")
	}
	if actualThirdTransaction.Type != bankstatement.TransactionTypeCredit {
		t.Fatal("third transaction type not match")
	}
	if actualThirdTransaction.Amount != float64(12000000) {
		t.Fatal("third transaction amount not match")
	}
	if actualThirdTransaction.Status != bankstatement.TransactionStatusSuccess {
		t.Fatal("third transaction status not match")
	}
	if actualThirdTransaction.Description != "salary" {
		t.Fatal("third transaction description not match")
	}

	actualFourthTransaction := transactions[3]
	if actualFourthTransaction.Timestamp != int64(1624615065) {
		t.Fatal("forth transaction timestamp not match")
	}
	if actualFourthTransaction.Name != "E-COMMERCE B" {
		t.Fatal("forth transaction name not match")
	}
	if actualFourthTransaction.Type != bankstatement.TransactionTypeDebit {
		t.Fatal("forth transaction type not match")
	}
	if actualFourthTransaction.Amount != float64(150000) {
		t.Fatal("forth transaction amount not match")
	}
	if actualFourthTransaction.Status != bankstatement.TransactionStatusPending {
		t.Fatal("forth transaction status not match")
	}
	if actualFourthTransaction.Description != "clothes" {
		t.Fatal("forth transaction description not match")
	}
}

func TestParseCSV_AllLines_Failed(t *testing.T) {
	text := "timestamp,name,type,amount,status,description"
	text += "\n1624507883,JOHN DOE,UNKNOWN,250000,SUCCESS,restaurant"
	text += "\nWIB,E-COMMERCE A,DEBIT,150000,FAILED,clothes"
	text += "\n1624512883,COMPANY A,CREDIT,DUA-BELAS,SUCCESS,salary"
	text += "\n1624615065,E-COMMERCE B,DEBIT,150000,UNKNOWN,clothes"
	text += "\n1624615066,,DEBIT,150000,SUCCESS,clothes"
	text += "\n1624615066,E-COMMERCE C,DEBIT,150000,SUCCESS,"

	transactions, err := bankstatement.ParseCSV(strings.NewReader(text))
	if len(transactions) != 0 {
		t.Fatal("transaction length not match")
	}
	if len(err) != 6 {
		t.Fatal("error length not match")
	}

	lineOneErrorMessages, ok := err["line[2]"].([]string)
	if !ok {
		t.Fatal("the first error key doesnt exists")
	}
	if lineOneErrorMessages[0] != "the posible value for 'type' is DEBIT, CREDIT" {
		t.Fatal("the frist error message not match")
	}

	lineTwoErrorMessages, ok := err["line[3]"].([]string)
	if !ok {
		t.Fatal("the second error key doesnt exists")
	}
	if lineTwoErrorMessages[0] != "invalid 'timestamp', expected to be integer" {
		t.Fatal("the second error message not match")
	}

	lineThreeErrorMessages, ok := err["line[4]"].([]string)
	if !ok {
		t.Fatal("the third error key doesnt exists")
	}
	if lineThreeErrorMessages[0] != "invalid 'amount', expected to be number" {
		t.Fatal("the third error message not match")
	}

	lineFourErrorMessages, ok := err["line[5]"].([]string)
	if !ok {
		t.Fatal("the forth error key doesnt exists")
	}
	if lineFourErrorMessages[0] != "the posible value for 'status' is SUCCESS, PENDING, FAILED" {
		t.Fatal("the forth error message not match")
	}

	lineFiveErrorMessages, ok := err["line[6]"].([]string)
	if !ok {
		t.Fatal("the fifth error key doesnt exists")
	}
	if lineFiveErrorMessages[0] != "the 'name' is required" {
		t.Fatal("the fifth error message not match")
	}

	lineSixErrorMessages, ok := err["line[7]"].([]string)
	if !ok {
		t.Fatal("the sixth error key doesnt exists")
	}
	if lineSixErrorMessages[0] != "the 'description' is required" {
		t.Fatal("the sixth error message not match")
	}
}
