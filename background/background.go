// background creates an image from the folder of preprocessed backgrounds.
// It picks 50 of those images and pastes them as is into a random place in an output image applying a random alpha blend

package main

import (
	FIU "FindIt/FIUtils"
    "os"
    "path/filepath"
    "image"
    "image/draw"
    //"image/color"
    "math/rand"
    "time"
    "strconv"
    //"strings"
    "fmt"
)

var alpha_path string
var paths FIU.OrigDest
var num_blits, num_files, num_alpha int

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))


func initVars() {
    //@TODO: Load this from a config file and check orig and dest are different
    alpha_path = filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/alpha")

    paths = FIU.OrigDest {
       filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/backgrounds"), 
       filepath.Join(FIU.FindIt_path, "/images_processed/backgrounds"), 
    }

    num_files = 1
    num_alpha = 1
    num_blits = 100
}

func main() {
    //FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

    // Load all the alphas we are going to user
	alphafilenames := getAlphaFiles( FIU.ListFilenames(alpha_path) )
	alphas := []image.Image{}
	alphas = make([]image.Image,0)
	for _, alpha := range alphafilenames {
		img := FIU.LoadImage(filepath.Join(alpha_path, alpha))
		alphas = append(alphas,  img)		
	}
	fmt.Println(alphas)

    // Generate the image
    output := image.NewRGBA( image.Rect(0,0,FIU.Out_width, FIU.Out_height) )
    files  := getBackgroundFiles( FIU.ListFilenames(paths.Orig) )

    for _, file := range files {
    	loaded := FIU.LoadImage(filepath.Join(paths.Orig, file))
    	b := loaded.Bounds()
    	// Load this image into an RGBA image
    	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    	draw.Draw(img, img.Bounds(), loaded, b.Min, draw.Src)

    	bounds := img.Bounds().Max
    	for i:=0; i<num_blits;  i++ {
    		// Random output rectangle
    		op := image.Point{ rng.Intn(FIU.Out_width) - bounds.X/2, rng.Intn(FIU.Out_height) - bounds.Y/2 }
    		r := image.Rectangle{op, op.Add(bounds)}

    		// Pick a random alpha channel
    		a := alphas[ rng.Intn( len(alphas) ) ]

    		draw.DrawMask(output, r, img, image.ZP, a, image.ZP, draw.Over) 
    	}
    }


    outfilename := strconv.FormatInt(time.Now().UnixNano(), 16)
    outfilename += ".png"

    // Save the image to the filename
    out_image := output.SubImage(image.Rectangle(output.Bounds()))
    FIU.SaveImage( &out_image, filepath.Join(paths.Dest, outfilename) )
}

func getBackgroundFiles(filelist []os.FileInfo) []string  {
    if len(filelist) < num_files {
        FIU.Error.Printf("Not enough normalised backgrounds. Found %d needed %d\n", len(filelist), num_files)
        os.Exit(1)
    }

    // Generate num_keys unique indexes
    p := rng.Perm(len(filelist))

    var files = []string{}
    files = make([]string, num_files)
    for i, v := range p[:num_files] {
        files[i] = filelist[v].Name()
    }

    return files
}

func getAlphaFiles(filelist []os.FileInfo) []string  {
    if len(filelist) < num_alpha {
        FIU.Error.Printf("Not enough alpha files. Found %d needed %d\n", len(filelist), num_alpha)
        os.Exit(1)
    }

    // Generate num_keys unique indexes
    p := rng.Perm(len(filelist))

    var files = []string{}
    files = make([]string, num_alpha)
    for i, v := range p[:num_alpha] {
        files[i] = filelist[v].Name()
    }

    return files
}