package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmic/pkg/api"
	"github.com/openconfig/ygot/ygot"
)

const (
	VRFPathSingle = "vrfs/vrf[vrf-name=%s]"
	VRFPathAll    = "vrfs/vrf"

	BGPVRFPathSingle = "router/bgp/as[as-number=49586]/vrfs/vrf[vrf-name=%s]"
	BGPVRFPathAll    = "router/bgp/as[as-number=49586]/vrfs"

	VRFModuleName    = "Cisco-IOS-XR-um-vrf-cfg"
	BGPVRFModuleName = "Cisco-IOS-XR-um-router-bgp-cfg"
)

type CLI struct {
	Host         string       `help:"Which host to benchmark" placeholder:"IP:PORT" required:""`
	Username     string       `help:"Username for authentication" required:""`
	Password     string       `help:"Password for authentication" required:""`
	CertPath     string       `help:"Path to certificate file"`
	Runs         int          `help:"Number of times to run" default:"1"`
	ResourceType ResourceType `help:"Type of resource to benchmark" enum:"vrf,bgpvrf" default:"vrf"`
	FetchType    FetchType    `help:"Type of fetch to perform" enum:"single,all" default:"single"`
	ResourceName string       `help:"Name of resource to benchmark" placeholder:"NAME"`
}

func (c *CLI) AfterApply() error {
	if c.FetchType == Single && c.ResourceName == "" {
		return fmt.Errorf("--resource-name is required when --fetch-type=single")
	}

	return nil
}

func main() {
	var cli CLI
	ctx := kong.Parse(&cli)

	config := getResourceConfig(cli.ResourceType)

	moduleName := config.ModuleName
	var path string
	switch cli.FetchType {
	case Single:
		path = fmt.Sprintf(config.SinglePath, cli.ResourceName)
	case All:
		path = config.AllPath
	}

	if err := runBenchmark(&cli, path, moduleName); err != nil {
		ctx.Fatalf("failed to run benchmark: %v", err)
	}
}

func runBenchmark(cli *CLI, path, moduleName string) error {
	ctx := context.Background()

	targetOptions := []api.TargetOption{
		api.Address(cli.Host),
		api.Username(cli.Username),
		api.Password(cli.Password),
		api.Insecure(true),
	}

	if cli.CertPath != "" {
		targetOptions = append(targetOptions, api.TLSCA(cli.CertPath))
	} else {
		targetOptions = append(targetOptions, api.Insecure(true))
	}

	target, err := api.NewTarget(targetOptions...)
	if err != nil {
		return fmt.Errorf("failed to create new target: %w", err)
	}

	if err := target.CreateGNMIClient(ctx); err != nil {
		return fmt.Errorf("failed to create new client from target: %w", err)
	}

	structuredPath, err := ygot.StringToStructuredPath(path)
	if err != nil {
		return fmt.Errorf("failed to create structured path from %s: %w", path, err)
	}
	structuredPath.Origin = moduleName

	getRequest := &gnmi.GetRequest{
		Type:     gnmi.GetRequest_ALL,
		Encoding: gnmi.Encoding_JSON_IETF,
		Path:     []*gnmi.Path{structuredPath},
	}

	start := time.Now()
	var totalTime time.Duration

	for range cli.Runs {
		iterStart := time.Now()

		_, err := target.Get(ctx, getRequest)
		if err != nil {
			return fmt.Errorf("failed to get resource(s) from target: %w", err)
		}

		iterElapsed := time.Since(iterStart)
		totalTime += iterElapsed
	}

	elapsed := time.Since(start)
	fmt.Printf("Performed %d Get operations in %s (avg: %s per call)\n", cli.Runs, elapsed, totalTime/time.Duration(cli.Runs))

	return nil
}
