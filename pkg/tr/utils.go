package tr

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func ExportTransactions(inputPath, outputPath string) {
	// Read the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	jsonParser := json.NewDecoder(inputFile)
	var timeline map[string]interface{}
	if err = jsonParser.Decode(&timeline); err != nil {
		log.Fatal(err)
	}
	logrus.Info("Write deposit and removal transactions to CSV file")
	// Write the csv output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	// Write the header
	outputFile.WriteString("Datum" + "," + "Type" + "," + "Werte" + "\n")
	// Write the transactions
	for _, event := range timeline {
		event = event.(map[string]interface{})["data"].(map[string]interface{})
		dateTime := event.(map[string]interface{})["timestamp"].(string)

		title := event.(map[string]interface{})["title"].(string)
		if v, ok := event.(map[string]interface{})["body"]; ok {
			if strings.Contains(v.(string), "storniert") {
				continue
			}
			if strings.Contains(title, "Einzahlung") || strings.Contains(title, "Bonuszahlung") {
				outputFile.WriteString(dateTime + "," + "deposit" + "," + event.(map[string]interface{})["cashChangeAmount"].(string) + "\n")
			} else if strings.Contains(title, "Auszahlung") {
				outputFile.WriteString(dateTime + "," + "removal" + "," + event.(map[string]interface{})["cashChangeAmount"].(string) + "\n")
			} else if strings.Contains(title, "Reinvestierung") {
				logrus.Warn("Reinvestierung not implemented")
			}

			logrus.Debug("Date: " + dateTime + " Title: " + title + " Value: " + event.(map[string]interface{})["cashChangeAmount"].(string))

		}

	}
	logrus.Info("Done!")
}

func filePathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getTimestamp(docMap map[string]interface{}) float64 {
	timestamp, ok := docMap["detail"].(float64)
	if !ok {
		timestamp = float64(time.Now().Unix())
	}
	return timestamp
}
