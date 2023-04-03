package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"example.com/controller"
	"github.com/gorilla/mux"
)

// Rest API application
type App struct {
	Router          *mux.Router
	DiComController controller.DiComController
	Config          Configuration
	Logger          *log.Logger
}

// Properties exist in config file
type Configuration struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	LogPath string `json:"logPath"`
}

func main() {
	// Reading configuration
	configFile, _ := ioutil.ReadFile("./config.json")
	config := Configuration{}
	err := json.Unmarshal([]byte(configFile), &config)
	if err != nil {
		log.Fatal(err)
	}
	a := App{Config: config}

	// setting up logger
	logOutput, err := os.OpenFile(a.Config.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logOutput.Close()
	log.SetOutput(logOutput)

	a.Initialize()
	a.Run()
}

// Runs the HTTP server using provided routes/configuration
func (a *App) Run() {
	runningTarget := fmt.Sprintf("%v:%v", a.Config.IP, a.Config.Port)
	msg := fmt.Sprintf("Server started on %v", runningTarget)
	fmt.Println(msg)
	log.Println(msg)
	log.Fatal(http.ListenAndServe(runningTarget, a.Router))
}

// Must be called to setup default values and routes of the Rest API application
func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.DiComController = controller.NewDiComController()
	a.initializeRoutes()
}

// defining routes and their corresponding handler methods
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.DiComController.Welcome).Methods("GET")
	a.Router.HandleFunc("/data", a.DiComController.Upload).Methods("POST")
	a.Router.HandleFunc("/data", a.DiComController.List).Methods("GET")
	a.Router.HandleFunc("/data/{id}", a.DiComController.Retrieve).Methods("GET")
	a.Router.HandleFunc("/data/{id}/image", a.DiComController.RetrieveImage).Methods("GET")
	a.Router.HandleFunc("/data/{id}/dicom", a.DiComController.RetrieveDiCom).Methods("GET")
	a.Router.HandleFunc("/data/{id}", a.DiComController.Delete).Methods("DELETE")
}
