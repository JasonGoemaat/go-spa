package mylog

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var LOG_FILE = "C:\\t\\local.log"

func Log(message string) {
	outfile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("ERROR OPENING LOG FILE C:\\t\\main.log")
		return
	}
	defer outfile.Close()
	dateString := time.Now().Format("2006-01-02 15:04:05.999") // interesting, this has to be the exact date used, i.e. '15' signifies 24hr, '03' would signify 12hr I think, so '01' is month, '02' is day, '03'/'15' is hour, '04' is minute, '05' is second, '06'/'2006' is year, 9s for partial seconds
	entireMessage := fmt.Sprintf("%s: %s", dateString, message)
	fmt.Fprintln(io.MultiWriter(outfile, os.Stdout), entireMessage)
}

func ShowEnvironmentVariables() {
	for _, key := range os.Environ() {
		value := os.Getenv(key)
		s := fmt.Sprintf("ENV: %s = %s", key, value)
		Log(s)
	}
}
