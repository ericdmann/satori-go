package rtm

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	channel       string
	RTMEndpoint   string
	RTMAppKey     string
	RTMRoleName   string
	RTMRoleSecret string
)

func LoadConfigFromEnvironment(t *testing.T) {
	if RTMEndpoint == "" || RTMAppKey == "" || RTMRoleSecret == "" {
		channel = os.Getenv("SATORI_CHANNEL_NAME")
		RTMEndpoint = os.Getenv("SATORI_ENDPOINT")
		RTMAppKey = os.Getenv("SATORI_APP_KEY")
		RTMRoleName = os.Getenv("SATORI_ROLE")
		RTMRoleSecret = os.Getenv("SATORI_SECRET")

		fmt.Println("SATORI ENV CONFIG:")
		fmt.Println("==================")
		fmt.Println("SATORI_CHANNEL_NAME:", channel)
		fmt.Println("SATORI_ENDPOINT:", RTMEndpoint)
		fmt.Println("SATORI_APP_KEY:", RTMAppKey)
		fmt.Println("SATORI_ROLE:", RTMRoleName)
		fmt.Println("SATORI_SECRET:", RTMRoleSecret)
		fmt.Println()
	}
}

func TestClientPublish(t *testing.T) {
	LoadConfigFromEnvironment(t)

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatal("Error creating RTM client: ", err)
	}

	msgStruct := struct {
		EventName string                 `json:"event_name"`
		Data      map[string]interface{} `json:"data"`
	}{
		EventName: "TestPublishEvent",
		Data: map[string]interface{}{
			"time":    time.Now().Unix(),
			"message": "Here is a really cool message...",
		},
	}
	rtWire, err := rtmClient.Publish(channel, msgStruct)


	if rtWire.Action != "rtm/publish/ok" {
		t.Fatalf("Unsuccessful publish: %s", rtWire.Action)
	} else if err != nil {
		t.Fatalf("Error encoutered publishing: %s", err.Error())
	}
}

func TestClientSubscribe(t *testing.T) {
	LoadConfigFromEnvironment(t)

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatal("Error creating RTM client: ", err)
	}
	err = rtmClient.Subscribe(channel)

	if err != nil {
		t.Fatalf("Unsuccessful subscribe: %s", err)
	}
}

func TestClientCancelSubscribe(t *testing.T) {
	LoadConfigFromEnvironment(t)

	rtmClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatal("Error creating RTM client: ", err)
	}
	rtmClient.CancelSubscription()

	if rtmClient.SubscriptionName != "" {
		t.Fatalf("Unsuccessful cancel subscribe: %s", rtmClient.SubscriptionName)
	}
}

func TestClientReadSubscriptionChannel(t *testing.T) {
	LoadConfigFromEnvironment(t)

	rPublishClient, err := NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)

	if err != nil {
		t.Fatal("Error creating RTM publish client: ", err)
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

	rtWire, err := rPublishClient.Publish(channel, map[string]interface{}{"time": time.Now().String()})
	if err != nil || rtWire.Action != "rtm/publish/ok" {
		t.Fatalf("Error publishing: %s", rtWire.Action)
	}

	item := <-rReadClient.Subscription

	fmt.Println(item)

}

func TestClientAuth(t *testing.T) {
	LoadConfigFromEnvironment(t)

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
		t.Fatal("Unable to connect: ", err)
	}
}
