package emitters

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	_ "github.com/mattn/go-sqlite3" // Per docs
	"github.com/myoung34/bluey/bluey"
	"log"
)

type SQLite struct {
	Enabled       bool
	File          string
	DeviceUUIDMap string `json:"device_uuids"`
}

func SQLiteEmit(payload blue.Payload, emitterConfig interface{}) (string, error) {
	sqlite := SQLite{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &sqlite)

	nickname := getNickNameFromUUIDs(sqlite.DeviceUUIDMap, payload.ID)

	db, err := sql.Open("sqlite3", sqlite.File)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
      CREATE TABLE IF NOT EXISTS data(
        id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
        minor INTEGER,
        major INTEGER,
        nickname VARCHAR(16),
        mac VARCHAR(17),
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)
       `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		level.Error(blue.Logger).Log("emitters.sqlite", err)
		return "", err
	}

	insertStmt := fmt.Sprintf(
		"insert into data (minor,major,nickname,mac) values (%d,%d,'%s','%s')",
		int(payload.Minor),
		int(payload.Major),
		nickname,
		payload.Mac,
	)
	level.Debug(blue.Logger).Log("emitters.sqlite", insertStmt)
	_, err = db.Exec(insertStmt)
	if err != nil {
		level.Error(blue.Logger).Log("emitters.sqlite", err)
		return "", err
	}

	return insertStmt, nil
}
