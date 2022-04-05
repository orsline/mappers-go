package mqttadapter

import (
	"context"
	"encoding/json"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/configmap"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/controller"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/instancepool"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/di"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
	"regexp"
)

// SyncInfo callback function of Mqtt subscribe message.
// The function will update device's value according to the message sent from the cloud
func SyncInfo(dic *di.Container, message mqtt.Message) {
	re := regexp.MustCompile(`hw/events/device/(.+)/twin/update/delta`)
	instanceID := re.FindStringSubmatch(message.Topic())[1]
	deviceInstances := instancepool.DeviceInstancesNameFrom(dic.Get)
	driver := instancepool.ProtocolDriverNameFrom(dic.Get)
	mapMutex := instancepool.DeviceLockNameFrom(dic.Get)
	if _, ok := deviceInstances[instanceID]; !ok {
		klog.Errorf("Instance :%s does not exist", instanceID)
		return
	} else {
		var delta DeviceTwinDelta
		if err := json.Unmarshal(message.Payload(), &delta); err != nil {
			klog.Errorf("Unmarshal message failed: %v", err)
			return
		}
		for twinName, twinValue := range delta.Delta {
			i := 0
			for i = 0; i < len(deviceInstances[instanceID].Twins); i++ {
				if twinName == deviceInstances[instanceID].Twins[i].PropertyName {
					break
				}
			}
			if i == len(deviceInstances[instanceID].Twins) {
				continue
			}
			// Desired value is not changed.
			if deviceInstances[instanceID].Twins[i].Desired.Value == twinValue {
				continue
			}
			klog.V(4).Infof("Set %s:%s value to %s", instanceID, twinName, twinValue)
			deviceInstances[instanceID].Twins[i].Desired.Value = twinValue
			err := controller.SetVisitor(instanceID, deviceInstances[instanceID].Twins[i], driver, mapMutex[instanceID], dic)
			if err != nil {
				klog.Error(err)
				return
			}
		}
	}

}

// UpdateDevice callback function of Mqtt subscribe message.
// The function support for dynamically adding/removing devices
func UpdateDevice(dic *di.Container, message mqtt.Message) {
	choices := make(map[string]string)
	if err := json.Unmarshal(message.Payload(), &choices); err != nil {
		klog.Errorf("Unmarshal message failed: %v", err)
		return
	}
	if choices["option"] == "add" {
		addDevice(dic, message)
	} else {
		removeDevice(dic, message)
	}

}

// removeDevice support for dynamically removing devices, delete only local memory data
func removeDevice(dic *di.Container, message mqtt.Message) {
	re := regexp.MustCompile(`hw/events/node/(.+)/membership/updated`)
	stopFunctions := instancepool.StopFunctionsNameFrom(dic.Get)
	deviceInstances := instancepool.DeviceInstancesNameFrom(dic.Get)
	deviceModels := instancepool.DeviceModelsNameFrom(dic.Get)
	protocol := instancepool.ProtocolNameFrom(dic.Get)
	instanceID := re.FindStringSubmatch(message.Topic())[1]
	mutex := instancepool.MutexNameFrom(dic.Get)
	mutex.Lock()
	defer mutex.Unlock()
	if cancelFunc, ok := stopFunctions[instanceID]; ok {
		cancelFunc()
		modelNameDeleted := deviceInstances[instanceID].Model
		protocolNameDeleted := deviceInstances[instanceID].ProtocolName
		if _, ok := deviceInstances[instanceID]; ok {
			delete(deviceInstances, instanceID)
		}
		modelFlag := true
		protocolFlag := true
		for k, v := range deviceInstances {
			if k != instanceID {
				if v.Model == modelNameDeleted {
					modelFlag = false
					break
				}
			}
		}
		if modelFlag {
			delete(deviceModels, modelNameDeleted)
		}
		for k, v := range deviceInstances {
			if k != instanceID {
				if v.ProtocolName == protocolNameDeleted {
					protocolFlag = false
					break
				}
			}
		}
		if protocolFlag {
			delete(protocol, protocolNameDeleted)
		}
		delete(stopFunctions, instanceID)
		klog.V(1).Infof("Remove %s successful\n", instanceID)
	} else {
		klog.V(1).Infof("Remove %s failed,there is no such instanceId\n", instanceID)
	}
}

// addDevice support for dynamically adding devices , delete only local memory data
func addDevice(dic *di.Container, message mqtt.Message) {
	re := regexp.MustCompile(`hw/events/node/(.+)/membership/updated`)
	instanceID := re.FindStringSubmatch(message.Topic())[1]
	configMap := instancepool.ConfigMapNameFrom(dic.Get)
	deviceInstances := instancepool.DeviceInstancesNameFrom(dic.Get)
	deviceModels := instancepool.DeviceModelsNameFrom(dic.Get)
	protocol := instancepool.ProtocolNameFrom(dic.Get)
	connectInfo := instancepool.ConnectInfoNameFrom(dic.Get)
	driver := instancepool.ProtocolDriverNameFrom(dic.Get)
	mqttClient := instancepool.MqttClientNameFrom(dic.Get)
	wg := instancepool.WgNameFrom(dic.Get)
	mapMutex := instancepool.DeviceLockNameFrom(dic.Get)
	stopFunctions := instancepool.StopFunctionsNameFrom(dic.Get)
	defaultConfigFile := configMap
	mutex := instancepool.MutexNameFrom(dic.Get)
	mutex.Lock()
	// parseConfigmap
	if err := configmap.ParseOdd(defaultConfigFile, deviceInstances, deviceModels, protocol, instanceID); err != nil {
		klog.Errorf("Please check you config-file %s,", err.Error())
		return
	}
	mutex.Unlock()
	configmap.GetConnectInfo(deviceInstances, connectInfo)
	go func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		err := SendTwin(instanceID, deviceInstances[instanceID], driver, mqttClient, wg, dic, mapMutex[instanceID], ctx)
		if err != nil {
			klog.Errorf("Failed to get %s %s:%v\n", instanceID, "twin", err)
		} else {
			err = SendData(instanceID, deviceInstances[instanceID], driver, mqttClient, wg, dic, mapMutex[instanceID], ctx)
			if err != nil {
				klog.Errorf("Failed to get %s %s:%v\n", instanceID, "data", err)
			}
			err = SendDeviceState(instanceID, deviceInstances[instanceID], driver, mqttClient, wg, dic, mapMutex[instanceID], ctx)
			if err != nil {
				klog.Errorf("Failed to get %s %s:%v\n", instanceID, "state", err)
			}
		}

		stopFunctions[instanceID] = cancelFunc
		klog.V(1).Infof("Add %s successful\n", instanceID)
	}()
}
