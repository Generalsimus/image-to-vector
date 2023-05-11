// package main

// import (
// 	"time"
// 	"vectoral/vector"

// 	"gocv.io/x/gocv"
// )

// func main() {
// 	webcam, _ := gocv.OpenVideoCapture(1)
// 	window := gocv.NewWindow("Hello")
// 	img := gocv.NewMat()

// 	///////////////////
// 	var fps time.Duration = 60
// 	framePerSecund := time.Second / fps
// 	for {
// 		time.Sleep(framePerSecund)
// 		webcam.Read(&img)
// 		image, _ := img.ToImage()
// 		vectorImg := vector.VectorImage{
// 			ColorDiffPercent: 0.1,
// 			Img:              image,
// 		}
// 		vectoredImg := vectorImg.ImageVector()

// 		updatedIMG, _ := gocv.ImageToMatRGBA(vectoredImg)
// 		// windoww.get
// 		// window.WaitKey(2)
// 		// fmt.Println("CLICK")
// 		window.IMShow(updatedIMG)

//			window.WaitKey(1)
//		}
//	}
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"vectoral/vector"
)

func main() {
	dirPath := "./images/"
	entries, _ := os.ReadDir(dirPath)

	for _, file := range entries {

		fileName := file.Name()
		path := dirPath + fileName

		fileInfo, _ := os.Stat(path)
		if fileInfo.IsDir() {
			continue
		}

		img, _ := getImageFromFilePath(path)
		vectorImg := vector.VectorImage{
			ColorDiffPercent: 0.5,
			Img:              img,
		}
		// vectorImg.ImageVector()
		fmt.Println("FILENAME: ", fileName)
		image, vector := vectorImg.ImageVector()
		vectorImg.SavePathsToSVGFile(vector, "./save/"+fileName+"A.svg")
		// fmt.Println("RES: ", len(vectorImg.ImageVector()))
		saveImageAt(image, "./save/"+fileName)
	}
	fmt.Println("PROCESS END")

}

func saveImageAt(image image.Image, path string) {
	var imageBuf bytes.Buffer
	png.Encode(&imageBuf, image)

	// Write to file.
	outfile, err := os.Create(path)
	if err != nil {
		// replace this with real error handling
		panic(err.Error())
	}
	defer outfile.Close()
	png.Encode(outfile, image)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _ := png.Decode(f)
	return image, err
}
