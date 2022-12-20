package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/kubeedge/mappers-go/mappers/mqtt/driver"
)

// main mqtt device program entry
func main() {
	mqtt := &driver.MQTT{}
	service.Bootstrap("MQTT", mqtt)
}
