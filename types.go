package main

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

type ResourceConfig struct {
	ModuleName string
	SinglePath string
	AllPath    string
}

func getResourceConfig(resourceType ResourceType) ResourceConfig {
	configs := map[ResourceType]ResourceConfig{
		VRF: {
			ModuleName: VRFModuleName,
			SinglePath: VRFPathSingle,
			AllPath:    VRFPathAll,
		},
		BGPVRF: {
			ModuleName: BGPVRFModuleName,
			SinglePath: BGPVRFPathSingle,
			AllPath:    BGPVRFPathAll,
		},
	}

	return configs[resourceType]
}
