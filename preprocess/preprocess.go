// Preprocess takes the normalised images and generates a set of images from them that
// have been processed i.e. rotated, skewed, scaled, etc
// It puts the processed images in the appropriate folder within processed_images

// Note: This code use imagemagik to play with the images.

package main

import (
	FIU "FindIt/FIUtils"
    "os"
    "os/exec"
    "path/filepath"
    "math"
    "math/rand"
    "time"
    "strings"
    "fmt"
    "strconv"
)

// List the files in images_processed/normalised/keys and images_processed/normalised/backgrounds
// For each image 
//     Load the image - If it fails to load then report a warning
//     Generate a defined number of new images
//		   rotate, skew, resize
//	       Save the image to the appropriate folder with a generated name 


// ENHANCEMENT: This can be better implemented by updating https://github.com/quirkey/magick/blob/master/magick.go
// and enabling that to give access to Image.affineTransform() 
// It would be much faster but is notably more complex so for speed of implementation will just call 
// out to the imagemagik command line

var generated_files int
var max_scale int
var max_rotation int
var max_skew int
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type IMImage struct {
	Width, Height int
	Files FIU.OrigDest
}

type Vertex struct {
	X, Y int
}

type Triangle struct {
	P [3]Vertex
}

func (t Triangle) String() string {
	return fmt.Sprintf("%d,%d  %d,%d  %d,%d", t.P[0].X, t.P[0].Y, t.P[1].X, t.P[1].Y, t.P[2].X, t.P[2].Y )
}

type Affine struct {
	src, dst Triangle
}

func (a Affine) String() string {
	return a.src.String() + " -> " + a.dst.String()
}

func initVars() {
	//@TODO: Load this from a config file and check orig and dest are different
	FIU.Paths = make(map[string]FIU.OrigDest)
	FIU.Paths["keys"] = FIU.OrigDest {
		filepath.Join(FIU.FindIt_path, "images_processed/normalised/keys"),
		filepath.Join(FIU.FindIt_path, "/images_processed/keys"), 
	}

	FIU.Paths["backgrounds"] = FIU.OrigDest {
		filepath.Join(FIU.FindIt_path, "/images_processed/normalised/backgrounds"),
		filepath.Join(FIU.FindIt_path, "/images_processed/backgrounds"), 
	}

	generated_files = 10
	max_scale = 80
	max_rotation = 80
	max_skew = 60
}

func main() {
    //FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

	FIU.Info.Println("Preprcessing all normalised images")

 	for _, folder := range []string{"keys", "backgrounds"} {
		files := FIU.ListFilenames(FIU.Paths[folder].Orig)
		for _, file := range files {
			img := IMImage {
				0, 0,
				FIU.OrigDest {
					filepath.Join(FIU.Paths[folder].Orig, file.Name()),
					FIU.Paths[folder].Dest,
				},			
			}

			if loadImageInfo(&img) != nil {
				continue
			}

			src := Triangle{[3]Vertex{{0,0}, {img.Width,0}, {img.Width,img.Height}}}
			
			for i:=0; i<generated_files; i++ {
				img_cpy := img
				transform := calculateTransform(&img_cpy)
				aff := Affine{ src, transform }
				FIU.Trace.Printf("%d %s:\t%s", i, img_cpy.Files.Dest, aff.String())
			}
		}
	}
}

// Use imagemagik identify to get the width and height of the image
func loadImageInfo(img *IMImage) error {
	cmdName := "identify"
	cmdArgs := []string{"-format", "'%w,%h'", img.Files.Orig}
	var cmdOut []byte
	FIU.Trace.Println("Exec: " + cmdName + " " + strings.Join(cmdArgs, " ") )

	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()

	if err != nil {
		FIU.Warning.Println(cmdName + " " + strings.Join(cmdArgs, " ") + " " + err.Error())
		return err
	}

	FIU.Trace.Printf("Width,Height = %s", string(cmdOut))

	// Not sure why but converting []bytes to string adds single quotes, need to strip them
	out := strings.Split(strings.Trim(string(cmdOut), "'"), ",")
	if img.Width, err = strconv.Atoi(out[0]); err != nil {
		FIU.Warning.Println("Invalid width '" + out[0] + "' :" + err.Error())
		return err
	}
	if img.Height, err = strconv.Atoi(out[1]); err != nil {
		FIU.Warning.Println("Invalid height '" + out[1] + "' :" + err.Error())
		return err
	}

	return nil
}


// Distort the image using imagemagik convert
// convert <INPUT> -alpha set -virtual-pixel transparent  +distort Affine 'iX1,iY1 oX1,oY1  iX2,iY2 oX2,oY2  iX3,iY3 oX3,oY3' <OUTPUT>
// iXn,iYn describe an input triangle, oXn,oYn describe an output triangle
// how input triangle maps to output triangle is the effect applied to the image i.e. scaled, rotated, skewed, translated
// The image can be considered as a right angle Triangle iT inscribed within a circle iC centered at iX3/2,iY3/2 where radius R(iC) = sqrt((iX3/2 * iX3/2) + (iY3/2 * iY3/2)) 
// Create a new circle oC of radius R such that R(oC) <= R(iC), this will allow for scale 
// Place the first virtex anywhere on the new circle, this rotates the image (probably +/-80' is good)
// Place the next points anywhere else on the circle that are not "too close", this skews and flips the image 

func calculateTransform(img *IMImage) Triangle {
	var t Triangle
	
	// iC
	center := Vertex{ img.Width / 2, img.Height / 2 }
	// R(iC)
	radius := math.Sqrt( float64(center.X * center.X + center.Y * center.Y) )
	// Determine the scale down factor (we never want the image larger than the normalised size)
	iscale := 100 - rng.Intn( max_scale )
	scale := float64(iscale) / 100
	new_radius := radius * scale 

	FIU.Trace.Printf("Radius: %f -> %f", radius, new_radius)
	
	// Angle from center to top left of image 0,0 (in radians)
	angle := math.Asin( float64(-center.X) / radius )
	// Adjust by the random number of degrees supplied from rand(max_rotation)
	irotation := rng.Intn( (2 * max_rotation + 1) ) - max_rotation
	rotation := float64( irotation )
	// Convert to radians
	rotation = (rotation * math.Pi) / 180
	new_angle := angle + rotation
	FIU.Trace.Printf("Angle to 0,0: %f -> %f", angle, new_angle)
	
	// Calculate the first point
	t.P[0].X = int(math.Sin(new_angle) * new_radius)
	t.P[0].Y = int(math.Cos(new_angle) * new_radius)
	// ----------------------------------------------------------------------

	// Angle from center to top right of image width,0 (in radians)
	angle := math.Asin( float64(center.X) / radius )


	

	iskew := rng.Intn( max_skew )

	// The destination filename of a preprocessed image is made from the original image name with the distortion parameters added
	dest_filename := filepath.Base(img.Files.Orig)
	dest_filename = strings.TrimSuffix(dest_filename, filepath.Ext(dest_filename))
	if irotation<0 {
		irotation +=360
	}
	img.Files.Dest = filepath.Join(img.Files.Dest, fmt.Sprintf("%s_%d_%d_%d.jpg", dest_filename, iscale, irotation, iskew)) 

	return t
}

func pointTransform(center, src Vertex, radius, random) Vertex {
	return p
}

/*func distortImage(img IMImage) bool {

}

func imageSize(img *IMImage) (int, int) {

}
*/