#pragma once
#include <rc_genicam_api/device.h>

/**
  The MyDevice struct encapsulates a Genicam device.
*/
struct MyDevice
{
	std::shared_ptr<rcg::Device> dev;
};
