package atox

import "testing"

func TestPrepare(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"", ""},
		{" ", ""},
		{"\t", ""},
		{"\r", ""},
		{"\n", ""},
		{"  \t\r\n  ", ""},
		{"0", "0"},
		{"1", "1"},
		{"+1", "+1"},
		{"-1", "-1"},
		{"0x1", "0x1"},
		{"0X1", "0X1"},
		{"0001", "1"},
		{"0010", "10"},
		{"+0010", "+10"},
		{"-0010", "-10"},
		{"0x0010", "0x10"},
		{"0X0010", "0X10"},
		{"0x1000", "0x1000"},
		{"0X1000", "0X1000"},
	}

	for _, test := range tests {
		if got := Prepare(test.input); got != test.want {
			t.Errorf("Prepare(%q) = %q, want %q", test.input, got, test.want)
		}
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		input string
		def   int64
		want  int64
	}{
		{"", 0, 0},
		{" ", 0, 0},
		{"\t", 0, 0},
		{"\r", 0, 0},
		{"\n", 0, 0},
		{"  \t\r\n  ", 0, 0},
		{"0", 0, 0},
		{"1", 0, 1},
		{"+1", 0, 1},
		{"-1", 0, -1},
		{"0x1", 0, 1},
		{"0X1", 0, 1},
		{"0001", 0, 1},
		{"0010", 0, 10},
		{"+0010", 0, 10},
		{"-0010", 0, -10},
		{"0x0010", 0, 16},
		{"0X0010", 0, 16},
		{"0x1000", 0, 4096},
		{"0X1000", 0, 4096},
		{"", 123, 123},
		{"abc", 123, 123},
	}

	for _, test := range tests {
		if got := Int64(test.input, test.def); got != test.want {
			t.Errorf("Int64(%q, %d) = %d, want %d", test.input, test.def, got, test.want)
		}
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		input string
		def   uint64
		want  uint64
	}{
		{"", 0, 0},
		{" ", 0, 0},
		{"\t", 0, 0},
		{"\r", 0, 0},
		{"\n", 0, 0},
		{"  \t\r\n  ", 0, 0},
		{"0", 0, 0},
		{"1", 0, 1},
		{"+1", 0, 1},
		{"-1", 0, 0},
		{"0x1", 0, 1},
		{"0X1", 0, 1},
		{"0001", 0, 1},
		{"0010", 0, 10},
		{"+0010", 0, 10},
		{"-0010", 0, 0},
		{"0x0010", 0, 16},
		{"0X0010", 0, 16},
		{"0x1000", 0, 4096},
		{"0X1000", 0, 4096},
		{"", 123, 123},
		{"abc", 123, 123},
	}

	for _, test := range tests {
		if got := Uint64(test.input, test.def); got != test.want {
			t.Errorf("Uint64(%q, %d) = %d, want %d", test.input, test.def, got, test.want)
		}
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		input string
		def   int
		want  int
	}{
		{"", 0, 0},
		{" ", 0, 0},
		{"\t", 0, 0},
		{"\r", 0, 0},
		{"\n", 0, 0},
		{"  \t\r\n  ", 0, 0},
		{"0", 0, 0},
		{"1", 0, 1},
		{"+1", 0, 1},
		{"-1", 0, -1},
		{"0x1", 0, 1},
		{"0X1", 0, 1},
		{"0001", 0, 1},
		{"0010", 0, 10},
		{"+0010", 0, 10},
		{"-0010", 0, -10},
		{"0x0010", 0, 16},
		{"0X0010", 0, 16},
		{"0x1000", 0, 4096},
		{"0X1000", 0, 4096},
		{"", 123, 123},
		{"abc", 123, 123},
	}

	for _, test := range tests {
		if got := Int(test.input, test.def); got != test.want {
			t.Errorf("Int(%q, %d) = %d, want %d", test.input, test.def, got, test.want)
		}
	}
}
