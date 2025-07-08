package scan

import (
	"fmt"
	"math/bits"
	"net/netip"
	"strconv"
	"strings"
)

func expandTargetsForScan(targets []string) ([]string, error) {
	expanded := make([]string, 0, len(targets))
	for _, target := range targets {
		value := strings.TrimSpace(target)
		if value == "" {
			continue
		}

		addrs, ok, err := parseIPv4RangePair(value)
		if err != nil {
			return nil, err
		}
		if ok {
			expanded = append(expanded, addrs...)
			continue
		}

		expanded = append(expanded, value)
	}

	return expanded, nil
}

// Parse masscan range formats.
func parseIPv4RangePair(target string) ([]string, bool, error) {
	parts := strings.SplitN(target, "-", 2)
	if len(parts) != 2 {
		return nil, false, nil
	}

	startValue := strings.TrimSpace(parts[0])
	endValue := strings.TrimSpace(parts[1])
	if startValue == "" || endValue == "" {
		return nil, false, nil
	}

	// Fall through if this is in nmap range format.
	if endIsOctet(endValue) {
		return nil, false, nil
	}

	startAddr, startOK := parseIPv4Addr(startValue)
	endAddr, endOK := parseIPv4Addr(endValue)
	if !startOK && !endOK { // Allows the case where the target is just a hostname with a dash.
		return nil, false, nil
	}
	if !startOK || !endOK { // Prevents the case where one is an address and the other part is not.
		return nil, false, fmt.Errorf("invalid range %q", target)
	}

	startAddr = startAddr.Unmap()
	endAddr = endAddr.Unmap()
	if !startAddr.Is4() || !endAddr.Is4() {
		return nil, true, fmt.Errorf("invalid range %q", target)
	}

	start := ipv4ToUint32(startAddr)
	end := ipv4ToUint32(endAddr)
	if start > end {
		return nil, true, fmt.Errorf("invalid range %q", target)
	}

	return expandIPv4RangeToTargets(start, end), true, nil
}

func parseIPv4Addr(value string) (netip.Addr, bool) {
	addr, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Addr{}, false
	}
	return addr, true
}

func endIsOctet(value string) bool {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return false
	}
	return parsed >= 0 && parsed <= 255
}

func expandIPv4RangeToTargets(start, end uint32) []string {
	if start > end {
		return nil
	}

	const maxUint32 = uint64(^uint32(0))
	remaining := uint64(end) - uint64(start) + 1
	results := make([]string, 0, 16)

	for current := uint64(start); remaining > 0; {
		if current > maxUint32 {
			return results
		}

		current32 := uint32(current)
		maxSize := uint64(1) << bits.TrailingZeros32(current32)
		for maxSize > remaining {
			maxSize >>= 1
		}

		prefixLen := 32 - (bits.Len64(maxSize) - 1)
		addr := uint32ToIPv4(current32)
		if maxSize == 1 {
			results = append(results, addr.String())
		} else {
			results = append(results, fmt.Sprintf("%s/%d", addr.String(), prefixLen))
		}

		current += maxSize
		remaining -= maxSize
	}

	return results
}

func ipv4ToUint32(addr netip.Addr) uint32 {
	value := addr.As4()
	return uint32(value[0])<<24 | uint32(value[1])<<16 | uint32(value[2])<<8 | uint32(value[3])
}

func uint32ToIPv4(value uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	})
}
