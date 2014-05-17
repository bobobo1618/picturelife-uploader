package main

import (
    "io/ioutil"
    "log"
    "net/url"
    "net/http"
    "encoding/json"
    "sync"
)

// Struct to hold a hash and a path.
type uploadJob struct {
    filePath string
    fileHash string
}

type signature_response struct {
    Status int
    Signatures map[string]interface{}
    ResponseTime int `json:"response_time"`
}

func processSignatureChecks(base_endpoint, access_token string, sigCheckChan chan uploadJob, uploadChan chan <- uploadJob, doneChan chan <- uploadJob){
    var wg sync.WaitGroup

    for job := range sigCheckChan {
        wg.Add(1)

        go func(job uploadJob){
            defer wg.Done()
            api_url, err := url.Parse(base_endpoint)

            if err != nil {
                panic(err)
            }

            api_url.Path = "/medias/check_signatures"

            queryValues := url.Values{}

            queryValues.Add("access_token", access_token)
            queryValues.Add("signatures", job.fileHash)

            queryString := queryValues.Encode()

            api_url.RawQuery = queryString

            func(req_url string){
                log.Printf("Checking %s\n", job.filePath)
                resp, err := http.Get(req_url)
                
                if err != nil || resp.StatusCode != 200{
                    log.Printf("Failed to get %s - %s\n", req_url, err)
                    return
                }

                defer resp.Body.Close()

                body, err := ioutil.ReadAll(resp.Body)

                if err != nil {
                    log.Printf("Failed to read %s - %s\n", req_url, err)
                    return
                }

                response := signature_response{}

                err = json.Unmarshal(body, &response)

                if err != nil {
                    log.Printf("Failed to parse %s - %s\n", req_url, err)
                    return
                }

                if response.Status != 20000 {
                    log.Printf("Bad response %s - %d\n", req_url, response.Status)
                    return
                }

                if response.Signatures[job.fileHash] == nil {
                    log.Printf("Queuing %s for upload.\n", job.filePath)
                    uploadChan <- job
                } else {
                    doneChan <- job
                }

            }(api_url.String())
        }(job)
    }
    wg.Wait()
    close(uploadChan)
}