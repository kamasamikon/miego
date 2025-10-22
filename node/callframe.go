package node

import (
	"fmt"
	"miego/klog"
	"strings"
	"time"

	"github.com/google/uuid"
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
		NewAt:      time.Now().UnixMicro(),
	}
}

func (f *KCallFrame) Dump() {
	x := ""
	w := 0
	ntWidth := 0 // NodeType
	dfWidth := 0 // DataFormat
	htWidth := 0 // Hint
	nmWidth := 0 // Name
	atWidth := 0 // NewAt

	var segs []string

	tmp := f
	for tmp != nil {
		x = tmp.Node.Type
		w = len(x)
		if w > ntWidth {
			ntWidth = w
		}
		segs = append(segs, x)

		x = NpStr(tmp.DataFormat)
		w = len(x)
		if w > dfWidth {
			dfWidth = w
		}
		segs = append(segs, x)

		x = NpStr(tmp.Hint)
		w = len(x)
		if w > htWidth {
			htWidth = w
		}
		segs = append(segs, x)

		x = fmt.Sprintf("%v", tmp.NewAt)
		w = len(x)
		if w > atWidth {
			atWidth = w
		}
		segs = append(segs, x)

		x = tmp.Node.Name
		w = len(x)
		if w > nmWidth {
			nmWidth = w
		}
		segs = append(segs, x)

		tmp = tmp.Caller
	}

	var ofmt string
	ofmt += "| "
	ofmt += fmt.Sprintf("%%-%ds", ntWidth)
	ofmt += " | "
	ofmt += fmt.Sprintf("%%%ds", dfWidth)
	ofmt += " | "
	ofmt += fmt.Sprintf("%%%ds", htWidth)
	ofmt += " | "
	ofmt += fmt.Sprintf("%%%ds", atWidth)
	ofmt += " | "
	ofmt += fmt.Sprintf("%%s")
	ofmt += "\n"

	var lines []string
	for i := range len(segs) / 5 {
		lines = append(
			lines,
			fmt.Sprintf(
				ofmt,
				segs[i*5+0],
				segs[i*5+1],
				segs[i*5+2],
				segs[i*5+3],
				segs[i*5+4],
			),
		)
	}

	klog.KLog(1, false, klog.ColorType_F, "D", "\n"+strings.Join(lines, ""))
}
