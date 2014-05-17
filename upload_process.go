package main

import (
    "net/http"
    "net/url"
    "log"
    "mime/multipart"
    "os"
    "io"
    "io/ioutil"
    "bytes"
    "sync"
)

func processUploads(cache_dir, base_endpoint, access_token string, uploadChan chan UploadJob, doneChan chan <- UploadJob){
    var wg sync.WaitGroup
    for job := range uploadChan {
        wg.Add(1)
        go func(job UploadJob){
            defer wg.Done()
            log.Printf("Uploading %s\n", job.filePath)
            api_url, err := url.Parse(base_endpoint)

            if err != nil {
                panic(err)
            }

            api_url.Path = "/medias/create"

            file, err := os.Open(job.filePath)
            if err != nil {
                log.Printf("Couldn't open %s for uploading: %s.\n", job.filePath, err)
                return
            }

            var formBuffer bytes.Buffer

            formWriter := multipart.NewWriter(&formBuffer)
            formWriter.WriteField("access_token", access_token)

            fileWriter, err := formWriter.CreateFormFile("file", job.filePath)

            if err != nil {
                log.Printf("Couldn't add a multipart form field for %s: %s\n", job.filePath, err)
                return
            }

            if _, err = io.Copy(fileWriter, file); err != nil {
                log.Printf("Error copying data between file and form: %s - %s\n", job.filePath, err)
                return
            }

            formWriter.Close()

            req, err := http.NewRequest("POST", api_url.String(), &formBuffer)
            if err != nil {
                log.Printf("Error forming upload request for %s - %s\n", job.filePath, err)
                return
            }

            req.Header.Set("Content-Type", formWriter.FormDataContentType())

            client := http.Client{}

            response, err := client.Do(req)
            if err != nil {
                log.Printf("Error uploading %s - %s\n", job.filePath, err)
                return
            }

            ioutil.ReadAll(response.Body)

            job.uploaded = true

            job.AddToCache(cache_dir)

            doneChan <- job
        }(job)
    }
    wg.Wait()
    close(doneChan)
}