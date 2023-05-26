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
	isUsed *bool
	Color  color.Color
	// StartLines *[]*[]*[2]int
	// EndLines   *[]*[]*[2]int
	//////////////////////////////////
	// CurrentLineIndex int
	// LastY            int
	///////////////////////////////////////
	// Lines     *[]*[]*[2]int
	// LineIndex int
	// CurrentY  int
	///////////////////////////////////////////
	LineChannel *func(Line *[]*[2]int) *[]*[2]int
	// GetPathLine *func(Line *[]*[2]int) *[]*[2]int
	CurrentY int
}

func (p *VectorPath) AddForkLine() *VectorPath {
	prevLineChannel := *p.LineChannel
	lineChannel := func(line *[]*[2]int) *[]*[2]int {

		return prevLineChannel(line)
	}
	// prevGetPathLine := *p.GetPathLine

	// *p.GetPathLine = func(Line *[]*[2]int) *[]*[2]int {
	// 	return prevGetPathLine(Line)
	// }

	return &VectorPath{
		LineChannel: &lineChannel,
		CurrentY:    -1,
	}

}
func (p *VectorPath) AddStart(columX int, rowY int) {
	prevLineChannel := *p.LineChannel
	*p.LineChannel = func(line *[]*[2]int) *[]*[2]int {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*line = append(*line, &move1, &move2)
		return prevLineChannel(line)
	}

}
func (p *VectorPath) AddEnd(columX int, rowY int) {
	prevLineChannel := *p.LineChannel
	*p.LineChannel = func(line *[]*[2]int) *[]*[2]int {
		prevLineChannel(line)
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*line = append(*line, &move1, &move2)
		return line
	}

}

// func (p *VectorPath) GetLinePath(rowY int) *VectorPath {
// 	if rowY == p.LastY {
// 		return &VectorPath{
// 			isUsed:           p.isUsed,
// 			Color:            p.Color,
// 			StartLines:       p.StartLines,
// 			EndLines:         p.EndLines,
// 			CurrentLineIndex: p.CurrentLineIndex + 1,
// 			LastY:            p.LastY,
// 		}
// 	} else {
// 		return p
// 	}

// }
// func (p *VectorPath) AddStart(columX int, rowY int) {
// 	startLines := *p.StartLines

// 	startLine := startLines[p.CurrentLineIndex]

// 	move1 := [2]int{columX, rowY}
// 	move2 := [2]int{columX, rowY + 1}

// 	*startLine = append(*startLine, &move1, &move2)
// 	p.LastY = rowY
// }
// func (p *VectorPath) AddEnd(columX int, rowY int) {
// 	endLines := *p.EndLines

// 	endLine := endLines[p.CurrentLineIndex]

// 	move1 := [2]int{columX, rowY}
// 	move2 := [2]int{columX, rowY + 1}

//		*endLine = append(*endLine, &move1, &move2)
//		p.LastY = rowY
//	}
func NewVectorPath(color color.Color) *VectorPath {

	isUsed := true

	lineChannel := func(line *[]*[2]int) *[]*[2]int {
		return line
	}
	return &VectorPath{
		isUsed:      &isUsed,
		Color:       color,
		LineChannel: &lineChannel,
		CurrentY:    -1,
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

//	func (v *VectorImage) MoveEnd(x int, y int) (float64, float64) {
//		return float64(x+1) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
//	}

func (v *VectorImage) ImageVector() (image.Image, []*VectorPath) {
	colorDiffNum := float64(255 * v.ColorDiffPercent)
	paths := []*VectorPath{}

	// jobChannel := &utils.JobChannel[func()]{}
	getPathValue, addPathModifier := utils.NewJobChannel[[]*VectorPath]()
	newImage := image.NewRGBA(image.Rect(0, 0, v.Widget, v.Height))
	pathShapes := make([]*VectorPath, v.Widget)
	// startChannel := func(paths []*VectorPath) []*VectorPath {
	// 	return paths
	// }
	// chanel := &startChannel

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
			if rowY > 6 {
				continue
			}

			leftOk, left := indexValue(pathShapes, columnX-1)
			curOk, current := indexValue(pathShapes, columnX)

			isColorCurrent := curOk && current.Color == pixelColor
			isColorLeft := leftOk && left.Color == pixelColor

			equal := curOk && leftOk && current.isUsed == left.isUsed
			//////////////////////////////////////////////////////////////////

			if !equal && isColorCurrent && isColorLeft {
				current.ConcatLine(left)
			}

			if isColorLeft {
				pathShapes[columnX] = left
				current = left
			} else if !isColorCurrent {
				current = NewVectorPath(pixelColor)
				isUsed := current.isUsed
				pathShapes[columnX] = current

				addPathModifier(func(paths []*VectorPath) []*VectorPath {

					col := color.RGBA{
						0,
						0,
						0,
						255,
					}

					if col == current.Color && *isUsed {
						return append(paths, current)
					}
					return paths
					// if *isUsed {
					// 	paths = append(paths, current)
					// }
				})
			}

			if !isColorLeft {
				if leftOk {
					left.AddEnd(columnX, rowY)
					left.AddCrossLine(columnX, rowY)
				}
				// if current.LastY == rowY && isColorCurrent {
				// if current.LastY == rowY {
				// 	current = current.DuplicateLine()
				// 	pathShapes[columnX] = current
				// 	// current.LastY = rowY
				// 	fmt.Println("AFTER: ", current.Color, current.CurrentLineIndex)
				// }
				current.AddCrossLine(columnX, rowY)
				current.AddStart(columnX, rowY)
			} else if columnX == (v.Widget - 1) {
				current.AddCrossLine(columnX+1, rowY)
				current.AddEnd(columnX+1, rowY)
			}
		}
	}
	// defer
	// wg.Done()
	// ProcessEnd.Wait()
	// paths :=

	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	// b, _ := json.MarshalIndent(&paths, "", "  ")
	// fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDD", len(paths), string(b))

	// for _, paa := range paths {
	// 	data, _ := json.Marshal(paa)
	// 	fmt.Println("RRRRRRRRRRRRRRRRRRR", len(paths), string(data))
	// 	// var prettyJSON bytes.Buffer
	// 	// json.Indent(&prettyJSON, data, "", "\t")
	// 	// fmt.Println("RRRRRRRRRRRRRRRRRRR", len(paths), string(prettyJSON.Bytes()))
	// }
	fmt.Println("ENDDDDDDDDDDDDDDDDDDDDDDDD", len(paths))
	return newImage, (*getPathValue)([]*VectorPath{})
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
		data3, _ := json.Marshal(path.Line)

		fmt.Println("LINE: ", string(data3))
		d := ""
		// startLines := *path.StartLines
		// endLines := *path.EndLines
		for _, XYPoint := range *path.Line {
			x := XYPoint[0]
			y := XYPoint[1]
			// 		// fmt.Println("X: ", x, "Y: ", y)
			d = d + fmt.Sprintf("L%v %v ", x, y)

		}
		// for index, startLine := range startLines {
		// fmt.Println("index: ", index)
		// if index > 1 {
		// 	continue
		// }
		// ok, pEndLineAddr := indexValue(endLines, index)

		// data2, _ := json.Marshal(startLine)
		// fmt.Println("START LINE: ", string(data2))
		// data1, _ := json.Marshal(*pEndLineAddr)
		// fmt.Println("End LINE: ", string(data1))
		// if ok {
		// 	endLine := *pEndLineAddr
		// 	for index, _ := range endLine {
		// 		XYPoint := endLine[len(endLine)-1-index]
		// 		x := XYPoint[0]
		// 		y := XYPoint[1]
		// 		// fmt.Println("X: ", x, "Y: ", y)
		// 		d = d + fmt.Sprintf("L%v %v ", x, y)
		// 	}
		// }

		// startLine := *startLine
		// for _, XYPoint := range startLine {
		// 	x := XYPoint[0]
		// 	y := XYPoint[1]
		// 	// fmt.Println("X: ", x, "Y: ", y)
		// 	// if index1 == 0 {

		// 	// 	d = fmt.Sprintf("M%v %v ", x, y) + d
		// 	// } else {

		// 	d = d + fmt.Sprintf("L%v %v ", x, y)
		// 	// }
		// }

		// }

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
