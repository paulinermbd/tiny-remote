package main

import (
	"time"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var advUntil = time.Now().Add(5 * time.Minute)
var advState = true

func main() {
	println("start process BLE")
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	// Define the peripheral device info.
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "tinyremote",
		ServiceUUIDs: []bluetooth.UUID{bluetooth.ServiceUUIDHumanInterfaceDevice},
	}))

	adapter.SetConnectHandler(func(device bluetooth.Device, connected bool) {
		if connected {
			println("connected, not advertising...")
			advState = false
		} else {
			println("disconnected, advertising...")
			advState = true
			advUntil = time.Now().Add(5 * time.Minute)
		}
	})
	must("start adv", adv.Start())

	/*must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: ,
				UUID:   bluetooth.CharacteristicUUIDMicrobitButtonAState,
				Value:  []byte{0, heartRate},
				Flags:  bluetooth.CharacteristicNotifyPermission,
			},
		},
	}))*/

	println("advertising...")
	address, _ := adapter.Address()
	for {
		if advState && time.Now().After(advUntil) {
			println("timeout, not advertising...")
			advState = false
			must("stop adv", adv.Stop())
		}
		println("Go Bluetooth /", address.MAC.String())
		time.Sleep(time.Second)
	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

/*func main() {
//TODO 1: pairing the device to the PC
fmt.Println("enable BLE...")
adapter.Enable()

fmt.Println("scanning...")
err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	println("found device:", device.Address.String(), device.RSSI, device.LocalName())
})
fmt.Errorf("failed error: %v", err)

//TODO 2: mapping between cardboard input and keyboard by using characteristic
/*bluetooth.CharacteristicConfig{
	Handle:     nil,
	UUID:       bluetooth.UUID{},
	Value:      nil,
	Flags:      0,
	WriteEvent: nil,
}*/

//TODO 3: check everything is ok on GoogleSlide

//TODO 4: test it on multiple devices
//}
