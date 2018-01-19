// Obfuscator takes sets of 3 keys and loads a number of their matching preprocessed 
// if then creates new images that are a combination of the 3 keys and pastes them randomly to an output image
// The output image is named using the 3 keys in alphabetical order


package main

import (
	FIU "FindIt/FIUtils"
    "os"
    "path/filepath"
    "image"
    "image/draw"
    "math/rand"
    "time"
    "sort"
    "strings"
    //"fmt"
)

var key_path string
var paths FIU.OrigDest
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))


var num_keys int
var num_blits int

func initVars() {
    //@TODO: Load this from a config file and check orig and dest are different
    key_path = filepath.Join(FIU.FindIt_path, "/images_processed/normalised/keys")

    paths = FIU.OrigDest {
       filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/keys"), 
       filepath.Join(FIU.FindIt_path, "/images_processed/obfuscator"), 
    }

    num_keys = 3
    num_blits = 20
}

func main() {
    //FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

    // Pick a random set of 3 keys
    files := getKeyFiles( FIU.ListFilenames(key_path) )
    var processed_keys []string
    processed_keys = make([]string, 0)
    for _, file := range files {
        file_glob := filepath.Base(file)
        file_glob = strings.TrimSuffix(file_glob, filepath.Ext(file_glob)) + "*"
        file_glob = filepath.Join(paths.Orig, file_glob)
        processed_files, _ := filepath.Glob(file_glob) // Ignore error as the only error is bad pattern
        processed_keys = append(processed_keys, processed_files...) 
    }

    // Generate the image
    output := image.NewRGBA( image.Rect(0,0,FIU.Out_width, FIU.Out_height) )

    rand_order := rng.Perm(len(processed_keys))
    for _, i := range rand_order {
        file := processed_keys[i]
        img := FIU.LoadImage(file)

        // A rectangle 1/4 the size of the image
        bounds := img.Bounds().Max
        rect_bounds := image.Point{ bounds.X / 4, bounds.Y / 4 }
      
        offset := image.Point{ bounds.X/10, bounds.Y/10 }
        sp := image.Point{ rng.Intn(bounds.X - rect_bounds.X - offset.X) + offset.X, rng.Intn(bounds.Y - rect_bounds.Y - offset.Y) + offset.Y }
        
        for i:=0; i<num_blits; i++ {
            // Take this random chunk of spurce imae and paste it into random points of the output image
            op := image.Point{ rng.Intn(FIU.Out_width) - rect_bounds.X/2, rng.Intn(FIU.Out_height) - rect_bounds.Y/2 }
            
            r := image.Rectangle{op, image.Point{op.X + rect_bounds.X, op.Y + rect_bounds.Y}}
            draw.Draw(output, r, img, sp, draw.Src)    
        } 
    }

    // Construct the output filename
    outfilename := ""
    for _, file := range files {
         filename := filepath.Base(file)
         filename = strings.TrimSuffix(filename, filepath.Ext(filename))
         outfilename += filename + "_" 
    }
    outfilename = strings.TrimSuffix(outfilename, "_")
    outfilename += ".png"
    outfilename = filepath.Join(paths.Dest, outfilename)

    // Save the image to the filename
    out_image := output.SubImage(image.Rectangle(output.Bounds()))
    FIU.SaveImage( &out_image, outfilename )
}

func getKeyFiles(filelist []os.FileInfo) []string {
    if len(filelist) < num_keys {
        FIU.Error.Printf("Not enough normalised keys. Found %d needed %d\n", len(filelist), num_keys)
        os.Exit(1)
    }

    // Generate num_keys unique indexes
    p := rng.Perm(len(filelist))

    var key_files = []string{}
    key_files = make([]string, num_keys)
    for i, v := range p[:num_keys] {
        key_files[i] = filelist[v].Name()
    }

    // Sort the key files
    sort.Slice(key_files, func(i, j int) bool { return key_files[i] < key_files[j] })
    FIU.Trace.Println("Keys: " + strings.Join(key_files, ", "))
    return key_files
}