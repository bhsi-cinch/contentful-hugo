package mapper

import (
	"encoding/json"
)

func MapTypes(rc []byte) (typeResult TypeResult, err error) {
	err = json.Unmarshal(rc, &typeResult)

	return
}

func MapItems(rc []byte) (itemResult ItemResult, err error) {
	err = json.Unmarshal(rc, &itemResult)

	return
}
