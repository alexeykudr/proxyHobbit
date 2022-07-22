package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// const requestURL = "http://46.227.245.119:8081/"
const requestURL = "http://127.0.0.1:8081/"

func HealthCheck() {
	start_url := fmt.Sprintf("%s%s", requestURL, "")
	req, err := http.NewRequest(http.MethodGet, start_url, nil)
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)

	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)

	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	log.Printf("client: response body: %s\n", string(resBody))
}
func Reboot(portId string) (string, error) {
	start_url := fmt.Sprintf("%s%s", requestURL, "reboot")
	req, err := http.NewRequest(http.MethodGet, start_url, nil)
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)

	}
	q := req.URL.Query()
	q.Add("id", portId)
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("client: error making http request: %s\n", err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	
	return string(resBody), err
}

func SetInterval(id, time string) {
	path := "updatePortInterval"
	start_url := fmt.Sprintf("%s%s", requestURL, path)
	req, err := http.NewRequest(http.MethodGet, start_url, nil)

	if err != nil {
		log.Printf("%s", err.Error())
	}

	q := req.URL.Query()
	q.Add("id", id)
	q.Add("interval", time)
	// q.Add("another_thing", "foo & bar")

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	fmt.Printf("client: response body: %s\n", string(resBody))
}

func DeleteInterval(id string) {
	path := "removePortInterval"
	start_url := fmt.Sprintf("%s%s", requestURL, path)
	req, err := http.NewRequest(http.MethodGet, start_url, nil)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	q := req.URL.Query()
	q.Add("id", id)
	// q.Add("another_thing", "foo & bar")

	req.URL.RawQuery = q.Encode()
	// fmt.Println(req.URL.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	fmt.Printf("client: response body: %s\n", string(resBody))
}

func SetConfig(id, login, password string)(error) {
	path := "setConfig"
	start_url := fmt.Sprintf("%s%s", requestURL, path)

	data := url.Values{
		"id":       {id},
		"username": {login},
		"password": {password},
	}

	resp, err := http.PostForm(start_url, data)

	if err != nil {
		log.Printf("%s", err.Error())
	}

	res := make(map[string]string)

	err_response := json.NewDecoder(resp.Body).Decode(&res)
	return err_response
}
