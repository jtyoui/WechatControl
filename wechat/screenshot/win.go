//go:build windows

package screenshot

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func Screenshot() (img image.Image, err error) {
	var rect win.RECT

	// 查找微信窗口
	name, _ := syscall.UTF16PtrFromString("微信")
	window := win.FindWindow(nil, name)
	if window == 0 {
		err = errors.New("未找到微信窗口")
		return
	}

	// 获取窗口在屏幕上的位置和大小（包含边框）
	if ok := win.GetWindowRect(window, &rect); !ok {
		err = errors.New("无法获取微信窗口位置")
		return
	}

	// 获取客户区位置（相对于窗口）
	var clientRect win.RECT
	if ok := win.GetClientRect(window, &clientRect); !ok {
		err = errors.New("无法获取客户区位置")
		return
	}

	clientRect.Left += 280
	clientRect.Top += 70
	clientRect.Right -= 15
	clientRect.Bottom -= 150

	// 计算客户区在屏幕上的绝对位置
	clientX := rect.Left + (clientRect.Left)
	clientY := rect.Top + (clientRect.Top)
	width := clientRect.Right - clientRect.Left
	height := clientRect.Bottom - clientRect.Top

	if width <= 0 || height <= 0 {
		err = errors.New("无效的客户区尺寸")
		return
	}

	// 获取整个屏幕的DC
	hdcScreen := win.GetDC(0)
	if hdcScreen == 0 {
		err = errors.New("无法获取屏幕设备上下文")
		return
	}
	defer win.ReleaseDC(0, hdcScreen)

	// 创建内存DC
	hdcMem := win.CreateCompatibleDC(hdcScreen)
	if hdcMem == 0 {
		err = errors.New("无法创建内存设备上下文")
		return
	}
	defer win.DeleteDC(hdcMem)

	// 创建兼容位图
	hBitmap := win.CreateCompatibleBitmap(hdcScreen, width, height)
	if hBitmap == 0 {
		err = errors.New("无法创建兼容位图")
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(hBitmap))

	// 将位图选入内存DC
	oldBmp := win.SelectObject(hdcMem, win.HGDIOBJ(hBitmap))
	defer win.SelectObject(hdcMem, oldBmp)

	// 直接从屏幕DC复制客户区内容（绕过窗口自绘问题）
	win.BitBlt(hdcMem, 0, 0, width, height, hdcScreen, clientX, clientY, win.SRCCOPY)

	// 定义位图信息头
	bmi := win.BITMAPINFO{
		BmiHeader: win.BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
			BiWidth:       width,
			BiHeight:      -height, // 负高度确保图像方向正确
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: win.BI_RGB,
		},
	}

	// 分配缓冲区
	buf := make([]byte, width*height*4)

	// 从内存DC获取位图数据
	if win.GetDIBits(hdcMem, hBitmap, 0, uint32(height), &buf[0], &bmi, win.DIB_RGB_COLORS) == 0 {
		err = errors.New("无法获取位图数据，错误码: " + fmt.Sprintf("%d", win.GetLastError()))
		return
	}

	if buf[0] != buf[1] && buf[1] != buf[2] && buf[0] != 237 {
		err = errors.New("不是有效的微信图片")
		return
	}

	// 转换为RGBA图像
	rgbaImg := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for y := int32(0); y < height; y++ {
		for x := int32(0); x < width; x++ {
			idx := (y*width + x) * 4
			b := buf[idx]
			g := buf[idx+1]
			r := buf[idx+2]
			a := buf[idx+3]
			rgbaImg.SetRGBA(int(x), int(y), color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	img = rgbaImg
	return
}
