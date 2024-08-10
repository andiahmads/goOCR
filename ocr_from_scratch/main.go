package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// Convert image to grayscale
func converTograyScale(img image.Image) *image.Gray {
	bound := img.Bounds()
	grayImage := image.NewGray(bound)

	for y := bound.Min.Y; y < bound.Max.Y; y++ {
		for x := bound.Min.X; x < bound.Max.X; x++ {
			originaColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originaColor)
			grayImage.Set(x, y, grayColor)
		}
	}
	return grayImage
}

// Binarize the grayscale image
func binarize(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	binaryImage := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := img.GrayAt(x, y)
			if grayColor.Y > threshold {
				binaryImage.SetGray(x, y, color.Gray{255})
			} else {
				binaryImage.SetGray(x, y, color.Gray{0})
			}
		}
	}
	return binaryImage
}

// Segment lines from the image
func segmentLines(img *image.Gray) []image.Rectangle {
	bounds := img.Bounds()
	var lineBounds []image.Rectangle
	inLine := false
	startY := 0

	minLineHeight := 5 // Minimum height to consider it a line of text
	// maxLineGap := 2    // Maximum gap between lines in pixels

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		isLine := false
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := img.GrayAt(x, y)
			if grayColor.Y < 128 { // Consider as part of text line
				isLine = true
				break
			}
		}

		if isLine && !inLine {
			// Start of a new line
			inLine = true
			startY = y
		} else if !isLine && inLine {
			// End of a line
			if y-startY >= minLineHeight {
				lineBounds = append(lineBounds, image.Rect(bounds.Min.X, startY, bounds.Max.X, y))
			}
			inLine = false
		}
	}

	return lineBounds
}

// Extract features from the character bound
func extractFeatures(img *image.Gray, charbound image.Rectangle) []float64 {
	features := []float64{}
	countBlackPixel := 0

	for y := charbound.Min.Y; y < charbound.Max.Y; y++ {
		for x := charbound.Min.X; x < charbound.Max.X; x++ {
			if img.GrayAt(x, y).Y == 0 {
				countBlackPixel++
			}
		}
	}
	aspectRatio := float64(charbound.Dx()) / float64(charbound.Dy())
	features = append(features, aspectRatio, float64(countBlackPixel))

	// Log extracted features
	log.Printf("Extracted features for rectangle %v: %v\n", charbound, features)
	return features
}

// Classify character based on features
// func classifyCharacter(features []float64) string {
// 	// Adjust classification logic as needed
// 	if features[0] > 1.5 {
// 		return "I"
// 	} else {
// 		return "O"
// 	}
// }

func classifyCharacter(features []float64) string {
	aspectRatio := features[0]
	countBlackPixel := features[1]

	if aspectRatio > 2.0 && countBlackPixel > 150 {
		return "I"
	} else if aspectRatio < 2.0 && countBlackPixel > 100 {
		return "O"
	} else {
		return "Unknown"
	}
}

// Use Tesseract to extract text from an image file
func extractTextUsingTesseract(filename string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage(filename)
	if err != nil {
		return "", err
	}

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}

func main() {
	file, err := os.Open("./testing.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	// Convert to grayscale and binarize
	grayImg := converTograyScale(img)
	binaryImage := binarize(grayImg, 128)

	// Save binary image for debugging
	outputFile, err := os.Create("binaryImage.png")
	if err != nil {
		log.Fatalf("Error creating image: %v", err)
	}
	defer outputFile.Close()
	png.Encode(outputFile, binaryImage)

	// Extract text from the binary image using Tesseract
	extractedText, err := extractTextUsingTesseract("binaryImage.png")
	if err != nil {
		log.Fatalf("Error extracting text: %v", err)
	}

	fmt.Println("Extracted Text:")
	fmt.Println(extractedText)

	// Segment lines
	lines := segmentLines(binaryImage)
	var recognizedText strings.Builder

	// Process each line
	// for _, line := range lines {
	// 	features := extractFeatures(binaryImage, line)
	// 	character := classifyCharacter(features)
	// 	log.Printf("Recognized character: %s", character)
	// }
	// log.Println("Processing complete")

	for _, line := range lines {
		lineImg := binaryImage.SubImage(line).(*image.Gray)
		x := line.Min.X
		for x < line.Max.X {
			charWidth, charHeight := 10, 20 // Assume fixed size for character bounding box
			charRect := image.Rect(x, line.Min.Y, x+charWidth, line.Min.Y+charHeight)

			// Ensure bounding box is within image bounds
			if charRect.Max.X > line.Max.X {
				charRect.Max.X = line.Max.X
			}
			if charRect.Max.Y > line.Max.Y {
				charRect.Max.Y = line.Max.Y
			}

			// Extract features and classify character
			features := extractFeatures(lineImg, charRect)
			fmt.Printf("Extracted features for rectangle %v: %v\n", charRect, features)
			character := classifyCharacter(features)
			fmt.Printf("Recognized character: %s\n", character)
			recognizedText.WriteString(character)

			x += charWidth // Move to the next character position
		}
		recognizedText.WriteString(" ") // Space between lines
	}

	log.Println("Complete:", recognizedText.String())
}
