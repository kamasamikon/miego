package node

// Data Format
// XMP=xmap, JSN=JSON, YML=YAML, XML=XML, BIN=bytes, NON=NO-DATA
var DF_XMP uint
var DF_XMP_NAME = "XMP"
var DF_JSN uint
var DF_JSN_NAME = "JSN"
var DF_YML uint
var DF_YML_NAME = "YML"
var DF_XML uint
var DF_XML_NAME = "XML"
var DF_STR uint
var DF_STR_NAME = "STR"
var DF_BIN uint
var DF_BIN_NAME = "BIN"
var DF_NON uint
var DF_NON_NAME = "NON"

// NT_ANY: Means match to ANY(*) node type.
// NT_BAD: Means an invalid node type
var NT_ANY uint
var NT_ANY_NAME = "NT_ANY"
var NT_BAD uint
var NT_BAD_NAME = "NT_BAD"

// HT_ANY: Match ANY(*) hint
// HT_BAD: An invalid hint
var HT_ANY uint
var HT_ANY_NAME = "HT_ANY"
var HT_BAD uint
var HT_BAD_NAME = "HT_BAD"

// HT_OK: Normal success
// HT_NG: Error: Normal failure, Please don't use the returned data.
// HT_IG: The processor do nothing about this data
// HT_CN: A CANCEL flag is set when process
// HT_BP: Error: Some BAD PARAMETER found.
var HT_OK uint
var HT_OK_NAME = "HT_OK"
var HT_NG uint
var HT_NG_NAME = "HT_NG"
var HT_IG uint
var HT_IG_NAME = "HT_IG"
var HT_CN uint
var HT_CN_NAME = "HT_CN"
var HT_BP uint
var HT_BP_NAME = "HT_BP"

func init() {
	NT_ANY = NpAdd(NT_ANY_NAME)
	NT_BAD = NpAdd(NT_BAD_NAME)

	HT_ANY = NpAdd(HT_ANY_NAME)
	HT_BAD = NpAdd(HT_BAD_NAME)

	HT_OK = NpAdd(HT_OK_NAME)
	HT_NG = NpAdd(HT_NG_NAME)
	HT_IG = NpAdd(HT_IG_NAME)
	HT_CN = NpAdd(HT_CN_NAME)
	HT_BP = NpAdd(HT_BP_NAME)

	DF_XMP = NpAdd(DF_XMP_NAME)
	DF_JSN = NpAdd(DF_JSN_NAME)
	DF_YML = NpAdd(DF_YML_NAME)
	DF_XML = NpAdd(DF_XML_NAME)
	DF_STR = NpAdd(DF_STR_NAME)
	DF_BIN = NpAdd(DF_BIN_NAME)
	DF_NON = NpAdd(DF_NON_NAME)
}
