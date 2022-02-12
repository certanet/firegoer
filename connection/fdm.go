package connection

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Fdm struct {
	Host     string
	Password string
	Verify   bool
}

func GetApiNoRead(fdm Fdm, url string) *http.Response {
	// Returns a Response so it can be read into a file or JSON etc.

	// For GetAPIVersions, the API version should not be sent, like this:
	// base_url := "https://" + fdm.Host + "/api/"
	base_url := "https://" + fdm.Host + "/api/fdm/latest/"
	full_url := base_url + url

	if !fdm.Verify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("GET", full_url, nil)

	// For GetAPIVersions, use basic auth instead of getting a bearer token
	// req.SetBasicAuth("admin", fdm.Password)

	// TODO check and store access_token for reuse
	token := getaccesstoken(fdm)
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	// req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	return res

}

func GetApi(fdm Fdm, url string) []byte {
	// Returns the body of the GET response
	res := GetApiNoRead(fdm, url)
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	// fmt.Printf("%s", body)
	return body
}

type fmdtoken struct {
	Access_Token string `access_token`
}

func getaccesstoken(fdm Fdm) string {
	base_url := "https://" + fdm.Host + "/api/fdm/latest/"
	full_url := base_url + "fdm/token"

	if !fdm.Verify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	// TODO clean this up
	var s1 = `{"grant_type": "password", "username": "admin", "password": "`
	var s2 = `"}`
	jsonAuth := make([]byte, 0)
	jsonAuth = append(jsonAuth, []byte(s1)...)
	jsonAuth = append(jsonAuth, []byte(fdm.Password)...)
	jsonAuth = append(jsonAuth, []byte(s2)...)

	req, err := http.NewRequest("POST", full_url, bytes.NewBuffer(jsonAuth))
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var token fmdtoken
	jsonErr := json.Unmarshal(body, &token)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// fmt.Printf("Got access token: %s", token.Access_Token)
	return token.Access_Token
}

func PostApi(fdm Fdm, url string, jsonData []byte) []byte {
	base_url := "https://" + fdm.Host + "/api/fdm/latest/"
	full_url := base_url + url

	if !fdm.Verify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest("POST", full_url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// TODO check and store access_token for reuse
	token := getaccesstoken(fdm)
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// fmt.Printf("%s", body)

	return body
}
