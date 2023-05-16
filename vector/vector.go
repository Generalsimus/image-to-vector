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
	isUsed         bool
	color          color.Color
	Lines          *[]*[]*[2]float64
	MoveLinesLeft  *[]*[2]float64
	MoveLinesRight *[]*[2]float64
	PosY               int
	PosYIndex        int
	// AssignLeftAt   *VectorPath
	// AssignRightAt  *VectorPath
}

//	func (p *VectorPath) Equal(el *VectorPath) bool {
//		return *p.index == *el.index
//	}
func (p *VectorPath) AddMove(x float64, y float64){ 
	if p.PosY != y { 
		p.PosYIndex = 0
	}
	line:=p.Lines[p.PosYIndex]
	if line == nil {
		
	}

	
	p.PosYIndex = p.PosYIndex+1
}
func (p *VectorPath) ColorEqual(color color.Color) bool {
	return p.color == color
}

func NewVectorPath(color color.Color) *VectorPath {

	moveLinesLeft := []*[2]float64{}
	moveLinesRight := []*[2]float64{}
	return &VectorPath{
		isUsed:         true,
		color:          color,
		MoveLinesLeft:  &moveLinesLeft,
		MoveLinesRight: &moveLinesRight,
	}
}

func (p *VectorPath) AddMoveLeft(x float64, y float64) {
	move := [2]float64{x, y}

	moveLinesLeft := append(*p.MoveLinesLeft, &move)
	*p.MoveLinesLeft = moveLinesLeft

	// move := [2]int{x, y}
	// moveLinesLeft := *p.MoveLinesLeft
	// lastIndex := len(moveLinesLeft) - 1

	// if lastIndex > 0 && moveLinesLeft[lastIndex][1] == y && x < moveLinesLeft[lastIndex][0] {
	// 	moveLinesLeft[lastIndex] = &move
	// } else {
	// 	moveLinesLeft = append(moveLinesLeft, &move)
	// }
	// *p.MoveLinesLeft = moveLinesLeft
	// p.PosY = y
}
func (p *VectorPath) AddMoveRight(x float64, y float64) {
	move := [2]float64{x, y}
	moveLinesRight := append(*p.MoveLinesRight, &move)
	*p.MoveLinesRight = moveLinesRight

	// move := [2]int{x, y}
	// moveLinesRight := *p.MoveLinesRight
	// lastIndex := len(moveLinesRight) - 1

	// if lastIndex > 0 && moveLinesRight[lastIndex][1] == y && x > moveLinesRight[lastIndex][0] {
	// 	moveLinesRight[lastIndex] = &move
	// } else {
	// 	moveLinesRight = append(moveLinesRight, &move)
	// }
	// *p.MoveLinesRight = moveLinesRight
	// p.PosY = y
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
	return float64(x) * v.OnePixelScaleX, float64(y) * v.OnePixelScaleY
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

	for row := 0; row < v.Height; row++ {
		for column := 0; column < v.Widget; column++ {
			r, g, b, a := v.Img.At(column, row).RGBA()

			red := math.Round(float64((r>>8))/colorDiffNum) * colorDiffNum
			green := math.Round(float64((g>>8))/colorDiffNum) * colorDiffNum
			blue := math.Round(float64((b>>8))/colorDiffNum) * colorDiffNum

			pixelColor := color.RGBA{
				uint8(red),
				uint8(green),
				uint8(blue),
				uint8(a),
			}
			newImage.Set(column, row, pixelColor)

			leftOk, left := PathInclude(pathShapes, column-1)
			curOk, current := PathInclude(pathShapes, column)

			equal := curOk && leftOk && *current == *left

			isColorCurrent := curOk && current.ColorEqual(pixelColor)
			isColorLeft := leftOk && left.ColorEqual(pixelColor)
			// equal := curOk && leftOk && *current == *left
			if isColorCurrent && isColorLeft && !equal {
				// fmt.Println("ASSIGN")
				// moveLinesLeft := append(*current.MoveLinesLeft, *left.MoveLinesLeft...)
				// moveLinesLeft = append(moveLinesLeft, *left.MoveLinesRight...)

				// moveLinesRight := append(*current.MoveLinesRight, *left.MoveLinesRight...)
				// *current.MoveLinesLeft = moveLinesLeft
				// *current.MoveLinesRight = moveLinesRight
				// index := *left.index

				// *left.MoveLinesLeSft = []*[2]int{}
				// *left.MoveLinesRight = []*[2]int{}

				// *left = *current
				// paths[index] = nil
				// equal = true
			}
			if isColorLeft {
				pathShapes[column] = left
				current = left
				equal = true
			} else if !isColorCurrent {
				current = NewVectorPath(pixelColor)
				pathShapes[column] = current
				jobChannel.AddJob(func() {
					if current.isUsed {
						paths = append(paths, current)
					}
			}

			if column == 0 {
				X, Y := v.MoveScale(column, row)
				current.AddMoveLeft(X, Y)
			} else if column == (v.Widget - 1) {
				X, Y := v.MoveScale(column, row)
				current.AddMoveRight(X, Y)
			} else if !equal {
				X, Y := v.MoveScale(column, row)
				if leftOk {
					left.AddMoveLeft(X, Y)
				}
				if curOk {
					current.AddMoveRight(X, Y)
				}
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

		moveLinesLeft := *path.MoveLinesLeft
		moveLinesRight := *path.MoveLinesRight
		sizeLeft := len(moveLinesLeft)
		sizeRight := len(moveLinesRight)
		if sizeLeft < 3 || sizeRight < 3 {
			continue
		}

		dLeft := ""
		for _, XYPoint := range moveLinesLeft {
			x := XYPoint[0] * float64(saveWidget)
			y := XYPoint[1] * float64(saveHeight)
			// strconv.FormatFloat(x, 'f', -1, 64)
			// strconv.FormatFloat(x, 'f', -1, 64)

			// fmt.Println("moveLinesLeft", x, y)
			dLeft = dLeft + fmt.Sprintf("L%v %v ", x, y)
		}

		dRight := ""
		for _, XYPoint := range moveLinesRight {
			x := XYPoint[0] * float64(saveWidget)
			y := XYPoint[1] * float64(saveHeight)

			// fmt.Println("moveLinesRight", x, y)
			dRight = fmt.Sprintf("L%v %v ", x, y) + dRight
		}

		// fmt.Println("D: ", dLeft)
		// fmt.Println("D: ", dRight)
		r, g, b, a := path.color.RGBA()
		color := fmt.Sprintf("rgba(%v,%v,%v,%v)", r>>8, g>>8, b>>8, a>>8)

		if _, err := f.Write([]byte(fmt.Sprintf("<path fill=\"%v\" d=\"M%v%vZ\" />\n", color, dLeft[1:], dRight))); err != nil {
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
