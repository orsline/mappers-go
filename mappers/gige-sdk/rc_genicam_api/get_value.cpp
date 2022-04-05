#include "get_value.h"

extern "C" int get_value(MyDevice myDevice, const char* key, char** value, char** err) {
	int ret = 0;
	*err = (char*)"";
	std::string v = "";
	try {
		std::shared_ptr<GenApi::CNodeMapRef> nodemap = myDevice.dev->getRemoteNodeMap();
		v = rcg::getString(nodemap, key);
		if (v.size() != 0) {
			*value = (char*)malloc(v.length());
			strcpy(*value, v.c_str());
		}
	}
	catch (const GENICAM_NAMESPACE::GenericException& ex) {
		ret = 1;
		//std::cout << ex.what() << std::endl;
		std::string e = ex.what();
		*err = (char*)malloc(e.length());
		strcpy(*err, ex.what());
	}
	catch (const std::exception& ex) {
		ret = 2;
		//std::cout << ex.what() << std::endl;
		std::string e = ex.what();
		*err = (char*)malloc(e.length());
		strcpy(*err, ex.what());
	}
	return ret;
}

