package systeminfo

import (
	"encoding/json"
	"log"

	"github.com/certanet/firegoer/connection"
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

func GetHostname(fdm connection.Fdm) string {
	var items HostnameItems

	resp := connection.GetApi(fdm, "devicesettings/default/devicehostnames")

	jsonErr := json.Unmarshal(resp, &items)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return items.Items[0].Hostname
}

func GetSystemInfo(fdm connection.Fdm) SystemInfo {
	var info SystemInfo

	resp := connection.GetApi(fdm, "/operational/systeminfo/default")

	jsonErr := json.Unmarshal(resp, &info)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return info
}
