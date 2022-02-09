package main

import (
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"mon-tool-be/models"
	"mon-tool-be/utils"
	"net/http"
	"strconv"
	"time"
)

type Counter struct {
	Current int
	Max     int
}

func InitCron(dbsql *gorm.DB) {
	db, err := scribble.New("store", nil)
	if err != nil {
		panic(err)
	}

	c := cron.New()
	c.AddFunc("* * * * *", func() {
		counter := Counter{}

		if err := db.Read("counter", "counter", &counter); err != nil {
			panic("please add starting point")
		}

		// do stuff

		var monitors []models.Monitor

		dbsql.Find(&monitors)

		// TODO:
		// This process should be in go routine.
		// I still don't have enough time to refactor this piece of code.
		for _, val := range monitors {
			if counter.Current%val.IntervalInMinute == 0 {
				var resp *http.Response
				var record models.Record

				client := http.Client{
					Timeout: time.Duration(val.TimeoutInSecond) * time.Second,
				}
				if val.RequestMethod == models.ReqGet {
					resp, err = client.Get(val.Type + "://" + val.Url)
				} else {
					resp, err = client.Post(val.Type+"://"+val.Url, "application/json", nil)
				}

				if err != nil {
					record.MonitorID = val.ID
					record.Status = "ERROR"
					record.Code = "9999"
					dbsql.Create(&record)
				} else {
					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						record.MonitorID = val.ID
						record.Status = "OK"
						record.Code = strconv.Itoa(resp.StatusCode)
						dbsql.Create(&record)
					} else {
						record.MonitorID = val.ID
						record.Status = "KO"
						record.Code = strconv.Itoa(resp.StatusCode)
						dbsql.Create(&record)

						var alerts []models.Alert
						dbsql.Where("monitor_id = ?", record.MonitorID).Find(&alerts)

						go utils.SendMessage(alerts, "Your server on address "+val.Type+"://"+val.Url+" has stopped working with error code "+record.Code+".")
					}
				}

			}
		}

		// end do stuff

		counter.Current = counter.Current + 1

		if err := db.Write("counter", "counter", counter); err != nil {
			fmt.Println("Error", err)
		}
	})
	c.Start()
}
