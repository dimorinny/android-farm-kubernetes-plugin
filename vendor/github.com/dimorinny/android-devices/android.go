package android

import (
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"
	"log"
)

const (
	usbInterfaceAdbClass    gousb.Class    = 0xFF
	usbInterfaceAdbSubclass gousb.Class    = 0x42
	usbInterfaceAdbProtocol gousb.Protocol = 0x1
)

func Devices() ([]*Device, error) {
	context := gousb.NewContext()
	defer func() {
		err := context.Close()
		if err != nil {
			log.Println("Android devices detector: failed to close context", err)
		}
	}()

	devices, err := getDevicesDescriptions(context)
	if err != nil {
		return nil, err
	}

	var result []*Device
	for _, device := range devices {
		if isAndroidDevice(device) {
			result = append(
				result,
				mapLibUsbDevicesToInternalModel(device),
			)
		}
	}

	return result, nil
}

func getDevicesDescriptions(context *gousb.Context) ([]*gousb.DeviceDesc, error) {
	var devices []*gousb.DeviceDesc

	_, err := context.OpenDevices(func(description *gousb.DeviceDesc) bool {
		devices = append(devices, description)
		// avoid device opening
		return false
	})

	return devices, err
}

func isAndroidDevice(description *gousb.DeviceDesc) bool {
	for _, configuration := range description.Configs {
		for _, usbInterface := range configuration.Interfaces {
			for _, interfaceConfiguration := range usbInterface.AltSettings {
				if interfaceConfiguration.Class == usbInterfaceAdbClass &&
					interfaceConfiguration.SubClass == usbInterfaceAdbSubclass &&
					interfaceConfiguration.Protocol == usbInterfaceAdbProtocol {

					return true
				}
			}
		}
	}

	return false
}

func mapLibUsbDevicesToInternalModel(description *gousb.DeviceDesc) *Device {
	return &Device{
		Description: usbid.Describe(description),

		Bus:     description.Bus,
		Address: description.Address,

		Vendor:  int(description.Vendor),
		Product: int(description.Product),
	}
}
