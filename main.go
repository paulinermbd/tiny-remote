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

	// Define the peripheral device info.
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "tinyremote",
		ServiceUUIDs: []bluetooth.UUID{
			bluetooth.ServiceUUIDDeviceInformation,
			bluetooth.ServiceUUIDBattery,
			bluetooth.ServiceUUIDHumanInterfaceDevice,
		},
	}))

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

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDDeviceInformation,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDManufacturerNameString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("Nice Keyboards"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDModelNumberString,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("nice!nano"),
			},
			{
				UUID:  bluetooth.CharacteristicUUIDPnPID,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{0x02, 0x8a, 0x24, 0x66, 0x82, 0x34, 0x36},
				//Value: []byte{0x02, uint8(0x10C4 >> 8), uint8(0x10C4 & 0xff), uint8(0x0001 >> 8), uint8(0x0001 & 0xff)},
			},
		},
	})
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDBattery,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDBatteryLevel,
				Value: []byte{80},
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
		},
	})

	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDGenericAccess,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				UUID:  bluetooth.CharacteristicUUIDDeviceName,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("tinygo-corne"),
			},
			{

				UUID:  bluetooth.New16BitUUID(0x2A01),
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte{uint8(0x03c4 >> 8), uint8(0x03c4 & 0xff)}, /// []byte(strconv.Itoa(961)),
			},
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDPeripheralPreferredConnectionParameters,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{0x02},
			// },

			// // 		//
		},
	})

	// hid
	adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.ServiceUUIDHumanInterfaceDevice,
		/*
			 - hid information r
			 - report map r
			 - report nr
			   - client charecteristic configuration
			   - report reference
			- report nr
			   - client charecteristic configuration
			   - report reference
			- hid control point wnr
		*/
		Characteristics: []bluetooth.CharacteristicConfig{
			// {
			// 	UUID:  bluetooth.CharacteristicUUIDHIDInformation,
			// 	Flags: bluetooth.CharacteristicReadPermission,
			// 	Value: []byte{uint8(0x0111 >> 8), uint8(0x0111 & 0xff), uint8(0x0002 >> 8), uint8(0x0002 & 0xff)},
			// },
			{
				//Handle: &reportmap,
				UUID:  bluetooth.CharacteristicUUIDReportMap,
				Flags: bluetooth.CharacteristicReadPermission,
				Value: []byte("reportMap"),
			},
			{

				//Handle: &reportIn,
				UUID:  bluetooth.CharacteristicUUIDReport,
				Value: []byte("test"),
				Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
			},
			{
				// protocl mode
				UUID:  bluetooth.New16BitUUID(0x2A4E),
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission | bluetooth.CharacteristicReadPermission,
				// Value: []byte{uint8(1)},
				// WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
				// 	print("protocol mode")
				// },
			},
			{
				UUID:  bluetooth.CharacteristicUUIDHIDControlPoint,
				Flags: bluetooth.CharacteristicWriteWithoutResponsePermission,
				//	Value: []byte{0x02},
			},
		},
	})

	println("advertising...")

	must("start adv", adv.Start())

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
