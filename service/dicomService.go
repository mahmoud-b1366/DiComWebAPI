package service

import (
	"bytes"
	"image/png"
	"io"
	"log"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// Service to provide functionalities used to extract information from DiCom files
type DiComService struct {
	dataset dicom.Dataset
}

// Constructor function return an instance of DiComService
func NewDiComService(filePath string) DiComService {
	s := DiComService{}
	s.Initialize(filePath)
	return s
}

// Sets default values of the object
// Needs to be called if new instance is not created using NewDiComService() method
func (s *DiComService) Initialize(filePath string) {
	s.dataset, _ = dicom.ParseFile(filePath, nil)
}

// Parses DiCom data and returns the image as PNG format
func (s *DiComService) GetImage() []byte {

	pixelDataElement, _ := s.dataset.FindElementByTag(tag.PixelData)
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	img, _ := pixelDataInfo.Frames[0].GetImage()
	buffer := new(bytes.Buffer)
	_ = png.Encode(io.Writer(buffer), img)
	return buffer.Bytes()
}

// Parses DiCom data and returns the value of the passed tag.
func (s *DiComService) GetTagValue(tagName string) any {
	tag, err := tag.FindByName(tagName)
	if err != nil {
		log.Print("Invalid tag name:", tagName)
		return "INVALID_TAG_NAME"
	}
	element, _ := s.dataset.FindElementByTag(tag.Tag)

	switch element.Value.ValueType() {
	case dicom.Strings:
		value := element.Value.GetValue().([]string)
		return value[0]
	case dicom.Floats:
		value := element.Value.GetValue().([]float64)
		return value[0]
	case dicom.Ints:
		value := element.Value.GetValue().([]int64)
		return value[0]
	case dicom.PixelData:
		return "Not Supported!"
	default:
		return element.Value.String()
	}
}
