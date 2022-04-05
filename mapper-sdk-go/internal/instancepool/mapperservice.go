package instancepool

import (
	"context"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/common"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/configmap"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/di"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/pkg/models"
	"sync"
)

// DeviceInstancesName contains the name of device service struct in the DIC.
var DeviceInstancesName = di.TypeInstanceToName(map[string]configmap.DeviceInstance{})

var DeviceModelsName = di.TypeInstanceToName(map[string]configmap.DeviceModel{})

var ProtocolName = di.TypeInstanceToName(map[string]configmap.Protocol{})

var ProtocolDriverName = di.TypeInstanceToName((*models.ProtocolDriver)(nil))

var WgName = di.TypeInstanceToName((*sync.WaitGroup)(nil))

var MutexName = di.TypeInstanceToName((*sync.Mutex)(nil))

var StopFunctionsName = di.TypeInstanceToName(map[string]context.CancelFunc(nil))

var ConnectInfoName = di.TypeInstanceToName(map[string]configmap.ConnectInfo{})

var DeviceLockName = di.TypeInstanceToName(map[string]common.Lock{})

// DeviceInstancesNameFrom helper function queries the DIC and returns device service struct.
func DeviceInstancesNameFrom(get di.Get) map[string]*configmap.DeviceInstance {
	return get(DeviceInstancesName).(map[string]*configmap.DeviceInstance)
}

func DeviceModelsNameFrom(get di.Get) map[string]*configmap.DeviceModel {
	return get(DeviceModelsName).(map[string]*configmap.DeviceModel)
}

func ProtocolNameFrom(get di.Get) map[string]*configmap.Protocol {
	return get(ProtocolName).(map[string]*configmap.Protocol)
}

func ProtocolDriverNameFrom(get di.Get) models.ProtocolDriver {
	return get(ProtocolDriverName).(models.ProtocolDriver)
}

func WgNameFrom(get di.Get) *sync.WaitGroup {
	return get(WgName).(*sync.WaitGroup)
}

func MutexNameFrom(get di.Get) *sync.Mutex {
	return get(MutexName).(*sync.Mutex)
}

func StopFunctionsNameFrom(get di.Get) map[string]context.CancelFunc {
	return get(StopFunctionsName).(map[string]context.CancelFunc)
}

// ConnectInfoNameFrom helper function queries the DIC and returns device service struct.
func ConnectInfoNameFrom(get di.Get) map[string]*configmap.ConnectInfo {
	return get(ConnectInfoName).(map[string]*configmap.ConnectInfo)
}

func DeviceLockNameFrom(get di.Get) map[string]*common.Lock {
	return get(DeviceLockName).(map[string]*common.Lock)
}
