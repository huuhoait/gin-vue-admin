package utils

import (
	"encoding/json"
	"strings"
)

func GetJSONKeys(jsonStr string) (keys []string, err error) {
	// Usejson.Decoder, ByConvenientAtParsePassFlowInRecordKeyofOrder
	dec := json.NewDecoder(strings.NewReader(jsonStr))
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	// ensureDataYesOnePieceObject
	if t != json.Delim('{') {
		return nil, err
	}
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return nil, err
		}
		keys = append(keys, t.(string))

		// ParseValue
		var value interface{}
		err = dec.Decode(&value)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
