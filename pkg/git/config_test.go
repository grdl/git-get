package git

import "testing"

func TestConfig(t *testing.T) {
	r := testRepoEmpty(t)

	value := r.GetCfg("gitget.host")
	if value != "" {
		t.Errorf("wrong value, expected %q; got %q", "", value)
	}

	r.setCfg("gitget.host", "gitlab.com")

	value = r.GetCfg("gitget.host")
	if value != "gitlab.com" {
		t.Errorf("wrong value, expected %q; got %q", "gitlab.com", value)
	}
}
