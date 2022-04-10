package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/kubeedge/mappers-go/mappers/modbus-sdk/driver"
)

// main Virtual device program entry
func main() {
	modbus := &driver.ModbusClient{}
	service.Bootstrap("modbus", modbus)
}
