package main

import (
	"fmt"
	"github.com/jwistrom/pellet-burner-monitor/internal/hardware"
	"github.com/jwistrom/pellet-burner-monitor/internal/persistence"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var htmlTemplate *template.Template

var burner hardware.Burner
var persistenceService persistence.Service

func main() {

	burner = &hardware.DummyBurner{}
	persistenceService = &persistence.DummyService{}

	startTemperatureCollection(time.Duration(5) * time.Minute)

	fmt.Println(os.Getwd())
	loadTemplate()

	http.HandleFunc("/", handleRoot)

	log.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadTemplate() {
	var err error
	htmlTemplate, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Template error %s", err)
	}
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	alarmStartTime := burner.ActiveAlarmStartTime().Format(time.RFC1123)

	tempHistory, err := persistenceService.GetTemperatureHistory()
	if err != nil {
		log.Printf("Failed to get temp history: %s", err)
		tempHistory = []persistence.TemperatureRecording{}
	}

	context := map[string]interface{}{
		"currentTemperature":   burner.GetCurrentTemperature(),
		"activeAlarm":          burner.AlarmIsActive(),
		"activeAlarmStartTime": alarmStartTime,
		"temperatureHistory":   tempHistory,
	}

	err = htmlTemplate.Execute(w, context)
	if err != nil {
		log.Printf("Template render error: %s", err)
		http.Error(w, err.Error(), 500)
	}
}

func startTemperatureCollection(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Collecting temperature")
				err := persistenceService.StoreTemperature(burner.GetCurrentTemperature())
				if err != nil {
					log.Fatalf("Failed to store temperature. Reason: %s", err)
				}
			}
		}
	}()

}
