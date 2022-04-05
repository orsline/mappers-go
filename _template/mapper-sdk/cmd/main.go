package main

import (
	"gitee.com/ascend/mapper-go-sdk/mapper-sdk-go/pkg/service"
	"gitee.com/ascend/mapper-go-sdk/mappers/Template/driver"
)

// main Template device program entry
func main() {
	d := &driver.Template{}
	service.Bootstrap("Template", d)
}
