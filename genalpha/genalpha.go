// genalpha creates one of a set of Alpha patterns
// and places them in the appropriate folder within processed_images

package main

import (
	FIU "FindIt/FIUtils"
    "image"
    "image/color"
    "os"
    "path/filepath"
    "math"
//    "strings"
    "strconv"
)

var alpha_path string
var width, height int

func initVars() {
	//@TODO: Load this from a config file and check orig and dest are different
	alpha_path = filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/alpha")

    width = 400
    height = 300
}

func main() {
	//FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
    initVars()

	//for _, ab := range []alphaBlend{ linear{"up"}, linear{"down"}, linear{"left"}, linear{"right"}} {
	for _, ab := range []alphaBlend{ circular{image.Pt(200,200)}, circular{image.Pt(0,0)}, circular{image.Pt(100,height)}, circular{image.Pt(width/2,height/2)} } {
		output := ab.getAlpha()
		outfilename := filepath.Join(alpha_path, ab.String() + ".png")

		// Save the image to the filename
	    out_image := output.SubImage(image.Rectangle(output.Bounds()))
	    FIU.SaveImage( &out_image, outfilename )
	}
}

type alphaBlend interface {
	getAlpha() *image.Alpha
	String() string
}

// Create the alpha bitmap and the closure that contains the nested for loop needed to run over it
func initAlpha() (*image.Alpha, func(max_inner, max_outer int, fn func(inner, outer int))) {
	alpha := image.NewAlpha(image.Rect(0, 0, width, height))

	loop := func( max_inner, max_outer int, fn func( inner, outer int) ) {
		for outer:=0; outer<max_outer; outer++ {
			for inner:=0; inner<max_inner; inner++ {
				fn(inner, outer)
			}
		}
	}

	return alpha, loop
}

// The linear alphaBlend is a simple linnear gradient fill in 1 of 4 directions up, down, left, right
type linear struct {
	direction string
}

func (l linear) getAlpha() *image.Alpha {
	alpha, loop := initAlpha()

	// use the loop closure to execute the set on each of the bytes in the alpha
	switch l.direction {
		case "up":
			loop( width, height, func( inner, outer int ) {
				col := 255 - (255 * outer) / height
				alpha.Set(inner, outer, color.RGBA{0, 0, 0, uint8(col)})
			})
			break
		case "down":
			loop( width, height, func( inner, outer int ) {
				col := (255 * outer) / height
				alpha.Set(inner, outer, color.RGBA{0, 0, 0, uint8(col)})
			})
			break
		case "left":
			loop( height, width, func( inner, outer int ) {
				col := 255 - (255 * outer) / width
				alpha.Set(outer, inner, color.RGBA{0, 0, 0, uint8(col)})
			})
			break
		case "right":
			loop( height, width, func( inner, outer int ) {
				col := (255 * outer) / width
				alpha.Set(outer, inner, color.RGBA{0, 0, 0, uint8(col)})
			})
			break
		default:
			panic("Invalid linear '" + l.direction + "', must be one of up, down, left, right")
	}	

	return alpha
}

func (l linear) String() string {
	return "linear_" + l.direction
}


// The linearFree alphaBlend is a simple gradient fill starting from 0 and transitioning evenly to 255 in a direction across the bitmap
// direction is the degree (0 is up) of the direction of the gradient. In the case of 0 degrees there would be 0 at the bottom and 255 at the top
type linearFree struct {
	direction int
}

func (l linearFree) getAlpha() *image.Alpha {
	alpha, _ := initAlpha()

	return alpha
}

func (l linearFree) String() string {
	return "linearFree_" + strconv.Itoa(l.direction)
}



// The circular alphaBland produces a circular pattern where 255 is at center and a linear gradient radiates to the longest edge at 0
type circular struct {
	center image.Point
}

func (c circular) getAlpha() *image.Alpha {
	alpha, loop := initAlpha()

	// Find the longest distance from the center to the edge
	gradient := 0
	for _, v := range []int{c.center.X, c.center.Y, width - c.center.X, height - c.center.Y} {
		if v > gradient {
			gradient = v
		}
	}
	
	loop(width, height, func( inner, outer int) {
		x := c.center.X - inner
		y := c.center.Y - outer
		dist := math.Sqrt(float64(x*x + y*y))
		col := 255 - (255 * dist) / float64(gradient)
		if col < 0 {
			col = 0
		} else if col > 255 {
			col = 255
		} 
		alpha.Set(inner, outer, color.RGBA{0, 0, 0, uint8(col)})
	})

	return alpha 
}

func (c circular) String() string {
	return "circular_" + strconv.Itoa(c.center.X) + "_" + strconv.Itoa(c.center.Y)
}


// The concentic alphaBlend is similar to the circular but has undulating ripples radiating from center at a rate of frequency.
// Technically circularis a subset of concentric, but circular is a simpler algorithm
type concentric struct {
	center image.Point
	frequency int
}

func (c concentric) getAlpha() *image.Alpha {
	alpha := image.NewAlpha(image.Rect(0, 0, width, height))

	return alpha 
}

func (c concentric) String() string {
	return "concentric: center " + c.center.String() + ", frequency " + strconv.Itoa(c.frequency)
}


// The spirla alphaBlend is similar to the concentric but radiates from center using a linear spiral pattern, frwuency defines the number of spirals
type spiral struct {
	center image.Point
	frequency int
}

func (s spiral) getAlpha() *image.Alpha {
	alpha := image.NewAlpha(image.Rect(0, 0, width, height))

	return alpha 
}

func (s spiral) String() string {
	return "spiral: center " + s.center.String() + ", frequency " + strconv.Itoa(s.frequency)
}


