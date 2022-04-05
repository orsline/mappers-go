#include "set_value.h"

extern "C" int set_value(MyDevice myDevice, const char* key, const char* value, char** err) {
	int ret = 0;
	*err = (char*)"";
	try {
		std::shared_ptr<GenApi::CNodeMapRef> nodemap = myDevice.dev->getRemoteNodeMap();
		ret = rcg::setString(nodemap, key, value, true);
	}
	catch (const GENICAM_NAMESPACE::GenericException& ex) {
		ret = 1;
		//std::cout << ret << std::endl;
		//std::cout << ex.what() << std::endl;
		std::string e = ex.what();
		*err = (char*)malloc(e.length());
		strcpy(*err, ex.what());
	}
	catch (const std::exception& ex) {
		ret = 2;
		//std::cout << ret << std::endl;
		//std::cout << ex.what() << std::endl;
		std::string e = ex.what();
		*err = (char*)malloc(e.length());
		strcpy(*err, ex.what());
	}
	return ret;
}