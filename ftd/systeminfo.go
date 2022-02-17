package ftd

import (
	"encoding/json"
	"log"
)

type HostnameItems struct {
	Items []Hostname
}

type Hostname struct {
	Hostname string `json:"hostname"`
}

type SystemInfo struct {
	Ipv4          string `json:"ipv4"`
	Ipv6          string `json:"ipv6"`
	Version       string `json:"softwareVersion"`
	Model         string `json:"platformModel"`
	Serial_Number string `json:"serialNumber"`
}

func (fdm *Fdm) GetHostname() string {
	var items HostnameItems

	resp := fdm.GetApi("devicesettings/default/devicehostnames")

	jsonErr := json.Unmarshal(resp, &items)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return items.Items[0].Hostname
}

func (fdm *Fdm) GetSystemInfo() SystemInfo {
	var info SystemInfo

	resp := fdm.GetApi("/operational/systeminfo/default")

	jsonErr := json.Unmarshal(resp, &info)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return info
}
