package wechat

import (
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/jtyoui/WechatControl/tool"
	"github.com/jtyoui/WechatControl/wechat/command"
	"github.com/jtyoui/WechatControl/wechat/ocr"
	"github.com/jtyoui/WechatControl/wechat/pusher"
	"github.com/jtyoui/WechatControl/wechat/screenshot"
)

// Print 每10秒打印一下任务
func Print() {
	for range time.Tick(10 * time.Second) {
		info := fmt.Sprintf("正常：%d任务，失败：%d任务", len(tool.WxNormal), len(tool.WxError))
		tool.Info(info)
	}
}

// LoopMonitor 开启一个循环监听服务
func LoopMonitor(ocr ocr.Func) {
	var (
		nextImg image.Image
		text    string
	)

	for range time.Tick(time.Second) {
		img, err := screenshot.Screenshot()
		if err != nil {
			_ = tool.Error(fmt.Sprintf("截图失败:%s", err))
			continue
		}

		if nextImg != nil {
			if !tool.ImagesSame(img, nextImg) {
				tool.Info("正在提取OCR文本...")
				if text, err = ocr(img); err != nil {
					_ = tool.Error(fmt.Sprintf("OCR识别失败：%s", err))
					continue
				}
				tool.Success("识别文本：" + text)
				text = string([]rune(text)[1:])
				if strings.HasPrefix(text, "夏目三千") {
					text = strings.ReplaceAll(text, "夏目三千", "")
					tool.Success("OCR识别成功！")
					tool.WxNormal <- tool.WxInfo{Data: text}
				}
			}
		}
		nextImg = img
	}
}

// Pusher 消息推送：用于通知用户是否成功。[可选]
func Pusher(push pusher.Func) {
	for wx := range tool.WxNormal {
		if wx.Over {
			continue
		}

		// 已经推送了当时还没有被执行
		if wx.Pusher.OK && !wx.Command.OK {
			tool.WxNormal <- wx
			continue
		}

		var text string

		if !wx.Pusher.OK {
			text = fmt.Sprintf("#正在执行：%s，请稍等！", wx.Data)
		}

		if wx.Command.OK {
			text = fmt.Sprintf("#执行问题：%s。答案：%s", wx.Data, wx.Command.Result)
		}

		tool.Info("正在消息推送...")

		if err := push(text); err != nil {
			err = tool.Error(fmt.Sprintf("消息推送服务失败：%s", err.Error()))
			if wx.Pusher.Stop(err) {
				tool.WxError <- wx
				continue
			}
			tool.Warn("消息推送正在重复执行中...")
		} else {
			tool.Success("消息推送成功！")
			wx.Pusher.OK = true

			// 全部结束
			if wx.Command.OK {
				wx.Over = true
				continue
			}
		}

		tool.WxNormal <- wx
	}
}

// Commander 执行命令
func Commander(command command.Func) {
	for wx := range tool.WxNormal {
		if wx.Over {
			continue
		}

		// 消息没有被推送 或者 命令已经执行过了 都跳过
		if !wx.Pusher.OK || wx.Command.OK {
			tool.WxNormal <- wx
			continue
		}

		tool.Info("正在执行命令...")

		if result, err := command(wx.Data); err != nil {
			err = tool.Error(fmt.Sprintf("执行命令失败：%s", err.Error()))
			if wx.Command.Stop(err) {
				tool.WxError <- wx
				continue
			}
			tool.Warn("命令正在重复执行中...")
		} else {
			wx.Command.OK = true
			wx.Command.Result = result
			tool.Success("执行命令成功！")
		}
		tool.WxNormal <- wx
	}
}

func Run() {
	go Print()
	go Pusher(pusher.Wechat)
	go Commander(command.OpenAI)
	LoopMonitor(ocr.TrWebOCR)
}
