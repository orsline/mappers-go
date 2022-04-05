# GPIO Mapper

## How to use this
0. Modify and apply crd [instance.yaml](../../build/crd-samples/devices/gpio-device-instance.yaml) and [model.yaml](../../build/crd-samples/devices/gpio-device-model.yaml)
1. Copy this Mapper to Raspberry PI
2. Connect the led to the corresponding pin port  
Yellow LED PIN :17  
Red    LED PIN :22   
Green  LED PIN :27 
3. `make build`
4. `./bin/main`