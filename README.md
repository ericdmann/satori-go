# satori-go

satori-go allows you to interact with the Satori platform to create server-based or mobile applications that use the RTM to publish and subscribe.

Sign up/sign in here: https://developer.satori.com/#/

## Setup/Usage

Publisher and subscription examples are located inside of `/examples`

### Create a client

```go
	rtmClient, err := rtm.NewClient(RTMEndpoint, RTMAppKey, RTMRoleName, RTMRoleSecret, true)
```

### Initialization properties
	RTMEndpoint   		//	See satori documentation for endpoint
	RTMAppKey     		//	See satori documentation for app key
	RTMRoleName   		//	See satori documentation for role name
	RTMRoleSecret 		//	See satori documentation for role secret
	Debug			//	Triggers logging on each RTM action


### Publish to a channel
Once a client has been created, publishing to a channel is quite simple. 

```go
	rtWire, err := rtmClient.Publish(channel, time.Now().String())
```

### Subscribe to a channel
Each client can subscribe to exacly one channel at a time. Subscriptions can be changed without reconnecting. 

```go
	err = rtmClient.Subscribe(channel)
```

### Cancelling a subscription
Because there is only one active subscription per client, you can cancel without providing a channel name.

```go
	err = rtmClient.CancelSubscription()
```

### Reading a subscription
Each client has a subscription channel `chan RTWire`. Trigger items to be pulled off of the socket and put on this channel by calling `ReadSubscription()`
```go
	go rtmClient.ReadSubscription()
	for {
		item := <-rtmClient.Subscription
		fmt.Println("(Received) : ", item.Body.Messages)
	}
```

### Example
Publisher and subscription examples are located inside of `/examples`

```go
	go run subscribe.go
```

###	Testing
Update `rtm_test.go` with your Satori account info, then just use `go test --cover`

```go
PASS
coverage: 82.2% of statements
ok	1.942s
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D


## TODO

