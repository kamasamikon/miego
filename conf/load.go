package conf

import (
	"fmt"
	"os"
	"strings"
)

func (cc *ConfCenter) LoadFromText(text string, overwrite bool) {
	Lines := strings.Split(strings.Replace(text, "\r", "\n", -1), "\n")
	size := len(Lines)
	i := 0
	for {
		if i >= size {
			break
		}

		Line := Lines[i]
		i++

		neat := strings.TrimSpace(Line)
		if neat == "" || neat[0] == '#' {
			continue
		}
		segs := strings.SplitN(neat, "=", 2)
		if len(segs) < 2 {
			continue
		}
		path, value := segs[0], segs[1]
		if len(path) < 4 || path[1] != ':' {
			continue
		}

		if strings.HasPrefix(value, "<<") {
			multiLineTag := value[2:]

			var sb strings.Builder
			for {
				if i >= size {
					break
				}

				Line = Lines[i]
				i++

				if Line == multiLineTag {
					break
				}

				sb.WriteString(Line)
				sb.WriteRune('\n')
			}

			cc.EntryAdd(path, sb.String(), overwrite)
		} else {
			cc.EntryAdd(path, value, overwrite)
		}
	}
}

// Load : configure from a file.
func (cc *ConfCenter) LoadFromFile(fileName string, overwrite bool) error {
	const (
		NGName = "s:/conf/Load/NG/%d/Name=%s"
		NGWhy  = "s:/conf/Load/NG/%d/Why=%s"
		OKName = "s:/conf/Load/OK/%d=%s"
	)

	data, err := os.ReadFile(fileName)
	if err != nil {
		cc.EntryAddByLine(fmt.Sprintf(NGName, cc.loadNGCount, fileName), false)
		cc.EntryAddByLine(fmt.Sprintf(NGWhy, cc.loadNGCount, err.Error()), false)
		cc.loadNGCount++
		dp("Error:'%s', fileName:'%s'", err.Error(), fileName)
		return err
	}

	cc.EntryAddByLine(fmt.Sprintf(OKName, cc.loadOKCount, fileName), false)
	cc.loadOKCount++

	cc.LoadFromText(string(data), overwrite)
	return nil
}

// cc.Set : Modify or Add conf entry
func (cc *ConfCenter) LoadFromEnv() {
	{
		cfgList := os.Getenv("KCFG_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				cc.LoadFromFile(f, true)
			}
		}
	}

	{
		cfgList := os.Getenv("KCFG_QQQ_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				cc.LoadFromFile(f, true)
				os.Remove(f)
			}
		}
	}
}

func (cc *ConfCenter) LoadFromArg() {
	argc := len(os.Args)

	// --kfg abc.cfg --kfg=xyz.cfg
	{
		for i, argv := range os.Args {
			if argv == "--kfg" {
				i++
				if i < argc {
					cc.LoadFromFile(os.Args[i], true)
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg=") {
				f := argv[6:]
				if f != "" {
					cc.LoadFromFile(f, true)
				}
				continue
			}
		}
	}

	// --kfg-qqq abc.cfg --kfg-qqq=xyz.cfg
	{
		for i, argv := range os.Args {
			if argv == "--kfg-qqq" {
				i++
				if i < argc {
					cc.LoadFromFile(os.Args[i], true)
					os.Remove(os.Args[i])
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg-qqq=") {
				f := argv[6:]
				if f != "" {
					cc.LoadFromFile(f, true)
					os.Remove(f)
				}
			}
		}
	}

	// --kfg-item i:/abc=777 --kfg-item=s:/xyz=abc
	for i, argv := range os.Args {
		if argv == "--kfg-item" {
			i++
			if i < argc {
				cc.EntryAddByLine(os.Args[i], true)
			}
			continue
		}
		if strings.HasPrefix(argv, "--kfg-item=") {
			item := argv[11:]
			cc.EntryAddByLine(item, true)
		}
	}
}
