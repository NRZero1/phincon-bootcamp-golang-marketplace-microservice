package utils

import (
	"fmt"
	"strconv"
	"time"
)

func GenerateNextTransactionID(lastID string) (string, error) {
	if lastID == "" {
		now := time.Now()
		currentDay := now.Format("060102") // YYMMDD format
		timePart := now.Format("150405")   // hhmmss format
		return fmt.Sprintf("%s%s%03d", currentDay, timePart, 1), nil
	}
	
	// Extract the date and last counter value
	_, lastCounter, err := ExtractCounterFromID(lastID)
	if err != nil {
		return "", err
	}

	// Increment the counter
	newCounter := lastCounter + 1

	// Get the current date
	now := time.Now()
	currentDay := now.Format("060102") // YYMMDD format
	timePart := now.Format("150405")   // hhmmss format

	// Format the new ID as YYMMDDhhmmssnnn
	nextID := fmt.Sprintf("%s%s%03d", currentDay, timePart, newCounter)
	return nextID, nil
}


func ExtractCounterFromID(transactionID string) (string, int, error) {
	if len(transactionID) < 15 {
		return "", 0, fmt.Errorf("invalid transaction ID length")
	}

	datePart := transactionID[:6]        // Extract YYMMDD part
	timePart := transactionID[6:12]      // Extract hhmmss part
	counterPartStr := transactionID[12:] // Extract nnn part

	// Convert the counter part to an integer
	counterPart, err := strconv.Atoi(counterPartStr)
	if err != nil {
		return "", 0, err
	}

	return datePart + timePart, counterPart, nil
}