package command

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jtyoui/WechatControl/tool"
)

func OpenAI(text string) (result string, err error) {
	params := map[string]any{
		"model": tool.Config.Command["model"],
		"messages": []map[string]any{
			{"role": "system", "content": "你是一个解决问题的专家。"},
			{"role": "user", "content": text},
		},
	}
	data, _ := json.Marshal(&params)

	request, err := http.NewRequest("POST", tool.Config.Command["url"], bytes.NewBuffer(data))
	if err != nil {
		return
	}

	request.Header.Set("Authorization", "Bearer "+tool.Config.Command["apikey"])
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	type Message struct {
		Content string `json:"content"`
	}

	type OpenAIResponse struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}

	var openAIResponse OpenAIResponse

	if err = json.NewDecoder(response.Body).Decode(&openAIResponse); err != nil {
		return
	}

	result = openAIResponse.Choices[0].Message.Content

	return
}
