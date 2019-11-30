package main

import (
	"encoding/json"
	"io/ioutil"
)

func (im *IncomingMessage) AddLine(logFile string) {
	file, _ := json.Marshal(im)
	_ = ioutil.WriteFile(logFile, file, 0644)
}
