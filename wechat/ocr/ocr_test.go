package ocr_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/jtyoui/WechatControl/wechat/ocr"
)

func TestTrWebOCR(t *testing.T) {
	imaPath := "C:\\Users\\jtyou\\Desktop\\img1.png"
	fs, _ := os.Open(imaPath)
	img, _ := png.Decode(fs)
	text, err := ocr.TrWebOCR(img)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(text)
}
