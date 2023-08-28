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


	/*1. Task to perform: Browse the pixels of the image and implement the formula for the conversion to grayscale.
      2. Note: It is advisable to browse the limits. Bounds.Min.Y and Bounds.Max.Y represent pixel width, while bouns.Min.X and bounds.Max.X represent pixel height.
	  3. To Do: Explore the pixels in the image and apply the method for converting to grayscale.
	  4. Advice: During the course, be sure to consider the terminals. The values Bounds.Min.Y and Bounds.Max.Y determine the width of the pixels, while bouns.Min.X and bounds.Max.X determine the height of the pixels.*/

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Get a pixel, use image.At(x, y).
			pixel := img.At(x, y)
			//To convert one pixel to gray: color.GrayModel.Convert(pixel)
			grayColor := color.GrayModel.Convert(pixel).(color.Gray)
			//To set every pixel on an image: gayImg.Set(x, y, pixel)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

/*
- Tips: Success or failure messages will be available in the results channel as strings.
- If successful: "The file was processed successfully: <file_path>"
- In the event of an error: "An error occurred while processing the file: <input_file>"
*/
func processImage(inputPath, outputPath string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	// Task to perform: Use the os.Open(filename) function to open the input image file.
	inputFile, err := os.Open(inputPath)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer inputFile.Close()

	// Task to perform: Use the image object by decoding the file. image.Decode(inputFile)
	img, _, err := image.Decode(inputFile)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	// Task to perform: Apply grayscale filter using the applyGrayscale function
	grayImg := applyGrayscale(img)

	// Task to perform: Create the output file
	// It should be in the same folder. os.Create(filepath.Join(outputPath, filepath.Base(inputPath)))

	outputFilePath := filepath.Join(outputPath, filepath.Base(inputPath))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}
	defer outputFile.Close()

	// Task to perform: Encode and save the grayscale image
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
	// Task to perform: Send a message to the results channel indicating success or failure
	if err != nil {
		results <- fmt.Sprintf("Error when processing %s: %s", inputPath, err)
		return
	}

	results <- fmt.Sprintf("Processed file at %s", inputPath)
}

func main() {
	inputPaths := []string{"file1.png", "file2.png", "file3.png"}
	outputPath := "output_pictures/"

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







