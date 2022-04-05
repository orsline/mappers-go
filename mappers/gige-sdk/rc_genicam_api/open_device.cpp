#include "open_device.h"

extern "C" int open_device(MyDevice * myDevice, const char* device_serial_number, char** err) {
	int ret = 0;
	*err = (char*)"";
	//std::cout << device_serial_number << std::endl;
	*myDevice = MyDevice();
	(*myDevice).dev = rcg::getDevice(device_serial_number);
	if ((*myDevice).dev != 0) {
		try {
			(*myDevice).dev->open(rcg::Device::CONTROL);
		}
		catch (const std::exception& ex) {
			ret = 1;
			std::string e = ex.what();
			*err = (char*)malloc(e.length());
			strcpy(*err, ex.what());
			//std::cout << ex.what() << std::endl;
		}
	}
	else {
		ret = 2;
		//std::cout << "Cannot find device: " << device_serial_number << std::endl;
		std::string e = "Cannot find device: " + (std::string)device_serial_number;
		*err = (char*)malloc(e.length());
		strcpy(*err, e.c_str());
	}
	return ret;
}
