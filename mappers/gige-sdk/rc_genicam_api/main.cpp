#include <iostream>
#include <time.h>
#include "find_device.h"
#include "get_value.h"
#include "set_value.h"
#include "get_image.h"
#include "close_device.h"
#include "open_device.h"
#include "const.h"

using namespace std;

int main() {
	time_t a, b;
	find_device();
	std::string key_value = "";
	const char* device_id = "";
	char* err = new char();
	int ret;
	MyDevice myDevice;
	cout << "device_id:";
	cin >> key_value;
	device_id = key_value.c_str();
	ret = open_device(&myDevice, device_id, &err);
	std::cout << ret << " " << err << endl;
	while (1) {
		unsigned char* image_buffer = NULL;
		int size = 0;
		cout << "key_value: ";
		cin >> key_value;
		try {
			if (key_value == "exit") {
				break;
			}
			else if (key_value == "png" || key_value == "pnm") {
				a = clock();
				const char* imgfmt = key_value.c_str();
				ret = get_image(myDevice, imgfmt, &image_buffer, &size, &err);
				cout << "err:" << err << endl;
				cout << "The size of image:" << size << endl;
				std::string name = "test." + key_value;
				std::ofstream out(name, std::ios::binary);
				std::streambuf* sb = out.rdbuf();
				for (int i = 0; i < size && out.good(); i++, image_buffer++) {
					sb->sputc(*image_buffer);
				}
				out.close();
				b = clock();
				std::cout << b - a << "  " << (double)((b - a) / CLOCKS_PER_SEC) << std::endl;
			}
			else {
				size_t k = key_value.find('=');
				char* value = NULL;
				get_value(myDevice, key_value.substr(0, k).c_str(), &value, &err);
				if (value != nullptr) {
					std::cout << key_value.substr(0, k).c_str() << ": " << value << endl;
				}
				std::cout << "1:" << err << std::endl;
				std::cout << set_value(myDevice, key_value.substr(0, k).c_str(), key_value.substr(k + 1).c_str(), &err) << " " << err << std::endl;
				get_value(myDevice, key_value.substr(0, k).c_str(), &value, &err);
				if (value != nullptr) {
					std::cout << key_value.substr(0, k).c_str() << ": " << value << endl;
				}
				std::cout << "2:" << err << std::endl;
			}
		}
		catch (const std::exception& ex) {
			std::cout << "Exception: " << ex.what() << std::endl;
		}
		catch (const GENICAM_NAMESPACE::GenericException& ex) {
			std::cout << "Exception: " << ex.what() << std::endl;
		}
		catch (...) {
			std::cout << "Unknown exception!" << std::endl;
		}

	}
	close_device(myDevice);
	rcg::System::clearSystems();

	return 0;
}