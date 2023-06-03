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
	login string `json: "login" binding:"required"`
	password string `json: "password" binding:"required"`
}

func LoginRequest() error {
    url := "http://localhost:8080/Login"
    data := Admin{login: "Admin", password: "Admin"}
    payload, err := json.Marshal(data)
    if err != nil {
        return errors.New("Json Encoding")
    }
    fmt.Println(payload)
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
func SignupRequest() error {
    url := "http://localhost:8080/Signup"
    data := Admin{login: "Admin", password: "Admin"}
    fmt.Println(data)
    payload, err := json.Marshal(data)
    if err != nil {
        return errors.New("Json Encoding")
    }
    fmt.Println(payload)
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
    // err := LoginRequest()
    err := SignupRequest();
	if (err != nil){
		fmt.Println("Error: ", err)
	}
}