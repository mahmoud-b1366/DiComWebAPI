module app

go 1.20

require (
	example.com/controller v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)

require (
	example.com/model v0.0.0-00010101000000-000000000000 // indirect
	example.com/service v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/suyashkumar/dicom v1.0.6 // indirect
	golang.org/x/text v0.3.8 // indirect
)

replace example.com/controller => ./controller

replace example.com/model => ./model

replace example.com/service => ./service
