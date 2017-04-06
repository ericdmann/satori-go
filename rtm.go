package rtm

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	subscriptionMutex sync.Mutex
)

const (
	RTMVersion = "v2"
)

type RTMWire struct {
	Action string      `json:"action"`
	Body   RTMWireBody `json:"body"`
	ID     string      `json:"id"`
}

type RTMWireBody struct {
	Credentials    RTMWireBodyCredentials `json:"credentials,omitempty"`
	Data           RTMWireBodyData        `json:"data,omitempty"`
	Method         string                 `json:"method,omitempty"`
	Channel        string                 `json:"channel,omitempty"`
	Message        string                 `json:"message,omitempty"`
	Prefix         string                 `json:"prefix,omitempty"`
	SubscriptionID string                 `json:"subscription_id,omitempty"`
	Error          string                 `json:"error,omitempty"`
	Reason         string                 `json:"reason,omitempty"`
	Position       string                 `json:"position,omitempty"`
	Messages       []string               `json:"messages,omitempty"`
}

type RTMWireBodyCredentials struct {
	Hash string `json:"hash,omitempty"`
}

type RTMWireBodyData struct {
	Nonce string `json:"nonce,omitempty"`
	Role  string `json:"role,omitempty"`
}

type RTMClient struct {
	Endpoint         string
	AppKey           string
	RoleName         string
	RoleSecret       string
	Debug            bool
	Subscription     chan RTMWire
	SubscriptionName string
	WSClient         *websocket.Conn
}

func NewClient(endpoint string, appkey string, rolename string, rolesecret string, debug bool) (RTMClient, error) {
	var rtm RTMClient

	if len(endpoint) < 15 || !strings.Contains(endpoint, "ws") {
		return rtm, errors.New("Invalid endpoint provided.")
	}
	if len(appkey) < 15 {
		return rtm, errors.New("Invalid appkey provided.")
	}
	if rolename == "" {
		return rtm, errors.New("Invalid rolename provided.")
	}
	if rolesecret == "" {
		return rtm, errors.New("Invalid rolesecret provided.")
	}

	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()
	rtm = RTMClient{
		Endpoint:     endpoint,
		AppKey:       appkey,
		RoleName:     rolename,
		RoleSecret:   rolesecret,
		Debug:        debug,
		Subscription: make(chan RTMWire),
	}

	err := rtm.Connect()
	return rtm, err
}

func computeHash(nonce string, state string) string {
	key := []byte(state)
	h := hmac.New(md5.New, key)
	h.Write([]byte(nonce))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
