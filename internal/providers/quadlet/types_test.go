package quadlet

import "testing"

func TestKind(t *testing.T) {
	tests := []struct {
		name string
		path string
		want Kind
	}{
		{"container", "grafana.container", ContainerKind},
		{"pod", "monitoring.pod", PodKind},
		{"volume", "grafana.volume", VolumeKind},
		{"unknown", "README.md", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kind(tt.path); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
