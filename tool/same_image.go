package tool

import (
	"image"
	"image/color"
)

// 比较两个颜色是否完全相同
func colorEq(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// ImagesSame 判断两个照片是否相同,返回 true表示相同
func ImagesSame(img1, img2 image.Image) (ok bool) {
	// 首先比较尺寸
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	// 如果边界不同，图像肯定不同
	if bounds1 != bounds2 {
		return
	}

	// 比较每一个像素点
	width := bounds1.Max.X - bounds1.Min.X
	height := bounds1.Max.Y - bounds1.Min.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素颜色
			color1 := img1.At(x+bounds1.Min.X, y+bounds1.Min.Y)
			color2 := img2.At(x+bounds1.Min.X, y+bounds1.Min.Y)

			// 比较颜色是否相同
			if !colorEq(color1, color2) {
				return
			}
		}
	}

	ok = true
	return
}
