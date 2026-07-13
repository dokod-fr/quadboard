package quadlet

import "testing"

func TestIsTemplate(t *testing.T) {
	if !isTemplate("redis@.container") {
		t.Fatal("expected template")
	}

	if isTemplate("redis@nextcloud.container") {
		t.Fatal("instance should not be a template")
	}

	if isTemplate("grafana.container") {
		t.Fatal("regular container should not be a template")
	}
}

func TestIsDropIn(t *testing.T) {
	if !isDropIn("/tmp/redis.container.d/10-env.conf") {
		t.Fatal("expected drop-in")
	}

	if isDropIn("/tmp/redis.container") {
		t.Fatal("regular container should not be a drop-in")
	}
}
