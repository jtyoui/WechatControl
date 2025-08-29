package ocr

import "image"

type Func func(img image.Image) (text string, err error)
