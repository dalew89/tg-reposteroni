package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func (im *IncomingMessage) AddLine(logFile string) {
	in, err := os.Open(logFile)
	logLine, _ := json.Marshal(im)
	out, err := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}
	out.Write(logLine)
	io.Copy(out, in)
}
