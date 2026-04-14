package utils

import (
	"fmt"
	"testing"
)

func TestGetJSONKeys(t *testing.T) {
	var jsonStr = `
	{
		"Name": "test",
		"TableName": "test",
		"TemplateID": "test",
		"TemplateInfo": "test",
		"Limit": 0
}`
	keys, err := GetJSONKeys(jsonStr)
	if err != nil {
		t.Errorf("GetJSONKeys failed: %v", err)
		return
	}
	if len(keys) != 5 {
		t.Errorf("GetJSONKeys failed: expected 5 keys, got %d", len(keys))
		return
	}
	if keys[0] != "Name" {
		t.Errorf("GetJSONKeys failed: keys[0]=%q", keys[0])

		return
	}
	if keys[1] != "TableName" {
		t.Errorf("GetJSONKeys failed: keys[1]=%q", keys[1])

		return
	}
	if keys[2] != "TemplateID" {
		t.Errorf("GetJSONKeys failed: keys[2]=%q", keys[2])

		return
	}
	if keys[3] != "TemplateInfo" {
		t.Errorf("GetJSONKeys failed: keys[3]=%q", keys[3])

		return
	}
	if keys[4] != "Limit" {
		t.Errorf("GetJSONKeys failed: keys[4]=%q", keys[4])

		return
	}

	fmt.Println(keys)
}
