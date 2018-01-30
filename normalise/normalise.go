// Normalise converts images in the images/keys and images/backgrounds 
// folders and places them in the appropriate folder within processed_images

// Note: This code relies on "go get github.com/nfnt/resize"

package main

import (
	FIU "FindIt/FIUtils"
    "github.com/nfnt/resize"
    "image"
    "os"
    "path/filepath"
    "strings"
)


// List the files in the images/keys and images/backgrounds folders
// For each image 
//     Load the image - If it fails to load then report a warning
//     Check the dimensions - scale to fit 400x400 maintaining aspect 
//	   Save the file as png 


var max_width uint
var max_height uint

func initVars() {
	//@TODO: Load this from a config file and check orig and dest are different
	FIU.Paths = make(map[string]FIU.OrigDest)
	FIU.Paths["keys"] = FIU.OrigDest {
		filepath.Join(FIU.FindIt_path, "/images/keys"),
		filepath.Join(FIU.FindIt_path, "/images_processed/normalised/keys"), 
	}

	FIU.Paths["backgrounds"] = FIU.OrigDest {
		filepath.Join(FIU.FindIt_path, "/images/backgrounds"),
		filepath.Join(FIU.FindIt_path, "/images_processed/normalised/backgrounds"),
	}

	max_width = 400
	max_height = 400
}

func main() {
    //FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

	FIU.Info.Println("Normalising all initial images")

 	for _, folder := range []string{"keys", "backgrounds"} {
		files := FIU.ListFilenames(FIU.Paths[folder].Orig)
		for _, file := range files {
			img := FIU.LoadImage( filepath.Join(FIU.Paths[folder].Orig, file.Name()) )
			if img != nil {
				//@TDOD: Replace underscore in source filenames with minus as we use underscore to separate the key files when making an obfuscator filename.
				img = resizeImage(&img)
				dest_filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				//@@TODO: Make sure the length of the filename is "sensible". Truncate it and make sure its unique. Alternatively just genetate all normalised images as unique IDs
				FIU.SaveImage(&img, filepath.Join(FIU.Paths[folder].Dest,  dest_filename + ".png"))
			    //@TODO: Count if there are enough key and backgrounds available 3 and 1
			}
		}
	}

}

// resize the image maintaining aspect ratio
func resizeImage(img *image.Image) image.Image {
	FIU.Trace.Println("resizeImage()")
	newimg := resize.Thumbnail(max_width, max_height, *img, resize.Lanczos3)
	FIU.Trace.Printf("New Bounds : %d, %d", newimg.Bounds().Max.X, newimg.Bounds().Max.Y)
	return newimg
}


