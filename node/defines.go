package node

// Data Format
// XMP=xmap, JSN=JSON, YML=YAML, XML=XML, BIN=bytes, NON=NO-DATA
var DT_ANY uint64
var DT_ANY_NAME = "DT_ANY"
var DT_XMP uint64
var DT_XMP_NAME = "DT_XMP"
var DT_JSN uint64
var DT_JSN_NAME = "DT_JSN"
var DT_YML uint64
var DT_YML_NAME = "DT_YML"
var DT_XML uint64
var DT_XML_NAME = "DT_XML"
var DT_STR uint64
var DT_STR_NAME = "DT_STR"
var DT_BIN uint64
var DT_BIN_NAME = "DT_BIN"
var DT_NON uint64
var DT_NON_NAME = "DT_NON"

// NT_ANY: Means match to ANY(*) node type.
// NT_BAD: Means an invalid node type
var NT_ANY uint64
var NT_ANY_NAME = "NT_ANY"
var NT_BAD uint64
var NT_BAD_NAME = "NT_BAD"

// HT_ANY: Match ANY(*) hint
// HT_BAD: An invalid hint
var HT_ANY uint64
var HT_ANY_NAME = "HT_ANY"
var HT_BAD uint64
var HT_BAD_NAME = "HT_BAD"

// HT_OK: Normal 成功，可以继续处理
// HT_NG: Error: 发生错误，就不要继续调用了？（饶是如此，Tracer的功能怎么办？）
// HT_IG: Ignore: 这个数据没有被Process处理。
// HT_CN: A CANCEL flag is set when process
// HT_BP: Error: Some BAD PARAMETER found.
var HT_OK uint64
var HT_OK_NAME = "HT_OK"
var HT_NG uint64
var HT_NG_NAME = "HT_NG"
var HT_IG uint64
var HT_IG_NAME = "HT_IG"
var HT_CN uint64
var HT_CN_NAME = "HT_CN"
var HT_BP uint64
var HT_BP_NAME = "HT_BP"

func init() {
	// NodeType
	NT_ANY = NpAdd(NT_ANY_NAME)
	NT_BAD = NpAdd(NT_BAD_NAME)

	// HintType
	HT_ANY = NpAdd(HT_ANY_NAME)
	HT_BAD = NpAdd(HT_BAD_NAME)

	HT_OK = NpAdd(HT_OK_NAME)
	HT_NG = NpAdd(HT_NG_NAME)
	HT_IG = NpAdd(HT_IG_NAME)
	HT_CN = NpAdd(HT_CN_NAME)
	HT_BP = NpAdd(HT_BP_NAME)

	// DataFormat
	DT_XMP = NpAdd(DT_XMP_NAME)
	DT_JSN = NpAdd(DT_JSN_NAME)
	DT_YML = NpAdd(DT_YML_NAME)
	DT_XML = NpAdd(DT_XML_NAME)
	DT_STR = NpAdd(DT_STR_NAME)
	DT_BIN = NpAdd(DT_BIN_NAME)
	DT_NON = NpAdd(DT_NON_NAME)
}
