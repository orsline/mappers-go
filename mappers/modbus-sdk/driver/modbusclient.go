package driver

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sailorvii/modbus"
)

var clients map[string]*modbus.Client

func (mbClient *ModbusClient) Set(registerType string, addr uint16, value uint16) (err error) {
	mbClient.mutex.Lock()
	defer mbClient.mutex.Unlock()
	fmt.Println("Set:", registerType, addr, value)
	switch registerType {
	case "CoilRegister":
		var valueSet uint16
		switch value {
		case 0:
			valueSet = 0x0000
		case 1:
			valueSet = 0xFF00
		default:
			return errors.New("Wrong value")
		}
		_, err = mbClient.client.WriteSingleCoil(addr, valueSet)
	case "HoldingRegister":
		_, err = mbClient.client.WriteSingleRegister(addr, value)
	default:
		return errors.New("Bad register type")
	}
	return nil
}

func (mbClient *ModbusClient) Get(registerType string, addr uint16, quantity uint16) (results []byte, err error) {
	switch registerType {
	case "CoilRegister":
		results, err = mbClient.client.ReadCoils(addr, quantity)
	case "DiscreteInputRegister":
		results, err = mbClient.client.ReadDiscreteInputs(addr, quantity)
	case "HoldingRegister":
		results, err = mbClient.client.ReadHoldingRegisters(addr, quantity)
	case "InputRegister":
		results, err = mbClient.client.ReadInputRegisters(addr, quantity)
	default:
		return nil, errors.New("Bad register type")
	}
	fmt.Println("Get result: ", results)
	return results, err
}

func (mbClient *ModbusClient) NewClient() (err error) {
	addr := mbClient.protocolCommonConfig.TCP.IP + ":" + strconv.Itoa(int(mbClient.protocolCommonConfig.TCP.Port))
	if _, ok := clients[addr]; ok {
		return nil
	}
	if clients == nil {
		clients = make(map[string]*modbus.Client)
	}
	handler := modbus.NewTCPClientHandler(addr)
	handler.Timeout = 1 * time.Second
	handler.IdleTimeout = 50
	handler.SlaveId = byte(mbClient.modbusProtocolConfig.SlaveID)
	mbClient.client = modbus.NewClient(handler)
	clients[addr] = &mbClient.client
	return nil
}
