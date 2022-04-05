// Package httpadapter is a package to process RESTful message
package httpadapter

import (
	"encoding/json"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/application"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/common"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/httpadapter/requests"
	"github.com/kubeedge/mappers-go/mapper-sdk-go/internal/httpadapter/response"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"strings"
)

// AddDevice Restful api to addDevice
func (c *RestController) AddDevice(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var addDeviceRequest requests.AddDeviceRequest
	err := json.NewDecoder(request.Body).Decode(&addDeviceRequest)
	if err != nil {
		klog.Error("Failed to decode JSON", err)
		c.sendMapperError(writer, request, err.Error(), common.ApiDeviceCallbackRoute)
		return
	}
	kind := application.AddDevice(addDeviceRequest, c.dic)
	if kind == "" {
		baseMessage := response.NewBaseResponse("", "", http.StatusOK)
		res := response.NewUpdateDeviceResponse(baseMessage, addDeviceRequest.DeviceInstance.ID, "add device", "Successful")
		c.sendResponse(writer, request, common.ApiDeviceCallbackRoute, res, http.StatusOK)
	} else {
		httpCode := response.CodeMapping(kind)
		baseMessage := response.NewBaseResponse("", "", httpCode)
		res := response.NewUpdateDeviceResponse(baseMessage, addDeviceRequest.DeviceInstance.ID, "add device", string(kind))
		c.sendResponse(writer, request, common.ApiDeviceCallbackRoute, res, httpCode)
	}
}

// RemoveDevice Restful api to remove device
func (c *RestController) RemoveDevice(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	urlItem := strings.Split(request.URL.Path, "/")
	itemLen := len(urlItem)
	instanceID := urlItem[itemLen-1]
	kind := application.DeleteDevice(instanceID, c.dic)
	if kind == "" {
		baseMessage := response.NewBaseResponse("", "", http.StatusOK)
		res := response.NewUpdateDeviceResponse(baseMessage, instanceID, "remove device", "Successful")
		c.sendResponse(writer, request, common.ApiDeviceCallbackIdRoute, res, http.StatusOK)
	} else {
		httpCode := response.CodeMapping(kind)
		baseMessage := response.NewBaseResponse("", "", httpCode)
		res := response.NewUpdateDeviceResponse(baseMessage, instanceID, "remove device", string(kind))
		c.sendResponse(writer, request, common.ApiDeviceCallbackIdRoute, res, httpCode)
	}
}

// WriteCommand  Restful api to write data to the device
func (c *RestController) WriteCommand(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	var reserved url.Values
	var err error
	_, reserved, err = filterQueryParams(request.URL.RawQuery)
	if err != nil {
		return
	}
	urlItem := strings.Split(request.URL.Path, "/")
	itemLen := len(urlItem)
	if len(reserved) != 1 {
		baseMessage := response.NewBaseResponse("", "Some errors have occurred", 500)
		c.sendResponse(writer, request, common.ApiDeviceWriteCommandByIdRoute, baseMessage, 500)
		return
	}
	kind := application.WriteDeviceData(urlItem[itemLen-1], reserved, c.dic)
	propertyName := ""
	for k := range reserved {
		propertyName = k
	}
	httpCode := response.CodeMapping(kind)
	baseMessage := response.NewBaseResponse("", "", httpCode)
	if httpCode < 300 {
		res := response.NewWriteCommandResponse(baseMessage, urlItem[itemLen-1], propertyName, "successful")
		c.sendResponse(writer, request, common.ApiDeviceWriteCommandByIdRoute, res, httpCode)
	} else {
		res := response.NewWriteCommandResponse(baseMessage, urlItem[itemLen-1], propertyName, "failed")
		c.sendResponse(writer, request, common.ApiDeviceWriteCommandByIdRoute, res, httpCode)
	}
}

// ReadCommand Restful api to read data from the device
func (c *RestController) ReadCommand(writer http.ResponseWriter, request *http.Request) {
	urlItem := strings.Split(request.URL.Path, "/")
	itemLen := len(urlItem)
	value, kind := application.ReadDeviceData(urlItem[itemLen-2], urlItem[itemLen-1], c.dic)
	httpCode := response.CodeMapping(kind)
	baseMessage := response.NewBaseResponse("", "", httpCode)
	if httpCode < 300 {
		res := response.NewReadCommandResponse(baseMessage, urlItem[itemLen-2], urlItem[itemLen-1], value)
		c.sendResponse(writer, request, common.ApiDeviceReadCommandByIdRoute, res, httpCode)
	} else {
		res := response.NewReadCommandResponse(baseMessage, urlItem[itemLen-2], urlItem[itemLen-1], string(kind))
		c.sendResponse(writer, request, common.ApiDeviceReadCommandByIdRoute, res, httpCode)
	}
}
