package unitfile

type File struct {
	Unit      Unit
	Service   Service
	Container Container
	Pod       Pod
	Install   Install
}

type Unit struct {
	Description   string
	Documentation string
}

type Service struct {
	Restart string
}

type Container struct {
	ContainerName string
	Image         string
	Pod           string

	Environment []string
	Volume      []string
	Label       []string
	Network     []string
}

type Pod struct {
	PodName string
	Network []string
}

type Install struct {
	WantedBy []string
}
