package main

const (
	VRFPathSingleSpecific = "vrfs/vrf[vrf-name=%s]/vrf-name"
	VRFPathAllSpecific    = "vrfs/vrf[vrf-name=*]/vrf-name"
	VRFPathSingle         = "vrfs/vrf[vrf-name=%s]"
	VRFPathAll            = "vrfs/vrf"

	BGPVRFPathSingleSpecific = "router/bgp/as[as-number=49586]/vrfs/vrf[vrf-name=%s]/vrf-name"
	BGPVRFPathAllSpecific    = "router/bgp/as[as-number=49586]/vrfs/vrf[vrf-name=*]/vrf-name"
	BGPVRFPathSingle         = "router/bgp/as[as-number=49586]/vrfs/vrf[vrf-name=%s]"
	BGPVRFPathAll            = "router/bgp/as[as-number=49586]/vrfs"

	VRFModuleName    = "Cisco-IOS-XR-um-vrf-cfg"
	BGPVRFModuleName = "Cisco-IOS-XR-um-router-bgp-cfg"
)

type ResourceType string

const (
	VRF    ResourceType = "vrf"
	BGPVRF ResourceType = "bgpvrf"
)

type FetchType string

const (
	Single FetchType = "single"
	All    FetchType = "all"
)

type QueryType string

const (
	Full     QueryType = "full"
	Specific QueryType = "specific"
)

type ResourceConfig struct {
	ModuleName         string
	SinglePath         string
	AllPath            string
	SinglePathSpecific string
	AllPathSpecific    string
}

func getResourceConfig(resourceType ResourceType) ResourceConfig {
	configs := map[ResourceType]ResourceConfig{
		VRF: {
			ModuleName:         VRFModuleName,
			SinglePath:         VRFPathSingle,
			AllPath:            VRFPathAll,
			SinglePathSpecific: VRFPathSingleSpecific,
			AllPathSpecific:    VRFPathAllSpecific,
		},
		BGPVRF: {
			ModuleName:         BGPVRFModuleName,
			SinglePath:         BGPVRFPathSingle,
			AllPath:            BGPVRFPathAll,
			SinglePathSpecific: BGPVRFPathSingleSpecific,
			AllPathSpecific:    BGPVRFPathAllSpecific,
		},
	}

	return configs[resourceType]
}
