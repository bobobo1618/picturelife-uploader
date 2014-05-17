package main

import (
    "path/filepath"
    "os"
)

func walk(scanPaths []string, hashChan chan <- string) {
    for _, filePath := range scanPaths {
        filepath.Walk(filePath, func(filePath string, info os.FileInfo, err error) error {
            if !info.IsDir() {
                hashChan <- filePath
            }
            return nil
        })
    }
    close(hashChan)
}