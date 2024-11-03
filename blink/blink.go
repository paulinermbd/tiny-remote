// blink.go

package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})

	for {
		fmt.Println("Blink")
		led.High()
		time.Sleep(time.Millisecond * 500)
		fmt.Println("Stop blinking")
		led.Low()
		time.Sleep(time.Millisecond * 500)
	}
}
