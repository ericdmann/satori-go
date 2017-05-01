package rtm

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Publish will publish a message to a specific RTM channel
func (r *RTMClient) Publish(channel string, message interface{}) (RTMWire, error) {
	var rtmResponse RTMWire

	rawMessage, err := r.ConvertToRawMessage(message)
	if err != nil {
		return rtmResponse, err
	}

	wire := RTMWire{
		Action: "rtm/publish",
		Body: RTMWireBody{
			Channel: channel,
			Message: rawMessage,
		},
	}

	jsonStr, err := json.Marshal(wire)
	if err != nil {
		return rtmResponse, err
	}

	if r.Debug {
		fmt.Printf("(Publish) Raw Message: %s\n", string(jsonStr))
	}

	r.WSClient.Write(jsonStr)
	retMsg := make([]byte, 512)
	n, err := r.WSClient.Read(retMsg)
	if err != nil {
		return rtmResponse, err
	}

	if r.Debug {
		fmt.Printf("(Publish): %s.\n", retMsg[:n])
	}
	err = json.Unmarshal(retMsg[:n], &rtmResponse)
	return rtmResponse, err
}

// Subscribe connects to the r.WSClient to the provided RTM channel for subscription events
func (r *RTMClient) Subscribe(channel string) error {

	wire := RTMWire{
		Action: "rtm/subscribe",
		Body: RTMWireBody{
			Channel:        channel,
			SubscriptionID: channel,
		},
	}
	jsonStr, err := json.Marshal(wire)
	if err != nil {
		return err
	}
	r.WSClient.Write(jsonStr)
	retMsg := make([]byte, 512)
	n, err := r.WSClient.Read(retMsg)
	if r.Debug {
		fmt.Printf("(Subscribe): %s.\n", retMsg[:n])
	}
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()
	r.SubscriptionName = channel
	return err
}

// CancelSubscription will cancel the active subscription
func (r *RTMClient) CancelSubscription() {
	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()
	r.SubscriptionName = ""
}

// ReadSubscription reads from the stream provided by r.WSClient if subscribed
func (r *RTMClient) ReadSubscription() error {
	for {
		//	Validate this subscription hasn't been cancelled
		subscriptionMutex.Lock()
		if r.SubscriptionName == "" {
			subscriptionMutex.Unlock()
			return errors.New("No active subscription")
		}
		subscriptionMutex.Unlock()

		var rtmResponse RTMWire
		retMsg := make([]byte, 512)
		n, err := r.WSClient.Read(retMsg)

		if err != nil {
			return err
		}

		if r.Debug {
			fmt.Printf("(Subscription): %s.\n", retMsg[:n])
		}

		err = json.Unmarshal(retMsg[:n], &rtmResponse)
		if err != nil {
			return err
		}
		r.Subscription <- rtmResponse
	}
}
