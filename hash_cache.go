package main

import (
    "crypto/sha256"
    "io"
    "path"
    "encoding/hex"
    "os"
    "fmt"
    "io/ioutil"
)

// Struct to hold a hash and a path.
type UploadJob struct {
    filePath string
    fileHash string
    lastModTime int64
    uploaded bool
}

func (job *UploadJob) AddToCache(cache_dir string) (err error) {
    cachePath := path.Join(cache_dir, job.filePath)
    dir, _ := path.Split(cachePath)

    // Make directories for the cache file.
    err = os.MkdirAll(dir, 0700)

    if err != nil {
        return
    }

    cacheFile, err := os.Create(cachePath)

    if err != nil {
        return
    } else {
        defer cacheFile.Close()

        _, err = fmt.Fprintf(cacheFile, "%t,%d,%s", job.uploaded, job.lastModTime, job.fileHash)

        if err != nil {
            return
        }

        return
    }
}

func (job *UploadJob) GetFromCache(cache_dir string) (err error) {
    cachePath := path.Join(cache_dir, job.filePath)

    fmt.Println(cachePath)

    cacheFile, err := os.Open(cachePath)

    if err != nil {
        return
    }

    defer cacheFile.Close()

    readBytes, err := ioutil.ReadAll(cacheFile)

    if err != nil {
        return
    }

    readSum := string(readBytes)
    fmt.Sscanf(readSum, "%t,%d,%s", &job.uploaded, &job.lastModTime, &job.fileHash)

    return nil
}

// Hashes the given file, caches the result and time in the given cache path directory.
func (job *UploadJob) HashAndCache(cache_dir string) (err error) {
    var file *os.File
    file, err = os.Open(job.filePath)

    if err != nil {
        return
    }

    defer file.Close()

    shaHash := sha256.New()
    io.Copy(shaHash, file)
    sum := shaHash.Sum(nil)

    job.fileHash = hex.EncodeToString(sum)
    
    err = job.AddToCache(cache_dir)

    if err != nil {
        return
    }

    return
}