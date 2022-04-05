# MQTT Mapper


## How to use this
0. Modify and apply crd [instance.yaml](../../build/crd-samples/devices/mqtt-gateway-instance.yaml) and [model.yaml](../../build/crd-samples/devices/mqtt-gateway-model.yaml)
1. Copy this Mapper to device that can subscribe mqtt broker
2. Connect the sensor device to the gateway
3. `make build`
4. `./bin/main`