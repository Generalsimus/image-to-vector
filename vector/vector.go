package vector

import (
	"image"
	"image/color"
	"math"
)

type VectorImage struct {
	ColorDiffPercent float64
	Img              image.Image
}

func (v VectorImage) RGBToNum(r uint8, g uint8, b uint8) uint8 {
	return (r << 16) + (g << 8) + (b)
}
func (v VectorImage) NumToRGB(num int) (uint8, uint8, uint8) {

	return uint8((num & 0xff0000) >> 16), uint8((num & 0x00ff00) >> 8), uint8(num & 0x0000ff)
}
func (v VectorImage) ImageVector() image.Image {
	bounds := v.Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	img := v.Img
	colorDiffNum := float64(255 * v.ColorDiffPercent)

	newImage := image.NewRGBA(image.Rect(0, 0, widget, height))
	for row := 0; row < height; row++ {
		for column := 0; column < widget; column++ {
			r, g, b, a := img.At(column, row).RGBA()

			red := math.Round(float64((r>>8))/colorDiffNum) * colorDiffNum
			green := math.Round(float64((g>>8))/colorDiffNum) * colorDiffNum
			blue := math.Round(float64((b>>8))/colorDiffNum) * colorDiffNum

			newImage.Set(column, row, color.RGBA{
				uint8(red),
				uint8(green),
				uint8(blue),
				uint8(a),
			})

		}
	}

	return newImage
}

// func rgbToNum[T int | uint | uint32](r T, g T, b T) T {
// 	rgb := r
// 	rgb = (rgb << 8) + g
// 	rgb = (rgb << 8) + b
// 	return rgb
// }

// func numToRgb[T int | uint | uint32](num T) (T, T, T) {

// 	red := (num >> 16) & 0xFF
// 	green := (num >> 8) & 0xFF
// 	blue := (num) & 0xFF

// 	return red, green, blue
// }
