package driver

import (
	"fmt"
	"github.com/kubeedge/mappers-go/mappers/common"
	"strings"
	"time"
)

// BaseMessage the base structure of mqttadapter message.
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

// DeviceTwinUpdate the structure of device twin update.
type DeviceTwinUpdate struct {
	BaseMessage
	Twin map[string]*MsgTwin `json:"twin"`
}

// TwinData the structure of device twin
type TwinData struct {
	Name       string
	Type       string
	Topic      string
	Value      string
	MqttClient mqttclient.MqttClient
	driverUnit DriverUnit
}

// CreateMessageTwinUpdate create twin update message.
func CreateMessageTwinUpdate(name string, valueType string, value string) (msg []byte, err error) {
	var updateMsg DeviceTwinUpdate

	updateMsg.BaseMessage.Timestamp = time.Now().UnixNano() / 1e6
	updateMsg.Twin = map[string]*MsgTwin{}
	updateMsg.Twin[name] = &MsgTwin{}
	updateMsg.Twin[name].Actual = &TwinValue{Value: &value}
	updateMsg.Twin[name].Metadata = &TypeMetadata{Type: valueType}

	msg, err = json.Marshal(updateMsg)
	return
}

// Run start timer function to get device's twin or data, and send it to mqtt broker
func (td *TwinData) sendDeviceTwinMsg() {
	var err error
	sData, err := controller.GetDeviceData(td.driverUnit.instanceID, td.driverUnit.twin, td.driverUnit.drivers, td.driverUnit.mutex, td.driverUnit.dic)
	if err != nil {
		klog.Errorf("Get %s data error:", td.driverUnit.instanceID, err.Error())
		return
	}
	td.Value = sData
	var payload []byte
	if strings.Contains(td.Topic, "$hw") {
		if payload, err = CreateMessageTwinUpdate(td.Name, td.Type, td.Value); err != nil {
			klog.Errorf("Create %s message twin update failed: %v", td.driverUnit.instanceID, err)
			return
		}
	} else {
		if payload, err = CreateMessageData(td.Name, td.Type, td.Value); err != nil {
			klog.Errorf("Create %s message data failed: %v", td.driverUnit.instanceID, err)
			return
		}
	}
	if err := td.MqttClient.Publish(td.Topic, payload); err != nil {
		klog.Errorf("Publish topic %v failed, err: %v", td.Topic, err)
	}
}
func (gigEClient *GigEVisionDevice) updateDeviceTwin(DeviceSN string) {

	twinData := TwinData{
		Name:       twinV.PropertyName,
		Type:       twinV.Desired.Metadatas.Type,
		Topic:      fmt.Sprintf(common.TopicTwinUpdate, id),
		MqttClient: mqttClient,
		driverUnit: DriverUnit{
			instanceID: id,
			twin:       twinV,
			drivers:    drivers,
			mutex:      mutex,
			dic:        dic,
		},
	}
	timer := common.Timer{Function: twinData.Run, Duration: collectCycle, Times: 0}
}
