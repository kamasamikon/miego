package node

// Hint: 类似于状态

var HT_ANY uint64
var HT_ANY_NAME = "HT_ANY"

var HT_NIL uint64
var HT_NIL_NAME = "HT_NIL"

var HT_OK uint64
var HT_OK_NAME = "HT_OK"

var HT_NG uint64
var HT_NG_NAME = "HT_NG"

var HT_IG uint64
var HT_IG_NAME = "HT_IG"

var HT_CN uint64
var HT_CN_NAME = "HT_CN"

func init() {
	HT_ANY = NpAdd(HT_ANY_NAME) // 匹配任意类型的数据，不要直接使用
	HT_NIL = NpAdd(HT_NIL_NAME) // 不存在的，非法的类型

	HT_OK = NpAdd(HT_OK_NAME) // HT_OK: Normal 成功，可以继续处理
	HT_NG = NpAdd(HT_NG_NAME) // HT_NG: Error: 发生错误，就不要继续调用了？（饶是如此，Tracer的功能怎么办？）
	HT_IG = NpAdd(HT_IG_NAME) // HT_IG: Ignore: 这个数据没有被Process处理。
	HT_CN = NpAdd(HT_CN_NAME) // HT_CN: A CANCEL flag is set when process
}
