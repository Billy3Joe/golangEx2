package main

/*
You are tasked with implementing a simple parallel image processing program. The program takes a list of image file paths as input and applies a grayscale filter to each image concurrently using Go routines and channels. The processed images are then saved to an output directory.

DETAILS:
In the applyGrayscale function:
- you need to iterate over the pixels of the input image and apply the grayscale formula to each pixel. You can use the color.Gray type to set the pixel values in the grayscale image.

In the processImage function:

- Open the input image file using os.Open.
- Apply the grayscale filter using the applyGrayscale function.
- Create the output file using os.Create.
- Encode and save the grayscale image to the output file using the appropriate image format package.

In the main function:
- Iterate over the results channel and print the processing status for each image.
*/

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"sync"
)

// Define a function to apply a grayscale filter to an image
func applyGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	// TODO: Iterate over the image pixels and apply the grayscale formula
	// Tips: iterate over the bounds. bounds.Min.Y and bounds.Max.Y indicate pixel-width and
	// bouns.Min.X and bounds.Max.X indicate pixel-height.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Tip: to get one pixel, img.At(x, y)
			pixel := img.At(x, y)
			// Tip: To convert one pixel to gray: color.GrayModel.Convert(pixel)
			grayColor := color.GrayModel.Convert(pixel).(color.Gray)
			// Tip: To set every pixel on an image: gayImg.Set(x, y, pixel)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

// Tips: results channel will have string messages indicating Success or failure
// Sucess: "Processed file at <filepath>"
// Error: "Error when processing <inputfile>"
func processImage(inputPath, outputPath string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	// TODO: Open the input image file using os.Open(filename)
	inputFile, err := os.Open(inputPath)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer inputFile.Close()

	// TODO: The the image object by decoding the file. image.Decode(inputFile)
	img, _, err := image.Decode(inputFile)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	// TODO: Apply grayscale filter using the applyGrayscale function
	grayImg := applyGrayscale(img)

	// TODO: Create the output file
	// It should be in the same folder. os.Create(filepath.Join(outputPath, filepath.Base(inputPath)))

	outputFilePath := filepath.Join(outputPath, filepath.Base(inputPath))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer outputFile.Close()

	// TODO: Encode and save the grayscale image
	// Check the extension of the image.
	switch filepath.Ext(inputPath) {
	// if filepath.Ext(inputFile) == ".png" -> png.Encode(outputfile, grayimg)
	case ".png":
		err = png.Encode(outputFile, grayImg)
	// if filepath.Ext(inputFile) == ".jpg" -> jpg.Encode(outputfile, grayimg)
	case ".jpg":
		err = jpeg.Encode(outputFile, grayImg, nil)
	default:
		results <- fmt.Sprintf("Error when processing %s: Unsupported image format", inputPath)
		return
	}
	// TODO: Send a message to the results channel indicating success or failure
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	results <- fmt.Sprintf("Processed file at %s", inputPath)
}

func main() {
	inputPaths := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
	outputPath := "output/"

	results := make(chan string)
	var wg sync.WaitGroup

	for _, inputPath := range inputPaths {
		wg.Add(1)
		go processImage(inputPath, outputPath, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	fmt.Println("Processing complete.")
}
