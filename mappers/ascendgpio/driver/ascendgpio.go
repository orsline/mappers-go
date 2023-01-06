package main

import (
	"fmt"
	"os"
	"time"
)

type Mode uint8
type Pin uint8
type State uint8

//
const (
	ascend_gpio_0_dir = "/sys/class/gpio/gpio504/direction"
	ascend_gpio_1_dir = "/sys/class/gpio/gpio444/direction"
	ascend_gpio_0_val = "/sys/class/gpio/gpio504/value"
	ascend_gpio_1_val = "/sys/class/gpio/gpio444/value"
)
const (
	Input Mode = iota
	Output
)

// State of pin, High / Low
const (
	Low uint8 = iota
	High
)

// setInput: Set pin as InputPin
func (pin Pin) SetInPut() {
	setPinMode(pin, Input)
}

// setOutput: Set pin as Output
func (pin Pin) SetOutPut() {
	setPinMode(pin, Output)
}

// setHight: Set pin Hight
func (pin Pin) SetHight() {
	gpioSetValue(pin, High)
}

// setLow: Set pin as Low
func (pin Pin) SetLow() {
	gpioSetValue(pin, Low)
}

// Write: Set pin state (high/low)
func (pin Pin) Write(val uint8) {
	WritePin(pin, val)
}

// Read pin state (high/low)
func (pin Pin) Read() uint8 {
	return ReadPin(pin)
}

func WritePin(pin Pin, val uint8) {
	gpioSetValue(pin, val)
}
func ReadPin(pin Pin) uint8 {
	var val uint8
	err := gpioGetValue(pin, &val)
	if err != nil {
		return 0
	}
	return val
}

// Spi mode should not be set by this directly, use SpiBegin instead.
func setPinMode(pin Pin, mode Mode) {
	f := uint8(0)
	const in uint8 = 0  // 000
	const out uint8 = 1 // 001

	switch mode {
	case Input:
		f = in
	case Output:
		f = out
	}

	GpioSetDirection(pin, f)

}

// Generic ioctl constants
const (
	IOC_NONE  = 0x0
	IOC_WRITE = 0x1
	IOC_READ  = 0x2

	IOC_NRBITS   = 8
	IOC_TYPEBITS = 8

	IOC_SIZEBITS = 14
	IOC_DIRBITS  = 2

	IOC_NRSHIFT   = 0
	IOC_TYPESHIFT = IOC_NRSHIFT + IOC_NRBITS     //8 + 0
	IOC_SIZESHIFT = IOC_TYPESHIFT + IOC_TYPEBITS //8 + 8
	IOC_DIRSHIFT  = IOC_SIZESHIFT + IOC_SIZEBITS //16 + 14

	IOC_NRMASK   = ((1 << IOC_NRBITS) - 1)
	IOC_TYPEMASK = ((1 << IOC_TYPEBITS) - 1)
	IOC_SIZEMASK = ((1 << IOC_SIZEBITS) - 1)
	IOC_DIRMASK  = ((1 << IOC_DIRBITS) - 1)
)

// Some useful additional ioctl constanst
const (
	IOC_IN        = IOC_WRITE << IOC_DIRSHIFT
	IOC_OUT       = IOC_READ << IOC_DIRSHIFT
	IOC_INOUT     = (IOC_WRITE | IOC_READ) << IOC_DIRSHIFT
	IOCSIZE_MASK  = IOC_SIZEMASK << IOC_SIZESHIFT
	IOCSIZE_SHIFT = IOC_SIZESHIFT
)

// IOC generate IOC
func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IOC_DIRSHIFT) | (t << IOC_TYPESHIFT) |
		(nr << IOC_NRSHIFT) | (size << IOC_SIZESHIFT)
}

// IOR generate IOR
func IOR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ, t, nr, size)
}

// IOW generate IOW
func IOW(t, nr, size uintptr) uintptr {
	return IOC(IOC_WRITE, t, nr, size)
}

// IOWR generate IOWR
func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IOC_READ|IOC_WRITE, t, nr, size)
}

// IO generate IO
func IO(t, nr uintptr) uintptr {
	return IOC(IOC_NONE, t, nr, 0)
}

func Open() (err error) {
	return nil
}
func Close() (err error) {
	return nil
}
func pca6416GpioSetValue(pin Pin, val uint8) (err error) {

	return nil
}
func pca6416GpioGetValue(pin Pin, val *uint8) (err error) {

	return nil
}

func GpioSetDirection(pin Pin, dir uint8) error {
	var fileName string
	var direction string
	var err error

	if pin == 0 {
		fileName = ascend_gpio_0_dir
	} else if pin == 1 {
		fileName = ascend_gpio_1_dir
	} else {
		err = fmt.Errorf("pin number is incorrect,must be 0 or 1")
		return err
	}
	direction = "out"
	if dir == 0 {
		direction = "in"
	}

	err = os.WriteFile(fileName, []byte(direction), 0666)
	if err != nil {
		return err
	}

	return nil
}

func AscendGpioSetValue(pin Pin, val uint8) error {
	var fileName string
	var err error

	if pin == 0 {
		fileName = ascend_gpio_0_val
	} else if pin == 1 {
		fileName = ascend_gpio_1_val
	} else {
		err = fmt.Errorf("pin number is incorrect,must be 0 or 1")
		return err
	}

	//buf := bytes.NewBuffer([]byte{})
	//binary.Write(buf, binary.LittleEndian, val)
	buff := []byte{val}
	err = os.WriteFile(fileName, buff, 0666)
	if err != nil {
		return err
	}
	return nil
}

func AscendGpioGetValue(pin Pin, val *uint8) error {
	var fileName string
	if pin == 0 {
		fileName = ascend_gpio_0_dir
	} else if pin == 1 {
		fileName = ascend_gpio_1_dir
	} else {
		err := fmt.Errorf("pin number is incorrect,the correct num is must be 0,1")
		return err
	}
	readFile, err := os.ReadFile(fileName)

	*val = readFile[0]
	return err
}
func gpioSetValue(pin Pin, val uint8) error {
	if pin == 0 || pin == 1 {
		return AscendGpioSetValue(pin, val)
	} else {
		return pca6416GpioSetValue(pin, val)
	}
}
func gpioGetValue(pin Pin, val *uint8) error {
	if pin == 0 || pin == 1 {
		return AscendGpioGetValue(pin, val)
	} else {
		return pca6416GpioGetValue(pin, val)
	}
}

func main() {
	var pin Pin

	pin = 0
	pinClient := Pin(pin)
	for i := 0; i < 10; i++ {
		pinClient.SetOutPut()
		pinClient.SetLow()
		fmt.Println("set outPut hight")
		time.Sleep(time.Second)
		pinClient.SetOutPut()
		pinClient.SetHight()
		fmt.Println("set outPut low")
		time.Sleep(time.Second)

	}
}

// IOCTL send ioctl
//func IOCTL(fd, name, data uintptr) error {
//	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, name, data)
//	if err != 0 {
//		return syscall.Errno(err)
//	}
//	return nil
//}
