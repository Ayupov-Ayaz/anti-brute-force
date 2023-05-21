package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var (
	errInvalidIP   = errors.New("invalid IP address")
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

func parseMask(mask string) (net.IPMask, error) {
	maskIP := net.IPMask(net.ParseIP(mask).To4())
	if maskIP == nil {
		return net.IPMask{}, fmt.Errorf("mask=%s: %w", mask, errInvalidMask)
	}

	return maskIP, nil
}

func joinIPMask(dirtyIP, dirtyMask string) (net.IP, error) {
	ip, err := parseIP(dirtyIP)
	if err != nil {
		return nil, err
	}

	mask, err := parseMask(dirtyMask)
	if err != nil {
		return nil, err
	}

	return ip.Mask(mask), nil
}

func (s *Service) ParseMaskedIP(ip, mask string) (net.IP, error) {
	return joinIPMask(ip, mask)
}

func (s *Service) ParseIP(ip string) (net.IP, error) {
	return parseIP(ip)
}

func (s *Service) IPToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}
