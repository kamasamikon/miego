package node

import (
	"fmt"
	"github.com/kamasamikon/miego/klog"
	"time"
)

type KCallFrame struct {
	// caller: upper layer frame
	// node: who process the data
	// data: data returned by upstreamNode.Processor()
	// dataFormat: uint
	// DataID: GUID of the data if the data been saved to database.
	// NewAt: Timestamp this frame created.
	Caller     *KCallFrame
	Node       *KNode
	Data       []byte
	DataFormat uint
	Hint       uint
	DataID     string
	NewAt      int64
}

func NewCallFrame(caller *KCallFrame, thisNode *KNode, data []byte, dataFormat uint, hint uint) *KCallFrame {
	f := &KCallFrame{
		Caller:     caller,
		Node:       thisNode,
		Data:       data,
		DataFormat: dataFormat,
		Hint:       hint,
		DataID:     "TODO: Id of DB?",
		NewAt:      time.Now().UnixNano(),
	}

	// TODO: Save the data to database;
	// TODO: Create DataID

	return f
}

func (f *KCallFrame) Dump() {
	var lines []string

	tmp := f
	fmtstr := " %-16s | %4s | %4s | %s | %lld"
	sp := fmt.Sprintf
	for tmp != nil {
		sDataFormat, sHint := NpStr(tmp.DataFormat), NpStr(tmp.Hint)

		line := sp(fmtstr, tmp.Node.Type, sDataFormat, sHint, tmp.Node.Name, tmp.NewAt)
		lines = append(lines, line)
		tmp = tmp.Caller
	}

	klog.D(">>>> DUMP (%d) >>>>", len(lines))
	cnt := len(lines)
	for i := range lines {
		klog.D("|%4d|%s", i, lines[cnt-i-1])
	}
	klog.D("<<<< DUMP (%d) <<<<", len(lines))
}
