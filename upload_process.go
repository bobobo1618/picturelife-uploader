package main

import (
    "net/http"
    "net/url"
    "log"
    "mime/multipart"
    "fmt"
    "os"
    "io"
    "io/ioutil"
    "bytes"
)

func processUploads(base_endpoint, access_token string, uploadChan chan uploadJob, doneChan chan <- uploadJob){
    for job := range uploadChan {
        api_url, err := url.Parse(base_endpoint)

        if err != nil {
            panic(err)
        }

        api_url.Path = "/medias/create"

        file, err := os.Open(job.filePath)
        if err != nil {
            log.Printf("Couldn't open %s for uploading: %s.\n", job.filePath, err)
            continue
        }

        var formBuffer bytes.Buffer

        formWriter := multipart.NewWriter(&formBuffer)
        formWriter.WriteField("access_token", access_token)

        fileWriter, err := formWriter.CreateFormFile("file", job.filePath)

        if err != nil {
            log.Printf("Couldn't add a multipart form field for %s: %s\n", job.filePath, err)
            continue
        }

        if _, err = io.Copy(fileWriter, file); err != nil {
            log.Printf("Error copying data between file and form: %s - %s\n", job.filePath, err)
        }

        formWriter.Close()

        req, err := http.NewRequest("POST", api_url.String(), &formBuffer)
        if err != nil {
            log.Printf("Error forming upload request for %s - %s\n", job.filePath, err)
        }

        req.Header.Set("Content-Type", formWriter.FormDataContentType())

        client := http.Client{}

        response, err := client.Do(req)
        if err != nil {
            log.Printf("Error uploading %s - %s\n", job.filePath, err)
        }

        body, _ := ioutil.ReadAll(response.Body)

        fmt.Println(response)
        fmt.Println(string(body))

        doneChan <- job
    }
    close(doneChan)
}