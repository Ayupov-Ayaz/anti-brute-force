package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var (
	errInvalidIP   = errors.New("invalid IPNet address")
	errInvalidMask = errors.New("invalid mask")
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func parseIP(ip string) (net.IP, error) {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return net.IP{}, fmt.Errorf("ip=%s: %w", ip, errInvalidIP)
	}

	return ipAddr, nil
}

func (s *Service) ParseIP(ip string) (net.IP, error) {
	return parseIP(ip)
}

func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func (s *Service) IPToUint32(ip net.IP) uint32 {
	return ipToUint32(ip)
}

func (s *Service) ParseCIDR(ip string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, fmt.Errorf("ip=%s: %w", ip, errInvalidMask)
	}

	return ipNet, nil
}
