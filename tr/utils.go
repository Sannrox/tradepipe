package tr

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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

func getTimestampAsString(docMap map[string]interface{}) string {
	timestamp, ok := docMap["detail"].(string)
	if !ok {
		timestamp = fmt.Sprint(time.Now().Unix())
	}
	return timestamp
}

func getTimestampAsInt64(docMap map[string]interface{}) int64 {
	timestamp := getTimestampAsString(docMap)

	timestampInt64, _ := strconv.ParseInt(timestamp, 10, 64)
	return timestampInt64
}

func isSavingsPlan(response TimelineDetail) bool {
	if response.SubtitleText == "Sparplan" || response.SubtitleText == "Savings plan" {
		return true
	} else {
		for _, section := range response.Sections {
			if section.Type == "actionButtons" {
				for _, button := range section.Data {
					if button.(map[string]interface{})["action"].(map[string]interface{})["type"] == "editSavingPlan" || button.(map[string]interface{})["action"].(map[string]interface{})["type"] == "deleteSavingPlan" {
						return true
					}
				}
			}
		}
	}
	return false
}

func getSavingsPlanFMT(response TimelineDetail, ifSavingPlan bool) string {
	if response.SubtitleText != "Sparplan" && ifSavingPlan {
		return " -- SPARPLAN"
	} else {
		return ""
	}
}
