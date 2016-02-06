package gotification

import (
	"errors"
	"github.com/alexjlockwood/gcm"
	"github.com/timehop/apns"
	"os"
	"sync"
)

// MUST: apn must be written constantry first

func checkErr(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

type Config struct {
	APN string
	GCM string
}

type Notification struct {
	Message      string
	APNReceivers []string
	GCMReceivers []string
}

var (
	apnClient apns.Client
	gcmClient *gcm.Sender
)

func (c *Config) Set() (err error) {
	// here apns certfile check
	apnClient, err = apns.NewClientWithCert(apns.SandboxGateway, *c.APN)
	checkErr(err)
	// here gcm
	if c.GCM == "" {
		err = errors.New("gcm apikey is not setted.")
		checkErr(err)
	} else {
		gcmClient = &gcm.Sender{ApiKey: c.GCM}
	}
	return
}

func (n *Notification) Send() (errList Notification, err error) {
	// GCM
	apnTask := make(chan bool)
	gcmTask := make(chan bool)
	var (
		apnWG sync.WaitGroup
		gcmWG sync.WaitGroup
	)
	go func() { // parallel between APN and GCM
		// APN
		p := apns.NewPayload()
		p.APS.Alert.Body = n.Message
		p.APS.ContentAvailable = 1
		m := apns.NewNotification()
		m.Payload = p
		m.Priority = apns.PriorityImmediate
		apnWG.Add(len(n.APNReceivers))
		for _, v := range APNReceivers {
			go func() {
				m.DeviceToken = v
				apnClient.Send(m)
				apnWG.Done()
			}() // parallel for APN
		}
		// HERE err handling for APN
		go func() {
			for _, f := range apnClient.FailedNotifs {
				errList.APNReceivers = append(errList.APNReceivers, f.Notif.ID)
			}
		}()
		apnWG.Wait()
		apnTask <- true
	}()
	// GCM
	go func() {
		gcmWG.Add(len(n.GCMReceivers))
		data := map[string]interface{}{"message": n.Message}
		for _, v := range GCMReceivers() {
			go func() {
				d := gcm.NewMessage(data, v)
				_, err := gcmClient.Send(d, 0)
				if err != nil {
					errList.GCMReceivers = append(errList.GCMReceivers, v)
				}
				gcmWG.Done()
			}()
		}
		gcmWG.Wait()
		gcmTask <- true
	}()
	<-apnTask
	<-gcmTask
	return
}
