module example.com/controller

go 1.20

replace example.com/model => ../model

replace example.com/service => ../service

require (
	example.com/model v0.0.0-00010101000000-000000000000
	example.com/service v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
)

require (
	github.com/suyashkumar/dicom v1.0.6 // indirect
	golang.org/x/text v0.3.8 // indirect
)
