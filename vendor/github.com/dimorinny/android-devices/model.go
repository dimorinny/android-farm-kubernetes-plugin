package android

type Device struct {
	Description string

	Bus     int
	Address int

	Vendor  int
	Product int
}
