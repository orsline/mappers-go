package driver

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"strconv"
	"sync"
)

var SensorData GatewayData
type MqttProtocolConfig struct {
	ProtocolName       string `json:"protocolName"`
	ProtocolConfigData `json:"configData"`
}

type ProtocolConfigData struct {

}

type MqttProtocolCommonConfig struct {
	CommonCustomizedValues `json:"customizedValues"`
}

type CommonCustomizedValues struct {
	IP string `json:"IP"`
	Port int `json:"port"`
	Topic string `json:"topic"`
}
type MqttVisitorConfig struct {
	ProtocolName      string `json:"protocolName"`
	VisitorConfigData `json:"configData"`
}

type VisitorConfigData struct {
	Feature string `json:"feature"`
}

type GatewayData struct {
	JsonData map[string]interface{}
}


// MQTT Realize the structure of random number
type MQTT struct {
	mutex                 sync.Mutex
	protocolConfig MqttProtocolConfig
	protocolCommonConfig  MqttProtocolCommonConfig
	visitorConfig         MqttVisitorConfig
	Client MqttClient
}

// InitDevice Sth that need to do in the first
// If you need mount a persistent connection, you should provide parameters in configmap's protocolCommon.
// and handle these parameters in the following function
func (d *MQTT) InitDevice(protocolCommon []byte) (err error) {
	err = json.Unmarshal(protocolCommon, &d.protocolCommonConfig)
	if err != nil {
		return err
	}
	d.newClient()
	return nil
}

// SetConfig Parse the configmap's raw json message
func (d *MQTT) SetConfig(protocolCommon, visitor, protocol []byte) (server string, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &d.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return  "", err
		}
	}
	if visitor != nil {
		if err = json.Unmarshal(visitor, &d.visitorConfig); err != nil {
			fmt.Printf("Unmarshal visitorConfig error: %v\n", err)
			return "", err
		}

	}
	if protocol != nil {
		if err = json.Unmarshal(protocol, &d.protocolConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolConfig error: %v\n", err)
			return "", err
		}
	}
	server = d.protocolCommonConfig.IP + ":" +strconv.Itoa(d.protocolCommonConfig.Port)
	return  server,nil

}

// ReadDeviceData  is an interface that reads data from a specific device, data is a type of string
func (d *MQTT) ReadDeviceData(protocolCommon, visitor, protocol []byte) (data interface{}, err error) {
	// Parse raw json message to get a virtualDevice instance
	_, err = d.SetConfig(protocolCommon, visitor, protocol)
	if err != nil {
		return nil, err
	}
	return SensorData.JsonData[d.visitorConfig.Feature],nil
}

// WriteDeviceData is an interface that write data to a specific device, data's DataType is Consistent with configmap
func (d *MQTT) WriteDeviceData(data interface{}, protocolCommon, visitor, protocol []byte) (err error) {
	return nil
}

// StopDevice is an interface to disconnect a specific device
// This function is called when mapper stops serving
func (d *MQTT) StopDevice() (err error) {
	// in this func, u can get ur device-instance in the client map, and give a safety exit
	fmt.Println("----------Stop Virtual Device Successful----------")
	return nil
}


// GetDeviceStatus is an interface to get the device status true is OK , false is DISCONNECTED
func (d *MQTT) GetDeviceStatus(protocolCommon, visitor, protocol []byte) (status bool) {
	err := d.Client.Publish("pulse", "")
	if err != nil {
		return false
	}
	return true
}


func (d *MQTT )newClient(){
	server := d.protocolCommonConfig.IP + ":" +strconv.Itoa(d.protocolCommonConfig.Port)
	d.Client = MqttClient{
		IP:         server,
		User:       "huawei",
		Passwd:     "",
		Cert:       "",
		PrivateKey: "",
	}
	err := d.Client.Connect()
	if err != nil {
		fmt.Printf("Failed to init mqtt client ,error:%v\n",err)
		os.Exit(1)
	}
	fmt.Println(server ," connect Successful")
	err = d.Client.Subscribe(d.protocolCommonConfig.Topic,onMessage)
	if err != nil {
		fmt.Printf("Failed to init mqtt client ,error:%v\n",err)
		os.Exit(1)
	}
	fmt.Println("Subscribe Successful")
}

// onMessage callback function of Mqtt subscribe message.
func onMessage(client mqtt.Client, message mqtt.Message) {
	err := json.Unmarshal(message.Payload(), &SensorData.JsonData)
	if err != nil {
		fmt.Println("json.Unmarshal error:", err.Error())
		return
	}
}