package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"errors"
    "io/ioutil"
)

type Admin struct {
	Login string `"json": "login" binding:"required"`
	Password string `"json": "password" binding:"required"`
}

func LoginRequest() error {
    url := "http://localhost:8080/Login"
    data := Admin{Login: "Admin", Password: "Admin"}
    payload, err := json.Marshal(data)
    if err != nil {
        return errors.New("Json Encoding")
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
    if err != nil {
        return errors.New("Setuping The Request")
	}

    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
		return errors.New("Sending The request")
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return errors.New("Reading The response body")
    }
    fmt.Println(resp.Status)
    fmt.Println(string(body))
	return nil
}
func main() {
    err := LoginRequest()
	if (err != nil){
		fmt.Println("Error: ", err)
	}
}