package main

import (
    "path/filepath"
    "os"
    "log"
    "sync"
)

func walk(scanPaths []string, hashChan chan <- string) {
    var waitGroup sync.WaitGroup

    for _, filePath := range scanPaths {
        waitGroup.Add(1)

        go func(filePath string){
            defer waitGroup.Done()
            filepath.Walk(filePath, func(filePath string, info os.FileInfo, err error) error {
                if !info.IsDir() {
                    absPath, err := filepath.Abs(filePath)
                    if err != nil {
                        log.Printf("Error getting absolute path for %s - %s\n", filePath, err)
                        return nil
                    }
                    hashChan <- absPath
                }
                return nil
            })
        }(filePath)
    }

    waitGroup.Wait()
    close(hashChan)
}