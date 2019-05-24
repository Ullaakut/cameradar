package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/ullaakut/cameradar"
	log "github.com/ullaakut/disgo"
	"github.com/ullaakut/disgo/style"
)

func parseArguments() error {
	viper.SetEnvPrefix("cameradar")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	pflag.StringSliceP("targets", "t", []string{}, "The targets on which to scan for open RTSP streams - required (ex: 172.16.100.0/24)")
	pflag.StringSliceP("ports", "p", []string{"554", "5554", "8554"}, "The ports on which to search for RTSP streams")
	pflag.StringP("custom-routes", "r", "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/routes", "The path on which to load a custom routes dictionary")
	pflag.StringP("custom-credentials", "c", "<GOPATH>/src/github.com/ullaakut/cameradar/dictionaries/credentials.json", "The path on which to load a custom credentials JSON dictionary")
	pflag.IntP("speed", "s", 4, "The nmap speed preset to use for discovery")
	pflag.IntP("timeout", "T", 2000, "The timeout in miliseconds to use for attack attempts")
	pflag.BoolP("log", "l", false, "Enable the logs for nmap's output to stdout")
	pflag.BoolP("help", "h", false, "displays this help message")

	viper.AutomaticEnv()

	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	if viper.GetBool("help") {
		pflag.Usage()
		fmt.Println("\nExamples of usage:")
		fmt.Println("\tScanning your home network for RTSP streams:\tcameradar -t 192.168.0.0/24")
		fmt.Println("\tScanning a remote camera on a specific port:\tcameradar -t 172.178.10.14 -p 18554 -s 2")
		fmt.Println("\tScanning an unstable remote network: \t\tcameradar -t 172.178.10.14/24 -s 1 --timeout 10000 -l")
		os.Exit(0)
	}

	if viper.GetStringSlice("targets") == nil {
		return errors.New("targets (-t, --targets) argument required\n    examples:\n      - 172.16.100.0/24\n      - localhost\n      - 8.8.8.8")
	}

	return nil
}

func main() {
	err := parseArguments()
	if err != nil {
		printErr(term, err)
	}

	c := cameradar.New(
		cameradar.WithTargets(viper.GetStringSlice("targets")),
		cameradar.WithPorts(viper.GetStringSlice("ports")),
		cameradar.WithDebug(viper.GetBool("log") || viper.GetBool("logging")),
		cameradar.WithCustomCredentials(viper.GetString("custom-credentials")),
		cameradar.WithCustomRoutes(viper.GetString("custom-routes")),
		cameradar.WithSpeed(viper.GetInt("speed")),
		cameradar.WithTimeout(viper.GetInt("timeout")*time.Millisecond),
	)

	scanResult, err := c.Scan()
	if err != nil {
		printErr(err)
	}

	streams, err := c.Attack(scanResult)
	if err != nil {
		printErr(err)
	}

	c.PrintStreams(streams)
}

func printErr(term *log.Terminal, err error) {
	term.Errorln(style.Failure(style.SymbolCross), err)
	os.Exit(1)
}
