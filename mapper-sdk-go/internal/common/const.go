// Package common used to store constants, data conversion functions, timers, etc
package common

// joint the topic like topic := fmt.Sprintf(TopicTwinUpdateDelta, deviceID)
const (
	TopicTwinUpdateDelta = "$hw/events/device/%s/twin/update/delta"
	TopicTwinUpdate      = "$hw/events/device/%s/twin/update"
	TopicStateUpdate     = "$hw/events/device/%s/state/update"
	TopicDataUpdate      = "$ke/events/device/%s/data/update"
	TopicDeviceUpdate    = "$hw/events/node/#"
)

// Device status definition.
const (
	DEVSTOK      = "OK"
	DEVSTDISCONN = "DISCONNECTED"
)

// joint x joint the instancepool like driverName :=  common.DriverPrefix+instanceID+twin.PropertyName
const (
	DriverPrefix = "Driver"
)

const (
	CorrelationHeader = "X-Correlation-ID"
)

const (
	ApiVersion = "v1"
	ApiBase    = "/api/v1"

	ApiDeviceRoute                 = ApiBase + "/device"
	ApiDeviceWriteCommandByIdRoute = ApiDeviceRoute + "/" + Id + "/{" + IdAndCommand + "}"
	ApiDeviceReadCommandByIdRoute  = ApiDeviceRoute + "/" + Id + "/{" + Id + "}" + "/{" + Command + "}"
	ApiDeviceCallbackRoute         = ApiBase + "/callback/device"
	ApiDeviceCallbackIdRoute       = ApiBase + "/callback/device/id/{id}"

	ApiPingRoute = ApiBase + "/ping"
)

const (
	Id           = "id"
	Command      = "command"
	IdAndCommand = "IdAndCommand"
)

// Constants related to the possible content types supported by the APIs
const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

type ErrKind string

// Constant Kind identifiers which can be used to label and group errors.
const (
	KindEntityDoesNotExist  ErrKind = "NotFound"
	KindServerError         ErrKind = "UnexpectedServerError"
	KindDuplicateName       ErrKind = "DuplicateName"
	KindInvalidId           ErrKind = "InvalidId"
	KindServiceUnavailable  ErrKind = "ServiceUnavailable"
	KindNotAllowed          ErrKind = "NotAllowed"
	KindServiceLocked       ErrKind = "ServiceLocked"
	KindNotImplemented      ErrKind = "NotImplemented"
	KindRangeNotSatisfiable ErrKind = "RangeNotSatisfiable"
	KindOverflowError       ErrKind = "OverflowError"
	KindNaNError            ErrKind = "NaNError"
)
