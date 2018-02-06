package main

import (
	FIU "FindIt/FIUtils"
	"flag"
    "github.com/go-martini/martini"
    "net/http"
    "html/template"
    "image"
    "image/png"
    "path/filepath"
    "encoding/json"
    "encoding/base64"
    "io/ioutil"
    "bytes"
    "os"
    "github.com/davecgh/go-spew/spew"
)

var revealed = false

func parseCommandLine() (filename *string, path *string, keypath *string) {
	// Get the filename of the image and metadata json file to use
	filename = flag.String("filename", "", "The base name of the FindIt png and json metadata files")
	path     = flag.String("filepath", ".", "Full path to the folder containing the FindIt files")
	keypath  = flag.String("keypath", ".", "Full path to the key files")
	flag.Parse()

	if *filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return filename, path, keypath
}

type MetaData struct {
	Key_rects []image.Rectangle
	Key_files []string
}

type Key struct {
	KeyName, KeyImage  string
}

type PageData struct {
	Title, MainImage string	
	Keys []Key
}

func main() {
	//FIU.initTrace(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
    FIU.InitTrace(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	filename, path, keypath := parseCommandLine()

	// Load the search image
	img := encodeImage(FIU.LoadImage(filepath.Join(*path, *filename + ".png")))

	pd := PageData {
		Title: "FindIt CAPTCHA",
		MainImage: img,
	}

	// Load the json metadata, then load and store each of the key images
	meta := loadJsonData(filepath.Join(*path, *filename + ".json"))
	for _, v := range meta.Key_files {
		key_img := encodeImage(FIU.LoadImage(filepath.Join(*keypath, v + ".png")))
		key := Key{
			KeyName: v,
			KeyImage: key_img,
		}
		pd.Keys = append(pd.Keys, key)
	}

	// Start martini
    m := martini.Classic()
    // Define the route to the css file
    m.Use(martini.Static("./static"))

    m.Get("/", func() string {
    	if revealed {
    		return "Sorry you've already tried this one. Try refreshing the page"
    	}

    	t := template.Must(template.New("FindIt.html").ParseFiles("FindIt.html"))
    	buf := new(bytes.Buffer)
    	err := t.Execute(buf, &pd)
    	if err != nil {
    		// TODO: Remove the error before deployement
        	return "It's all gone wrong: " + err.Error()
    	}
        return buf.String()
    })

    m.Post("/verify/", func(req *http.Request) string {
    	if revealed {
    		return "Sorry you've already tried this one. Try refreshing the page"
    	}
    	revealed = true
    	spew.Dump(req.URL.Query())
    	out := "Verify keys:"
    	for _, v := range req.URL.Query() {
    		out += "\n" + v[0]
    	}
    	return out
     })

    m.Get("/reveal/", func(params martini.Params) string {
		revealed = true
		b, err := json.Marshal( meta )
		if err != nil {
			// TODO: Remove the error before deployement
			return "JSON encode error: " + err.Error()
		}
    	return string(b[:]);
    })

    m.Run()
}

func loadJsonData(filename string) *MetaData {
	md := MetaData {}
    meta, err := ioutil.ReadFile(filename)

    if err != nil {
        FIU.Error.Println("Unable to load '" + filename + "': " + err.Error())
        os.Exit(1)  
    }

    json.Unmarshal(meta, &md)
    return &md
}

func encodeImage(img image.Image) string {
    buffer := new(bytes.Buffer)
    if err := png.Encode(buffer, img); err != nil {
        FIU.Error.Println("unable to encode image")
        os.Exit(1)
    }

    return base64.StdEncoding.EncodeToString(buffer.Bytes())
}