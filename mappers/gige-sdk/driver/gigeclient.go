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
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
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
	if strings.EqualFold(FeatureName, "ImageFormat") {
		if strings.EqualFold(convert, "png") {
			gigEClient.deviceMeta[DeviceSN].imageFormat = "png"
		} else if strings.EqualFold(convert, "pnm") {
			gigEClient.deviceMeta[DeviceSN].imageFormat = "pnm"
		} else if strings.EqualFold(convert, "jpeg") {
			gigEClient.deviceMeta[DeviceSN].imageFormat = "jpeg"
		} else {
			errorMsg := fmt.Sprintf("Set %s's image format failed, it only support format jpeg、 png or pnm.", DeviceSN)
			err = errors.New(errorMsg)
			return err
		}
		return nil
	} else if strings.EqualFold(FeatureName, "ImageURL") {
		convert = strings.TrimSpace(convert)
		_, err := url.Parse(convert)
		if err != nil {
			errorMsg := fmt.Sprintf("Set imageURL failed because of incorrect format, message: %s", err)
			err = errors.New(errorMsg)
			return err
		}
		gigEClient.deviceMeta[DeviceSN].imageURL = convert
		gigEClient.PostImage(DeviceSN)
		return nil
	} else {
		var msg *C.char
		defer C.free(unsafe.Pointer(msg))
		signal := C.set_value(gigEClient.deviceMeta[DeviceSN].dev, C.CString(FeatureName), C.CString(convert), &msg)
		if signal != 0 {
			errorMsg := fmt.Sprintf("Set command from device %s type failed: %s", DeviceSN, (string)(C.GoString(msg)))
			err = errors.New(errorMsg)
			if signal == 1 || signal == 2 {
				gigEClient.deviceMeta[DeviceSN].deviceStatus = false
				go gigEClient.ReConnectDevice(DeviceSN)
			}
			return err
		}
		return nil
	}
}

func (gigEClient *GigEVisionDevice) Get(DeviceSN string, FeatureName string) (results string, err error) {
	if strings.EqualFold(FeatureName, "Image") {
		var imageBuffer *byte
		var size int
		var p = &imageBuffer
		var msg *C.char
		defer C.free(unsafe.Pointer(msg))
		signal := C.get_image(gigEClient.deviceMeta[DeviceSN].dev, C.CString(gigEClient.deviceMeta[DeviceSN].imageFormat), (**C.char)(unsafe.Pointer(p)), (*C.int)(unsafe.Pointer(&size)), &msg)
		if signal != 0 {
			errorMsg := fmt.Sprintf("Failed to get %s's images: %s", DeviceSN, (string)(C.GoString(msg)))
			err = errors.New(errorMsg)
			if signal != 1 {
				gigEClient.deviceMeta[DeviceSN].deviceStatus = false
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
	} else if strings.EqualFold(FeatureName, "ImageFormat") {
		if gigEClient.deviceMeta[DeviceSN].imageFormat != "" {
			results = gigEClient.deviceMeta[DeviceSN].imageFormat
		} else {
			errorMsg := fmt.Sprintf("Maybe init %s's image format failed, it only support format png or pnm.", DeviceSN)
			err = errors.New(errorMsg)
			return "", err
		}
	} else if strings.EqualFold(FeatureName, "ImageURL") {
		if gigEClient.deviceMeta[DeviceSN].imageURL != "" {
			results = gigEClient.deviceMeta[DeviceSN].imageURL
		} else {
			errorMsg := fmt.Sprintf("Maybe init %s's image format failed, it only support format png or pnm.", DeviceSN)
			err = errors.New(errorMsg)
			return "", err
		}
		gigEClient.PostImage(DeviceSN)
		return results, nil
	} else {
		var msg *C.char
		var value *C.char
		signal := C.get_value(gigEClient.deviceMeta[DeviceSN].dev, C.CString(FeatureName), &value, &msg)
		if signal != 0 {
			errorMsg := fmt.Sprintf("Get command from device %s's failed: %s", DeviceSN, (string)(C.GoString(msg)))
			err = errors.New(errorMsg)
			if signal == 1 {
				gigEClient.deviceMeta[DeviceSN].deviceStatus = false
				go gigEClient.ReConnectDevice(DeviceSN)
			}
			return "", err
		}
		results = C.GoString(value)
	}
	return results, err
}

func (gigEClient *GigEVisionDevice) NewClient() (err error) {
	var msg *C.char
	var dev *C.uint
	if gigEClient.deviceMeta == nil {
		gigEClient.deviceMeta = make(map[string]*DeviceMeta)
	}
	if _, ok := gigEClient.deviceMeta[gigEClient.protocolCommonConfig.DeviceSN]; !ok {
		if gigEClient.deviceMeta[gigEClient.protocolCommonConfig.DeviceSN] == nil {
			signal := C.open_device(&dev, C.CString(gigEClient.protocolCommonConfig.DeviceSN), &msg)
			if signal != 0 {
				klog.Infof("Failed to open device %s failed: %s", gigEClient.protocolCommonConfig.DeviceSN, (string)(C.GoString(msg)))
				gigEClient.deviceMeta[gigEClient.protocolCommonConfig.DeviceSN] = &DeviceMeta{
					dev:          nil,
					deviceStatus: false,
					imageFormat:  "jpeg",
					imageURL:     "http://192.168.137.61:8081/image_infer",
				}
				go gigEClient.ReConnectDevice(gigEClient.protocolCommonConfig.DeviceSN)
				return nil
			}
			gigEClient.deviceMeta[gigEClient.protocolCommonConfig.DeviceSN] = &DeviceMeta{
				dev:          dev,
				deviceStatus: true,
				imageFormat:  "jpeg",
				imageURL:     "http://192.168.137.61:8081/image_infer",
			}
		}
	}
	return nil
}

func (gigEClient *GigEVisionDevice) PostImage(DeviceSN string) {
	var imageBuffer *byte
	var size int
	var p = &imageBuffer
	var msg *C.char
	defer C.free(unsafe.Pointer(msg))
	signal := C.get_image(gigEClient.deviceMeta[DeviceSN].dev, C.CString(gigEClient.deviceMeta[DeviceSN].imageFormat), (**C.char)(unsafe.Pointer(p)), (*C.int)(unsafe.Pointer(&size)), &msg)
	if signal != 0 {
		fmt.Printf("Failed to get %s's images: %s\n", DeviceSN, (string)(C.GoString(msg)))
		if signal != 1 {
			gigEClient.deviceMeta[DeviceSN].deviceStatus = false
			go gigEClient.ReConnectDevice(DeviceSN)
		}
		return
	}
	go func() {
		var buffer []byte
		var bufferHdr = (*reflect.SliceHeader)(unsafe.Pointer(&buffer))
		bufferHdr.Data = uintptr(unsafe.Pointer(imageBuffer))
		bufferHdr.Len = size
		bufferHdr.Cap = size
		klog.V(4).Infof("buffer:%s", buffer[:100])
		postStr := base64.URLEncoding.EncodeToString(buffer)
		klog.V(4).Infof("poststr:%s", postStr[:100])
		v := url.Values{}
		v.Set("gigEImage", postStr)
		body := ioutil.NopCloser(strings.NewReader(v.Encode()))
		req, _ := http.NewRequest(http.MethodPost, gigEClient.deviceMeta[DeviceSN].imageURL, body)
		if req == nil {
			fmt.Printf("Failed to post %s's images: URL can't POST\n", DeviceSN)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		client := &http.Client{}
		resp, _ := client.Do(req)
		if resp == nil {
			fmt.Printf("Failed to post %s's images: URL no reaction\n", DeviceSN)
			return
		}
		defer resp.Body.Close()
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(data))
	}()
}
