package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	subClient mqtt.Client
	pubClient mqtt.Client
	topic     = "#"
)

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func init() {
	localMQTT, err := url.Parse(os.Getenv("LOCALMQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}
	subClient = connect("sub", localMQTT)

	cloudMQTT, err := url.Parse(os.Getenv("CLOUDMQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}
	pubClient = connect("pub", cloudMQTT)
}

func processMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("recieved message from topic: %s\n", msg.Topic())
	pubClient.Publish(msg.Topic(), 0, false, msg.Payload())
}

func main() {
	subClient.Subscribe(topic, 0, processMessage)

	// timer := time.NewTicker(1 * time.Second)
	// for t := range timer.C {
	// 	subClient.Publish(topic, 0, false, t.String())
	// }
	select {}
}
