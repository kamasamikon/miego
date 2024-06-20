package node

// NodeType:

var NT_ANY uint64
var NT_ANY_NAME = "NT_ANY"

var NT_NIL uint64
var NT_NIL_NAME = "NT_NIL"

func init() {
	NT_ANY = NpAdd(NT_ANY_NAME) // 匹配任意类型的节点，不要直接使用
	NT_NIL = NpAdd(NT_NIL_NAME) // 不存在的，非法的类型
}
