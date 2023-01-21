package persistence

import (
	"math/rand"
	"time"
)

type Service interface {
	StoreTemperature(temperature int) error
	GetTemperatureHistory() ([]TemperatureRecording, error)
}

type TemperatureRecording struct {
	Time        time.Time `json:"time"`
	Temperature int       `json:"temperature"`
}

type DummyService struct {
	store []int
}

func (d *DummyService) GetTemperatureHistory() (recs []TemperatureRecording, e error) {
	for i := 0; i < 1000; i++ {
		temp := rand.Intn(50) + 50
		recs = append(recs, TemperatureRecording{time.Now().Add(-time.Duration(i*5) * time.Minute), temp})
	}
	return
}

func (d *DummyService) StoreTemperature(temperature int) error {
	d.store = append(d.store, temperature)
	return nil
}
