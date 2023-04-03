package service_test

import (
	"testing"

	"example.com/service"
)

func TestExistingTag(t *testing.T) {
	fileName := "../dicomFiles/a377afb7-3b0d-47e8-987f-bec71755718b"
	tagName := "PatientName"
	svc := service.NewDiComService(fileName)
	value := svc.GetTagValue(tagName)
	if value != "NAYYAR^HARSH" {
		t.Errorf("Retrieved value for tag %v of file %v is invalid ", tagName, fileName)
	}
}

func TestInvalidTag(t *testing.T) {
	fileName := "../dicomFiles/a377afb7-3b0d-47e8-987f-bec71755718b"
	tagName := "InvalidTag"
	svc := service.NewDiComService(fileName)
	value := svc.GetTagValue(tagName)
	if value != "INVALID_TAG_NAME" {
		t.Errorf("Retrieved value for invalid tag is not 'INVALID_TAG_NAME'")
	}
}
