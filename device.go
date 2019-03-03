package main

import "time"

type Device struct {
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

		l.devices <- devices

		time.Sleep(time.Second)
	}
}

func (l *UsbAndroidDevicesListener) Devices() chan []*Device {
	return l.devices
}

func (l *UsbAndroidDevicesListener) Errors() chan error {
	return l.errors
}

func (l *UsbAndroidDevicesListener) getDevices() ([]*Device, error) {
	return []*Device{
		{
			devicePath: "device1",
		},
		{
			devicePath: "device2",
		},
	}, nil
}

func (l *UsbAndroidDevicesListener) abort(err error) {
	l.errors <- err

	close(l.errors)
	close(l.devices)
}
