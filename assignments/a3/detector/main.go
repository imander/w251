package main

import (
	"fmt"
	"image"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gocv.io/x/gocv"
)

var (
	mqttClient mqtt.Client
	topic      string
	// get device index of the webcam
	deviceID int = getWebcam()

	// get the wait interval in ms between captures
	waitInterval int = getWaitInterval()
)

func init() {
	mqtt, err := url.Parse(os.Getenv("MQTT_URL"))
	if err != nil {
		log.Fatal(err)
	}
	mqttClient = connect("pub", mqtt)

	topic = os.Getenv("MQTT_TOPIC")
	if topic == "" {
		topic, err = os.Hostname()
		if err != nil {
			topic = "default"
		}
	}
	topic += "_png"
}

func main() {
	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	var window *gocv.Window
	if os.Getenv("DISPLAY") != "" {
		window = gocv.NewWindow("Face Detect")
		defer window.Close()
	}

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	// boxColor := color.RGBA{0, 255, 0, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load("/data/haarcascade_frontalface_default.xml") {
		fmt.Println("Error reading cascade file: /data/haarcascade_frontalface_default.xml")
		os.Exit(1)
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image
		for _, r := range rects {
			// gocv.Rectangle(&img, r, boxColor, 3)
			if window != nil {
				window.IMShow(cropImage(img, r))
			}
			croppedImg := cropImage(img, r)
			if err := publishImage(croppedImg); err != nil {
				fmt.Println(err.Error)
				os.Exit(1)
			}
		}

		// show the image in the window
		if window != nil {
			window.IMShow(img)
		}
		fmt.Printf("sleeping %d seconds\n", waitInterval/1000)
		time.Sleep(time.Duration(waitInterval) * time.Millisecond)
	}
}

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

func publishImage(img gocv.Mat) error {
	fmt.Println("publishing image")
	encoded, err := gocv.IMEncode(gocv.PNGFileExt, img)
	if err != nil {
		return fmt.Errorf("error encoding image: %s", err.Error())
	}

	mqttClient.Publish(topic, 0, false, encoded)
	return nil
}

func saveImage(img gocv.Mat) error {
	fmt.Println("saving image")
	encoded, err := gocv.IMEncode(gocv.PNGFileExt, img)
	if err != nil {
		return fmt.Errorf("error encoding image: %s", err.Error())
	}

	file, err := os.OpenFile(
		"data/pic.png",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return fmt.Errorf("error obtaining file handle: %s", err.Error())
	}
	defer file.Close()
	_, err = file.Write(encoded)
	if err != nil {
		return fmt.Errorf("error writing to file: %s", err.Error())
	}
	return nil
}

func cropImage(img gocv.Mat, rect image.Rectangle) gocv.Mat {
	croppedMat := img.Region(rect)
	resultMat := croppedMat.Clone()
	return resultMat
}

func getWebcam() int {
	i := os.Getenv("CAMERA_INDEX")
	if i == "" {
		return 0
	}
	idx, err := strconv.Atoi(i)
	if err != nil {
		fmt.Printf("invalid camera index: %s", i)
		os.Exit(1)
	}
	return idx
}

func getWaitInterval() int {
	i := os.Getenv("WAIT_INTERVAL")
	if i == "" {
		return 1
	}
	interval, err := strconv.Atoi(i)
	if err != nil {
		fmt.Printf("invalid wait interval: %s\n", i)
		os.Exit(1)
	}
	fmt.Printf("setting wait interval: %d\n", interval)
	return interval
}
