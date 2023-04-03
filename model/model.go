package model

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Model object used to store/retrieve and delete DiCom files.
type DiComModel struct {
	Path    string
	baseDir string
}

// Constructor function return an instance of NewDiComModel
func NewDiComModel(fileStoragePath string) DiComModel {
	m := DiComModel{Path: fileStoragePath}
	m.Initialize()
	return m
}

// Sets default values of the object
// Needs to be called if new instance is not created using NewDiComModel() method
func (m *DiComModel) Initialize() {
	location, err := filepath.Abs(m.Path)
	if err != nil {
		log.Fatal(err)
	}
	m.baseDir = location
}

// Stores DiCom data for the passed ID
func (m *DiComModel) Create(data []byte, id string) bool {
	log.Printf("Creating record ID: %v", id)
	location := filepath.Join(m.baseDir, id)
	err := os.WriteFile(location, data, 0644)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Deletes the DiCom data for the passed ID
func (m *DiComModel) Delete(id string) bool {
	log.Printf("Deleting record ID: %v", id)
	location := filepath.Join(m.baseDir, id)
	err := os.Remove(location)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Retrieves DiCom data associated with the passed ID
func (m *DiComModel) Retrieve(id string) ([]byte, error) {
	log.Printf("Retrieving record ID: %v", id)
	location := filepath.Join(m.baseDir, id)
	data, err := os.ReadFile(location)
	if err != nil {
		log.Println(err)
	}

	return data, err
}

// Returns location of the DiCom file stored which corresponds to the passed ID
func (m *DiComModel) RetrieveAsFile(id string) (string, error) {
	log.Printf("Retrieving record ID as file: %v", id)
	location := filepath.Join(m.baseDir, id)
	_, err := os.Stat(location)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("Error retrieving file for ID:", id)
	}
	return location, err
}

// Returns key list of available DiCom files
func (m *DiComModel) List() []string {
	log.Printf("Retrieving List of records")

	files, err := ioutil.ReadDir(m.baseDir)
	list := make([]string, 0)
	if err != nil {
		log.Print(err)
	} else {
		for _, file := range files {
			list = append(list, file.Name())
		}
	}

	return list
}
