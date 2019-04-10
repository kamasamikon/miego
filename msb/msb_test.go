package msb

import (
	"encoding/json"
	"testing"
)

func TestMain(t *testing.T) {
	aa := KService{
		ServiceName: "a",
		Version:     "v2",
		IPAddr:      "aa.aa.aa.aa",
		Port:        11111,
	}

	ab := KService{
		ServiceName: "a",
		Version:     "v2",
		IPAddr:      "ab.ab.ab.ab",
		Port:        22222,
	}

	b := KService{
		ServiceName: "b",
		Version:     "va",
		IPAddr:      "b.b.b.b",
		Port:        22222,
	}

	var bin []byte
	bin, _ = json.Marshal(aa)
	msSet(bin)

	bin, _ = json.Marshal(ab)
	msSet(bin)

	bin, _ = json.Marshal(b)
	msSet(bin)

	nginxConfWrite()
}
