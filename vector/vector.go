package vector

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"reflect"
	"vectoral/utils"
)

type VectorPath struct {
	isUsed     *bool
	Color      color.Color
	StartLines *[]*[]*[2]int
	EndLines   *[]*[]*[2]int
	//////////////////////////////////
	// CurrentLine             int
	CurrentStartY           int
	CurrentStartYLindeIndex int
	CurrentEndY             int
	CurrentEndYLineIndex    int
}

func (p *VectorPath) AddStart(columX int, rowY int) {
	if p.CurrentStartY == rowY {
		p.CurrentStartYLindeIndex = p.CurrentStartYLindeIndex + 1
	} else {
		p.CurrentStartYLindeIndex = 0
	}
	startLines := *p.StartLines
	var startLine *[]*[2]int

	if len(startLines) <= p.CurrentStartYLindeIndex {
		startLine = &[]*[2]int{}
		*p.StartLines = append(*p.StartLines, startLine)
	} else {
		startLine = startLines[p.CurrentStartYLindeIndex]
	}

	// if p.CurrentStartY != rowY {
	move1 := [2]int{columX, rowY}
	move2 := [2]int{columX, rowY + 1}

	*startLine = append(*startLine, &move1, &move2)
	p.CurrentStartY = rowY
	// }
}
func (p *VectorPath) AddEnd(columX int, rowY int) {
	if p.CurrentEndY == rowY {
		p.CurrentEndYLineIndex = p.CurrentEndYLineIndex + 1
	} else {
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

	// if p.CurrentEndY == rowY {
	// 	endLineVal := *endLine

	// 	move1 := [2]int{columX, rowY}
	// 	move2 := [2]int{columX, rowY + 1}
	// 	endLineVal = append(endLineVal[:len(endLineVal)-2], &move1, &move2)
	// 	*endLine = endLineVal
	// } else {
	move1 := [2]int{columX, rowY}
	move2 := [2]int{columX, rowY + 1}

	*endLine = append(*endLine, &move1, &move2)
	p.CurrentEndY = rowY
	// }
}

func (p *VectorPath) Concat(p2 *VectorPath) {
	// return
	// p2 არის სტარტი
	// p დასასრული
	*p2.isUsed = false
	// pEndLines := *p.EndLines
	// pLastLine := pEndLines[len(pEndLines)-1]
	// pStartLine := pStartLines[p.CurrentStartYLindeIndex]
	p2StartLines := *p2.StartLines
	p2EndLines := *p2.EndLines
	//////////////////////////// START LINE
	startLine := new([]*[2]int) // p2
	for index := len(p2StartLines) - 1; index >= 0; index-- {
		p2StartLine := *p2StartLines[index]

		ok, p2EndLineAddr := indexValue(p2EndLines, index)
		if ok {
			p2EndLine := *p2EndLineAddr
			// p2EndLine := *p2EndLines[index]

			for i := len(p2EndLine) - 1; i >= 0; i-- {
				p2EndLineItem := p2EndLine[i]
				*startLine = append(*startLine, p2EndLineItem)
				// 	// 	*pStartLine = append(*pStartLine, p2EndLineV[i])
			}
		}
		// line = append(line, p2EndLine...)
		*startLine = append(*startLine, p2StartLine...)

	}
	//////////////////////////// END LINE
	endLine := new([]*[2]int) // p
	pStartLines := *p.EndLines
	pEndLines := *p.StartLines
	for index := len(pStartLines) - 1; index >= 0; index-- {
		pStartLine := *pStartLines[index]
		// pEndLine := *pEndLines[index]
		ok, pEndLineAddr := indexValue(pEndLines, index)
		if ok {
			pEndLine := *pEndLineAddr
			for i := len(pEndLine) - 1; i >= 0; i-- {
				p2EndLineItem := pEndLine[i]
				*endLine = append(*endLine, p2EndLineItem)
				// 	// 	*pStartLine = append(*pStartLine, p2EndLineV[i])
			}

		}

		// line = append(line, p2EndLine...)
		*endLine = append(*endLine, pStartLine...)

	}
	// *[]*[]*[2]int
	// eee :=
	p.CurrentStartY = -1
	// p.CurrentStartYLindeIndex = 0
	p.CurrentEndY = -1
	// p.CurrentEndYLineIndex = 0
	*p.StartLines = []*[]*[2]int{startLine}
	*p.EndLines = []*[]*[2]int{endLine}

	// startLineEl := *startLine
	// endLineEl := *endLine
	// data1, _ := json.Marshal(startLineEl[0])
	// data2, _ := json.Marshal(endLineEl[0])
	// fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	// fmt.Println(data1)
	// fmt.Println("EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	// fmt.Println(data2)
	*p2 = *p
}

func NewVectorPath(color color.Color) *VectorPath {
	startLines := []*[]*[2]int{}
	endLines := []*[]*[2]int{}
	isUsed := true
	return &VectorPath{
		isUsed:                  &isUsed,
		Color:                   color,
		StartLines:              &startLines,
		EndLines:                &endLines,
		CurrentStartY:           -1,
		CurrentStartYLindeIndex: 0,
		CurrentEndY:             -1,
		CurrentEndYLineIndex:    0,
	}
}

// func contains(elems []T, v T) bool {
func indexValue[T comparable](vectorPaths []T, index int) (bool, T) {
	count := len(vectorPaths)

	if index > -1 && index < count {
		el := vectorPaths[index]

		return !reflect.ValueOf(el).IsNil(), el
	}
	var e T
	return false, e
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
			// if rowY > 4 {
			// 	continue
			// }

			leftOk, left := indexValue(pathShapes, columnX-1)
			curOk, current := indexValue(pathShapes, columnX)

			equal := curOk && leftOk && current.isUsed == left.isUsed

			isColorCurrent := curOk && current.Color == pixelColor
			isColorLeft := leftOk && left.Color == pixelColor

			//////////////////////////////////////////////////////////////////
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

					// if col == current.Color && *isUsed {
					// 	paths = append(paths, current)
					// }
					if *isUsed {
						paths = append(paths, current)
					}
				})
			}

			if columnX == (v.Widget - 1) {
				current.AddEnd(columnX+1, rowY)
			}
			if !isColorLeft {
				if leftOk {
					left.AddEnd(columnX, rowY)
				}
				current.AddStart(columnX, rowY)
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

	for _, paa := range paths {
		data, _ := json.Marshal(paa)
		fmt.Println("RRRRRRRRRRRRRRRRRRR", len(paths), string(data))
		// var prettyJSON bytes.Buffer
		// json.Indent(&prettyJSON, data, "", "\t")
		// fmt.Println("RRRRRRRRRRRRRRRRRRR", len(paths), string(prettyJSON.Bytes()))
	}
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
		startLines := *path.StartLines
		endLines := *path.EndLines
		for index, startLine := range startLines {
			// if index > 0 {
			// 	continue
			// }
			startLine := *startLine
			for _, XYPoint := range startLine {
				x := XYPoint[0]
				y := XYPoint[1]
				// fmt.Println("X: ", x, "Y: ", y)
				// if index1 == 0 {

				// 	d = fmt.Sprintf("M%v %v ", x, y) + d
				// } else {

				d = d + fmt.Sprintf("L%v %v ", x, y)
				// }
			}
			endLine := *endLines[index]
			for index, _ := range endLine {
				XYPoint := endLine[len(endLine)-1-index]
				x := XYPoint[0]
				y := XYPoint[1]
				// fmt.Println("X: ", x, "Y: ", y)
				d = d + fmt.Sprintf("L%v %v ", x, y)
			}

		}

		r, g, b, a := path.Color.RGBA()
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
