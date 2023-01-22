package main

import (
	"github.com/jwistrom/pellet-burner-monitor/internal/hardware"
	"github.com/jwistrom/pellet-burner-monitor/internal/notification"
	"github.com/jwistrom/pellet-burner-monitor/internal/persistence"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var htmlTemplate *template.Template

var burner hardware.Burner
var persistenceService persistence.Service
var notificationService notification.Service

func main() {

	burner = &hardware.BurnerImpl{}
	persistenceService = &persistence.ServiceImpl{}
	notificationService = setupNotificationService()

	burner.AddAlarmListener(func() {
		err := persistenceService.StoreAlarm(time.Now())
		if err != nil {
			log.Println("Failed to store alarm", err)
			return
		}
	})

	burner.AddAlarmListener(func() {
		notificationService.SendNotification("Ett alarm har triggats av brännaren!!")
	})

	startTemperatureCollection(time.Duration(5) * time.Minute)

	loadTemplate()

	http.HandleFunc("/", handleRoot)

	log.Println("Serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupNotificationService() notification.Service {
	pwd := os.Getenv("EMAIL_PWD")
	if len(pwd) == 0 {
		log.Fatalln("Password needs to be entered")
	}

	host := os.Getenv("EMAIL_HOST")
	if len(host) == 0 {
		host = "smtp.gmail.com"
	}

	port := os.Getenv("EMAIL_PORT")
	if len(port) == 0 {
		port = "587"
	}

	from := os.Getenv("EMAIL_FROM")
	if len(from) == 0 {
		from = "johan.wistroem@gmail.com"
	}

	recipients := os.Getenv("EMAIL_RECIPIENTS")
	if len(recipients) == 0 {
		recipients = "johan.wistroem@gmail.com;blomkvistlisa@hotmail.com"
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln("Cannot parse port " + port)
	}
	parsedRecipients := strings.Split(recipients, ";")

	return notification.NewEmailService(host, parsedPort, from, pwd, parsedRecipients)
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

	tempHistory, err := persistenceService.GetTemperatureRecordings()
	if err != nil {
		log.Printf("Failed to get temp history: %s", err)
		tempHistory = []persistence.TemperatureRecording{}
	}

	alarmHistory, err := persistenceService.GetAlarmRecordings()
	if err != nil {
		log.Printf("Failed to get alarm history: %s", err)
		alarmHistory = []persistence.AlarmRecording{}
	}

	context := map[string]interface{}{
		"currentTemperature":   burner.GetCurrentTemperature(),
		"activeAlarm":          burner.AlarmIsActive(),
		"activeAlarmStartTime": alarmStartTime,
		"temperatureHistory":   tempHistory,
		"alarmHistory":         alarmHistory,
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
