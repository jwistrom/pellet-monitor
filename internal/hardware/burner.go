package hardware

import (
	"log"
	"time"
)

type Burner interface {
	// GetCurrentTemperature gets current temperature or -1 in case of error
	GetCurrentTemperature() int

	// AddAlarmListener registers a listener for any type of alarm from the burner
	AddAlarmListener(alarmListener AlarmListener)

	// AlarmIsActive returns whether any alarm is currently active
	AlarmIsActive() bool

	// ActiveAlarmStartTime returns start time of active alarm. If no active alarm it returns nil
	ActiveAlarmStartTime() time.Time
}

type AlarmListener func()

type BurnerImpl struct {
	listeners []AlarmListener
}

func (d *BurnerImpl) ActiveAlarmStartTime() time.Time {
	return time.Now().Add(-time.Duration(30) * time.Minute)
}

func (d *BurnerImpl) AlarmIsActive() bool {
	return true
}

func (d *BurnerImpl) GetCurrentTemperature() int {
	return 76
}

func (d *BurnerImpl) AddAlarmListener(alarmListener AlarmListener) {
	log.Println("Listener added")
	d.listeners = append(d.listeners, alarmListener)
}
