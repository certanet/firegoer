package configexport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/certanet/firegoer/connection"
)

type ConfigExportReq struct {
	File_Name    string `json:"diskFileName"`
	Dont_Encrypt bool   `json:"doNotEncrypt"`
	Type         string `json:"type"`
}

type ConfigExportResp struct {
	File_Name string `json:"diskFileName"`
	Job_Id    string `json:"jobHistoryUuid"`
}

type ConfigExportStatus struct {
	File_Name string `json:"diskFileName"`
	Status    string `json:"status"`
}

func ExportConfig(fdm connection.Fdm, config_name string) ConfigExportResp {
	// Submits a ConfigExport job
	var export ConfigExportResp

	body := ConfigExportReq{
		File_Name:    config_name,
		Dont_Encrypt: true,
		Type:         "scheduleconfigexport"}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	resp := connection.PostApi(fdm, "action/configexport", payloadBuf.Bytes())

	jsonErr := json.Unmarshal(resp, &export)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return export
}

func GetConfigExportStatus(fdm connection.Fdm, export_job_id string) ConfigExportStatus {
	// Gets the status of the given ConfigExport job
	var status ConfigExportStatus

	resp := connection.GetApi(fdm, "jobs/configexportstatus/"+export_job_id)

	jsonErr := json.Unmarshal(resp, &status)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return status
}

func DownloadConfigFile(fdm connection.Fdm, remote_filename string, local_filename string) {
	// Downloads the given remote exported config file

	// Create the zip file locally
	out, err := os.Create(local_filename)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	defer out.Close()

	dl_url := "action/downloadconfigfile/" + remote_filename

	// Get the http.Response without reading it's contents and close when done
	res := connection.GetApiNoRead(fdm, dl_url)
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Read the reponse body in chunks directly to the file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
}
