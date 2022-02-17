package ftd

import (
	"encoding/json"
	"log"
)

type ApiVersion struct {
	Supported_Versions []string `json:"supportedVersions"`
}

func (fdm *Fdm) GetApiVers() ApiVersion {
	var vers ApiVersion

	body := fdm.GetApi("versions")

	jsonErr := json.Unmarshal(body, &vers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return vers
}
