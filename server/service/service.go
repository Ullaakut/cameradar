package service

import (
	"fmt"
	"time"

	"github.com/EtixLabs/cameradar"
	"github.com/pkg/errors"
)

// Cameradar is the service in charge of communicating with the GUI
type Cameradar struct {
	Streams []cmrdr.Stream

	options *cmrdr.Options

	toClient   chan<- string
	fromClient <-chan string
}

// New instanciates a new Cameradar service
func New(routesFilePath, credentialsFilePath string, fromClient <-chan string, toClient chan<- string) (*Cameradar, error) {
	routes, err := cmrdr.LoadRoutes(routesFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't load routes dictionary")
	}

	credentials, err := cmrdr.LoadCredentials(credentialsFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't load credentials dictionary")
	}

	cameradar := &Cameradar{
		Streams: nil,
		options: &cmrdr.Options{
			Ports:       "554,8554",
			Routes:      routes,
			Credentials: credentials,
			OutputFile:  "/tmp/cameradar_nmap_result.xml",
			Speed:       4,
			Timeout:     2000,
		},

		fromClient: fromClient,
		toClient:   toClient,
	}

	go cameradar.Run()
	return cameradar, nil
}

// Run launches the service that will automatically call the service methods
// using the instructions received over websocket
func (c *Cameradar) Run() {
	for {
		msg := <-c.fromClient
		go c.handleRequest(msg)
	}
}

// Discover launches a Cameradar scan using the service's options
func (c *Cameradar) Discover() ([]cmrdr.Stream, error) {
	streams, err := cmrdr.Discover(c.options.Target, c.options.Ports, c.options.OutputFile, c.options.Speed, true)
	if err != nil {
		return streams, errors.Wrap(err, "could not discover streams")
	}
	c.Streams = streams
	return streams, nil
}

// AttackRoute launches a Cameradar route attack using the service's options
func (c *Cameradar) AttackRoute() ([]cmrdr.Stream, error) {
	streams, err := cmrdr.AttackRoute(c.Streams, c.options.Routes, c.options.Timeout, true)
	if err != nil {
		return streams, errors.Wrap(err, "could not discover streams")
	}
	c.Streams = streams
	return streams, nil
}

// AttackCredentials launches a Cameradar credential attack using the service's options
func (c *Cameradar) AttackCredentials() ([]cmrdr.Stream, error) {
	streams, err := cmrdr.AttackCredentials(c.Streams, c.options.Credentials, c.options.Timeout, true)
	if err != nil {
		return streams, errors.Wrap(err, "could not discover streams")
	}
	c.Streams = streams
	return streams, nil
}

// SetOptions sets all options using an option structure
func (c *Cameradar) SetOptions(options cmrdr.Options) {
	c.options.Target = options.Target

	if len(options.Ports) > 0 {
		c.options.Ports = options.Ports
	}

	if len(options.OutputFile) > 0 {
		c.options.OutputFile = options.OutputFile
	}

	// TODO: Add custom dictionary support through ws

	c.SetSpeed(options.Speed)

	c.SetTimeout(options.Timeout)
}

// SetNmapOutputFile sets the OutputFile option
func (c *Cameradar) SetNmapOutputFile(path string) {
	c.options.OutputFile = path
}

// SetRoutes overwrites the routes dictionary with new values
func (c *Cameradar) SetRoutes(routes string) {
	c.options.Routes = cmrdr.ParseRoutesFromString(routes)
}

// SetCredentials overwrites the routes dictionary with new values
func (c *Cameradar) SetCredentials(credentials string) error {
	newCredentials, err := cmrdr.ParseCredentialsFromString(credentials)
	if err != nil {
		return errors.Wrap(err, "could not decode credentials")
	}
	c.options.Credentials = newCredentials
	return nil
}

// SetSpeed sets the Speed option
func (c *Cameradar) SetSpeed(speed int) error {
	if speed < cmrdr.PARANOIAC || speed > cmrdr.INSANE {
		return fmt.Errorf("invalid speed value '%d'. should be between '%d' and '%d'", speed, cmrdr.PARANOIAC, cmrdr.INSANE)
	}
	c.options.Speed = speed
	return nil
}

// SetTimeout sets the Timeout option
func (c *Cameradar) SetTimeout(timeout time.Duration) error {
	if timeout < 0 {
		return fmt.Errorf("invalid timeout value '%d'. should be superior to 0", timeout)
	}
	c.options.Timeout = time.Millisecond * time.Duration(timeout)
	return nil
}
