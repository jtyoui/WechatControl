package ocr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/jtyoui/WechatControl/tool"
)

type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		RawOut [][]any `json:"raw_out"`
	} `json:"data"`
}

func TrWebOCR(img image.Image) (text string, err error) {
	var (
		buf      bytes.Buffer
		response *http.Response
		value    result
	)

	if err = png.Encode(&buf, img); err != nil {
		return
	}

	payload := url.Values{
		"img":      {base64.StdEncoding.EncodeToString(buf.Bytes())},
		"compress": {"0"},
		"is_draw":  {"0"},
	}

	if response, err = http.PostForm(tool.Config.OCR["url"], payload); err != nil {
		return
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&value); err != nil {
		return
	}

	if value.Code != 200 {
		return
	}

	rawOut := value.Data.RawOut

	rawOut = slices.DeleteFunc(rawOut, func(out []interface{}) bool {
		return out[1].(string) == ""
	})

	for i, out := range rawOut {
		text += out[1].(string)
		if out[1] == "" || i >= len(rawOut)-1 {
			continue
		}

		nextPoint := rawOut[i+1][0].([]interface{})[1].(float64)
		point := out[0].([]interface{})[1].(float64)
		width := out[0].([]interface{})[3].(float64)
		if nextPoint-point > width+5 {
			text += "\n"
		}
	}

	text = strings.TrimSpace(text)

	if text != "" {
		// 获取最后一段
		texts := strings.Split(text, "\n")
		text = texts[len(texts)-1]
		return
	}

	err = errors.New("ocr识别为空")
	return
}
