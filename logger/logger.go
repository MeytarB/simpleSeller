package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/prebid/openrtb/v20/openrtb2"
)

type Logger struct {
	logFile *os.File
	mtx     sync.Mutex
}

func NewLogger() *Logger {
	newFile, err := os.OpenFile("simpleServerLogger", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error creating log file:", err)
	}
	fileHeader := "logger:" + (time.Now().String()) + "\n"
	newFile.Write([]byte(fileHeader))
	return &Logger{logFile: newFile}
}

func (l *Logger) Close() {
	l.logFile.Close()
}

type bidReqLogInfo struct {
	Id     string
	Target string
	AdType string
}

func LogBidRequest(l *Logger, breq *openrtb2.BidRequest) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	var target, adType string
	switch {
	case breq.Site != nil:
		target = "site"
	case breq.App != nil:
		target = "app"
	}

	switch {
	case breq.Imp[0].Banner != nil:

		adType = "banner"
	case breq.Imp[0].Video != nil:
		target = "video"
	}

	breqInfo := bidReqLogInfo{Id: breq.ID, Target: target, AdType: adType}
	fmt.Println(breqInfo)
	breqJSON, _ := json.Marshal(breqInfo)
	l.logFile.Write(append([]byte("INFO: ")))
	l.logFile.Write(append(breqJSON))
	l.logFile.Write(append([]byte(time.Now().String()), '\n'))

}
