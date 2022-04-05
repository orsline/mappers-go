package driver

/*
#include <dlfcn.h>
int open_device(unsigned int** device,char* deviceId,char** error);
int get_value (unsigned int* device, char* feature, char** value,char** error);
int set_value (unsigned int* device, char* feature, char* value,char** error);
void close_device (unsigned int* device);
//链接dl库
#cgo LDFLAGS: -ldl
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	"sync"
)

type GigEVisionDeviceProtocolConfig struct {
}

type GigEVisionDeviceProtocolCommonConfig struct {
	DeviceSN string `json:"deviceSN"`
}

type GigEVisionDeviceVisitorConfig struct {
	FeatureName string `json:"FeatureName"`
}

type GigEVisionDevice struct {
	mutex                sync.RWMutex
	GigEprotocolConfig   GigEVisionDeviceProtocolConfig
	ProtocolCommonConfig GigEVisionDeviceProtocolCommonConfig
	visitorConfig        GigEVisionDeviceVisitorConfig
	dev                  *C.uint
}

func (geClient *GigEVisionDevice) InitDevice(protocolCommon []byte) (err error) {
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &geClient.ProtocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return err
		}
	}
	err = geClient.NewClient()
	if err != nil {
		fmt.Printf("Failed to new a GigE client: %v\n", err)
		return err
	}
	return nil
}

func (geClient *GigEVisionDevice) SetConfig(protocolCommon, visitor, protocol []byte) (err error) {
	//geClient.NewClient()
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &geClient.ProtocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return err
		}
		//klog.V(0).Info(geClient.GigEprotocolConfig)
	}
	if visitor != nil {
		if err = json.Unmarshal(visitor, &geClient.visitorConfig); err != nil {
			fmt.Printf("Unmarshal visitorConfig error: %v\n", err)
			return err
		}

	}

	if protocol != nil {
		if err = json.Unmarshal(protocol, &geClient.GigEprotocolConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolConfig error: %v\n", err)
			return err
		}
		//klog.V(0).Info(geClient.ProtocolCommonConfig)
	}

	return nil
}

// ReadDeviceData  is an interface that reads data from a specific device, data is a type of string
func (geClient *GigEVisionDevice) ReadDeviceData(protocolCommon, visitor, protocol []byte) (data interface{}, err error) {
	geClient.NewClient()
	err = geClient.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return nil, err
	}
	data, err = geClient.Get(geClient.visitorConfig.FeatureName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return data, err
}

// WriteDeviceData is an interface that write data to a specific device, data's DataType is Consistent with configmap
func (geClient *GigEVisionDevice) WriteDeviceData(data interface{}, protocolCommon, visitor, protocol []byte) (err error) {
	geClient.NewClient()
	err = geClient.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return err
	}
	//geClient.mutex.Lock()
	//defer geClient.mutex.Unlock()
	err = geClient.Set(geClient.visitorConfig.FeatureName, data)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// StopDevice is an interface to disconnect a specific device
func (geClient *GigEVisionDevice) StopDevice() (err error) {
	geClient.mutex.Lock()
	defer geClient.mutex.Unlock()
	C.close_device(geClient.dev)
	geClient.dev = nil
	fmt.Println("----------Stop GigE Device Successful----------")
	return nil
}

// GetDeviceStatus is an interface to get the device status true is OK , false is DISCONNECTED
func (geClient *GigEVisionDevice) GetDeviceStatus(protocolCommon, visitor, protocol []byte) (status bool) {
	geClient.mutex.Lock()
	defer geClient.mutex.Unlock()
	var msg *C.char
	var value *C.char
	err := geClient.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return false
	}
	signal := C.get_value(geClient.dev, C.CString(geClient.visitorConfig.FeatureName), &value, &msg)
	if signal != 0 {
		klog.Errorf("Device %s unconnected", geClient.ProtocolCommonConfig.DeviceSN)
		go geClient.ReConnectDevice()
		return false
	}
	signal = C.set_value(geClient.dev, C.CString(geClient.visitorConfig.FeatureName), value, &msg)
	if signal != 0 {
		klog.Errorf("Device %s unconnected", geClient.ProtocolCommonConfig.DeviceSN)
		go geClient.ReConnectDevice()
		return false
	}
	klog.Infof("Device %s is connected", geClient.ProtocolCommonConfig.DeviceSN)
	return true
}

func (geClient *GigEVisionDevice) ReConnectDevice() {
	var dev *C.uint
	var msg *C.char
	if geClient.dev != nil {
		C.close_device(geClient.dev)
		geClient.dev = nil
	}
	for {
		signal := C.open_device(&dev, C.CString(geClient.ProtocolCommonConfig.DeviceSN), &msg)
		if signal != 0 {
			klog.Errorf("Failed to restart device %s : %s Please check your camera link", geClient.ProtocolCommonConfig.DeviceSN, C.GoString(msg))
		} else {
			geClient.dev = dev
			break
		}
	}
	klog.Infof("Device %s restsrt success!", geClient.ProtocolCommonConfig.DeviceSN)
	return
}
