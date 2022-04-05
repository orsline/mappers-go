package configmap

import (
	"encoding/json"
	"errors"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/common"
	"io/ioutil"
	"k8s.io/klog/v2"
)

// Parse is a method to parse the configmap.
func Parse(path string,
	devices map[string]*DeviceInstance,
	dms map[string]*DeviceModel,
	protocols map[string]*Protocol) error {
	var deviceProfile DeviceProfile
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New("failed to read " + path + " file")
		return err
	}
	//Parse the JSON file and convert it into the data structure of DeviceProfile
	if err = json.Unmarshal(jsonFile, &deviceProfile); err != nil {
		return err
	}
	for i := 0; i < len(deviceProfile.DeviceInstances); i++ {
		instance := deviceProfile.DeviceInstances[i]
		j := 0
		for j = 0; j < len(deviceProfile.Protocols); j++ {
			if instance.ProtocolName == deviceProfile.Protocols[j].Name {
				instance.PProtocol = deviceProfile.Protocols[j]
				break
			}
		}
		if j == len(deviceProfile.Protocols) {
			err = errors.New("protocol mismatch")
			return err
		}
		for k := 0; k < len(instance.PropertyVisitors); k++ {
			modelName := instance.PropertyVisitors[k].ModelName
			propertyName := instance.PropertyVisitors[k].PropertyName
			l := 0
			for l = 0; l < len(deviceProfile.DeviceModels); l++ {
				if modelName == deviceProfile.DeviceModels[l].Name {
					m := 0
					for m = 0; m < len(deviceProfile.DeviceModels[l].Properties); m++ {
						if propertyName == deviceProfile.DeviceModels[l].Properties[m].Name {
							instance.PropertyVisitors[k].PProperty = deviceProfile.DeviceModels[l].Properties[m]
							break
						}
					}
					if m == len(deviceProfile.DeviceModels[l].Properties) {
						err = errors.New("property mismatch")
						return err
					}
					break
				}
			}
			if l == len(deviceProfile.DeviceModels) {
				err = errors.New("device model mismatch")
				return err
			}
		}
		for k := 0; k < len(instance.Twins); k++ {
			name := instance.Twins[k].PropertyName
			l := 0
			for l = 0; l < len(instance.PropertyVisitors); l++ {
				if name == instance.PropertyVisitors[l].PropertyName {
					instance.Twins[k].PVisitor = &instance.PropertyVisitors[l]
					break
				}
			}
			if l == len(instance.PropertyVisitors) {
				err = errors.New("propertyVisitor mismatch")
				return err
			}
		}
		for k := 0; k < len(instance.Datas.Properties); k++ {
			name := instance.Datas.Properties[k].PropertyName
			l := 0
			for l = 0; l < len(instance.PropertyVisitors); l++ {
				if name == instance.PropertyVisitors[l].PropertyName {
					instance.Datas.Properties[k].PVisitor = &instance.PropertyVisitors[l]
					break
				}
			}
			if l == len(instance.PropertyVisitors) {
				err = errors.New("propertyVisitor mismatch")
				return err
			}
		}
		devices[instance.ID] = new(DeviceInstance)
		devices[instance.ID] = &instance
		klog.V(4).Infof("Instance:%s Successfully registered", instance.ID)
	}
	for i := 0; i < len(deviceProfile.DeviceModels); i++ {
		dms[deviceProfile.DeviceModels[i].Name] = new(DeviceModel)
		dms[deviceProfile.DeviceModels[i].Name] = &deviceProfile.DeviceModels[i]
	}
	for i := 0; i < len(deviceProfile.Protocols); i++ {
		protocols[deviceProfile.Protocols[i].Name] = new(Protocol)
		protocols[deviceProfile.Protocols[i].Name] = &deviceProfile.Protocols[i]
	}
	return nil
}

// GetConnectInfo is a method to generate link information for each attribute
func GetConnectInfo(
	devices map[string]*DeviceInstance,
	connectInfo map[string]*ConnectInfo) {
	for id, instance := range devices {
		tempId := id
		tempInstance := instance
		for _, visitorV := range tempInstance.PropertyVisitors {
			tempVisitorV := visitorV
			driverName := common.DriverPrefix + tempId + visitorV.PropertyName
			connectInfo[driverName] = &ConnectInfo{
				ProtocolCommonConfig: tempInstance.PProtocol.ProtocolCommonConfig,
				VisitorConfig:        tempVisitorV.VisitorConfig,
				ProtocolConfig:       tempInstance.PProtocol.ProtocolConfigs,
			}
		}
	}
}

// ParseOdd is a method to parse the configmap.
func ParseOdd(path string,
	devices map[string]*DeviceInstance,
	dms map[string]*DeviceModel,
	protocols map[string]*Protocol,
	id string) error {
	var deviceProfile DeviceProfile
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		err = errors.New("failed to read " + path + " file")
		return err
	}
	//Parse the JSON file and convert it into the data structure of DeviceProfile
	if err = json.Unmarshal(jsonFile, &deviceProfile); err != nil {
		return err
	}
	for i := 0; i < len(deviceProfile.DeviceInstances); i++ {
		instance := deviceProfile.DeviceInstances[i]
		if instance.ID == id {
			j := 0
			for j = 0; j < len(deviceProfile.Protocols); j++ {
				if instance.ProtocolName == deviceProfile.Protocols[j].Name {
					instance.PProtocol = deviceProfile.Protocols[j]
					break
				}
			}
			if j == len(deviceProfile.Protocols) {
				err = errors.New("protocol not found")
				return err
			}
			for k := 0; k < len(instance.PropertyVisitors); k++ {
				modelName := instance.PropertyVisitors[k].ModelName
				propertyName := instance.PropertyVisitors[k].PropertyName
				l := 0
				for l = 0; l < len(deviceProfile.DeviceModels); l++ {
					if modelName == deviceProfile.DeviceModels[l].Name {
						m := 0
						for m = 0; m < len(deviceProfile.DeviceModels[l].Properties); m++ {
							if propertyName == deviceProfile.DeviceModels[l].Properties[m].Name {
								instance.PropertyVisitors[k].PProperty = deviceProfile.DeviceModels[l].Properties[m]
								break
							}
						}
						if m == len(deviceProfile.DeviceModels[l].Properties) {
							err = errors.New("property not found")
							return err
						}
						break
					}
				}
				if l == len(deviceProfile.DeviceModels) {
					err = errors.New("device model not found")
					return err
				}
			}
			for k := 0; k < len(instance.Twins); k++ {
				name := instance.Twins[k].PropertyName
				l := 0
				for l = 0; l < len(instance.PropertyVisitors); l++ {
					if name == instance.PropertyVisitors[l].PropertyName {
						instance.Twins[k].PVisitor = &instance.PropertyVisitors[l]
						break
					}
				}
				if l == len(instance.PropertyVisitors) {
					return errors.New("propertyVisitor not found")
				}
			}
			for k := 0; k < len(instance.Datas.Properties); k++ {
				name := instance.Datas.Properties[k].PropertyName
				l := 0
				for l = 0; l < len(instance.PropertyVisitors); l++ {
					if name == instance.PropertyVisitors[l].PropertyName {
						instance.Datas.Properties[k].PVisitor = &instance.PropertyVisitors[l]
						break
					}
				}
				if l == len(instance.PropertyVisitors) {
					return errors.New("propertyVisitor mismatch")
				}
			}
			if _, ok := devices[instance.ID]; !ok {
				devices[instance.ID] = new(DeviceInstance)
				devices[instance.ID] = &instance
			} else {
				return errors.New(instance.ID + " already in the device list")
			}
			for i := 0; i < len(deviceProfile.DeviceModels); i++ {
				if _, ok := dms[deviceProfile.DeviceModels[i].Name]; !ok {
					dms[deviceProfile.DeviceModels[i].Name] = &deviceProfile.DeviceModels[i]
				}
			}
			for i := 0; i < len(deviceProfile.Protocols); i++ {
				protocols[deviceProfile.Protocols[i].Name] = &deviceProfile.Protocols[i]
			}
			return nil
		}
	}

	return errors.New("can't find the device in profile.json")
}
