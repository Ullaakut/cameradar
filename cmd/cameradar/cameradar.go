package main

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/Ullaakut/cameradar/v5"
	"github.com/urfave/cli/v3"
)

func runCameradar(ctx context.Context, cmd *cli.Command) error {
	var targets []netip.Prefix
	for _, tgt := range cmd.StringSlice(flagTargets) {
		target, err := netip.ParsePrefix(tgt)
		if err != nil {
			return fmt.Errorf("invalid target %q: %w", tgt, err)
		}

		targets = append(targets, target)
	}

	var ports []uint16
	for _, port := range cmd.UintSlice(flagPorts) {
		if port > 65535 {
			return fmt.Errorf("%d is out of range (0-65535)", port)
		}

		ports = append(ports, uint16(port))
	}

	var opts []cameradar.Option
	if cmd.IsSet(flagCustomRoutes) {
		opts = append(opts, cameradar.WithCustomRoutes(cmd.String(flagCustomRoutes)))
	}
	if cmd.IsSet(flagCustomCredentials) {
		opts = append(opts, cameradar.WithCustomCredentials(cmd.String(flagCustomCredentials)))
	}
	if cmd.IsSet(flagScanSpeed) {
		opts = append(opts, cameradar.WithScanSpeed(int(cmd.Int(flagScanSpeed))))
	}
	if cmd.IsSet(flagAttackInterval) {
		opts = append(opts, cameradar.WithAttackInterval(cmd.Duration(flagAttackInterval)))
	}
	if cmd.IsSet(flagTimeout) {
		opts = append(opts, cameradar.WithTimeout(cmd.Duration(flagTimeout)))
	}

	c, err := cameradar.New(targets, ports, opts...)
	if err != nil {
		return fmt.Errorf("creating scanner: %w", err)
	}

	return c.Run(ctx)
}
