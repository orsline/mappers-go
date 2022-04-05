.PHONY: build help clean docker
help:
	#
	# Usage:
	#   make build  :  compile main.go
	#   make clean  :  clean binaries
	#   make docker :  pack binaries to image
	@echo
build:
	export GENICAM_GENTL64_PATH=$(abspath $(lastword $(MAKEFILE_LIST)))/baumer/Ubuntu-16.04/arm64
	sed -i '1a /bin/genicam/Linux64_ARM' /etc/ld.so.conf
	ldconfig
	sysctl -w net.core.rmem_max=33554432
	sysctl -w net.core.netdev_max_backlog=2000
	sysctl -w net.core.netdev_budget=600
	go build -o ./bin/main ./cmd/main.go
	@echo  "[INFO] go build successful"
	@echo  "[INFO] you can cd ./bin to execute it"
clean:
	rm -f ./bin/main
docker:
	docker build -f ./Dockerfile \
		-t mapper-go-gigevision:v1.0.0 .