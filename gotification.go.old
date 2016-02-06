package gotification

import (
	"fmt"
	"github.com/alexjlockwood/gcm"
	"github.com/mattprice/Go-APNs"
	"io/ioutil"
	"sync"
	"time"
)

var (
	apnauth string
	gcmauth string
)

type Config struct {
	APNCertFile string
	GCMAPIKey   string
}

type Notification struct {
	Message          string
	AndroidReceivers []string
	IOSReceivers     []string
}

func (c *Config) Set() {
	apnauth = c.APNCertFile
	gcmauth = c.GCMAPIKey
}

func (n *Notification) Notify() bool {
	task1 := make(chan bool)
	task2 := make(chan bool)
	err1 := []string{}
	err2 := []string{}
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	go func() {
		wg1.Add(len(n.AndroidReceivers))
		for _, v := range n.AndroidReceivers {
			// go gcmNotify(v, n.Message, &wg1)
			go gcmNotify(v, n.Message, &err1, &wg1)
		}
		wg1.Wait()
		task1 <- true
	}()
	go func() {
		wg2.Add(len(n.IOSReceivers))
		for _, v := range n.IOSReceivers {
			// go apnNotify(v, n.Message, &wg2)
			go apnNotify(v, n.Message, &err2, &wg2)
		}
		wg2.Wait()
		task2 <- true
	}()
	<-task1
	<-task2
	for _, v := range err1 {
		println(v)
	}
	for _, v := range err2 {
		println(v)
	}
	return true
}

func test(att string, msg string, i int, wg *sync.WaitGroup, err *[]int) (status bool) {
	content := []byte(fmt.Sprintf("%s %s %s\n", time.Now(), att, msg))
	ioutil.WriteFile("/Users/QsF/tmp/go-testfile.log", content, 0755)
	fmt.Println(time.Now(), att, msg, i)
	time.Sleep(1 * time.Second)
	if i == 32 || i == 33 || i == 564 || i == 921 || i == 612 || i == 741 || i == 742 || i == 743 {
		status = false
		*err = append(*err, i)
	} else {
		status = true
	}
	wg.Done()
	return
}

func apnNotify(att string, msg string, errHandle *[]string, wg *sync.WaitGroup) (status bool) {
	status = true
	err := apns.LoadCertificateFile(false, apnauth)
	if err != nil {
		status = false
	}
	payload := &apns.Notification{Alert: fmt.Sprintf(msg), Badge: 0, Sandbox: true}
	payload.SetExpiryDuration(24 * time.Hour)
	err = payload.SendTo(att)
	if err != nil {
		status = false
	}
	wg.Done()
	if status == false {
		*errHandle = append(*errHandle, att)
	}
	return
}

func gcmNotify(att string, msg string, errHandle *[]string, wg *sync.WaitGroup) (status bool) {
	status = true
	d := map[string]interface{}{"message": msg}
	sender := &gcm.Sender{ApiKey: gcmauth}
	data := gcm.NewMessage(d, att)
	_, err := sender.Send(data, 0)
	if err != nil {
		status = false
	}
	wg.Done()
	if status == false {
		*errHandle = append(*errHandle, att)
	}
	return
}
