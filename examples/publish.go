package main

import (
	"fmt"
	"time"

	"github.com/ericdmann/satori-go"
)

const (
	RTMEndpoint   = "<RTMEndpoint>"
	RTMAppKey     = "<RTMAppKey>"
	RTMRoleName   = "<RTMRoleName>"
	RTMRoleSecret = "<RTMRoleSecret>"
)

func main() {
	rClient, err := rtm.NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		fmt.Println("Error creating RTM client: ", err)
		return
	}

	fmt.Println("RTM client created successfully")

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {

		rtmWire, err := rClient.Publish("test", map[string]interface{}{"message": "hi"})
		if err != nil {
			fmt.Println("Error publishing message: ", err)
			return
		}

		fmt.Println("Message published: ", rtmWire.Action)
	}
}
