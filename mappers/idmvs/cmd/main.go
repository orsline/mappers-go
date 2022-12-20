package main

import (
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/service"
	"github.com/kubeedge/mappers-go/mappers/idmvs-sdk/driver"
)

// main IDMVS device program entry
func main() {
	d := &driver.IDMVS{}
	service.Bootstrap("IDMVS", d)
}
