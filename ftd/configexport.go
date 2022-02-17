package ftd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

func (fdm *Fdm) ExportConfig(config_name string) ConfigExportResp {
	// Submits a ConfigExport job
	var export ConfigExportResp

	body := ConfigExportReq{
		File_Name:    config_name,
		Dont_Encrypt: true,
		Type:         "scheduleconfigexport",
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	resp := fdm.PostApi("action/configexport", payloadBuf.Bytes())

	jsonErr := json.Unmarshal(resp, &export)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return export
}

func (fdm *Fdm) GetConfigExportStatus(export_job_id string) ConfigExportStatus {
	// Gets the status of the given ConfigExport job
	var status ConfigExportStatus

	resp := fdm.GetApi("jobs/configexportstatus/" + export_job_id)

	jsonErr := json.Unmarshal(resp, &status)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return status
}

func (fdm *Fdm) DownloadConfigFile(remote_filename string, local_filename string) {
	// Downloads the given remote exported config file

	// Create the zip file locally
	out, err := os.Create(local_filename)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	defer out.Close()

	dl_url := "action/downloadconfigfile/" + remote_filename

	// Get the http.Response without reading it's contents and close when done
	res := fdm.GetApiNoRead(dl_url)
	if res.Body != nil {
		defer res.Body.Close()
	}

	// Read the reponse body in chunks directly to the file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
}

func (fdm *Fdm) DeleteConfigExport(remote_filename string) {
	// Deletes the given config file
	_ = fdm.DeleteApi("action/configfiles/"+remote_filename, nil)
}
