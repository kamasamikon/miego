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

func IndexOf(s string) int {
	for i, x := range os.Args {
		if x == s {
			return i
		}
	}
	return -1
}

// Parse 解析 os.Args
// flags: -flag 形式，值为 true
// opts:  --opt 形式，同名 opt 的值追加到切片
// args:  其余参数，含 "--" 之后的所有内容
func Parse() (flags map[string]bool, opts map[string][]string, args []string) {
	flags = make(map[string]bool)
	opts = make(map[string][]string)
	args = []string{}

	argv := os.Args[1:]
	afterDoubleDash := false

	for i := 0; i < len(argv); i++ {
		arg := argv[i]

		// 已进入 "--" 之后，全部是 args
		if afterDoubleDash {
			args = append(args, arg)
			continue
		}

		// 分隔符 "--"
		if arg == "--" {
			afterDoubleDash = true
			continue
		}

		// opt：-- 开头
		if strings.HasPrefix(arg, "--") {
			body := arg[2:]

			if body == "" {
				// 理论上不会到这里（"--" 已处理），保险起见丢弃
				continue
			}

			if idx := strings.Index(body, "="); idx >= 0 {
				// --optA=optValueA 形式
				name := body[:idx]
				value := body[idx+1:]
				if name == "" {
					// --=value 形式：opt 名为空，丢弃
					continue
				}
				opts[name] = append(opts[name], value)
			} else {
				// --optA optValueA 形式
				name := body
				if i+1 < len(argv) {
					opts[name] = append(opts[name], argv[i+1])
					i++ // 消耗下一个参数
				}
				// 没有下一个参数：非法值，丢弃（不加任何值）
			}
			continue
		}

		// flag：- 开头且长度 > 1
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			flags[arg[1:]] = true
			continue
		}

		// 其余（含单独的 "-"）都是 args
		args = append(args, arg)
	}

	return
}
