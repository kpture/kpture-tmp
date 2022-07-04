package capture

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type LogLine struct {
	ts  time.Time
	msg string
}
type Logs []LogLine

func (p Logs) Len() int {
	return len(p)
}

func (p Logs) Less(i, j int) bool {
	return p[i].ts.Before(p[j].ts)
}

func (p Logs) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (k *Kpture) SortLogs() (*Logs, error) {
	kptureLogs := Logs{}
	err := filepath.Walk(k.basePath,
		func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".log") {
				readFile, err := os.Open(path)
				if err != nil {
					return errors.WithMessage(err, "could not open file for logs")
				}
				fileScanner := bufio.NewScanner(readFile)

				fileScanner.Split(bufio.ScanLines)

				for fileScanner.Scan() {
					line := fileScanner.Text()
					timeStamp := strings.Split(line, " ")
					if len(timeStamp) >= 1 {
						k8sTimestamp, e := time.Parse(
							time.RFC3339,
							timeStamp[0])
						if e == nil {
							kptureLogs = append(kptureLogs, LogLine{
								ts:  k8sTimestamp,
								msg: strings.Join(timeStamp[1:], ""),
							})
						}
					}
				}
				readFile.Close()
			}
			return nil
		})

	sort.Sort(kptureLogs)

	return &kptureLogs, errors.WithMessage(err, "error creating kpture log")
}

func (k *Kpture) WriteLogs(logs *Logs) error {
	fileToWrite, err := os.OpenFile(filepath.Join(k.basePath, "kpture.log"), os.O_CREATE|os.O_RDWR, fs.ModePerm)
	if err != nil {
		k.logger.Error(err)
		return errors.WithMessage(err, "could not openfile for logs")
	}

	for _, log := range *logs {
		fileToWrite.WriteString(log.msg + "\n")
	}
	fileToWrite.Close()
	return nil
}
