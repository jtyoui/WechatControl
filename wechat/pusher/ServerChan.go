package pusher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jtyoui/WechatControl/tool"
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

	link := fmt.Sprintf("https://1315.push.ft07.com/send/%s.send", tool.Config.Pusher["sendKey"])

	params := url.Values{}

	params.Set("title", "微信推送")
	params.Set("desp", text)
	params.Set("tags", "微信")

	if response, err = http.PostForm(link, params); err != nil {
		return
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&result)

	if result.Message != "SUCCESS" {
		err = errors.New(result.Message)
	}

	return
}
