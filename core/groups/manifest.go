package groups

type GroupManifest struct {
	path string
	data *GroupManifestData
}

type GroupManifestData struct {
	Version string            `yaml:"version"`
	Groups  map[string]*Group `yaml:"groups"`
}
