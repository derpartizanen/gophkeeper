package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/config"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func unsetGophEnv() {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GOPH_") {
			pair := strings.SplitN(env, "=", 2)

			_ = os.Unsetenv(pair[0])
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	unsetGophEnv()

	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}

func TestConfigFromEnv(t *testing.T) {
	_ = os.Setenv("GOPH_USERNAME", gophtest.Username)
	_ = os.Setenv("GOPH_PASSWORD", string(gophtest.Password))
	_ = os.Setenv("GOPH_ADDRESS", "192.168.0.10:8080")
	_ = os.Setenv("GOPH_CA_PATH", "/etc/ssl/root.crt")
	_ = os.Setenv("GOPH_VERBOSE", "1")

	t.Cleanup(unsetGophEnv)

	sat := config.New()

	snaps.MatchSnapshot(t, sat.String())
}
