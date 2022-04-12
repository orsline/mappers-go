package driver

/*
#include <dlfcn.h>
#include <stdlib.h>
int open_device(unsigned int** device, char* deviceSN, char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int**, char*, char**);
    handle = dlopen("../bin/librcapi_arm64.so", 1);
    if(handle == NULL){
        *error = (char *)dlerror();
        return -1;
    }
    FPTR fptr = (FPTR)dlsym(handle, "open_device");
    if(fptr == NULL){
        *error = (char *)dlerror();
		return -1;
    }
    int result = (*fptr)(device, deviceSN, error);
	dlclose(handle);
    return result;
}

int set_value (unsigned int* device, char* feature, char* value, char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int*, char*, char*, char**);
    handle = dlopen("../bin/librcapi_arm64.so", 1);
    if(handle == NULL){
        *error = (char *)dlerror();
        return -1;
    }
    FPTR fptr = (FPTR)dlsym(handle, "set_value");
    if(fptr == NULL){
        *error = (char *)dlerror();
		return -1;
    }
    int result = (*fptr)(device, feature, value, error);
	dlclose(handle);
    return result;
}

int get_value (unsigned int* device, char* feature, char** value, char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int*, char*, char**, char**);
	handle = dlopen("../bin/librcapi_arm64.so", 1);
    if(handle == NULL){
        *error = (char *)dlerror();
        return -1;
    }
    FPTR fptr = (FPTR)dlsym(handle, "get_value");
    if(fptr == NULL){
        *error = (char *)dlerror();
		return -1;
    }
    int result = (*fptr)(device, feature, value, error);
	dlclose(handle);
    return result;
}

int get_image (unsigned int* device, char* type, char** bufferPointer, int* size, char** error)
{
    typedef int (*FPTR)(unsigned int*, char*, char**, int*, char**);
	void* handle;
	handle = dlopen("../bin/librcapi_arm64.so", 1);
    if(handle == NULL){
        *error = (char *)dlerror();
		return -1;
    }
    FPTR fptr = (FPTR)dlsym(handle, "get_image");
    if(fptr == NULL){
        *error = (char *)dlerror();
		return -1;
    }
    int result = (*fptr)(device, type, bufferPointer, size, error);
	dlclose(handle);
    return result;
}

int close_device (unsigned int* device)
{
    void* handle;
    typedef void (*FPTR)(unsigned int*);
    handle = dlopen("../bin/librcapi_arm64.so", 1);
    if(handle == NULL){
		return -1;
    }
    FPTR fptr = (FPTR)dlsym(handle, "close_device");
    if(fptr == NULL){
		return -1;
    }
    (*fptr)(device);
	int result = dlclose(handle);
    return result;
}
#cgo LDFLAGS: -ldl
*/
import "C"
import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func (gigEClient *GigEVisionDevice) Set(DeviceSN, FeatureName string, value interface{}) (err error) {
	var convert string
	switch value := value.(type) {
	case float64:
		convert = strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		convert = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case int:
		convert = strconv.Itoa(value)
	case uint:
		convert = strconv.Itoa(int(value))
	case int8:
		convert = strconv.Itoa(int(value))
	case uint8:
		convert = strconv.Itoa(int(value))
	case int16:
		convert = strconv.Itoa(int(value))
	case uint16:
		convert = strconv.Itoa(int(value))
	case int32:
		convert = strconv.Itoa(int(value))
	case uint32:
		convert = strconv.Itoa(int(value))
	case int64:
		convert = strconv.FormatInt(value, 10)
	case uint64:
		convert = strconv.FormatUint(value, 10)
	case string:
		convert = value
	case bool:
		convert = strconv.FormatBool(value)
	case []byte:
		convert = string(value)
	default:
		if err == nil {
			errorMsg := fmt.Sprintf("Assertion is not supported for %v type ", value)
			err = errors.New(errorMsg)
		}
		return err
	}
	var msg *C.char
	defer C.free(unsafe.Pointer(msg))
	signal := C.set_value(gigEClient.dev[DeviceSN], C.CString(FeatureName), C.CString(convert), &msg)
	if signal != 0 {
		errorMsg := fmt.Sprintf("Set command from device %s type failed: %s", DeviceSN, (string)(C.GoString(msg)))
		err = errors.New(errorMsg)
		if signal == 1 || signal == 2 {
			gigEClient.dev[DeviceSN] = nil
			go gigEClient.ReConnectDevice(DeviceSN)
		}
		return err
	}
	return nil
}

func (gigEClient *GigEVisionDevice) Get(DeviceSN string, FeatureName string) (results string, err error) {
	var imageFormat = C.CString("png") //png or pnm
	defer C.free(unsafe.Pointer(imageFormat))
	if strings.EqualFold(FeatureName, "image") {
		var imageBuffer *byte
		var size int
		var p = &imageBuffer
		var msg *C.char
		defer C.free(unsafe.Pointer(msg))
		signal := C.get_image(gigEClient.dev[DeviceSN], imageFormat, (**C.char)(unsafe.Pointer(p)), (*C.int)(unsafe.Pointer(&size)), &msg)
		if signal != 0 {
			errorMsg := fmt.Sprintf("Failed to get %s's images: %s", DeviceSN, (string)(C.GoString(msg)))
			err = errors.New(errorMsg)
			if signal != 1 {
				gigEClient.dev[DeviceSN] = nil
				go gigEClient.ReConnectDevice(DeviceSN)
			}
			return "", err
		}
		var buffer []byte
		var bufferHdr = (*reflect.SliceHeader)(unsafe.Pointer(&buffer))
		bufferHdr.Data = uintptr(unsafe.Pointer(imageBuffer))
		bufferHdr.Len = size
		bufferHdr.Cap = size
		results = base64.StdEncoding.EncodeToString(buffer)
	} else {
		var msg *C.char
		var value *C.char
		signal := C.get_value(gigEClient.dev[DeviceSN], C.CString(FeatureName), &value, &msg)
		if signal != 0 {
			errorMsg := fmt.Sprintf("Get command from device %s's failed: %s", DeviceSN, (string)(C.GoString(msg)))
			err = errors.New(errorMsg)
			gigEClient.dev[DeviceSN] = nil
			go gigEClient.ReConnectDevice(DeviceSN)
			return "", err
		}
		results = C.GoString(value)
	}
	return results, err
}

func (gigEClient *GigEVisionDevice) NewClient() (err error) {
	var msg *C.char
	var dev *C.uint
	if gigEClient.dev == nil {
		gigEClient.dev = make(map[string]*C.uint)
	}
	if _, ok := gigEClient.dev[gigEClient.protocolCommonConfig.DeviceSN]; !ok {
		if gigEClient.dev[gigEClient.protocolCommonConfig.DeviceSN] == nil {
			signal := C.open_device(&dev, C.CString(gigEClient.protocolCommonConfig.DeviceSN), &msg)
			if signal != 0 {
				errorMsg := fmt.Sprintf("Failed to open device %s failed: %s", gigEClient.protocolCommonConfig.DeviceSN, (string)(C.GoString(msg)))
				err = errors.New(errorMsg)
				return err
			}
			gigEClient.dev[gigEClient.protocolCommonConfig.DeviceSN] = dev
		}
	}
	return nil
}
