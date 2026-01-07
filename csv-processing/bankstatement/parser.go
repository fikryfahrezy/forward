package bankstatement

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type columnName string

var (
	timestampColumn   columnName = "timestamp"
	nameColumn        columnName = "name"
	typeColumn        columnName = "type"
	amountColumn      columnName = "amount"
	statusColumn      columnName = "status"
	descriptionColumn columnName = "description"
)

type column struct {
	name   columnName
	action func(record []string, index int) (any, error)
}

var columns = []column{
	{
		name:   timestampColumn,
		action: parseTimestamp,
	},
	{
		name:   nameColumn,
		action: parseName,
	},
	{
		name:   typeColumn,
		action: parseType,
	},
	{
		name:   amountColumn,
		action: parseAmount,
	},
	{
		name:   statusColumn,
		action: parseStatus,
	},
	{
		name:   descriptionColumn,
		action: parseDescription,
	},
}

func parseTimestamp(record []string, index int) (any, error) {
	rawTimestamp := strings.TrimSpace(record[index])
	timestamp, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid 'timestamp', expected to be integer")
	}

	return timestamp, nil
}

func parseName(record []string, index int) (any, error) {
	rawName := strings.TrimSpace(record[index])
	if rawName == "" {
		return "", fmt.Errorf("the 'name' is required")
	}

	return rawName, nil
}

func parseType(record []string, index int) (any, error) {
	rawType := TransactionType(strings.TrimSpace(record[index]))
	if !rawType.Valid() {
		return TransactionTypeUnknown, fmt.Errorf("the posible value for 'type' is %s, %s", TransactionTypeDebit, TransactionTypeCredit)
	}

	return rawType, nil
}

func parseAmount(record []string, index int) (any, error) {
	rawAmount := strings.TrimSpace(record[index])
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid 'amount', expected to be number")
	}

	return amount, nil
}

func parseStatus(record []string, index int) (any, error) {
	rawStatus := TransactionStatus(strings.TrimSpace(record[index]))
	if !rawStatus.Valid() {
		return TransactionStatusUnknown, fmt.Errorf("the posible value for 'status' is %s, %s, %s", TransactionStatusSuccess, TransactionStatusPending, TransactionStatusFailed)
	}

	return rawStatus, nil
}

func parseDescription(record []string, index int) (any, error) {
	rawDescription := strings.TrimSpace(record[index])
	if rawDescription == "" {
		return "", fmt.Errorf("the 'description' is required")
	}

	return rawDescription, nil
}

func parseRowData(record []string, columnIndex map[columnName]int) (Transaction, []string) {
	tx := Transaction{}
	errors := []string{}

	for _, col := range columns {
		result, err := col.action(record, columnIndex[col.name])
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		switch col.name {
		case timestampColumn:
			tx.Timestamp = result.(int64)
		case nameColumn:
			tx.Name = result.(string)
		case typeColumn:
			tx.Type = result.(TransactionType)
		case amountColumn:
			tx.Amount = result.(float64)
		case statusColumn:
			tx.Status = result.(TransactionStatus)
		case descriptionColumn:
			tx.Description = result.(string)
		}
	}

	return tx, errors
}

func parseHeader(record []string) map[columnName]int {
	columnIndex := map[columnName]int{}
	for i, header := range record {
		for _, col := range columns {
			if header == string(col.name) {
				columnIndex[col.name] = i
			}
		}
	}
	return columnIndex
}

func ParseCSV(file io.Reader) ([]Transaction, map[string]any) {
	reader := csv.NewReader(file)

	keyPrefix := "line"
	errorMap := map[string]any{}
	transactions := []Transaction{}
	line := 0

	var columnIndex map[columnName]int
	for {
		line++
		errorKey := createCSVErrorKey(keyPrefix, line)
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			errorMap[errorKey] = []string{"the content is invalid"}
			continue
		}

		// Validate the header
		if columnIndex == nil {
			columnIndex = parseHeader(record)
			if len(columnIndex) != len(columns) {
				errorMap[errorKey] = []string{
					fmt.Sprintf("expected 6 required columns, got %d", len(columnIndex)),
				}
				break
			}
			continue
		}

		transaction, errors := parseRowData(record, columnIndex)
		if len(errors) == 0 {
			transactions = append(transactions, transaction)
		} else {
			errorMap[errorKey] = errors
		}
	}

	return transactions, errorMap
}

func createCSVErrorKey(prefix string, line int) string {
	return fmt.Sprintf("%s[%d]", prefix, line)
}
