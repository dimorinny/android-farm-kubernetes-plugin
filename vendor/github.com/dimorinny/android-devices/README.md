## Detect connected android devices

```go
package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dimorinny/android-devices"
	"log"
)

func main() {
	devices, err := android.Devices()
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(devices)
}

```

Example output:
```
([]*android.Device) (len=1 cap=1) {
 (*android.Device)(0xc0004a18f0)({
  Description: (string) (len=43) "Galaxy (MTP) (Samsung Electronics Co., Ltd)",
  Bus: (int) 1,
  Address: (int) 9,
  Vendor: (int) 1256,
  Product: (int) 26720
 })
}
```