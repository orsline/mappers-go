#include "close_device.h"

extern "C" int close_device(MyDevice myDevice) {
	int ret = 0;
	if (myDevice.dev != nullptr) {
		myDevice.dev->close();
	}
	else {
		ret = 1;
	}
	return ret;
}