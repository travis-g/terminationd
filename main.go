package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// TerminationTimeURL is the AWS EC2 Instance meta-data URL for checking Spot
// Instance termination times.
const TerminationTimeURL = `http://169.254.169.254/latest/meta-data/spot/termination-time`

var (
	// Wait is the time (in seconds) between consecutive checks
	Wait = time.Second * 5

	// Client is a non-default HTTP client to use for polling
	Client = &http.Client{
		Timeout: time.Second * 3,
	}
)

// IsTerminating returns true when querying the instance's termination time
// returns a Time-like result, or false and any received errors if a Time-like
// value was not returned.
func IsTerminating() (bool, error) {
	resp, err := Client.Get(TerminationTimeURL)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}

	t, err := time.Parse(time.RFC3339, string(body))
	if err == nil {
		// Successfully parsed a UTC time string:
		log.Printf("Received termination notice for %v\n", t)
		return true, nil
	}

	// Response was unparsable as a time, but could be a full HTTP 404
	return false, nil
}

func main() {
	c := make(chan os.Signal, 1)

	// tick asynchronously in case there is a timeout
	ticker := time.NewTicker(Wait)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			terminating, err := IsTerminating()
			if err != nil {
				log.Println(err)
			}
			if terminating {
				signal.Notify(c)
			}
		}
	}()

	s := <-c
	log.Printf("%v received\n", s)
	os.Exit(1)
}
