package driver

/*
#include <dlfcn.h>
#include <stdlib.h>
int open_device(unsigned int** device,char* deviceId,char** error);
int get_value (unsigned int* device, char* feature, char** value,char** error);
int set_value (unsigned int* device, char* feature, char* value,char** error);
int close_device (unsigned int* device);
//链接dl库
#cgo LDFLAGS: -ldl
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type GigEVisionDeviceProtocolCommonConfig struct {
	CommonCustomizedValues `json:"customizedValues"`
}

type CommonCustomizedValues struct {
	DeviceSN string `json:"deviceSN"`
}

type GigEVisionDeviceVisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	FeatureName string `json:"FeatureName"`
}

type GigEVisionDevice struct {
	mutex                sync.RWMutex
	protocolCommonConfig GigEVisionDeviceProtocolCommonConfig
	visitorConfig        GigEVisionDeviceVisitorConfig
	deviceMeta           map[string]*DeviceMeta
}

type DeviceMeta struct {
	dev          *C.uint
	deviceStatus bool
	imageFormat  string
	imageURL     string
}

func (gigEClient *GigEVisionDevice) InitDevice(protocolCommon []byte) (err error) {
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &gigEClient.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return err
		}
	}
	err = gigEClient.NewClient()
	if err != nil {
		fmt.Printf("Failed to new a GigE client: %v\n", err)
		return err
	}
	return nil
}

func (gigEClient *GigEVisionDevice) SetConfig(protocolCommon, visitor []byte) (featureName string, deviceSN string, err error) {
	gigEClient.mutex.Lock()
	defer gigEClient.mutex.Unlock()
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &gigEClient.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return "", "", err
		}
	}
	if visitor != nil {
		if err = json.Unmarshal(visitor, &gigEClient.visitorConfig); err != nil {
			fmt.Printf("Unmarshal visitorConfig error: %v\n", err)
			return "", "", err
		}
	}
	return gigEClient.visitorConfig.FeatureName, gigEClient.protocolCommonConfig.DeviceSN, nil
}

// ReadDeviceData  is an interface that reads data from a specific device, data is a type of string
func (gigEClient *GigEVisionDevice) ReadDeviceData(protocolCommon, visitor, protocol []byte) (data interface{}, err error) {
	featureName, deviceSN, err := gigEClient.SetConfig(protocolCommon, visitor)
	if err != nil {
		return nil, err
	}
	if !gigEClient.deviceMeta[deviceSN].deviceStatus {
		errorMsg := fmt.Sprintf("Device %s is unreachable and failed to read device data.", deviceSN)
		err = errors.New(errorMsg)
		return nil, err
	}
	data, err = gigEClient.Get(deviceSN, featureName)
	if err != nil {
		return nil, err
	}
	return data, err
}

// WriteDeviceData is an interface that write data to a specific device, data's DataType is Consistent with configmap
func (gigEClient *GigEVisionDevice) WriteDeviceData(data interface{}, protocolCommon, visitor, protocol []byte) (err error) {
	featureName, deviceSN, err := gigEClient.SetConfig(protocolCommon, visitor)
	if err != nil {
		return err
	}
	if !gigEClient.deviceMeta[deviceSN].deviceStatus {
		errorMsg := fmt.Sprintf("Device %s is unreachable and failed to get.", deviceSN)
		err = errors.New(errorMsg)
		return err
	}
	err = gigEClient.Set(deviceSN, featureName, data)
	if err != nil {
		return err
	}
	return nil
}

// StopDevice is an interface to disconnect a specific device
func (gigEClient *GigEVisionDevice) StopDevice() (err error) {
	for s := range gigEClient.deviceMeta {
		if gigEClient.deviceMeta[s].deviceStatus {
			C.close_device(gigEClient.deviceMeta[s].dev)
			gigEClient.deviceMeta[s].dev = nil
			gigEClient.deviceMeta[s].deviceStatus = false
		}
	}
	fmt.Println("----------Stop GigE Device Successful----------")
	return nil
}

// GetDeviceStatus is an interface to get the device status true is OK , false is DISCONNECTED
func (gigEClient *GigEVisionDevice) GetDeviceStatus(protocolCommon, visitor, protocol []byte) (status bool) {
	_, deviceSN, err := gigEClient.SetConfig(protocolCommon, visitor)
	if err == nil {
		return false
	}
	return gigEClient.deviceMeta[deviceSN].deviceStatus
}

func (gigEClient *GigEVisionDevice) ReConnectDevice(DeviceSN string) {
	var msg *C.char
	var dev *C.uint
	gigEClient.deviceMeta[DeviceSN].dev = nil
	for {
		signal := C.open_device(&dev, C.CString(DeviceSN), &msg)
		if signal != 0 {
			fmt.Printf("Failed to restart device %s: %s\n", DeviceSN, (string)(C.GoString(msg)))
			time.Sleep(5 * time.Second)
		} else {
			gigEClient.deviceMeta[DeviceSN].dev = dev
			gigEClient.deviceMeta[DeviceSN].deviceStatus = true
			break
		}
	}
	fmt.Printf("Device %s restart success!\n", DeviceSN)
}
