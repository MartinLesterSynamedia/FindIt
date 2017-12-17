// Normalise converts images in the images/keys and images/backgrounds 
// folders and places them in the appropriate folder within processed_images

// Note: This code relies on "go get github.com/nfnt/resize"

package main

import (
	"image"
    "image/jpeg"
    _ "image/png"
    _ "image/gif"
    "github.com/nfnt/resize"
    "bytes"
	"log"
	"io"
    "io/ioutil"  
    "os"
)


// List the files in the images/keys and images/backgrounds folders
// For each image 
//     Load the image - If it fails to load then report a warning
//     Check the dimensions - scale to fit 400x400 maintaining aspect 
//	   Save the file as jpeg with moderate compression to 

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func initTrace(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

type OrigDest struct {
	orig, dest string
}

var paths map[string]OrigDest

var max_width uint
var max_height uint

func initVars() {
	//@TODO: Load this from a config file and check orig and dest are different
	paths = make(map[string]OrigDest)
	paths["keys"] = OrigDest {
		"../images/keys", "../images_processed/normalised/keys", 
	}

	paths["backgrounds"] = OrigDest {
		"../images/backgrounds", "../images_processed/normalised/backgrounds",
	}

	max_width = 400
	max_height = 400
}

func main() {
    //initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    initTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

	Info.Println("Normalising all initial images")

 	for _, folder := range []string{"keys", "backgrounds"} {
		files := listFilenames(paths[folder].orig)
		for _, file := range files {
			img := loadImage(paths[folder].orig + "/" + file.Name())
			if img != nil {
				img = resizeImage(&img)
				saveImageAsJpg(&img, paths[folder].dest + "/" + file.Name())
			    //@TODO: Count if there are enough key and backgrounds available 3 and 1
			}
		}
	}

}

// Use os.readdir to list the files in the folders
func listFilenames(path string) []os.FileInfo {
	Trace.Println("listFilenames(" + path + ")")

	files, err := ioutil.ReadDir(path)
    if err != nil {
        Error.Println(err)
        os.Exit(1)
    }

	return files
}

// Load an image irrespective of format. Images that fail to load produce warnings
func loadImage(filename string) image.Image {
	Trace.Println("loadImage(" + filename + ")")

	imgBuffer, err := ioutil.ReadFile(filename)

    if err != nil {
        Warning.Println("Unable to load '" + filename + "': " + err.Error())
 		return nil  
    }

    reader := bytes.NewReader(imgBuffer)

    img, formatname, err := image.Decode(reader) // <--- here

    if err != nil {
        Warning.Println("Unable to read '" + filename + "' of type " + formatname + ": " + err.Error())
 		return nil 
 	}

    Trace.Printf("Bounds : %d, %d", img.Bounds().Max.X, img.Bounds().Max.Y)

	return img
}

// resize the image maintaining aspect ratio
func resizeImage(img *image.Image) image.Image {
	Trace.Println("resizeImage()")
	newimg := resize.Thumbnail(max_width, max_height, *img, resize.Lanczos3)
	Trace.Printf("New Bounds : %d, %d", newimg.Bounds().Max.X, newimg.Bounds().Max.Y)
	return newimg
}

// Save the image as a jpeg and save a bit more space
// TODO: Fix the bug here that means the file is saved with the wrong extension
func saveImageAsJpg(img *image.Image, filename string) {
	Trace.Println("saveImageAsJpg(" + filename + ")")

    out, err := os.Create(filename)
    if err != nil {
    	Warning.Println("Unable to create file '" + filename + "': " + err.Error())
    }
    var opt jpeg.Options
    opt.Quality = 80

	err = jpeg.Encode(out, *img, &opt)
	if err != nil {
    	Warning.Println("Unable to write normalised jpeg '" + filename + "': " + err.Error())
    }
}
