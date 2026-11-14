package accord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var logMessageEvNum = 0

func logEvent(subj string, ts time.Time, ev any) {
	if ev == nil {
		return
	}

	jsn, _ := json.MarshalIndent(ev, "", "  ")

	now := time.Now()
	name := fmt.Sprintf(
		"logs/%s-%d-%d-%s.json",
		now.Format("2006-01-02T15-04-05-000"),
		ts.Unix(),
		logMessageEvNum,
		subj,
	)

	logMessageEvNum++

	log.Printf("[LOG MESSAGE EVENT] name=%s", name)
	err := os.WriteFile(name, jsn, 0664)
	if err != nil {
		log.Printf("[LOG MESSAGE EVENT] failed to write to file '%s': %v", name, err)
	}
}
