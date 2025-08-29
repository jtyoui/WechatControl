//go:build windows

package screenshot

import (
	"errors"
	"image"
	"syscall"

	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

func Screenshot() (img image.Image, err error) {
	var rect win.RECT

	// 查找包含微信标题的窗口
	name, _ := syscall.UTF16PtrFromString("微信")
	window := win.FindWindow(nil, name)

	// 等待窗口激活
	flag := win.SWP_NOMOVE | win.SWP_NOSIZE
	win.SetWindowPos(window, win.HWND_TOP, 0, 0, 0, 0, uint32(flag))
	win.SetWindowPos(window, win.HWND_NOTOPMOST, 0, 0, 0, 0, uint32(flag))
	win.SetForegroundWindow(window)

	// 获取窗口位置和大小
	if ok := win.GetWindowRect(window, &rect); !ok {
		err = errors.New("没有找到微信窗口")
		return
	}

	// 调整坐标（与原Python代码的调整值对应）
	x1 := int(rect.Left) + 280
	y1 := int(rect.Top) + 70
	x2 := int(rect.Right) - 15
	y2 := int(rect.Bottom) - 150

	// 截图
	if img, err = robotgo.CaptureImg(x1, y1, x2-x1, y2-y1); err != nil {
		return
	}

	return
}
