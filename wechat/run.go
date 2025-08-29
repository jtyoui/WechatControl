package wechat

import (
	"fmt"
	"image"
	"time"

	"github.com/jtyoui/WechatControl/tool"
	"github.com/jtyoui/WechatControl/wechat/ocr"
	"github.com/jtyoui/WechatControl/wechat/pusher"
	"github.com/jtyoui/WechatControl/wechat/screenshot"
)

var (
	WxInfo = make(chan string, 100)
)

// LoopMonitor 开启一个循环监听服务
func LoopMonitor(ocr ocr.Func) {
	var (
		nextImg image.Image
		text    string
	)

	for range time.Tick(2 * time.Second) {
		img, err := screenshot.Screenshot()
		if err != nil {
			tool.Error(fmt.Sprintf("截图失败:%s", err))
			continue
		}

		if nextImg != nil {
			if !tool.ImagesSame(img, nextImg) {
				if text, err = ocr(img); err != nil {
					tool.Error(fmt.Sprintf("OCR识别失败:%s", err))
					continue
				}
				tool.Success("OCR识别成功！")
				WxInfo <- text
			}
		}
		nextImg = img
	}
}

// Pusher 消息推送：用于通知用户是否成功。[可选]
func Pusher(push pusher.Func) {
	for text := range WxInfo {
		tool.Info("正在消息推送...")

		if err := push(text); err != nil {
			tool.Error(fmt.Sprintf("消息推送服务失败：%s", err.Error()))
			continue
		}

		tool.Success("消息推送成功！")
	}
}

func Run() {
	go Pusher(pusher.ServerChan)
	LoopMonitor(ocr.TrWebOCR)
}
