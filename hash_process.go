package main

import (
    "path"
    "os"
    "log"
    "sync"
)

func processHashChan(cache_dir string, hashChan chan string, sigCheckChan chan <- UploadJob){
    var wg sync.WaitGroup
    for filePath := range hashChan {
        wg.Add(1)
        go func(filePath string){
            defer wg.Done()
            log.Printf("Processing %s\n", filePath)
            // Look at the cache
            cachePath := path.Join(cache_dir, filePath)

            fileInfo, err := os.Stat(filePath)
            if err != nil {
                log.Printf("Failed to get information about %s - %s\n", filePath, err)
                return
            }

            lastModTime := fileInfo.ModTime().Unix()

            currentJob := UploadJob{
                lastModTime: lastModTime,
                filePath: filePath,
                uploaded: false,
            }

            cacheJob := UploadJob{
                filePath: filePath,
            }

            err = cacheJob.GetFromCache(cache_dir)

            if err != nil {
                // If the cache file doesn't exist, create it.
                if os.IsNotExist(err) {
                    err = currentJob.HashAndCache(cache_dir)
                    if err != nil {
                        log.Printf("Error hashing and caching: %s - %s", filePath, err)
                        return
                    }
                } else {
                    log.Printf("Had an issue with the cache file: %s - %s\n", cachePath, err)
                    return
                }
            } else {
                // If the file has changed, recache it.
                if cacheJob.lastModTime < currentJob.lastModTime {
                    err = currentJob.HashAndCache(cache_dir)
                    if err != nil {
                        log.Printf("Error hashing and caching: %s - %s", filePath, err)
                        return
                    }
                    currentJob.uploaded = false
                } else {
                    currentJob = cacheJob
                }
            }

            if !currentJob.uploaded {
                sigCheckChan <- currentJob
            }
        }(filePath)
    }
    wg.Wait()
    close(sigCheckChan)
}