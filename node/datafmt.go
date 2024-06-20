package node

// Data Format: 和DataType不同，Format是装扮，Type是内在。

var DF_ANY uint64
var DF_ANY_NAME = "DF_ANY"

var DF_NIL uint64
var DF_NIL_NAME = "DF_NIL"

var DF_XMP uint64
var DF_XMP_NAME = "DF_XMP"

var DF_JSN uint64
var DF_JSN_NAME = "DF_JSN"

var DF_YML uint64
var DF_YML_NAME = "DF_YML"

var DF_XML uint64
var DF_XML_NAME = "DF_XML"

var DF_STR uint64
var DF_STR_NAME = "DF_STR"

var DF_BIN uint64
var DF_BIN_NAME = "DF_BIN"

func init() {
	// DataFormat
	DF_ANY = NpAdd(DF_ANY_NAME) // 匹配任意类型的数据，不要直接使用
	DF_NIL = NpAdd(DF_NIL_NAME) // 不存在的，非法的类型

	DF_XMP = NpAdd(DF_XMP_NAME) // miego.xmap，就是map
	DF_JSN = NpAdd(DF_JSN_NAME) // JSON 字符串
	DF_YML = NpAdd(DF_YML_NAME) // YAML 字符串
	DF_XML = NpAdd(DF_XML_NAME) // XML 字符串
	DF_STR = NpAdd(DF_STR_NAME) // 普通字符串
	DF_BIN = NpAdd(DF_BIN_NAME) // 二进制数据，等于bytes[]
}
