package screenshot_test

import (
	"github.com/jtyoui/WechatControl/wechat/screenshot"
	"image/jpeg"
	"os"
	"testing"
	"time"
)

func TestScreenshot(t *testing.T) {
	img, err := screenshot.Screenshot()
	if err != nil {
		t.Error(err)
		return
	}
	name := time.Now().Format("20060102_150405.jpg")
	file, _ := os.Create(name)
	if err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Error(err)
	}
}
