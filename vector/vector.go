package vector

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"vectoral/utils"
)

type VectorPath struct {
	isUsed    bool
	color     color.Color
	StartLine *[]*[2]int
	EndLine   *[]*[2]int
}

func (p *VectorPath) AddMoveStart(columX int, rowY int) {
	move1 := [2]int{columX, rowY}
	move2 := [2]int{columX, rowY + 1}

	*p.StartLine = append(*p.StartLine, &move1, &move2)

}
func (p *VectorPath) AddMoveEnd(columX int, rowY int) {
	move1 := [2]int{columX, rowY}
	move2 := [2]int{columX, rowY + 1}

	*p.EndLine = append(*p.EndLine, &move1, &move2)

}

func NewVectorPath(color color.Color) *VectorPath {
	startLine := []*[2]int{}
	endLine := []*[2]int{}
	return &VectorPath{
		isUsed:    true,
		color:     color,
		StartLine: &startLine,
		EndLine:   &endLine,
	}
}

func PathInclude(vectorPaths []*VectorPath, index int) (bool, *VectorPath) {
	count := len(vectorPaths)

	if index > -1 && index < count {
		el := vectorPaths[index]

		return el != nil, el
	}

	return false, nil
}

type VectorImage struct {
	ColorDiffPercent float64
	Img              image.Image
	Widget           int
	Height           int
	Bounds           image.Rectangle
	OnePixelScaleX   float64
	OnePixelScaleY   float64
}

// func (v *VectorImage) MoveScale(x int, y int) (float64, float64) {
// 	return float64(x) / float64(v.Widget), float64(y) / float64(v.Height)
// }

func (v *VectorImage) MoveScale(x int, y int) (float64, float64) {
	// return float64(x) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
	return float64(x), float64(y)
}

// func (v *VectorImage) MoveEnd(x int, y int) (float64, float64) {
// 	return float64(x+1) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
// }

func (v *VectorImage) ImageVector() (image.Image, []*VectorPath) {
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := []*VectorPath{}

	jobChannel := &utils.JobChannel[func()]{}
	newImage := image.NewRGBA(image.Rect(0, 0, v.Widget, v.Height))
	pathShapes := make([]*VectorPath, v.Widget)

	for rowY := 0; rowY < v.Height; rowY++ {
		for columnX := 0; columnX < v.Widget; columnX++ {
			r, g, b, a := v.Img.At(columnX, rowY).RGBA()

			red := math.Round(float64((r>>8))/colorDiffNum) * colorDiffNum
			green := math.Round(float64((g>>8))/colorDiffNum) * colorDiffNum
			blue := math.Round(float64((b>>8))/colorDiffNum) * colorDiffNum

			pixelColor := color.RGBA{
				uint8(red),
				uint8(green),
				uint8(blue),
				uint8(a),
			}
			newImage.Set(columnX, rowY, pixelColor)

			leftOk, left := PathInclude(pathShapes, columnX-1)
			curOk, current := PathInclude(pathShapes, columnX)

			// equal := curOk && leftOk && *current == *left

			isColorCurrent := curOk && current.color == pixelColor
			isColorLeft := leftOk && left.color == pixelColor

			if isColorLeft {
				pathShapes[columnX] = left
				current = left
				curOk = true
				// equal = true
			} else if !isColorCurrent {
				current = NewVectorPath(pixelColor)
				curOk = true
				// equal = false
				pathShapes[columnX] = current
				jobChannel.AddJob(func() {
					// col := color.RGBA{0, 0, 0, 255}
					// if current.isUsed && col == current.color {
					// 	paths = append(paths, current)
					// }
					if current.isUsed {
						paths = append(paths, current)
					}
				})
			}

			if columnX == (v.Widget - 1) {
				current.AddMoveEnd(columnX+1, rowY)
			}
			if !isColorLeft {
				if leftOk {
					left.AddMoveEnd(columnX, rowY)
				}
				current.AddMoveStart(columnX, rowY)
			}
		}
	}
	jobChannel.Run()
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDD", len(paths))
	return newImage, paths
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

func (v VectorImage) SavePathsToSVGFile(paths []*VectorPath, fileName string, saveWidget int, saveHeight int) {
	os.Remove(fileName)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// bounds := v.Img.Bounds()
	// widget := bounds.Max.X
	// height := bounds.Max.Y
	// if _, err := f.Write([]byte(fmt.Printf(""))); err != nil {
	// 	log.Fatal(err)
	// }
	// viewBox="0 0 %v %v"
	if _, err := f.Write([]byte(fmt.Sprintf(
		"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"%vpx\" height=\"%vpx\">\n",
		saveWidget,
		saveHeight,
	))); err != nil {
		log.Fatal(err)
	}
	// <svg height="210" width="400">
	// fmt.Println("paths: ", paths)</svg></svg>

	for _, path := range paths {
		if path == nil {
			continue
		}

		// moveLinesLeft := *path.MoveLinesLeft
		// moveLinesRight := *path.MoveLinesRight
		// sizeLeft := len(moveLinesLeft)
		// sizeRight := len(moveLinesRight)
		// if sizeLeft < 3 || sizeRight < 3 {
		//გამოიწვიე საკუთარი თავი
		// 	continue
		// }
		d := ""
		for _, XYPoint := range *path.StartLine {
			x := XYPoint[0]
			// * float64(saveWidget)
			y := XYPoint[1]
			// * float64(saveHeight)
			d = d + fmt.Sprintf("L%v %v ", x, y)
		}
		for _, XYPoint := range *path.EndLine {
			x := XYPoint[0]
			// * float64(saveWidget)
			y := XYPoint[1]
			// * float64(saveHeight)
			d = fmt.Sprintf("L%v %v ", x, y) + d
		}
		// }
		// fmt.Println("d:", d)
		// dLeft := ""
		// for _, XYPoint := range moveLinesLeft {
		// 	x := XYPoint[0] * float64(saveWidget)
		// 	y := XYPoint[1] * float64(saveHeight)
		// 	dLeft = dLeft + fmt.Sprintf("L%v %v ", x, y)
		// }

		// dRight := ""
		// for _, XYPoint := range moveLinesRight {
		// 	x := XYPoint[0] * float64(saveWidget)
		// 	y := XYPoint[1] * float64(saveHeight)

		// 	dRight = fmt.Sprintf("L%v %v ", x, y) + dRight
		// }

		// fmt.Println("D: ", dLeft)
		// fmt.Println("D: ", dRight)
		r, g, b, a := path.color.RGBA()
		color := fmt.Sprintf("rgba(%v,%v,%v,%v)", r>>8, g>>8, b>>8, a>>8)

		if _, err := f.Write([]byte(fmt.Sprintf("<path fill=\"%v\" d=\"M%vZ\" />\n", color, d[1:]))); err != nil {
			log.Fatal(err)
		}

	}
	if _, err := f.Write([]byte("</svg>")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
func NewVectorImage(Img image.Image, ColorDiffPercent float64) *VectorImage {
	bounds := Img.Bounds()
	widget := bounds.Max.X
	height := bounds.Max.Y
	OnePixelScaleX := 1 / float64(widget)
	OnePixelScaleY := 1 / float64(height)

	return &VectorImage{
		ColorDiffPercent: ColorDiffPercent,
		Img:              Img,
		Widget:           widget,
		Height:           height,
		Bounds:           bounds,
		OnePixelScaleX:   OnePixelScaleX,
		OnePixelScaleY:   OnePixelScaleY,
	}
}
