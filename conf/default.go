package conf

//////////////////////////////////////////////////////////////////////////
// Values

// Bool

func BFlip(key string) {
	Default.BFlip(key)
}
func BGet(key string, vdef bool) bool {
	return Default.BGet(key, vdef)
}
func B(key string, vdef bool) bool {
	return Default.B(key, vdef)
}
func BTrue(key string) bool {
	return Default.BTrue(key)
}
func BFalse(key string) bool {
	return Default.BFalse(key)
}
func BGetb(key string) (bool, bool) {
	return Default.BGetb(key)
}
func BHas(key string) bool {
	return Default.BHas(key)
}
func BMonitorAdd(key string, cbName string, cb bMonitor) string {
	return Default.BMonitorAdd(key, cbName, cb)
}
func BMonitorRem(key string, cbName string) {
	Default.BMonitorRem(key, cbName)
}
func BRem(key string) {
	Default.BRem(key)
}
func BSet(key string, val bool) {
	Default.BSet(key, val)
}
func BSetf(key string, val bool) {
	Default.BSetf(key, val)
}

// Event

func EHas(key string) bool {
	return Default.EHas(key)
}
func EAdd(key string) {
	Default.EAdd(key)
}
func EListenerAdd(key string, cbName string, cb eListener) string {
	return Default.EListenerAdd(key, cbName, cb)
}
func EListenerRem(key string, cbName string) {
	Default.EListenerRem(key, cbName)
}
func ERem(key string) {
	Default.ERem(key)
}
func ESend(key string, arg any) {
	Default.ESend(key, arg)
}
func ESendf(key string, arg any) {
	Default.ESendf(key, arg)
}

// Int

func IGet(key string, vdef int64) int64 {
	return Default.IGet(key, vdef)
}
func I(key string, vdef int64) int64 {
	return Default.I(key, vdef)
}
func IGetb(key string) (int64, bool) {
	return Default.IGetb(key)
}
func IHas(key string) bool {
	return Default.IHas(key)
}
func IInc(key string, inc int64) int64 {
	return Default.IInc(key, inc)
}
func IMonitorAdd(key string, cbName string, cb iMonitor) string {
	return Default.IMonitorAdd(key, cbName, cb)
}
func IMonitorRem(key string, cbName string) {
	Default.IMonitorRem(key, cbName)
}
func IRem(key string) {
	Default.IRem(key)
}
func ISet(key string, val any) {
	Default.ISet(key, val)
}
func ISetf(key string, val any) {
	Default.ISetf(key, val)
}

// String

func SGet(key string, vdef string) string {
	return Default.SGet(key, vdef)
}
func S(key string) string {
	return Default.S(key)
}
func SGetb(key string) (string, bool) {
	return Default.SGetb(key)
}
func SHas(key string) bool {
	return Default.SHas(key)
}
func SMonitorAdd(key string, cbName string, cb sMonitor) string {
	return Default.SMonitorAdd(key, cbName, cb)
}
func SMonitorRem(key string, cbName string) {
	Default.SMonitorRem(key, cbName)
}
func SRem(key string) {
	Default.SRem(key)
}
func SSet(key string, val string) {
	Default.SSet(key, val)
}
func SSetf(key string, val string) {
	Default.SSetf(key, val)
}
func SSplit(key string, sep string) []string {
	return Default.SSplit(key, sep)
}

// Execture

func XHas(key string) bool {
	return Default.XHas(key)
}
func XAdd(key string) {
	Default.XAdd(key)
}
func XSet(key string, arg any) {
	Default.XSet(key, arg)
}
func XGet(key string) any {
	return Default.XGet(key)
}
func XRem(key string) {
	Default.XRem(key)
}
func XSetSetter(key string, setter xSetter) {
	Default.XSetSetter(key, setter)
}
func XSetGetter(key string, getter xGetter) {
	Default.XSetGetter(key, getter)
}

// ///////////////////////////////////////////////////////////////////////
// Entries

func EntryRem(path string) {
	Default.EntryRem(path)
}

func EntryAddByLine(line string, overwrite bool) {
	Default.EntryAddByLine(line, overwrite)
}

func EntryAdd(path string, value string, overwrite bool) {
	Default.EntryAdd(path, value, overwrite)
}

// ///////////////////////////////////////////////////////////////////////
// Load

func LoadFromText(text string, overwrite bool) {
	Default.LoadFromText(text, overwrite)
}
func LoadFromFile(fileName string, overwrite bool) error {
	return Default.LoadFromFile(fileName, overwrite)
}
func LoadFromEnv() {
	Default.LoadFromEnv()
}
func LoadFromArg() {
	Default.LoadFromArg()
}

// ///////////////////////////////////////////////////////////////////////
// Dump

func Dump(joinBy string) string {
	return Default.Dump(joinBy)
}
func DumpMap() map[string]string {
	return Default.DumpMap()
}
func DumpRaw(joinBy string) string {
	return Default.DumpRaw(joinBy)
}

func Raw(path string) (string, bool) {
	return Default.Raw(path)
}

// ///////////////////////////////////////////////////////////////////////
// Others
func Names() []string {
	return Default.Names()
}
func OnReady(cb func()) {
	Default.OnReady(cb)
}
func Ready() {
	Default.Ready()
}
func Go() {
	Default.Go()
}
