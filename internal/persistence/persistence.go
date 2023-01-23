package persistence

import (
	"math/rand"
	"time"
)

type Service interface {
	StoreTemperature(temperature int) error
	GetTemperatureRecordings() ([]TemperatureRecording, error)
	StoreAlarm(time time.Time) error
	GetAlarmRecordings() ([]time.Time, error)
}

type TemperatureRecording struct {
	Time        time.Time `json:"time"`
	Temperature int       `json:"temperature"`
}

type ServiceImpl struct {
	tempStore  []TemperatureRecording
	alarmStore []time.Time
}

func (d *ServiceImpl) GetTemperatureRecordings() (recs []TemperatureRecording, e error) {
	for i := 0; i < 1000; i++ {
		temp := rand.Intn(50) + 50
		recs = append(recs, TemperatureRecording{time.Now().Add(-time.Duration(i*5) * time.Minute), temp})
	}
	return
}

func (d *ServiceImpl) StoreAlarm(time time.Time) error {
	d.alarmStore = append(d.alarmStore, time)
	return nil
}

func (d *ServiceImpl) GetAlarmRecordings() ([]time.Time, error) {
	recs := []time.Time{
		time.Now().Add(-time.Duration(2) * time.Hour),
		time.Now().Add(-time.Duration(8) * time.Hour),
		time.Now().Add(-time.Duration(24) * time.Hour),
		time.Now().Add(-time.Duration(36) * time.Hour),
		time.Now().Add(-time.Duration(40) * time.Hour),
		time.Now().Add(-time.Duration(168) * time.Hour),
		time.Now().Add(-time.Duration(336) * time.Hour),
		time.Now().Add(-time.Duration(360) * time.Hour),
		time.Now().Add(-time.Duration(362) * time.Hour),
		time.Now().Add(-time.Duration(365) * time.Hour),
	}
	return recs, nil
}

func (d *ServiceImpl) StoreTemperature(temperature int) error {
	d.tempStore = append(d.tempStore, TemperatureRecording{
		Time:        time.Now(),
		Temperature: temperature,
	})
	return nil
}
