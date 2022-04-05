package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/kubeedge/mappers-go/mappers/virtualdevice-sdk/driver"
)

// main Virtual device program entry
func main() {
	vd := &driver.VirtualDevice{}
	service.Bootstrap("RandomNumber", vd)
}
