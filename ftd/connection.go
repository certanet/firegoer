package ftd

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

func CreateFdmConnection(host, password string, verify bool) *Fdm {
	return &Fdm{
		Host:     host,
		Password: password,
		Verify:   verify,
	}
}

func (fdm *Fdm) GetApiNoRead(url string) *http.Response {
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
	token := fdm.getaccesstoken()
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)
	// req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	return res
}

type fmdtoken struct {
	Access_Token string `access_token`
}

func (fdm *Fdm) getaccesstoken() string {
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

func (fdm *Fdm) apiCall(url string, method string, reqData interface{}) []byte {
	base_url := "https://" + fdm.Host + "/api/fdm/latest/"
	full_url := base_url + url

	if !fdm.Verify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	// Check if there's a reuest body to send and covnert to bytes
	bod := func(reqData interface{}) (b []byte) {
		if reqData != nil {
			return reqData.([]byte)
		}
		return b
	}(reqData)

	req, err := http.NewRequest(method, full_url, bytes.NewBuffer(bod))

	// TODO check and store access_token for reuse
	token := fdm.getaccesstoken()
	bearer := "Bearer " + token
	req.Header.Set("Authorization", bearer)

	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		if res.StatusCode == 422 {
			fmt.Println("Resource does not exist!")
		} else {
			fmt.Println("Non-OK HTTP status:", res.StatusCode)
			panic(res.StatusCode)
		}
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// fmt.Printf("%s", body)

	return body
}

func (fdm *Fdm) GetApi(url string) []byte {
	// Returns the body of the GET response
	return fdm.apiCall(url, "GET", nil)
}

func (fdm *Fdm) PostApi(url string, reqData interface{}) []byte {
	// Returns the body of the GET response
	return fdm.apiCall(url, "POST", reqData)
}

func (fdm *Fdm) DeleteApi(url string, reqData interface{}) []byte {
	return fdm.apiCall(url, "DELETE", reqData)
}
