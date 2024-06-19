package node

import (
	"miego/xmap"
)

type KNode struct {
	// Point back to Manager
	nm *KNodeManager

	//
	// User defined
	//
	Type   string // "NT_XXX"
	TypeId uint64 // NpAdd("NT_XXX")

	Name string
	Desc string

	// item: "hint@node;hint@node"
	// if match any, use * instead.
	Follows string

	// Processor will process the given data, and return the result
	// If no data generated, set output to nil.
	// Processor func(f *KCallFrame) (output []byte, datafmt uint, dataType uint)

	// result: output of this function, please convert according to format
	// format: the format of the result
	// hint: Why return this data?
	Processor func(f *KCallFrame) (result []byte, datafmt uint64, hint uint64)

	//
	// OnXXX
	//
	// OnBeforeReg   func(nm *KNodeManager, self *KNode) // Not in NM
	// OnAfterReg    func(nm *KNodeManager, self *KNode) // In NM
	// OnBeforeDereg func(nm *KNodeManager, self *KNode) // In NM
	// OnAfterDereg  func(nm *KNodeManager, self *KNode) // Not in NM
	OnStart func(nm *KNodeManager, self *KNode)
	// OnStop        func(nm *KNodeManager, self *KNode)

	UserData xmap.Map
}

func (n *KNode) SendToSubs(data []byte, datafmt uint64, hint uint64) {
	n.nm.sendtoSubs(nil, n, data, datafmt, hint)
}
