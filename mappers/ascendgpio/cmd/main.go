package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/kubeedge/mappers-go/mappers/ascend-gpio/driver"
)

// main gpio device program entry
func main() {
	gpio := &driver.GPIO{}
	service.Bootstrap("GPIO", gpio)
}
