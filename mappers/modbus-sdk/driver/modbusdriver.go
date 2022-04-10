package driver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

func (mbClient *ModbusClient) ReadDeviceData() (data interface{}, err error) {
	data, err = mbClient.Get(mbClient.visitorConfig.Register, mbClient.visitorConfig.Offset, uint16(mbClient.visitorConfig.Limit))
	if err != nil {
		fmt.Println("ModbusTcp Get error", err.Error())
		return nil, err
	}
	dataValue := data.([]byte)
	data = float32(binary.BigEndian.Uint16(dataValue)) / 10
	fmt.Printf("ReadDeviceData from Device %v\n", data)
	return data, err
}

func (mbClient *ModbusClient) WriteDeviceData(data interface{}) (err error) {
	err = mbClient.Set(mbClient.visitorConfig.Register, mbClient.visitorConfig.Offset, uint16(data.(int64)))
	if err != nil {
		fmt.Println("Set Err", err.Error())
		return err
	}
	return nil
}

func (mbClient *ModbusClient) StopDevice() (err error) {
	fmt.Println("---------Stop Modbus Successful---------")
	return nil
}

func (mbClient *ModbusClient) InitDevice(protocolCommon, visitor, protocol []byte) (err error) {
	if protocolCommon != nil {
		if err = json.Unmarshal(protocolCommon, &mbClient.protocolCommonConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolCommonConfig error: %v\n", err)
			return err
		}
	}
	if visitor != nil {
		if err = json.Unmarshal(visitor, &mbClient.visitorConfig); err != nil {
			fmt.Printf("Unmarshal visitorConfig error: %v\n", err)
			return err
		}
	}

	if protocol != nil {
		if err = json.Unmarshal(protocol, &mbClient.modbusProtocolConfig); err != nil {
			fmt.Printf("Unmarshal ProtocolConfig error: %v\n", err)
			return err
		}
	}
	err = mbClient.NewClient()
	if err != nil {
		fmt.Printf("NewClient error: %v\n", err)
		return err
	}
	return nil
}

func (mbClient *ModbusClient) GetDeviceStatus() (status bool) {
	return true
}