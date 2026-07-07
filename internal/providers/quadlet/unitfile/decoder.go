package unitfile

func parseKeyValue(f *File, section, key, value string) {
	switch section {

	case "Unit":

		switch key {
		case "Description":
			f.Unit.Description = value

		case "Documentation":
			f.Unit.Documentation = value
		}

	case "Container":

		switch key {

		case "ContainerName":
			f.Container.ContainerName = value

		case "Image":
			f.Container.Image = value

		case "Pod":
			f.Container.Pod = value

		case "Environment":
			f.Container.Environment = append(f.Container.Environment, value)

		case "Volume":
			f.Container.Volume = append(f.Container.Volume, value)

		case "Label":
			f.Container.Label = append(f.Container.Label, value)

		case "Network":
			f.Container.Network = append(f.Container.Network, value)
		}

	case "Pod":

		switch key {

		case "PodName":
			f.Pod.PodName = value

		case "Network":
			f.Pod.Network = append(f.Pod.Network, value)
		}
	}
}
