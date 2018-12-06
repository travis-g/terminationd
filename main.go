package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// InstanceActionURL is the AWS EC2 Instance meta-data URL for checking Spot
// Instance termination time.
const InstanceActionURL = `http://169.254.169.254/latest/meta-data/spot/instance-action`

var (
	// Wait is the time (in seconds) between consecutive checks
	Wait = time.Second * 5

	// Client is a non-default HTTP client to use for polling
	Client = &http.Client{
		Timeout: time.Second * 3,
	}

	nilTime = (time.Time{}).UnixNano()
)

// InstanceAction is the EC2 metadata response for Spot instance termination
// actions. If a Spot Instance is not about to be terminated the response will
// be an HTTP 404, and this object will consist of empty or zero-like values.
type InstanceAction struct {
	Action string    `json:"action"`
	Time   time.Time `json:"time"`
}

// GetInstanceAction queries the requested URL for an EC2 instance-action
// result.
func GetInstanceAction(url string) (InstanceAction, error) {
	instanceAction := &InstanceAction{}
	resp, err := Client.Get(url)
	if err != nil {
		return *instanceAction, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return *instanceAction, err
	}

	err = json.Unmarshal(body, instanceAction)
	return *instanceAction, err
}

// IsTerminating returns true when the instance action's termination time
// returns a non-default Time-like result.
func (ia *InstanceAction) IsTerminating() bool {
	return ia.Time.UnixNano() != nilTime
}

func main() {
	c := make(chan os.Signal, 1)

	// tick asynchronously in case there is a timeout
	ticker := time.NewTicker(Wait)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			ia, err := GetInstanceAction(InstanceActionURL)
			if err != nil {
				log.Println(err)
			}
			if ia.IsTerminating() {
				log.Println(ia.Time, ia.Action)
				signal.Notify(c)
			}
		}
	}()

	log.Println("Initialized")

	s := <-c
	log.Printf("%v received\n", s)
	os.Exit(1)
}
