// Package response used to implement the responses structure
package response

import "github.com/kubeedge/mappers-go/mapper-sdk-go/internal/common"

type BaseResponse struct {
	Version    string
	RequestId  string `json:"requestId,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode"`
}

type ReadCommandResponse struct {
	BaseResponse
	DeviceID     string
	PropertyName string
	Value        string
}

type WriteCommandResponse struct {
	BaseResponse
	DeviceID     string
	PropertyName string
	Status       string
}

type UpdateDeviceResponse struct {
	BaseResponse
	DeviceID  string
	Operation string
	Status    string
}

func NewBaseResponse(requestId string, message string, statusCode int) BaseResponse {
	return BaseResponse{
		Version:    common.ApiVersion,
		RequestId:  requestId,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewReadCommandResponse(response BaseResponse, deviceId, propertyName, value string) ReadCommandResponse {
	return ReadCommandResponse{
		response,
		deviceId,
		propertyName,
		value,
	}
}

func NewWriteCommandResponse(response BaseResponse, deviceId, propertyName, status string) WriteCommandResponse {
	return WriteCommandResponse{
		response,
		deviceId,
		propertyName,
		status,
	}
}

func NewUpdateDeviceResponse(response BaseResponse, deviceId, operation, status string) UpdateDeviceResponse {
	return UpdateDeviceResponse{
		response,
		deviceId,
		operation,
		status,
	}
}
