package tr

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type Downloader struct {
	TR             APIClient
	OutputPath     string
	HistoryFile    string
	FilenameFmt    string
	FilePaths      []string
	DocUrls        []string
	DocUrlsHistory []string
	TimeLine       *TimeLine
}

func NewDownloader(client APIClient) *Downloader {
	return &Downloader{TR: client, HistoryFile: "history.txt", FilenameFmt: "%s_%s_%s_%s_%s.pdf", OutputPath: "./tmp/"}
}
func (d *Downloader) DownloadDocument(docUrl string, filePath string) error {
	resp, err := d.TR.Client.Get(docUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (d *Downloader) Download(doc Doc, titleText string, subTitleText string, subFolder string) error {
	docUrl := doc.Action.Payload
	date := doc.Detail
	docId := doc.ID
	re := regexp.MustCompile(`\d+:\d+`)
	matches := re.FindStringSubmatch(subTitleText)
	var time string
	if len(matches) == 0 {
		time = "00:00"
	} else {
		time = matches[0]
	}
	var dir string
	if subFolder != "" {
		dir = d.OutputPath + subFolder
	} else {
		dir = d.OutputPath
	}

	docType := strings.Split(doc.Title, " ")
	var docTypeNumber string
	if isNumeric(docType[len(docType)-1]) {
		docTypeNumber = docType[len(docType)-1]
		docType[len(docType)-1] = ""

	} else {
		docTypeNumber = ""
	}
	docTypeJoined := strings.Join(docType, " ")
	docTypeJoined = strings.TrimSpace(docTypeJoined)
	titleText = strings.Replace(titleText, "\n", "", -1)
	titleText = strings.Replace(titleText, "/", "-", -1)
	filename := fmt.Sprintf(d.FilenameFmt, date, time, titleText, docTypeNumber, docId)
	var filePath string
	if strings.Contains(docTypeJoined, "Kontoauszug") || strings.Contains(docTypeJoined, "Depotauszug") {
		filePath = filepath.Join(dir, "Abschluesse", filename, docTypeJoined)
	} else {
		filePath = filepath.Join(dir, docTypeJoined, filename)
	}
	if contains(d.FilePaths, filePath) {
		logrus.Debug("File already downloaded: ", filePath)
		return nil
	} else {
		d.FilePaths = append(d.FilePaths, filePath)
	}
	logrus.Infof("Trying to download to: %s", filePath)
	// check if file already exists
	if !filePathExists(filePath) {
		docURLBase := strings.Split(docUrl, "?")[0]
		if contains(d.DocUrls, docURLBase) {
			logrus.Debug("File already downloaded: ", filePath)
		} else if contains(d.DocUrlsHistory, docURLBase) {
			logrus.Debug("File already downloaded: ", filePath)
		} else {
			d.DocUrls = append(d.DocUrls, docURLBase)
		}
		err := os.MkdirAll(path.Dir(filePath), 0o755)
		if err != nil {
			return err
		}

		err = d.DownloadDocument(docUrl, filePath)
		if err != nil {
			logrus.Error("Error downloading document: ", err)
			return err
		}
		err = d.WriteHistoryFile(docURLBase)
		if err != nil {
			logrus.Error("Error writing history file: ", err)
			return err
		}

	} else {
		logrus.Debug("File already exists: ", filePath)
	}
	return nil
}

func (d *Downloader) DownloadAll(ctx context.Context, data chan Message) error {

	if d.HistoryFile != "" {
		err := d.ReadHistoryFile()
		if err != nil {
			logrus.Error("Error reading history file: ", err)
			return err
		}
	}
	err := d.CreateOutputDir()
	if err != nil {
		logrus.Error("Error creating output directory: ", err)
		return err
	}

	d.TimeLine = NewTimeLine(&d.TR)

	if err := d.TimeLine.LoadTimeLine(ctx, data); err != nil {
		logrus.Error("Error loading timeline: ", err)
		return err
	}

	if err := d.TimeLine.LoadTimeLineDetails(ctx, data); err != nil {
		logrus.Error("Error loading timeline details: ", err)
		return err
	}
	for _, detail := range d.TimeLine.TimelineDetails {
		isSavingsPlan := isSavingsPlan(detail)
		savingsPlan := getSavingsPlanFMT(detail, isSavingsPlan)
		logrus.Infof("%d/%d: %s -- %s%s", d.TimeLine.ReceivedDetail, d.TimeLine.NumberofTimlineDetails, detail.TitleText, detail.SubtitleText, savingsPlan)

		if err := d.DownloadTimeLineDetails(detail, d.TimeLine.SinceTimestamp, isSavingsPlan); err != nil {
			return err
		}
	}
	if d.TimeLine.ReceivedDetail == d.TimeLine.NumberofTimlineDetails {
		if err := d.CreateOutputDir(); err != nil {
			return err
		}
		d.WriteFile(d.TimeLine.TimeLineEventsWithDocs, "timelineEventsWithDocs.json")
		d.WriteFile(d.TimeLine.TimelineEventsWithoutDocs, "timelineEventsWithoutDocs.json")

		// ExportTransactions("", d.OutputPath)
		// }
	}
	return nil
}

func (d *Downloader) DownloadTimeLineDetails(response TimelineDetail, maxAgeTimestamp int64, isSavingsPlan bool) error {
	for _, sec := range response.Sections {
		if sec.Type == "documents" {
			for _, doc := range sec.Documents {
				b, err := json.Marshal(doc)
				if err != nil {
					logrus.Error("Error marshalling document: ", err)
					return err
				}
				var f map[string]interface{}
				err = json.Unmarshal(b, &f)
				if err != nil {
					logrus.Error("Error unmarshalling document: ", err)
					return err
				}
				timestamp := getTimestampAsInt64(f)
				if maxAgeTimestamp == 0 || timestamp > maxAgeTimestamp {
					if isSavingsPlan {
						err := d.Download(doc, response.TitleText, response.SubtitleText, "SavingsPlan")
						if err != nil {
							logrus.Error("Error downloading document: ", err)
							return err
						}
					} else {
						err := d.Download(doc, response.TitleText, response.SubtitleText, "")
						if err != nil {
							logrus.Error("Error downloading document: ", err)
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (d *Downloader) WriteHistoryFile(docUrl string) error {
	logrus.Debug("Writing history file: ", d.HistoryFile)
	f, err := os.OpenFile(d.HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		logrus.Error("Error opening history file: ", err)
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(docUrl + "\n"); err != nil {
		logrus.Error("Error writing history file: ", err)
		return err
	}
	return nil
}

func (d *Downloader) CreateOutputDir() error {
	if !filePathExists(d.OutputPath) {
		err := os.MkdirAll(d.OutputPath, 0o755)
		if err != nil {
			logrus.Error("Error creating output directory: ", err)
			return err
		}
	}
	return nil
}

func (d *Downloader) WriteFile(data []TimeLineEvent, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		logrus.Error("Error creating file: ", err)
		return err
	}
	defer f.Close()
	for _, v := range data {
		_, err := f.WriteString(fmt.Sprintf("%v\n", v))
		if err != nil {
			logrus.Error("Error writing file: ", err)
			return err
		}
	}

	return nil
}

func (d *Downloader) ReadHistoryFile() error {
	d.HistoryFile = filepath.Join(d.OutputPath, "history.txt")
	if filePathExists(d.HistoryFile) {
		file, err := os.Open(d.HistoryFile)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			d.DocUrlsHistory = append(d.DocUrlsHistory, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
