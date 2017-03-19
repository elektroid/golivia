package imgResize

import (
	"gopkg.in/gographics/imagick.v1/imagick"
	"os"
	"io"
"log"
)



// file copy done with io.Copy
func copy(input string, output string) error {
	// open files r and w
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(output)
	if err != nil {
		return err
	}
	defer w.Close()

	// do the actual work
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

// ugly map of input/output formats
func outputFormat(inputFormat string) string {
	if inputFormat == "png" {
		return "png"
	}
	if inputFormat == "gif" {
		return "gif"
	}
	return "jpeg"
}

func MakeMini(targetWidth uint, targetHeight uint, quality int, inputFile string, outputFile string) error {
	
	imagick.Initialize()
	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()

	err = mw.ReadImage(inputFile)
	if err != nil {
		panic(err)
	}

	// Get original logo size
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	baseProportion := float64(width) / float64(height)

	log.Printf("proportion: %f - %+v", baseProportion, baseProportion)
	log.Printf("target w %d h %d\n", targetWidth, targetHeight)
	log.Printf("real w %d h %d\n", width, height)

	if width<targetWidth{
		targetWidth=width
	}

	if height<targetHeight{
		targetHeight=height
	}

	log.Printf("final target iwth p %f w %d h %d\n", baseProportion, targetWidth, targetHeight) 
	targetWidth = uint(float64(baseProportion) * float64(targetHeight))
	log.Printf("final corrected target w %d h %d\n", targetWidth, targetHeight) 

	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(targetWidth, targetHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		panic(err)
	}

	 // open output file
    fo, err := os.Create(outputFile)
    if err != nil {
        panic(err)
    }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()


	err = mw.WriteImageFile(fo)
	return err
}
