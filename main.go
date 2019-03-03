package main

func main() {
	devicesListener := NewUsbAndroidDevicesListener()
	devicePlugin := NewAndroidDevicesPlugin(devicesListener)

	devicePlugin.Start()
}
