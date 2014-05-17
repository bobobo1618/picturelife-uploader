package main

import (
    "path"
    "os"
    "log"
    "io/ioutil"
    "fmt"
    "sync"
)

func processHashChan(cache_dir string, hashChan chan string, sigCheckChan chan <- uploadJob){
    var wg sync.WaitGroup
    for filePath := range hashChan {
        wg.Add(1)
        go func(filePath string){
            defer wg.Done()
            log.Printf("Processing %s\n", filePath)
            // Look at the cache
            cachePath := path.Join(cache_dir, filePath)

            var baseSum string;

            fileInfo, err := os.Stat(filePath)
            if err != nil {
                log.Printf("Failed to get information about %s - %s\n", filePath, err)
                return
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
                        return
                    }
                } else {
                    log.Printf("Had an issue with the cache file: %s - %s\n", cachePath, err)
                    return
                }
            } else {
                defer cacheFile.Close()
                // Read the cache.
                readBytes, err := ioutil.ReadAll(cacheFile)

                if err != nil {
                    log.Printf("Failed to read contents of cache file: %s - %s\n", cachePath, err)
                    return
                }

                readSum := string(readBytes)
                var readModTime int64
                fmt.Sscanf(readSum, "%d,%s", &readModTime, &baseSum)

                // If the file has changed, recache it.
                if readModTime != lastModTime {
                    baseSum, err = hashAndCache(cachePath, filePath, lastModTime)
                    if err != nil {
                        log.Printf("Error hashing and caching: %s - %s", filePath, err)
                        return
                    }
                }
            }

            sigCheckChan <- uploadJob{
                filePath: filePath,
                fileHash: baseSum,
            }
        }(filePath)
    }
    wg.Wait()
    close(sigCheckChan)
}