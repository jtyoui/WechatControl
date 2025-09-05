package tool

import (
	"encoding/json"
	"os"
)

type config struct {
	Pusher  map[string]string `json:"pusher"`
	OCR     map[string]string `json:"ocr"`
	Command map[string]string `json:"command"`
}

var Config config

func init() {
	fs, err := os.Open("config.json")
	if err != nil {
		panic(Error("没有找到配置文件"))
		return
	}
	if err = json.NewDecoder(fs).Decode(&Config); err != nil {
		panic(err)
	}
}
