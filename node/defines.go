package node

// Data Format
var DF_JSN int
var DF_JSN_NAME = "JSN"
var DF_YML int
var DF_YML_NAME = "YML"
var DF_XML int
var DF_XML_NAME = "XML"
var DF_STR int
var DF_STR_NAME = "STR"
var DF_BIN int
var DF_BIN_NAME = "BIN"
var DF_NON int
var DF_NON_NAME = "NON"

// NT_ANY: Means match to ANY(*) node type.
// NT_BAD: Means an invalid node type
var NT_ANY int
var NT_ANY_NAME = "NT_ANY"
var NT_BAD int
var NT_BAD_NAME = "NT_BAD"

// HT_ANY: Match ANY(*) hint
// HT_BAD: An invalid hint
var HT_ANY int
var HT_ANY_NAME = "HT_ANY"
var HT_BAD int
var HT_BAD_NAME = "HT_BAD"

// HT_OK: Normal success
// HT_NG: Error: Normal failure, Please don't use the returned data.
// HT_IG: The processor do nothing about this data
// HT_CN: A CANCEL flag is set when process
// HT_BP: Error: Some BAD PARAMETER found.
var HT_OK int
var HT_OK_NAME = "HT_OK"
var HT_NG int
var HT_NG_NAME = "HT_NG"
var HT_IG int
var HT_IG_NAME = "HT_IG"
var HT_CN int
var HT_CN_NAME = "HT_CN"
var HT_BP int
var HT_BP_NAME = "HT_BP"

func init() {
	NT_ANY = npAdd(NT_ANY_NAME)
	NT_BAD = npAdd(NT_BAD_NAME)

	HT_ANY = npAdd(HT_ANY_NAME)
	HT_BAD = npAdd(HT_BAD_NAME)

	HT_OK = npAdd(HT_OK_NAME)
	HT_NG = npAdd(HT_NG_NAME)
	HT_IG = npAdd(HT_IG_NAME)
	HT_CN = npAdd(HT_CN_NAME)
	HT_BP = npAdd(HT_BP_NAME)

	DF_JSN = npAdd(DF_JSN_NAME)
	DF_YML = npAdd(DF_YML_NAME)
	DF_XML = npAdd(DF_XML_NAME)
	DF_STR = npAdd(DF_STR_NAME)
	DF_BIN = npAdd(DF_BIN_NAME)
	DF_NON = npAdd(DF_NON_NAME)
}
