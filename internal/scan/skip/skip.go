package skip

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
)

// Scanner is a stream scanner that skips discovery and treats every target/port as a stream.
type Scanner struct {
	targets []string
	ports   []string
}

// New builds a scanner that skips discovery and treats every target/port as a stream.
func New(targets, ports []string) *Scanner {
	return &Scanner{
		targets: targets,
		ports:   ports,
	}
}

// Scan returns the precomputed list of streams.
func (s *Scanner) Scan(ctx context.Context) ([]cameradar.Stream, error) {
	return buildStreamsFromTargets(ctx, s.targets, s.ports)
}

func buildStreamsFromTargets(ctx context.Context, targets, ports []string) ([]cameradar.Stream, error) {
	resolvedPorts, err := parsePorts(ctx, ports)
	if err != nil {
		return nil, err
	}
	if len(resolvedPorts) == 0 {
		return nil, errors.New("no valid ports provided")
	}

	resolvedTargets, err := expandTargets(ctx, targets)
	if err != nil {
		return nil, err
	}
	if len(resolvedTargets) == 0 {
		return nil, errors.New("no valid target addresses resolved")
	}

	streams := make([]cameradar.Stream, 0, len(resolvedTargets)*len(resolvedPorts))
	for _, addr := range resolvedTargets {
		for _, port := range resolvedPorts {
			streams = append(streams, cameradar.Stream{
				Address: addr,
				Port:    port,
			})
		}
	}

	return streams, nil
}

func parsePorts(ctx context.Context, ports []string) ([]uint16, error) {
	seen := make(map[uint16]struct{})
	resolved := make([]uint16, 0, len(ports))

	for _, entry := range ports {
		for raw := range strings.SplitSeq(entry, ",") {
			value := strings.TrimSpace(raw)
			if value == "" {
				continue
			}

			values, err := parsePortValue(ctx, value)
			if err != nil {
				return nil, err
			}

			for _, port := range values {
				if _, exists := seen[port]; exists {
					continue
				}
				seen[port] = struct{}{}
				resolved = append(resolved, port)
			}
		}
	}

	return resolved, nil
}

func parsePortValue(ctx context.Context, value string) ([]uint16, error) {
	if strings.Contains(value, "-") {
		parts := strings.SplitN(value, "-", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port range %q", value)
		}

		start, err := parsePortNumber(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid port range %q: %w", value, err)
		}
		end, err := parsePortNumber(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid port range %q: %w", value, err)
		}
		if start > end {
			return nil, fmt.Errorf("invalid port range %q", value)
		}

		ports := make([]uint16, 0, end-start+1)
		for port := start; port <= end; port++ {
			ports = append(ports, port)
		}
		return ports, nil
	}

	port, err := parsePortNumber(value)
	if err == nil {
		return []uint16{port}, nil
	}

	servicePort, lookupErr := net.DefaultResolver.LookupPort(ctx, "tcp", value)
	if lookupErr != nil {
		return nil, fmt.Errorf("invalid port %q", value)
	}
	if servicePort < 1 || servicePort > 65535 {
		return nil, fmt.Errorf("port %d out of range", servicePort)
	}
	return []uint16{uint16(servicePort)}, nil
}

func parsePortNumber(value string) (uint16, error) {
	port, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port %d out of range", port)
	}
	return uint16(port), nil
}

func expandTargets(ctx context.Context, targets []string) ([]netip.Addr, error) {
	seen := make(map[netip.Addr]struct{})
	resolved := make([]netip.Addr, 0, len(targets))

	for _, target := range targets {
		value := strings.TrimSpace(target)
		if value == "" {
			continue
		}

		addrs, err := parseTargetAddrs(ctx, value)
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if !addr.IsValid() {
				continue
			}
			if _, exists := seen[addr]; exists {
				continue
			}
			seen[addr] = struct{}{}
			resolved = append(resolved, addr)
		}
	}

	return resolved, nil
}

func parseTargetAddrs(ctx context.Context, target string) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(target)
	if err == nil { // Return early.
		return expandPrefix(prefix), nil
	}

	if strings.Contains(target, "-") {
		addrs, ok, err := parseIPv4Range(target)
		if ok {
			return addrs, err
		}
	}

	addr, err := netip.ParseAddr(target)
	if err == nil { // Return early.
		return []netip.Addr{addr}, nil
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, target)
	if err != nil {
		return nil, fmt.Errorf("resolving hostname %q: %w", target, err)
	}

	addrs := make([]netip.Addr, 0, len(ips))
	for _, ip := range ips {
		addr, ok := netip.AddrFromSlice(ip.IP)
		if !ok {
			continue
		}
		addrs = append(addrs, addr.Unmap())
	}

	if len(addrs) == 0 {
		return nil, fmt.Errorf("no ip addresses found for hostname %q", target)
	}

	return addrs, nil
}

func expandPrefix(prefix netip.Prefix) []netip.Addr {
	if !prefix.IsValid() {
		return nil
	}

	prefix = prefix.Masked()
	addr := prefix.Addr()
	addrs := make([]netip.Addr, 0, 16)

	for current := addr; prefix.Contains(current); {
		addrs = append(addrs, current)
		next := current.Next()
		if !next.IsValid() {
			break
		}
		current = next
	}

	return addrs
}

type octetRange struct {
	start int
	end   int
}

func parseIPv4Range(target string) ([]netip.Addr, bool, error) {
	parts := strings.Split(target, ".")
	if len(parts) != 4 {
		return nil, false, nil
	}

	ranges := make([]octetRange, 4)
	for i, part := range parts {
		parsed, ok, err := parseOctetRange(part)
		if err != nil {
			return nil, true, err
		}

		if !ok {
			return nil, false, nil
		}
		ranges[i] = parsed
	}

	addrs := make([]netip.Addr, 0, 16)
	for first := ranges[0].start; first <= ranges[0].end; first++ {
		for second := ranges[1].start; second <= ranges[1].end; second++ {
			for third := ranges[2].start; third <= ranges[2].end; third++ {
				for fourth := ranges[3].start; fourth <= ranges[3].end; fourth++ {
					addrs = append(addrs, netip.AddrFrom4([4]byte{
						byte(first),
						byte(second),
						byte(third),
						byte(fourth),
					}))
				}
			}
		}
	}

	return addrs, true, nil
}

func parseOctetRange(value string) (octetRange, bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return octetRange{}, false, nil
	}

	if strings.Contains(value, "-") {
		parts := strings.SplitN(value, "-", 2)
		if len(parts) != 2 {
			return octetRange{}, true, fmt.Errorf("invalid range %q", value)
		}

		start, err := parseOctetValue(strings.TrimSpace(parts[0]))
		if err != nil {
			return octetRange{}, true, err
		}
		end, err := parseOctetValue(strings.TrimSpace(parts[1]))
		if err != nil {
			return octetRange{}, true, err
		}
		if start > end {
			return octetRange{}, true, fmt.Errorf("invalid range %q", value)
		}

		return octetRange{start: start, end: end}, true, nil
	}

	if !isDigits(value) {
		return octetRange{}, false, nil
	}

	octet, err := parseOctetValue(value)
	if err != nil {
		return octetRange{}, true, err
	}

	return octetRange{start: octet, end: octet}, true, nil
}

func parseOctetValue(value string) (int, error) {
	if !isDigits(value) {
		return 0, fmt.Errorf("invalid octet %q", value)
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid octet %q", value)
	}
	if parsed < 0 || parsed > 255 {
		return 0, fmt.Errorf("octet %d out of range", parsed)
	}
	return parsed, nil
}

func isDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return value != ""
}
