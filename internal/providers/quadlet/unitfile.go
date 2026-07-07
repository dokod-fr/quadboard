package quadlet

type UnitFile map[string]Section

type Section map[string][]string

func (u UnitFile) Section(name string) Section {
	return u[name]
}

func (s Section) First(key string) string {
	values := s[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
