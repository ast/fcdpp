package main

import (
	"bytes"
	"fmt"
	"github.com/karalabe/hid"
	"log"
)

const (
	FCD_HID_CMD_QUERY uint8 = 1 // Returns string with "FCDAPP version"

	FCD_HID_CMD_SET_FREQUENCY_KHZ = 100 // Send with 3 byte unsigned little endian frequency in kHz.
	FCD_HID_CMD_SET_FREQUENCY_HZ  = 101 // Send with 4 byte unsigned little endian frequency in Hz, returns with actual frequency set in Hz
	FCD_HID_CMD_GET_FREQUENCY_HZ  = 102 // Returns 4 byte unsigned little endian frequency in Hz.

	FCD_HID_CMD_SET_LNA_GAIN   = 110 // Send one byte, 1 on, 0 off
	FCD_HID_CMD_SET_RF_FILTER  = 113 // Send one byte enum, see TUNERRFFILTERENUM
	FCD_HID_CMD_SET_MIXER_GAIN = 114 // Send one byte, 1 on, 0 off
	FCD_HID_CMD_SET_IF_GAIN    = 117 // Send one byte value, valid value 0 to 59 (dB)
	FCD_HID_CMD_SET_IF_FILTER  = 122 // Send one byte enum, see TUNERIFFILTERENUM
	FCD_HID_CMD_SET_BIAS_TEE   = 126 // Send one byte, 1 on, 0 off

	FCD_HID_CMD_GET_LNA_GAIN   = 150 // Returns one byte, 1 on, 0 off
	FCD_HID_CMD_GET_RF_FILTER  = 153 // Returns one byte enum, see TUNERRFFILTERENUM
	FCD_HID_CMD_GET_MIXER_GAIN = 154 // Returns one byte, 1 on, 0 off
	FCD_HID_CMD_GET_IF_GAIN    = 157 // Returns one byte value, valid value 0 to 59 (dB)
	FCD_HID_CMD_GET_IF_FILTER  = 162 // Returns one byte enum, see TUNERIFFILTERENUM
	FCD_HID_CMD_GET_BIAS_TEE   = 166 // Returns one byte, 1 on, 0 off
)

type RFFilter byte

func (filter RFFilter) String() string {
	switch filter {
	case TRFE_0_4:
		return "0-4MHz"
	case TRFE_4_8:
		return "4-8MHz"
	case TRFE_8_16:
		return "8-16MHz"
	case TRFE_16_32:
		return "16-32MHz"
	case TRFE_32_75:
		return "32-75MHz"
	case TRFE_75_125:
		return "75-125MHz"
	case TRFE_125_250:
		return "125-250MHz"
	case TRFE_145:
		return "145MHz"
	case TRFE_410_875:
		return "410-875MHz"
	case TRFE_435:
		return "435MHz"
	case TRFE_875_2000:
		return "875-2000MHz"
	}

	return "unknown RF filter"
}

const (
	TRFE_0_4 = RFFilter(iota)
	TRFE_4_8
	TRFE_8_16
	TRFE_16_32
	TRFE_32_75
	TRFE_75_125
	TRFE_125_250
	TRFE_145
	TRFE_410_875
	TRFE_435
	TRFE_875_2000
)

type IFFilter byte

func (filter IFFilter) String() string {
	switch filter {
	case TIFE_200KHZ:
		return "200kHz"
	case TIFE_300KHZ:
		return "300kHz"
	case TIFE_600KHZ:
		return "600kHz"
	case TIFE_1536KHZ:
		return "1536kHz"
	case TIFE_5MHZ:
		return "5MHz"
	case TIFE_6MHZ:
		return "6MHz"
	case TIFE_7MHZ:
		return "7MHz"
	case TIFE_8MHZ:
		return "8MHz"
	}

	return "unknown IF filter"
}

const (
	TIFE_200KHZ IFFilter = IFFilter(iota)
	TIFE_300KHZ
	TIFE_600KHZ
	TIFE_1536KHZ
	TIFE_5MHZ
	TIFE_6MHZ
	TIFE_7MHZ
	TIFE_8MHZ
)

type FCDPP struct {
	dev *hid.Device
	buf []byte
}

func NewFCDPP(dev *hid.Device) *FCDPP {
	fcdpp := &FCDPP{
		dev: dev,
		buf: make([]byte, 65, 65),
	}
	return fcdpp
}

func (fcdpp *FCDPP) clearBuf() {
	for i := range fcdpp.buf {
		fcdpp.buf[i] = 0x00
	}
}

func (fcdpp *FCDPP) writeRead() {
	if _, err := fcdpp.dev.Write(fcdpp.buf); err != nil {
		log.Fatal(err)
	}
	if _, err := fcdpp.dev.Read(fcdpp.buf); err != nil {
		log.Fatal(err)
	}
}

func (fcdpp *FCDPP) Close() {
	fcdpp.dev.Close()
}

func (fcdpp *FCDPP) Query() string {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_QUERY
	fcdpp.writeRead()
	n := bytes.IndexByte(fcdpp.buf, 0)
	return string(fcdpp.buf[2:n])
}

func (fcdpp *FCDPP) Frequency() uint32 {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_FREQUENCY_HZ
	fcdpp.writeRead()

	freq := uint32(fcdpp.buf[2])
	freq |= uint32(fcdpp.buf[3]) << 8
	freq |= uint32(fcdpp.buf[4]) << 16
	freq |= uint32(fcdpp.buf[5]) << 24

	return freq
}

func (fcdpp *FCDPP) LNAGain() bool {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_LNA_GAIN
	fcdpp.writeRead()

	if fcdpp.buf[2] == 1 {
		return true
	}

	return false
}

func (fcdpp *FCDPP) RFFilter() RFFilter {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_RF_FILTER
	fcdpp.writeRead()
	return RFFilter(fcdpp.buf[2])
}

func (fcdpp *FCDPP) MixerGain() bool {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_MIXER_GAIN
	fcdpp.writeRead()

	if fcdpp.buf[2] == 1 {
		return true
	}

	return false
}

func (fcdpp *FCDPP) IFGain() uint8 {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_IF_GAIN
	fcdpp.writeRead()
	return fcdpp.buf[2]
}

func (fcdpp *FCDPP) IFFilter() IFFilter {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_IF_FILTER
	fcdpp.writeRead()
	return IFFilter(fcdpp.buf[2])
}

func (fcdpp *FCDPP) BiasTee() bool {
	fcdpp.clearBuf()
	fcdpp.buf[1] = FCD_HID_CMD_GET_IF_GAIN
	fcdpp.writeRead()

	if fcdpp.buf[2] == 1 {
		return true
	}

	return false
}

func main() {
	// get device
	devs := hid.Enumerate(0x04d8, 0xfb31)
	if len(devs) != 1 {
		log.Fatal("could not open device")
	}

	dev, err := devs[0].Open()
	if err != nil {
		log.Fatal(err)
	}
	fcdpp := NewFCDPP(dev)
	defer fcdpp.Close()

	fmt.Printf("%12s %v\n", "Query:", fcdpp.Query())
	fmt.Printf("%12s %v Hz\n", "Freq:", fcdpp.Frequency())
	fmt.Printf("%12s %v\n", "LNA gain:", fcdpp.LNAGain())
	fmt.Printf("%12s %v\n", "RF filter:", fcdpp.RFFilter())
	fmt.Printf("%12s %v\n", "Mixer gain:", fcdpp.MixerGain())
	fmt.Printf("%12s %v dB\n", "IF gain:", fcdpp.IFGain())
	fmt.Printf("%12s %v\n", "IF filter:", fcdpp.IFFilter())
	fmt.Printf("%12s %v\n", "Bias tee:", fcdpp.BiasTee())
}
