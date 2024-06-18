package node

import (
	"fmt"
	"miego/klog"
	"time"
)

type KCallFrame struct {
	Caller     *KCallFrame // caller: upper layer frame
	Node       *KNode      // node: who process the data
	Data       []byte      // data: data returned by upstreamNode.Processor()
	DataFormat uint        // datafmt: uint
	Hint       uint        // XXX
	DataID     string      // DataID: GUID of the data if the data been saved to database.
	NewAt      int64       // NewAt: Timestamp this frame created.
}

func NewCallFrame(caller *KCallFrame, this *KNode, data []byte, datafmt uint, hint uint) *KCallFrame {
	f := &KCallFrame{
		Caller:     caller,
		Node:       this,
		Data:       data,
		DataFormat: datafmt,
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
