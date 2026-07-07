package quadlet

type Model struct {
	Pods       []Pod
	Containers []Container
	Volumes    []Volume
}

type Pod struct {
	Name        string // Nom logique
	Filename    string // nextcloud
	Path        string
	Description string
}

type Container struct {
	Name        string // Nom logique
	Filename    string // nextcloud-app
	Path        string
	Description string
	Pod         string

	File UnitFile
}

type Volume struct {
	Name string
	Path string
}
