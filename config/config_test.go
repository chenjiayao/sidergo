package config

import (
	"bytes"
	"testing"
)

func Test_parseConfig(t *testing.T) {
	config := `

bind 0.0.0.0
port 6399
maxclients 128

appendonly yes
appendfilename appendonly.aof 
`
	var buf bytes.Buffer
	buf.Write([]byte(config))

	c := parseConfig(&buf)

	gotBind := c.Bind
	wantBind := "0.0.0.0"
	if gotBind != wantBind {
		t.Errorf("loadConfig bind = %s, want = %s", gotBind, wantBind)
	}

	gotAppendonly := c.Appendonly
	if !gotAppendonly {
		t.Errorf("loadConfig bind = %t, want = %t", gotAppendonly, true)
	}
}

func Test_loadConfig(t *testing.T) {

	config := `

bind 0.0.0.0
port 6399
maxclients 128

appendonly yes
appendfilename appendonly.aof 
`
	var buf bytes.Buffer
	buf.Write([]byte(config))

	configMap := loadConfig(&buf)

	if len(configMap) != 5 {
		t.Errorf("len(loadConfig) = %d, want = %d", len(configMap), 5)

	}

	gotBind := configMap["bind"]
	wantBind := "0.0.0.0"
	if gotBind != wantBind {
		t.Errorf("loadConfig bind = %s, want = %s", gotBind, wantBind)
	}

	gotAppendonly := configMap["appendonly"]
	wantAppendonly := "yes"
	if gotAppendonly != wantAppendonly {
		t.Errorf("loadConfig bind = %s, want = %s", gotAppendonly, wantAppendonly)
	}
}
