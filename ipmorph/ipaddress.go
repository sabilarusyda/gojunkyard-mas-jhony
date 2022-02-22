package ipmorph

import (
	"encoding/binary"
	"errors"
	"net"
)

func ConvertNetIPToDecimal(IPAddress string) (uint32, error) {
	ipNetIP := net.ParseIP(IPAddress)
	if ipNetIP == nil {
		return 0, errors.New("Invalid ip address value")
	}
	return binary.BigEndian.Uint32(ipNetIP[12:16]), nil
}

func ConvertDecimalToNetIP(IPAddress uint32) net.IP {
	ipByteForm := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByteForm, IPAddress)
	return net.IPv4(ipByteForm[0], ipByteForm[1], ipByteForm[2], ipByteForm[3])
}
