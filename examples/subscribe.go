package main

import (
	"fmt"
	"time"

	"github.com/ericdmann/satori-go"
)

const (
	channel       = "<ChannelName>"
	RTMEndpoint   = "<RTMEndpoint>"
	RTMAppKey     = "<RTMAppKey>"
	RTMRoleName   = "<RTMRoleName>"
	RTMRoleSecret = "<RTMRoleSecret>"
)

func main() {
	rPublishClient, err := rtm.NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		fmt.Println("Error creating RTM client: ", err)
		return
	}

	go sendMessages(rPublishClient, channel)

	rReadClient, err := rtm.NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		fmt.Println("Error creating RTM client: ", err)
		return
	}

	err = rReadClient.Subscribe(channel)
	if err != nil {
		fmt.Println("Error subscribing to channel: ", err)
		return
	}

	go rReadClient.ReadSubscription()
	for {
		item := <-rReadClient.Subscription
		fmt.Println("(Received) : ", item.Body.Messages)
	}
}

func sendMessages(rClient rtm.RTMClient, channel string) {
	ticker := time.NewTicker(time.Second * 1)

	for range ticker.C {

		rtWire, err := rClient.Publish(channel, time.Now().String())
		if err != nil {
			fmt.Println("Error publishing message: ", err)
			continue
		}
		fmt.Println("(Published): ", rtWire)
	}
}
