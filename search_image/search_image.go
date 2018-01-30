// search_image takes an obfuscation layer, a background and the 3 key images (from the pre-processed folder) that make up the obfuscation then
// It builds a new single image starting with the obfuscation, then it transparently blits the background over the top 
// Then it transparently blits the 3 keys over that ensuring they do not overlap

// Note: This code relies on "go get github.com/nfnt/resize"

package main

import (
	FIU "FindIt/FIUtils"
    "github.com/nfnt/resize"
    "math/rand"
    "time"
    "image"
    "image/draw"
    "image/color"
    "os"
    "path/filepath"
    "strings"
    "strconv"
)


var key_path, backgrounds_path string
var paths FIU.OrigDest
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var num_keys int
var bg_transparency uint8

func initVars() {
    //@TODO: Load this from a config file and check orig and dest are different
    key_path = filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/keys")
    backgrounds_path = filepath.Join(FIU.FindIt_path, "/images_processed/backgrounds")

    paths = FIU.OrigDest {
       filepath.Join(FIU.FindIt_path, "/images_processed/obfuscator"), 
       filepath.Join(FIU.FindIt_path, "/images_processed/search"), 
    }

    num_keys = 3
    bg_transparency = 150
}

func main() {
    //FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

    // Load a random obfuscator
    obfuscator_files := FIU.ListFilenames(paths.Orig)
    obfuscator_file := obfuscator_files[rng.Intn(len(obfuscator_files))].Name()
    obfuscator_img := FIU.LoadImage( filepath.Join(paths.Orig, obfuscator_file) )
    FIU.Trace.Println("Obfuscator: " + obfuscator_file)

    // Load obfuscator into RGBA image
    r := image.Rect(0, 0, FIU.Out_width, FIU.Out_height)
    output := image.NewRGBA(r)
    draw.Draw(output, r, obfuscator_img, image.ZP, draw.Src)

    // Load a random background
    background_files := FIU.ListFilenames(backgrounds_path)
    background_file := background_files[rng.Intn(len(background_files))].Name()
	background_img := FIU.LoadImage( filepath.Join(backgrounds_path, background_file))
	FIU.Trace.Println("Background: " + background_file)

	// Transparent blit the background over the obfuscator
	draw.DrawMask(output, r, background_img, image.ZP, &image.Uniform{color.RGBA{0, 0, 0, bg_transparency}}, image.ZP, draw.Over)

	// Split the chosen obfuscator filename to select the 3 key files
	file_glob := filepath.Base(obfuscator_file)
    file_glob = strings.TrimSuffix(file_glob, filepath.Ext(file_glob))    
    key_files := strings.Split(file_glob, "_")
	FIU.Trace.Println("Keys: " + strings.Join(key_files, " "))

	var key_rects []image.Rectangle
	key_rects = make([]image.Rectangle, 0, num_keys)
    // Load 1 normailised version of each key file
    for i, v := range key_files {
    	keys, _ := filepath.Glob( filepath.Join(key_path, v + "_*") ) 
    	key := keys[rng.Intn(len(keys))]
    	// Load the preprocessed key file
    	key_img := FIU.LoadImage(key)

    	// If the key is too large scale it down 
    	bounds := key_img.Bounds().Max
    	if bounds.X > FIU.Out_width/4 || bounds.Y > FIU.Out_height/4 {
    		key_img = resize.Thumbnail(FIU.Out_width/4, FIU.Out_height/4, key_img, resize.Lanczos3)
    		bounds = key_img.Bounds().Max
    	}

    	// Keep randomly placing the rect till is doesn't overlap any others
    	// TODO: This is okay with small numbers of keys but may be impossible to solve with larger numbers and so hang
	    RetryRect:
    	// Random output rectangle
		op := image.Point{ rng.Intn(FIU.Out_width - bounds.X), rng.Intn(FIU.Out_height - bounds.Y) }
		r := image.Rectangle{op, op.Add(bounds)}
		for _, rect := range key_rects {
			if r.Overlaps(rect) {
				FIU.Trace.Println("Rect overlaped RetryRect")
				goto RetryRect
			}
		}		
		key_rects = append(key_rects, r)
		FIU.Trace.Println("key " + strconv.Itoa(i+1) + ": " + r.String())
    	
    	// Use an alphablend to transparent blit the key into the image
		draw.DrawMask(output, r, key_img, image.ZP, &image.Uniform{color.RGBA{0, 0, 0, 170}}, image.ZP, draw.Over)
    }

    // Save the image to the filename
    outfilename := filepath.Join(paths.Dest, filepath.Base(obfuscator_file))
    out_image := output.SubImage(image.Rectangle(output.Bounds()))
    FIU.SaveImage( &out_image, outfilename )
}
