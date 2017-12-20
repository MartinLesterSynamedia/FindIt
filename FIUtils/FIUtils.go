// Simple trace copied from somewhere on the web

package FIUtils

import (
	"log"
	"io"
    "os"
    "io/ioutil"
    "bytes"
    "image"
    _ "image/jpeg"
    _ "image/png"
    _ "image/gif"
)

var FindIt_path string = os.Getenv("GOPATH") + "/src/FindIt"

type OrigDest struct {
    Orig, Dest string
}

var Paths map[string]OrigDest

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func InitTrace(
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

// Use os.readdir to list the files in the folders
func ListFilenames(path string) []os.FileInfo {
    Trace.Println("listFilenames(" + path + ")")

    files, err := ioutil.ReadDir(path)
    if err != nil {
        Error.Println(err)
        os.Exit(1)
    }

    return files
}

// Load an image irrespective of format. Images that fail to load produce warnings
func LoadImage(filename string) image.Image {
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

