package main

import (
	"fmt"
	"github.com/tinygo-org/tinygo/src/machine"
	"time"
	"tinygo.org/x/drivers/dht"
)

func main() {
	fmt.Println("Start measurements...")
	pin := machine.IO17
	dhtSensor := dht.New(pin, dht.DHT11)
	for {
		temp, hum, err := dhtSensor.Measurements()
		if err == nil {
			fmt.Printf("Temperature: %02d.%d°C, Humidity: %02d.%d%%\n", temp/10, temp%10, hum/10, hum%10)
		} else {
			fmt.Printf("Could not take measurements from the sensor: %s\n", error.Error(err))
		}
		// Measurements cannot be updated only 2 seconds. More frequent calls will return the same value
		time.Sleep(time.Second * 2)
	}
}
