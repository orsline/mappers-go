package driver

/*
#include <dlfcn.h>

void find_device()
{
    void* handle;
    typedef void (*FPTR)();

    handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "find_device");
    (*fptr)();
	dlclose(handle);
    return;
}

int open_device(unsigned int** device,char* deviceId,char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int**,char*,char**);

    handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "open_device");

    int result = (*fptr)(device,deviceId,error);
	dlclose(handle);
    return result;
}

int set_value (unsigned int* device, char* feature, char* value,char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int*,char*,char*,char**);

    handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "set_value");

    int result = (*fptr)(device,feature,value,error);
	dlclose(handle);
    return result;
}

int get_value (unsigned int* device, char* feature, char** value,char** error)
{
    void* handle;
    typedef int (*FPTR)(unsigned int*,char*,char**,char**);

	handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "get_value");

    int result = (*fptr)(device,feature,value,error);
	dlclose(handle);
    return result;
}

int get_image (unsigned int* device, char* type, char** bufferPointer,int* size,char** error)
{
    typedef int (*FPTR)(unsigned int*, char*, char**,int*,char**);
	void* handle;
	handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "get_image");

    int result = (*fptr)(device,type,bufferPointer,size,error);
	dlclose(handle);
    return result;
}

int close_device (unsigned int* device)
{
    void* handle;
    typedef void (*FPTR)(unsigned int*);

    handle = dlopen("../bin/librcapi_arm64.so", 1);
    FPTR fptr = (FPTR)dlsym(handle, "close_device");

    (*fptr)(device);
	int result=dlclose(handle);
    return result;
}//链接dl库
#cgo LDFLAGS: -ldl
*/
import "C"
import (
	"encoding/json"
	"errors"
	"k8s.io/klog/v2"
	"reflect"
	"strconv"
	"unsafe"
)

var clients map[string]*C.uint

func (geClient *GigEVisionDevice) Set(FeatureName string, value interface{}) (err error) {
	geClient.mutex.Lock()
	if geClient.dev == nil {
		var errorMsg = "Device " + geClient.ProtocolCommonConfig.DeviceSN + " unreachabled " + " failed to set"
		err = errors.New(errorMsg)
		return err
	}
	var convert string
	switch value.(type) {
	case float64:
		ft := value.(float64)
		convert = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		convert = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		convert = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		convert = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		convert = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		convert = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		convert = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		convert = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		convert = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		convert = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		convert = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		convert = strconv.FormatUint(it, 10)
	case string:
		convert = value.(string)
	case bool:
		it := value.(bool)
		convert = strconv.FormatBool(it)
	case []byte:
		convert = string(value.([]byte))
	default:
		_, err := json.Marshal(value)
		errors.As(err, "断言失败")
		return err
	}
	defer geClient.mutex.Unlock()
	var msg *C.char
	signal := C.set_value(geClient.dev, C.CString(FeatureName), C.CString(convert), &msg)
	if signal != 0 {
		var errorMsg = "Set command from device " + geClient.ProtocolCommonConfig.DeviceSN + " failed : " + C.GoString(msg)
		err = errors.New(errorMsg)
		if signal == 1|2 {
			go geClient.ReConnectDevice()
		}
		return err
	}
	klog.Infof("Set command success from device %s set %s : %s ", geClient.ProtocolCommonConfig.DeviceSN, FeatureName, convert)
	return nil
}

func (geClient *GigEVisionDevice) Get(FeatureName string) (results []byte, err error) {
	geClient.mutex.RLock()
	defer geClient.mutex.RUnlock()
	var imageForm = C.CString("png")
	if geClient.dev == nil {
		var errorMsg = "Device " + geClient.ProtocolCommonConfig.DeviceSN + " unreachabled " + " failed to get"
		err = errors.New(errorMsg)
		return nil, err
	}
	if FeatureName == "image" {
		var imageBuffer *byte
		var size int
		var p = &imageBuffer
		var msg *C.char
		signal := C.get_image(geClient.dev, imageForm, (**C.char)(unsafe.Pointer(p)), (*C.int)(unsafe.Pointer(&size)), &msg)
		if signal != 0 {
			var errorMsg = "Failed to get " + geClient.ProtocolCommonConfig.DeviceSN + "'s images : " + C.GoString(msg)
			err = errors.New(errorMsg)
			if signal == 2|3|4|5 {
				go geClient.ReConnectDevice()
			}
			return nil, err
		}
		klog.Infof("Get images success from device %s , image form is %s", geClient.ProtocolCommonConfig.DeviceSN, imageForm)
		var bufferHdr = (*reflect.SliceHeader)(unsafe.Pointer(&results))
		bufferHdr.Data = uintptr(unsafe.Pointer(imageBuffer))
		bufferHdr.Len = size
		bufferHdr.Cap = size
	} else {
		var msg *C.char
		var value *C.char
		signal := C.get_value(geClient.dev, C.CString(FeatureName), &value, &msg)
		if signal != 0 {
			var errorMsg = "Get command from device " + geClient.ProtocolCommonConfig.DeviceSN + " failed : " + C.GoString(msg)
			err = errors.New(errorMsg)
			return nil, err
		}
		klog.Infof("Get command success from device %s get %s : %s ", geClient.ProtocolCommonConfig.DeviceSN, FeatureName, C.GoString(value))
		results = []byte(C.GoString(value))
	}
	return results, err
}

func (geClient *GigEVisionDevice) NewClient() (err error) {
	var msg *C.char
	var dev *C.uint
	addr := geClient.ProtocolCommonConfig.DeviceSN
	if _, ok := clients[addr]; ok {
		return nil
	}
	if clients == nil {
		clients = make(map[string]*C.uint)
	}
	signal := C.open_device(&dev, C.CString(geClient.ProtocolCommonConfig.DeviceSN), &msg)
	if signal != 0 {
		klog.Errorf("Failed to open device %s : %s Please check your camera link", geClient.ProtocolCommonConfig.DeviceSN, C.GoString(msg))
		go geClient.ReConnectDevice()
		var errorMsg = "Failed to open device " + geClient.ProtocolCommonConfig.DeviceSN + " : " + C.GoString(msg)
		err = errors.New(errorMsg)
		return err
	}
	geClient.dev = dev
	clients[addr] = geClient.dev
	return nil
}
