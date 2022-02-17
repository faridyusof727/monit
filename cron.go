package main

import (
	"context"
	"crypto/tls"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"mon-tool-be/models"
	"mon-tool-be/utils"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

				startTime := time.Now()
				client := http.Client{
					Timeout: time.Duration(val.TimeoutInSecond) * time.Second,
				}
				if val.RequestMethod == models.ReqGet {
					resp, err = client.Get(val.Type + "://" + strings.TrimSpace(val.Url))
				} else {
					resp, err = client.Post(val.Type+"://"+strings.TrimSpace(val.Url), "application/json", nil)
				}

				if val.Type == models.TypeHttps {

					record.SSLStatus = "OK"

					u, _ := url.Parse(val.Type + "://" + strings.TrimSpace(val.Url))

					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

					conf := &tls.Config{
						InsecureSkipVerify: false,
					}

					dialer := tls.Dialer{
						Config: conf,
					}

					tlsConn, err := dialer.DialContext(ctx, "tcp", u.Host+":443")
					cancel()
					if err != nil {
						record.SSLStatus = "KO"
					}

					var conn *tls.Conn
					if record.SSLStatus == "OK" {
						conn = tlsConn.(*tls.Conn)
					}

					if record.SSLStatus == "OK" {
						err = conn.VerifyHostname(u.Host)
						if err != nil {
							record.SSLStatus = "KO"
							_ = conn.Close()
						}
					}

					if record.SSLStatus == "OK" {
						certs := conn.ConnectionState().PeerCertificates
						for _, cert := range certs {
							record.SSLIssuer = fmt.Sprint(cert.Issuer)
							record.SSLExpiry = fmt.Sprint(cert.NotAfter.Format("2006-January-02"))
							record.SSLCommonName = fmt.Sprint(cert.Issuer.CommonName)
						}
					}

					if record.SSLStatus == "OK" {
						_ = conn.Close()
					}
				}

				duration := time.Now().Sub(startTime)
				if err != nil {
					record.MonitorID = val.ID
					record.Status = "ERROR"
					record.Code = "9999"
					record.ResponseTime = duration.Milliseconds()
					dbsql.Create(&record)
				} else {
					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						record.MonitorID = val.ID
						record.Status = "OK"
						record.ResponseTime = duration.Milliseconds()
						record.Code = strconv.Itoa(resp.StatusCode)
						dbsql.Create(&record)
					} else {
						record.MonitorID = val.ID
						record.Status = "KO"
						record.ResponseTime = duration.Milliseconds()
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
