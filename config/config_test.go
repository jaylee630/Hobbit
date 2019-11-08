package config

import (
	"os"
	"testing"
)

func TestLoadYAMLConfig(t *testing.T) {

	var data = `
host: 192.168.102.42
port: 4242
dev_mode: true
mysql:
  host: 192.168.102.43
  port: 4343
`
	cfg := &BasicConfig{}
	err := loadYAMLConfig([]byte(data), cfg)
	if err != nil {
		t.Fatal(err.Error())
	}

	if cfg.Host != "192.168.102.42" {
		t.Fatalf("fail to load configs, excepted: host=192.168.102.42, result: host=%s", cfg.Host)
	}

	if cfg.MySQL.Port != 4343 {
		t.Fatalf("fail to load configs, excepted: mysql.port=4343, result: mysql.port=%d", cfg.MySQL.Port)
	}
}

func TestLoadEnvConfig(t *testing.T) {

	cfg := &BasicConfig{}
	prefix := "OWLNEST_UNITTEST"

	name1 := "_HOST"
	name2 := "_MYSQL_HOST"

	_ = os.Setenv(prefix+name1, "192.168.102.42")
	defer os.Unsetenv(prefix + name1)

	_ = os.Setenv(prefix+name2, "192.168.102.42")
	defer os.Unsetenv(prefix + name2)

	v, ok := os.LookupEnv(prefix + name1)
	t.Logf("set env %s%s=%s, state:%t", prefix, name1, v, ok)

	v, ok = os.LookupEnv(prefix + name2)
	t.Logf("set env %s%s=%s, state:%t", prefix, name2, v, ok)

	err := loadEnvConfig(prefix, cfg)
	if err != nil {
		t.Fatal(err.Error())
	}

	if cfg.Host != "192.168.102.42" {
		t.Fatalf("failed to read env, excepted: host=192.168.102.42, result: host=%s ", cfg.Host)
	}

	if cfg.MySQL.Host != "192.168.102.42" {
		t.Fatalf("failed to read env, excepted: mysql.host=192.168.102.42, result: host=%s ", cfg.Host)
	}

}
