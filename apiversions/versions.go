package apiversions

import (
	"encoding/json"
	"log"

	"github.com/certanet/firegoer/connection"
)

type ApiVersion struct {
	Supported_Versions []string `json:"supportedVersions"`
}

func GetApiVers(fdm connection.Fdm) ApiVersion {
	var vers ApiVersion

	body := connection.GetApi(fdm, "versions")

	jsonErr := json.Unmarshal(body, &vers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return vers
}
