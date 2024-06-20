package node

import (
	"fmt"
	"miego/klog"
	"time"

	"github.com/twinj/uuid"
)

type KCallFrame struct {
	Caller     *KCallFrame // caller: upper layer frame
	Node       *KNode      // node: who process the data
	Data       interface{} // data: data returned by upstreamNode.Processor()
	DataFormat uint64      // datafmt: uint
	Hint       uint64      // XXX
	DataID     string      // DataID: GUID of the data if the data been saved to database.
	NewAt      int64       // NewAt: Timestamp this frame created.
}

func NewCallFrame(caller *KCallFrame, this *KNode, data interface{}, datafmt uint64, hint uint64) *KCallFrame {
	return &KCallFrame{
		Caller:     caller,
		Node:       this,
		Data:       data,
		DataFormat: datafmt,
		Hint:       hint,
		DataID:     uuid.NewV4().String(),
		NewAt:      time.Now().UnixNano(),
	}
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
