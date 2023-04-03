package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"example.com/model"
	"example.com/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ApiMessageResponse struct {
	Success  bool   `json:"success"`
	Feedback string `json:"feedback"`
}
type ApiUploadResponse struct {
	Success  bool   `json:"success"`
	Feedback string `json:"feedback"`
	ID       string `json:"id"`
}

// Controller to handle CRUD operations for DiCom files
type DiComController struct {
	model model.DiComModel
}

const fileStoragePath = "./dicomFiles"

// Constructor function return an instance of DiComController
func NewDiComController() DiComController {
	c := DiComController{}
	c.Initialize()
	return c
}

// Sets default values of the object
// Needs to be called if new instance is not created using NewDiComController() method
func (c *DiComController) Initialize() {
	log.Println("Initializing controller...")
	c.model = model.NewDiComModel(fileStoragePath)
}

// Rest Api's home page
func (c *DiComController) Welcome(w http.ResponseWriter, r *http.Request) {
	// TODO: replace with useful help document
	fmt.Fprintf(w, "Welcome to Pocket Health API")
}

// Retrieves tags information for the given ID
func (c *DiComController) Retrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	service, err := c.getService(id)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Sprintf("Error retrieving record:%v", id))
	} else {
		result := make(map[string]any)
		tagsString := r.URL.Query().Get("tags")
		if tagsString != "" {
			tags := strings.Split(tagsString, ",")
			for _, t := range tags {
				value := service.GetTagValue(t)
				result[t] = value
			}

		}

		dataResponse(w, 200, result)
	}
}

// Returns list of IDs of the uploaded decom files
func (c *DiComController) List(w http.ResponseWriter, r *http.Request) {
	list := c.model.List()
	dataResponse(w, 200, list)
}

// Deletes specific record from uploaded decom files
func (c *DiComController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	success := c.model.Delete(id)
	if success {
		messageResponse(w, fmt.Sprintf("Successfully deleted record:%v", id))
	} else {
		errorResponse(w, fmt.Sprintf("Error deleting record:%v", id))
	}
}

// Saves uploaded DiCom file to defined storage folder
func (c *DiComController) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(25)
	file, handler, err := r.FormFile("dicomFile")
	if err != nil {
		errorResponse(w, "Error receiving uploaded file, try again")
		log.Println(err)
		return
	}
	defer file.Close()

	id := uuid.New().String()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	success := c.model.Create(bytes, id)
	if success {
		uploadResponse(w, handler.Filename, handler.Size, id)
	} else {
		errorResponse(w, "Error Storing uploaded file, try again")
	}
}

// Parses target DiCom file and returns the result as image/png
func (c *DiComController) RetrieveImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	service, err := c.getService(id)
	if err != nil {
		log.Println(err)
		errorResponse(w, fmt.Sprintf("Error retrieving image for:%v", id))
	} else {
		data := service.GetImage()
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		if _, err := w.Write(data); err != nil {
			log.Println("unable to send image.")
		}
	}
}

// Restores uploaded DiCom file and sends it to caller
func (c *DiComController) RetrieveDiCom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	data, err := c.model.Retrieve(id)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Error retrieving dicom file for:%v", id))
	} else {
		w.Header().Set("Content-Type", "application/dicom")
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		if _, err := w.Write(data); err != nil {
			log.Println("unable to send image.")
		}
	}
}

// Helper method to create DiComService for given ID
func (c *DiComController) getService(id string) (service.DiComService, error) {
	var svc service.DiComService
	filePath, err := c.model.RetrieveAsFile(id)
	if err == nil {
		svc = service.NewDiComService(filePath)
	}
	return svc, err
}

// Helper method used to send error messages to caller
func errorResponse(w http.ResponseWriter, message string) {
	result := ApiMessageResponse{
		Success:  false,
		Feedback: message,
	}
	dataResponse(w, 400, result)
}

// Helper method used to send feedbacks to caller when operation is successful
func messageResponse(w http.ResponseWriter, message string) {
	result := ApiMessageResponse{
		Success:  true,
		Feedback: message,
	}
	dataResponse(w, 200, result)
}

// Helper method used to prepare response when upload operation is successful
func uploadResponse(w http.ResponseWriter, fileName string, size int64, id string) {
	result := ApiUploadResponse{
		Success:  true,
		Feedback: fmt.Sprintf("%v with Size of %v uploaded successfully", fileName, size),
		ID:       id,
	}
	dataResponse(w, 201, result)
}

// Helper method used to send JSON response to the caller
func dataResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
