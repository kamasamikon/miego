package conf

func EntryRem(path string) {
	Default.EntryRem(path)
}

func EntryAddByLine(line string, overwrite bool) {
	Default.EntryAddByLine(line, overwrite)
}

func EntryAdd(path string, value string, overwrite bool) {
	Default.EntryAdd(path, value, overwrite)
}

func SetSetter(path string, setter setter) {
	Default.SetSetter(path, setter)
}

func SetGetter(path string, getter getter) {
	Default.SetGetter(path, getter)
}

func LoadFromText(text string, overwrite bool) {
	Default.LoadFromText(text, overwrite)
}

func LoadFromFile(fileName string, overwrite bool) error {
	return Default.LoadFromFile(fileName, overwrite)
}

func Ref(path string) (int64, int64) {
	return Default.Ref(path)
}

func Has(path string) bool {
	return Default.Has(path)
}

func Raw(name string) (string, bool) {
	return Default.Raw(name)
}

func Flip(path string) bool {
	return Default.Flip(path)
}

func Str(defval string, paths ...string) string {
	return Default.Str(defval, paths...)
}

func StrX(paths ...string) (string, bool) {
	return Default.StrX(paths...)
}

func Bool(defval bool, paths ...string) bool {
	return Default.Bool(defval, paths...)
}

func BoolX(paths ...string) (bool, bool) {
	return Default.BoolX(paths...)
}

func Obj(defval any, paths ...string) any {
	return Default.Obj(defval, paths...)
}

func ObjX(paths ...string) (any, bool) {
	return Default.ObjX(paths...)
}

func List(sep string, paths ...string) []string {
	return Default.List(sep, paths...)
}

func Names() []string {
	return Default.Names()
}

func SafeNames() []string {
	return Default.SafeNames()
}

func Set(path string, value any, create bool) {
	Default.Set(path, value, create)
}

func Ready() {
	Default.Ready()
}

func LoadFromEnv() {
	Default.LoadFromEnv()
}

func LoadFromArg() {
	Default.LoadFromArg()
}

func OnReady(cb func()) {
	Default.OnReady(cb)
}

func Go() {
	Default.Go()
}

func Dump(safeMode bool, joinBy string) string {
	return Default.Dump(safeMode, joinBy)
}
func DumpJson(safeMode bool) map[string]string {
	return Default.DumpJson(safeMode)
}
func DumpRaw(safeMode bool, group bool, joinBy string) string {
	return Default.DumpRaw(safeMode, group, joinBy)
}
