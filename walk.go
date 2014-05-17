package main

import (
    "path/filepath"
    "os"
    "log"
)

func walk(scanPaths []string, hashChan chan <- string) {
    for _, filePath := range scanPaths {
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
    }
    close(hashChan)
}