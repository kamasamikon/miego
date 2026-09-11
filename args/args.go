package args

import (
	"os"
	"strings"
)

// argGet -flagX -flagY -flagZ --keyX valueX --keyY valueY --- argX argY
// argGet(os.Args, "-flagX")
// argGet(os.Args, "--keyX")
// argGet(os.Args, "---")
//
// argGetFlag(args, "flagX")
func Flag(flag string) bool {
	key := "-" + flag

	i := 0
	size := len(os.Args)
	for {
		if i >= size {
			break
		}

		arg := os.Args[i]
		i++

		if arg == key {
			return true
		}
	}

	return false
}

func Opts(name string) (opts []string) {
	key := "--" + name
	key2 := "--" + name + "="

	hit := false
	i := 0
	size := len(os.Args)
	for {
		if i >= size {
			break
		}

		arg := os.Args[i]
		i++

		if hit {
			hit = false
			opts = append(opts, arg)
			continue
		}

		if arg == key {
			hit = true
			continue
		}

		if strings.HasPrefix(arg, key2) {
			_, after, _ := strings.Cut(arg, "=")
			opts = append(opts, after)
		}
	}

	return opts
}

// 我们只认第一个
func Opt(name string) (opt string, ok bool) {
	key := "--" + name
	key2 := "--" + name + "="

	ok = true

	hit := false
	i := 0
	size := len(os.Args)
	for {
		if i >= size {
			break
		}

		arg := os.Args[i]
		i++

		if hit {
			opt = arg
			return
		}

		if arg == key {
			hit = true
			continue
		}

		if strings.HasPrefix(arg, key2) {
			_, after, _ := strings.Cut(arg, "=")
			opt = after
			return
		}
	}

	ok = false
	return
}

// -f arg -key val arg -f -- arg arg arg
func Args() (args []string) {
	i := 1
	size := len(os.Args)
	for {
		if i >= size {
			break
		}

		arg := os.Args[i]
		i++

		if arg == "--" {
			args = append(args, os.Args[i:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			i += 1
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		args = append(args, arg)
	}

	return
}

// -f arg -key val arg -f -- arg arg arg
// 性能问题：先不管了
func Arg(index int) (arg string, ok bool) {
	var args []string
	i := 1
	size := len(os.Args)
	for {
		if i >= size {
			break
		}

		arg := os.Args[i]
		i++

		if arg == "--" {
			args = append(args, os.Args[i:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			i += 1
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		args = append(args, arg)
	}

	if len(args) == 0 {
		return
	}

	if index < 0 {
		index = len(args) + index
	}
	if index < 0 || index >= len(args) {
		return
	}

	arg = args[index]
	ok = true
	return
}

// os.Args[index]
func At(index int) (arg string, ok bool) {
	size := len(os.Args)

	if index < 0 {
		index = size + index
	}
	if index < 0 || index >= size {
		return
	}

	arg = os.Args[index]
	ok = true
	return
}
