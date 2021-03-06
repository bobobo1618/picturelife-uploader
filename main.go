package main

import (
    "flag"
    "log"
)

func main() {
    // Command line parsing.
    var base_endpoint, access_token, cache_dir string;

    flag.StringVar(&base_endpoint, "base_endpoint", "https://api.picturelife.com/", "API base endpoint location.")
    flag.StringVar(&access_token, "token", "", "API access token.")
    flag.StringVar(&cache_dir, "cache_dir", "./", "Path to where to store hash cache.")

    flag.Parse()

    if !(len(access_token)>0) {
        panic("Required arguments were not supplied.")
    }

    // Paths to walk.
    scanPaths := flag.Args()

    // Channel for files to be hashed to go into.
    hashChan := make(chan string)

    // Channel to take upload jobs.
    sigCheckChan := make(chan UploadJob)
    uploadChan := make(chan UploadJob)
    doneChan := make(chan UploadJob)

    // Walks directories in one goroutine.
    go func(){
        walk(scanPaths, hashChan)
    }()

    // Hashes files in another goroutine.
    go func(){
        processHashChan(cache_dir, hashChan, sigCheckChan)
    }()

    go func(){
        processSignatureChecks(cache_dir, base_endpoint, access_token, sigCheckChan, uploadChan, doneChan)
    }()

    go func(){
        processUploads(cache_dir, base_endpoint, access_token, uploadChan, doneChan)
    }()

    for job := range doneChan {
        log.Printf("Job %s done!\n", job.filePath)
    }
}