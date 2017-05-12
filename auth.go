package rtm

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

// Connect completes the handshake and authentication via role_secret to the RTM endpoint
func (r *RTMClient) Connect() error {
	var err error

	config, err := websocket.NewConfig(r.Endpoint+"/"+RTMVersion+"?appkey="+r.AppKey, "http://localhost")
	if err != nil {
		return err
	}

	r.WSClient, err = websocket.DialConfig(config)
	if err != nil {
		return err
	}

	//	handshake
	wire := RTMWire{
		Action: "auth/handshake",
		Body: RTMWireBody{
			Method: "role_secret",
			Data: RTMWireBodyData{
				Role: r.RoleName,
			},
		},
	}

	jsonStr, err := json.Marshal(wire)
	if err != nil {
		return err
	}
	r.WSClient.Write(jsonStr)
	var retMsg = make([]byte, 512)

	n, err := r.WSClient.Read(retMsg)
	if err != nil {
		return err
	}

	var rtmResponse RTMWire
	if r.Debug {
		fmt.Printf("(Handshake): %s.\n", retMsg[:n])
	}

	err = json.Unmarshal(retMsg[:n], &rtmResponse)
	if err != nil {
		return err
	}

	//	authenticate
	hash := computeHash(rtmResponse.Body.Data.Nonce, r.RoleSecret)

	wire = RTMWire{
		Action: "auth/authenticate",
		Body: RTMWireBody{
			Method: "role_secret",
			Data: RTMWireBodyData{
				Role: r.RoleName,
			},
			Credentials: RTMWireBodyCredentials{
				Hash: hash,
			},
		},
	}
	jsonStr, _ = json.Marshal(wire)
	r.WSClient.Write(jsonStr)
	retMsg = make([]byte, 512)
	n, err = r.WSClient.Read(retMsg)
	if err != nil {
		return err
	}

	if r.Debug {
		fmt.Printf("(Authentication): %s.\n", retMsg[:n])
	}

	err = json.Unmarshal(retMsg[:n], &rtmResponse)
	if err != nil {
		return err
	}
	return err
}
