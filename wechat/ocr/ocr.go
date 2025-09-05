package ocr

import "image"

// Func 输入一张照片 返回照片里面的文本信息
type Func func(img image.Image) (text string, err error)
