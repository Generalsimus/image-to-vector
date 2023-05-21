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
	isUsed        *bool
	color         color.Color
	StartLine     *[]*[2]int
	EndLine       *[]*[2]int
	StartLines    *[]*[]*[2]int
	EndLines      *[]*[]*[2]int
	CurrentStartY int
	// CurrentStartYLineIndex int
	CurrentEndY int
	LineIndex   int
	// CurrentEndYLineIndex   int
}

func (p *VectorPath) AddStart1(columX int, rowY int) {
	if 

}

func (p *VectorPath) AddStart(columX int, rowY int) {
	if rowY != p.CurrentStartY {
		p.CurrentStartY = rowY
		p.CurrentStartYLineIndex = 0
	}

	startLines := *p.StartLines
	var startLine *[]*[2]int

	if len(startLines) <= p.CurrentStartYLineIndex {
		startLine = &[]*[2]int{}
		*p.StartLines = append(*p.StartLines, startLine)
	} else {
		startLine = startLines[p.CurrentStartYLineIndex]
	}

	if p.CurrentStartY != rowY {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*startLine = append(*startLine, &move1, &move2)
		p.CurrentStartY = rowY
	}
	*p.StartLines = startLines
	p.CurrentStartYLineIndex = p.CurrentStartYLineIndex + 1
}

func (p *VectorPath) AddEnd(columX int, rowY int) {
	if rowY != p.CurrentStartY {
		p.CurrentEndYLineIndex = 0
	}
	endLines := *p.EndLines
	var endLine *[]*[2]int
	if len(endLines) <= p.CurrentEndYLineIndex {
		endLine = &[]*[2]int{}
		*p.EndLines = append(*p.EndLines, endLine)
	} else {
		endLine = endLines[p.CurrentEndYLineIndex]
	}
	if p.CurrentStartY != rowY {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*endLine = append(*endLine, &move1, &move2)
		p.CurrentStartY = rowY
	}
	*p.EndLines = endLines

	p.CurrentEndY = rowY
	p.CurrentEndYLineIndex = p.CurrentEndYLineIndex + 1
}

//	func (p *VectorPath) String() string {
//		b, err := json.MarshalIndent(p, "", "  ")
//		if err != nil {
//			fmt.Println(err)
//			return ""
//		}
//		return string(b)
//	}
func (p *VectorPath) AddMoveStart(columX int, rowY int) {
	if p.CurrentStartY != rowY {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*p.StartLine = append(*p.StartLine, &move1, &move2)
		p.CurrentStartY = rowY
	}
}
func (p *VectorPath) AddMoveEnd(columX int, rowY int) {
	if p.CurrentEndY == rowY {
		endLine := *p.EndLine

		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}
		*p.EndLine = append(endLine[:len(endLine)-2], &move1, &move2)

	} else {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*p.EndLine = append(*p.EndLine, &move1, &move2)
		p.CurrentEndY = rowY
	}
}
func (p *VectorPath) Concat(p2 *VectorPath) {
	// return
	*p2.isUsed = false

	startLineP2 := *p2.StartLine
	endLineP2 := *p2.EndLine
	startLineP := *p.StartLine

	for _, point := range endLineP2 {
		startLineP2 = append([]*[2]int{point}, startLineP2...)
	}

	for i := len(startLineP) - 1; i >= 0; i-- {
		startLineP2 = append([]*[2]int{startLineP[i]}, startLineP2...)
	}
	*p.StartLine = startLineP2
	/////////////////////////////////////////////////////////////////////////////////

	*p2 = *p
}

func NewVectorPath(color color.Color) *VectorPath {
	startLine := []*[2]int{}
	endLine := []*[2]int{}
	isUsed := true
	return &VectorPath{
		isUsed:        &isUsed,
		color:         color,
		StartLine:     &startLine,
		EndLine:       &endLine,
		CurrentStartY: -1,
		CurrentEndY:   -1,
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

// func (v *VectorImage) MoveScale(x int, y int) (float64, float64) {
// 	// return float64(x) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
// 	return float64(x), float64(y)
// }

// func (v *VectorImage) MoveEnd(x int, y int) (float64, float64) {
// 	return float64(x+1) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
// }

func (v *VectorImage) ImageVector() (image.Image, []*VectorPath) {
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := []*VectorPath{}

	jobChannel := &utils.JobChannel[func()]{}
	newImage := image.NewRGBA(image.Rect(0, 0, v.Widget, v.Height))
	pathShapes := make([]*VectorPath, v.Widget)
	////////////////////////////////////////////////
	// var wg sync.WaitGroup
	// wg.Add(1)
	// var ProcessEnd sync.WaitGroup

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
			// if rowY > 3 {
			// 	continue
			// }

			leftOk, left := PathInclude(pathShapes, columnX-1)
			curOk, current := PathInclude(pathShapes, columnX)

			equal := curOk && leftOk && current.isUsed == left.isUsed

			isColorCurrent := curOk && current.color == pixelColor
			isColorLeft := leftOk && left.color == pixelColor
			if !equal && isColorCurrent && isColorLeft {
				// fmt.Println(columnX, rowY, current, left, current.isUsed)
				current.Concat(left)
			}

			if isColorLeft {
				pathShapes[columnX] = left
				current = left
				curOk = true
			} else if !isColorCurrent {
				current = NewVectorPath(pixelColor)
				curOk = true
				isUsed := current.isUsed
				pathShapes[columnX] = current

				jobChannel.AddJob(func() {

					// col := color.RGBA{
					// 	255,
					// 	51,
					// 	51,
					// 	255,
					// }

					// col := color.RGBA{
					// 	0,
					// 	0,
					// 	0,
					// 	255,
					// }

					// if col == current.color && *isUsed {
					// 	paths = append(paths, current)
					// }
					if *isUsed {
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
	// defer
	// wg.Done()
	// ProcessEnd.Wait()
	jobChannel.Run()
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	// b, _ := json.MarshalIndent(&paths, "", "  ")
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDD", len(paths), string(b))
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
		"<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"%vpx\" height=\"%vpx\" viewBox=\"0 0 %v %v\">\n",
		saveWidget,
		saveHeight,
		saveWidget,
		saveHeight,
	))); err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		if path == nil {
			continue
		}
		d := ""
		// fmt.Println("StartLine", path.color)
		for _, XYPoint := range *path.StartLine {
			x := XYPoint[0]
			y := XYPoint[1]
			// fmt.Println("X: ", x, "Y: ", y)
			d = d + fmt.Sprintf("L%v %v ", x, y)
		}
		// fmt.Println("EndLine", path.color)
		for _, XYPoint := range *path.EndLine {
			x := XYPoint[0]
			y := XYPoint[1]
			// fmt.Println("X: ", x, "Y: ", y)
			d = fmt.Sprintf("L%v %v ", x, y) + d
		}

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
