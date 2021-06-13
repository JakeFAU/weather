package main

import (
	"log"
	"net/url"
	"os"
	"time"

	aio "github.com/jakefau/goAdafruit"
	"github.com/jakefau/rpi-devices/dev"
	"golang.org/x/exp/io/i2c"
)

// connect to the API
func connect() aio.Client {
	// basic stuff
	username := "JakeFau"
	baseURL := "https://io.adafruit.com/"
	// Hide your key
	key := os.Getenv("ADAFRUIT_IO_KEY")
	// get a client
	client := aio.NewClient(key, username)
	// set the base url, aka the host
	client.BaseURL, _ = url.Parse(baseURL)
	return *client
}

func getFeed(feedKey string, client aio.Client) *aio.Feed {
	feed, _, err := client.Feed.Get(feedKey)
	if err != nil {
		log.Fatal(err)
	}
	return feed
}

func main() {
	io := os.Args[1]
	client := connect()
	//get the feed
	tempFeed := getFeed("weather.temperature", client)
	humidFeed := getFeed("weather.humidity", client)
	pressFeed := getFeed("weather.pressure", client)
	oTempFeed := getFeed("outdoor.temperature", client)
	oHumidFeed := getFeed("outdoor.humidity", client)

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x77)
	if err != nil {
		panic(err)
	}
	b := dev.New(d)
	err = b.Init()
	if err != nil {
		log.Fatal(err)
	}

	dht := dev.NewDHT11()

	if io == "indoor" {
		for {
			t, p, h, _ := b.EnvData()
			client.SetFeed(tempFeed)
			client.Data.Create(&aio.Data{Value: convert64(toFahrenheit(t))})
			client.SetFeed(humidFeed)
			client.Data.Create(&aio.Data{Value: convert64(toMercury(h))})
			client.SetFeed(pressFeed)
			client.Data.Create(&aio.Data{Value: convert64(p)})
			log.Printf("Temp: %fF, Press: %f, Hum: %f%%\n", toFahrenheit(t), toMercury(p), h)
			time.Sleep(time.Second * 20)
		}
	} else {
		for {
			t, h, _ := dht.TempHumidity()
			client.SetFeed(oTempFeed)
			client.Data.Create(&aio.Data{Value: convert64(toFahrenheit(t))})
			client.SetFeed(oHumidFeed)
			client.Data.Create(&aio.Data{Value: convert64(h)})
			log.Printf("Temp %v  Humid %v", toFahrenheit(t), h)
			time.Sleep(time.Second * 20)
		}
	}
}
