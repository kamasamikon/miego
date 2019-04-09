package node

import (
	"fmt"
	"github.com/kamasamikon/miego/klog"
	"time"
)

// "github.com/satori/go.uuid"

type KCallFrame struct {
	//
	// nodex.Processor(data, dateType, dataId)
	//

	// caller: upper layer frame
	// node: who process the data
	// data: data returned by upstreamNode.Processor()
	// dataFormat: int
	// dataId: GUID of the data if the data been saved to database.
	// createdAt: Timestamp this frame created.
	Caller     *KCallFrame
	Node       *KNode
	Data       []byte
	DataFormat int
	Hint       int
	DataId     string
	CreatedAt  int
}

func dumpFrames(f *KCallFrame) {
	var lines []string

	tmp := f
	for tmp != nil {
		lines = append(lines, fmt.Sprintf(" %-16s | %4s | %4s | %s", tmp.Node.Type, npName(tmp.DataFormat), npName(tmp.Hint), tmp.Node.Name))
		tmp = tmp.Caller
	}

	klog.D(">>>> DUMP (%d) >>>>", len(lines))
	cnt := len(lines)
	for i := range lines {
		klog.D("|%4d|%s", i, lines[cnt-i-1])
	}
	klog.D("<<<< DUMP (%d) <<<<", len(lines))
}

func NewCallFrame(caller *KCallFrame, thisNode *KNode, data []byte, dataFormat int, hint int) *KCallFrame {
	f := &KCallFrame{
		Caller:     caller,
		Node:       thisNode,
		Data:       data,
		DataFormat: dataFormat,
		Hint:       hint,
		DataId:     "TODO: Id of DB?",
		CreatedAt:  time.Now().Nanosecond(),
	}

	// TODO: Save the data to database;
	// TODO: Create dataId

	return f
}

type KNode struct {
	//
	// User defined
	//
	Type   string // "NT_XXX"
	typeId int    // npAdd("NT_XXX")

	Name string
	Desc string

	// item: "hint@node;hint@node"
	// if match any, use * instead.
	Follows string

	// Processor will process the given data, and return the result
	// If no data generated, set output to nil.
	// Processor func(f *KCallFrame) (output []byte, dataFormat int, dataType int)

	// result: output of this function, please convert according to format
	// format: the format of the result
	// hint: Why return this data?
	Processor func(f *KCallFrame) (result []byte, dataFormat int, hint int)

	//
	// OnXXX
	//
	// OnBeforeReg   func(nm *KNodeManager, self *KNode) // Not in NM
	// OnAfterReg    func(nm *KNodeManager, self *KNode) // In NM
	// OnBeforeDereg func(nm *KNodeManager, self *KNode) // In NM
	// OnAfterDereg  func(nm *KNodeManager, self *KNode) // Not in NM
	OnStart func(nm *KNodeManager, self *KNode)
	// OnStop        func(nm *KNodeManager, self *KNode)

	UserDataType string
	UserData     interface{}
}

func init() {
}
