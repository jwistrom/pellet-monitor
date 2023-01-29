package persistence

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresqlServiceImpl struct {
	db *sql.DB
}

func (p *PostgresqlServiceImpl) StoreTemperature(temperature int) error {
	_, err := p.db.Exec("INSERT INTO temperature (time, temperature) VALUES($1, $2)", time.Now(), temperature)
	p.deleteOldTemperatureRecordings()
	return err
}

func (p *PostgresqlServiceImpl) GetTemperatureRecordings() ([]TemperatureRecording, error) {
	rows, err := p.db.Query("SELECT time, temperature FROM temperature")

	recordings := make([]TemperatureRecording, 0)
	if err == nil {
		for rows.Next() {
			var rec TemperatureRecording
			err := rows.Scan(&rec.Time, &rec.Temperature)
			if err != nil {
				return nil, err
			}
			recordings = append(recordings, rec)
		}
	}

	return recordings, err
}

func (p *PostgresqlServiceImpl) StoreAlarm(time time.Time) error {
	_, err := p.db.Exec("INSERT INTO alarm (timestamp) VALUES($1)", time)
	p.deleteOldAlarmRecordings()
	return err
}

func (p *PostgresqlServiceImpl) GetAlarmRecordings() ([]time.Time, error) {
	rows, err := p.db.Query("SELECT timestamp FROM alarm")
	recordings := make([]time.Time, 0)
	if err == nil {
		for rows.Next() {
			var t time.Time
			err := rows.Scan(&t)
			if err != nil {
				return nil, err
			}
			recordings = append(recordings, t)
		}
	}
	return recordings, err
}

func (p *PostgresqlServiceImpl) deleteOldTemperatureRecordings() {
	_, err := p.db.Exec("DELETE FROM temperature where time < $1", time.Now().AddDate(0, 0, -7))
	if err != nil {
		log.Println("Failed to delete old temp recs")
	}
}

func (p *PostgresqlServiceImpl) deleteOldAlarmRecordings() {
	_, err := p.db.Exec("DELETE FROM alarm where timestamp < $1", time.Now().AddDate(0, -1, 0))
	if err != nil {
		log.Println("Failed to delete old alarm recs")
	}
}

func tableExists(tableName string, db *sql.DB) bool {
	rows, err := db.Query("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename  = $1)", tableName)
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}

	var res bool
	for rows.Next() {
		if rows.Scan(&res) != nil {
			log.Fatalln("Failed to read table exists-query")
		}
	}
	return res
}

func NewPostgresService(user string, pwd string, dbName string) *PostgresqlServiceImpl {
	connStr := "postgresql://" + user + ":" + pwd + "@localhost/" + dbName + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Failed to connect to db", err)
	}

	if !tableExists("alarm", db) {
		log.Fatalln("Table 'alarm' does not exist")
	}
	if !tableExists("temperature", db) {
		log.Fatalln("Table 'temperature' does not exist")
	}

	return &PostgresqlServiceImpl{db: db}
}
