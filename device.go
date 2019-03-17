package main

import (
	"fmt"
	"github.com/dimorinny/android-devices"
	"log"
	"os"
	"time"
)

type Device struct {
	deviceName string
	devicePath string
}

type DevicesListener interface {
	Listen()
	Devices() chan []*Device
	Errors() chan error
}

type UsbAndroidDevicesListener struct {
	devices chan []*Device
	errors  chan error
}

func NewUsbAndroidDevicesListener() DevicesListener {
	return &UsbAndroidDevicesListener{}
}

func (l *UsbAndroidDevicesListener) Listen() {
	l.devices = make(chan []*Device)
	l.errors = make(chan error)

	for {
		devices, err := l.getDevices()
		if err != nil {
			l.abort(err)
			break
		}

		log.Println(
			fmt.Sprintf("Found %d devices", len(devices)),
		)

		l.devices <- devices

		time.Sleep(time.Second * 5)
	}
}

func (l *UsbAndroidDevicesListener) Devices() chan []*Device {
	return l.devices
}

func (l *UsbAndroidDevicesListener) Errors() chan error {
	return l.errors
}

func (l *UsbAndroidDevicesListener) getDevices() ([]*Device, error) {
	devices, err := android.Devices()
	if err != nil {
		return nil, err
	}

	var resultDevices []*Device
	for _, device := range devices {
		linuxDevicePath := fmt.Sprintf("/dev/bus/usb/%03d/%03d", device.Bus, device.Address)

		if _, err := os.Stat(linuxDevicePath); os.IsNotExist(err) {
			log.Println(
				fmt.Sprintf(
					"Linux device file: %s not found (for device: %s)",
					linuxDevicePath,
					device.Description,
				),
			)

			continue
		}

		resultDevices = append(
			resultDevices,
			&Device{
				deviceName: device.Description,
				devicePath: linuxDevicePath,
			},
		)
	}

	return resultDevices, nil
}

func (l *UsbAndroidDevicesListener) abort(err error) {
	l.errors <- err

	close(l.errors)
	close(l.devices)
}
