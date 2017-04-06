package rtm

import (
	"fmt"
	"testing"
	"time"
)

const (
	channel       = "<ChannelName>"
	RTMEndpoint   = "<RTMEndpoint>"
	RTMAppKey     = "<RTMAppKey>"
	RTMRoleName   = "<RTMRoleName>"
	RTMRoleSecret = "<RTMRoleSecret>"
)

func TestClientPublish(t *testing.T) {

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatalf("Error creating RTM client: ", err)
	}

	timeInABottle := time.Now().String()
	rtWire, err := rtmClient.Publish(channel, timeInABottle)

	if rtWire.Action != "rtm/publish/ok" {
		t.Fatalf("Unsuccessful publish: %s", rtWire.Action)
	}
}

func TestClientSubscribe(t *testing.T) {

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatalf("Error creating RTM client: ", err)
	}
	err = rtmClient.Subscribe(channel)

	if err != nil {
		t.Fatalf("Unsuccessful subscribe: %s", err)
	}
}

func TestClientCancelSubscribe(t *testing.T) {

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatalf("Error creating RTM client: ", err)
	}
	rtmClient.CancelSubscription()

	if rtmClient.SubscriptionName != "" {
		t.Fatalf("Unsuccessful cancel subscribe: %s", rtmClient.SubscriptionName)
	}
}

func TestClientReadSubscriptionChannel(t *testing.T) {
	rPublishClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatalf("Error creating RTM publish client: ", err)
	}

	rReadClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)
	if err != nil {
		t.Fatalf("Error creating RTM read client: %s", err)
	}

	err = rReadClient.Subscribe(channel)
	if err != nil {
		t.Fatalf("Error creating subscription: %s", err)
	}

	go rReadClient.ReadSubscription()

	rtWire, err := rPublishClient.Publish(channel, time.Now().String())
	if err != nil || rtWire.Action != "rtm/publish/ok" {
		t.Fatalf("Error publishing: %s", rtWire.Action)
	}

	item := <-rReadClient.Subscription

	fmt.Println(item)

}

func TestClientAuth(t *testing.T) {
	rtm := RTMClient{
		Endpoint:     RTMEndpoint,
		AppKey:       RTMAppKey,
		RoleName:     RTMRoleName,
		RoleSecret:   RTMRoleSecret,
		Debug:        true,
		Subscription: make(chan RTMWire),
	}

	err := rtm.Connect()
	if err != nil {
		t.Fatalf("Unable to connect: ", err)
	}
}
