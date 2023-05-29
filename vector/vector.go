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
	isUsed              *bool
	Color               color.Color
	FillPathLine        *func(Line *[]*[2]int) *[]*[2]int
	GetPreviousEndPoint func() (int, int)
	StartLineChannel    *func(Line *[]*[2]int) *[]*[2]int
	EndLineChannel      *func(Line *[]*[2]int) *[]*[2]int

	CurrentY int
}

func (p *VectorPath) ForkLine() *VectorPath {

	// PrevColumX, PrevRowY := p.GetPreviousEndPoint()
	prevStartLineChannel := p.StartLineChannel
	prevEndLineChannel := p.EndLineChannel

	startChannel := func(line *[]*[2]int) *[]*[2]int {
		// *line = append(*line, &[2]int{PrevColumX, PrevRowY}, &[2]int{PrevColumX, PrevRowY + 1})

		(*prevStartLineChannel)(line)
		(*prevEndLineChannel)(line)
		return line
	}
	endChannel := func(line *[]*[2]int) *[]*[2]int {
		return line
	}

	startChannelAddr := &startChannel
	endChannelAddr := &endChannel
	prevFillPathLine := *p.FillPathLine

	*p.FillPathLine = func(line *[]*[2]int) *[]*[2]int {
		prevFillPathLine(line)
		(*startChannelAddr)(line)
		(*endChannelAddr)(line)
		return line
	}

	return &VectorPath{
		isUsed:           p.isUsed,
		Color:            p.Color,
		StartLineChannel: startChannelAddr,
		EndLineChannel:   endChannelAddr,
		// LineChannel:     lineChannelAddr,/
		FillPathLine: p.FillPathLine,
		CurrentY:     -1,
	}

}
func (p *VectorPath) ConcatLine(p2 *VectorPath) {
	return
	*p2.isUsed = false
	// pStartLineChannel := *p.StartLineChannel
	// pEndLineChannel := *p.EndLineChannel
	// pFillPathLine := *p.FillPathLine
	//////

	/////////////////////////////////////////////////////////////////////////////////////
	// pStartLineChannel := *p.StartLineChannel // start
	// pStartLineChannel := *p.EndLineChannel  // End
	// p2StartLineChannel := *p2.StartLineChannel // start
	// p2StartLineChannel := *p2.EndLineChannel  // End

	pStartLineChannel := *p.StartLineChannel   // start
	pEndLineChannel := *p.EndLineChannel       // End
	p2StartLineChannel := *p2.StartLineChannel // start
	p2EndLineChannel := *p2.EndLineChannel     // End
	////
	startLineChannel := func(line *[]*[2]int) *[]*[2]int {
		pEndLineChannel(line)
		p2EndLineChannel(line)
		pStartLineChannel(line)
		p2StartLineChannel(line)
		data, _ := json.Marshal(line)
		fmt.Println(p.Color, string(data))

		// pEndLineChannel(line)

		// pStartLineChannel(line)
		// p2StartLineChannel(line)

		// p2EndLineChannel(line)

		return line
	}
	endLineChannel := func(line *[]*[2]int) *[]*[2]int {
		return line
	}
	// startLineChannelAddr := &startLineChannel
	*p.StartLineChannel = startLineChannel
	*p.EndLineChannel = endLineChannel
	// pStartLineChannel := *p.StartLineChannel
	// p2EndLineChannel := *p2.EndLineChannel
	*p.FillPathLine = func(line *[]*[2]int) *[]*[2]int {
		(*p.StartLineChannel)(line)
		(*p.EndLineChannel)(line)
		return line
	}
	// *p.StartLineChannel = *p2.StartLineChannel
	// e2 :=
	fillPathLine := func(line *[]*[2]int) *[]*[2]int {
		data, _ := json.Marshal(line)
		fmt.Println("math.Random(1)", string(data))
		return line
	}
	*p2.FillPathLine = fillPathLine
	// *p2.StartLineChannel = fillPathLine
	// *p2.EndLineChannel = fillPathLine
	// *p2.EndLineChannel = func(line *[]*[2]int) *[]*[2]int {
	// 	data, _ := json.Marshal(line)
	// 	fmt.Println("math.Random(2)", string(data))
	// 	return line
	// }
	*p2 = *p
}

func (p *VectorPath) AddStart(columX int, rowY int) {
	prevStartLineChannel := *p.StartLineChannel

	*p.StartLineChannel = func(line *[]*[2]int) *[]*[2]int {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}
		prevStartLineChannel(line)

		*line = append(*line, &move1, &move2)

		return line
	}

	p.CurrentY = rowY
}
func (p *VectorPath) AddEnd(columX int, rowY int) {
	prevEndLineChannel := *p.EndLineChannel

	*p.EndLineChannel = func(line *[]*[2]int) *[]*[2]int {
		move1 := [2]int{columX, rowY}
		move2 := [2]int{columX, rowY + 1}

		*line = append(*line, &move2, &move1)
		prevEndLineChannel(line)

		return line
	}

	p.GetPreviousEndPoint = func() (int, int) {
		*p.EndLineChannel = prevEndLineChannel
		return columX, rowY
	}

	p.CurrentY = rowY
}

func NewVectorPath(color color.Color) *VectorPath {
	isUsed := true

	startChannel := func(line *[]*[2]int) *[]*[2]int {
		return line
	}
	endChannel := func(line *[]*[2]int) *[]*[2]int {
		return line
	}

	startChannelAddr := &startChannel
	endChannelAddr := &endChannel
	fillPathLine := func(line *[]*[2]int) *[]*[2]int {
		(*startChannelAddr)(line)
		(*endChannelAddr)(line)
		return line
	}
	fillPathLineAddr := &fillPathLine
	return &VectorPath{
		isUsed:           &isUsed,
		Color:            color,
		FillPathLine:     fillPathLineAddr,
		StartLineChannel: startChannelAddr,
		EndLineChannel:   endChannelAddr,
		CurrentY:         -1,
	}
}

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
			// if rowY > 6 {
			// 	continue
			// }

			leftOk, left := indexValue(pathShapes, columnX-1)
			curOk, current := indexValue(pathShapes, columnX)

			isColorCurrent := curOk && current.Color == pixelColor
			isColorLeft := leftOk && left.Color == pixelColor

			equal := curOk && leftOk && current.FillPathLine == left.FillPathLine
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

					if *isUsed {
						paths = append(paths, current)
					}
					return paths
				})
			}

			if !isColorLeft {
				if leftOk {
					left.AddEnd(columnX, rowY)
				}
				if current.CurrentY == rowY {
					current = current.ForkLine()
					pathShapes[columnX] = current

				}

				current.AddStart(columnX, rowY)
			}
			if columnX == (v.Widget - 1) {
				// fmt.Println(current.Color, columnX, rowY)
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
	fmt.Println("paths: ", len(paths))
	for _, path := range paths {
		if path == nil {
			continue
		}
		line := &[]*[2]int{}

		(*path.FillPathLine)(line)

		if len(*line) == 0 {
			continue
		}

		d := ""

		for _, XYPoint := range *line {
			x := XYPoint[0]
			y := XYPoint[1]
			// 		// fmt.Println("X: ", x, "Y: ", y)
			d = d + fmt.Sprintf("L%v %v ", x, y)

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
