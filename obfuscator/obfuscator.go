// Obfuscator takes sets of 3 keys and loads a number of their matching preprocessed 
// if then creates new images that are a combination of the 3 keys and pastes them randomly to an output image
// The output image is named using the 3 keys in alphabetical order


package main

import (
	FIU "FindIt/FIUtils"
    "os"
    // "os/exec"
    "path/filepath"
    // "image"
    // "image/jpeg"
    // "math"
    "math/rand"
    "time"
    "sort"
    "strings"
    "fmt"
    // "strconv"
)

var key_path string
var paths FIU.OrigDest
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

var num_keys int

func initVars() {
    //@TODO: Load this from a config file and check orig and dest are different
    key_path = filepath.Join(FIU.FindIt_path, "/images_processed/normalised/keys")

    paths = FIU.OrigDest {
       filepath.Join(FIU.FindIt_path, "/images_processed/preprocessed/keys"), 
       filepath.Join(FIU.FindIt_path, "/images_processed/obfuscator"), 
    }

    num_keys = 3
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

    fmt.Println(processed_keys)
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