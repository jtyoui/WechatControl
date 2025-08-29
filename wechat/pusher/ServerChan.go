package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type sendResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ServerChan Server酱推送服务 https://sct.ftqq.com/
func ServerChan(text string) (err error) {
	var (
		response *http.Response
		result   sendResponse
	)

	sendKey := ""
	link := fmt.Sprintf("https://1315.push.ft07.com/send/%s.send", sendKey)

	params := url.Values{}

	texts := strings.Split(strings.TrimSpace(text), "\n")
	text = "# 【OCR识别结果】\n" + texts[len(texts)-1]
	params.Set("title", "微信OCR推送")
	params.Set("desp", text)
	params.Set("tags", "微信|OCR")

	response, err = http.Post(link, "application/x-www-form-urlencoded", strings.NewReader(params.Encode()))
	if err != nil {
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)

	if result.Message != "SUCCESS" {
		err = errors.New(result.Message)
	}

	return
}
