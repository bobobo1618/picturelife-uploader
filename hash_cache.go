package main

import (
    "crypto/sha256"
    "io"
    "path"
    "encoding/hex"
    "os"
    "fmt"
)

// Hashes the given file, caches the result and time in the given cache path directory.
func hashAndCache(cachePath, filePath string, lastModTime int64) (hash string, err error) {
    dir, _ := path.Split(cachePath)

    // Make directories for the cache file.
    hash = ""
    err = os.MkdirAll(dir, 0700)

    if err != nil {
        return
    }

    cacheFile, err := os.Create(cachePath)

    if err != nil {
        return
    } else {
        defer cacheFile.Close()

        var file *os.File
        file, err = os.Open(filePath)

        if err != nil {
            return
        }

        defer file.Close()

        shaHash := sha256.New()
        io.Copy(shaHash, file)
        sum := shaHash.Sum(nil)

        hash = hex.EncodeToString(sum)
        writtenSum := fmt.Sprintf("%d,%s", lastModTime, hash)

        cacheFile.Write([]byte(writtenSum))

        return
    }
}