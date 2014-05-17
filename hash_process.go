package main

import (
    "path"
    "os"
    "log"
    "io/ioutil"
    "fmt"
)

func processHashChan(cache_dir string, hashChan chan string, sigCheckChan chan <- uploadJob){
    for filePath := range hashChan {

        // Look at the cache
        cachePath := path.Join(cache_dir, filePath)

        var baseSum string;

        fileInfo, err := os.Stat(filePath)
        if err != nil {
            log.Printf("Failed to get information about %s - %s\n", filePath, err)
        }

        lastModTime := fileInfo.ModTime().Unix()

        // Check the cache
        cacheFile, err := os.Open(cachePath)

        if err != nil {
            // If the cache file doesn't exist, create it.
            if os.IsNotExist(err) {
                baseSum, err = hashAndCache(cachePath, filePath, lastModTime)
                if err != nil {
                    log.Printf("Error hashing and caching: %s - %s", filePath, err)
                }
            } else {
                log.Printf("Had an issue with the cache file: %s - %s\n", cachePath, err)
                continue
            }
        } else {
            defer cacheFile.Close()
            // Read the cache.
            readBytes, err := ioutil.ReadAll(cacheFile)

            if err != nil {
                log.Printf("Failed to read contents of cache file: %s - %s\n", cachePath, err)
                continue
            }

            readSum := string(readBytes)
            var readModTime int64
            fmt.Sscanf(readSum, "%d,%s", &readModTime, &baseSum)

            // If the file has changed, recache it.
            if readModTime != lastModTime {
                baseSum, err = hashAndCache(cachePath, filePath, lastModTime)
                if err != nil {
                    log.Printf("Error hashing and caching: %s - %s", filePath, err)
                }
            }
        }

        sigCheckChan <- uploadJob{
            filePath: filePath,
            fileHash: baseSum,
        }
    }
    close(sigCheckChan)
}