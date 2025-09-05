//go:build windows

package pusher

import (
	"errors"
	"syscall"
	"time"

	"github.com/lxn/win"
)

func Wechat(text string) (err error) {
	// 查找微信窗口
	name, _ := syscall.UTF16PtrFromString("微信")
	window := win.FindWindow(nil, name)
	if window == 0 {
		err = errors.New("未找到微信窗口")
		return
	}

	win.SetFocus(window)

	win.ShowWindow(window, win.SW_RESTORE)
	if !win.SetForegroundWindow(window) {
		err = errors.New("没有获取到光标")
		return
	}

	// 发送文本
	for _, c := range text {
		win.SendMessage(window, win.WM_CHAR, uintptr(c), 0)
		time.Sleep(10 * time.Millisecond)
	}

	// 发送回车
	win.SendMessage(window, win.WM_KEYDOWN, win.VK_RETURN, 0)
	time.Sleep(100 * time.Millisecond)
	win.SendMessage(window, win.WM_KEYUP, win.VK_RETURN, 0)
	return
}
