package node

import (
	"fmt"
	"strings"

	"miego/klog"
)

type KNodeManager struct {
	//
	// ALL nodes
	//
	nodes []*KNode

	// nodeTypeA <> downstreamNodes
	// dataReady(n, data) => lbNode := nodeTopoType[n.Type]
	//
	// XXX: nm.NodeAdd(node, NodeType, DataFormat, DataType)
	// "NodeType*hint" == (NodeA, NodeB, ..., NodeX)
	followMap map[uint][]*KNode
}

func (nm *KNodeManager) NodeAdd(node *KNode) {
	node.TypeId = NpAdd(node.Type)
	node.nm = nm
	nm.nodes = append(nm.nodes, node)

	s := node.Follows
	if s == "" {
		return
	}

	var iHint, iNode uint
	for _, m := range strings.Split(s, ";") {
		segs := strings.Split(m, "@")
		h, n := segs[0], segs[1]

		if h == "" || h == "*" {
			iHint = HT_ANY
		} else {
			iHint = NpNum(h)
		}

		if n == "" || n == "*" {
			iNode = NT_ANY
		} else {
			iNode = NpNum(n)
		}

		nm.LinkAfter(node, iNode, iHint)
	}
}

func catchError() {
	if err := recover(); err != nil {
		klog.F("----------")
		klog.E("%s\n", fmt.Sprint(err))
		klog.F("----------")
	}
}

func (nm *KNodeManager) LinkAfter(node *KNode, nodeType uint, hint uint) {
	hashkey := uint(nodeType) * uint(hint)
	klog.F("nodeType:%d, hint:%d", nodeType, hint)

	subNodes, ok := nm.followMap[hashkey]
	if !ok {
		subNodes = make([]*KNode, 0)
	}
	subNodes = append(subNodes, node)
	nm.followMap[hashkey] = subNodes
}

func (nm *KNodeManager) startNodes() {
	for _, n := range nm.nodes {
		if n.OnStart != nil {
			n.OnStart(nm, n)
		}
	}
}

// endPipeline: The pipeline at this node and this data.
func (nm *KNodeManager) endPipeline(node *KNode, data []byte, dataFormat uint, hint uint) {
	// TODO: Save the frame to database.
	klog.D("TODO: %s.", node.Name)
}

// Create new frame base on @f and send @data to @nodeDst
//
// @f: Parent frame
//
// XXX: Caller MUST ensure the nodeDst.Processor exists.
func (nm *KNodeManager) callNode(caller *KCallFrame, nodeDst *KNode, data []byte, dataFormat uint, hint uint) {
	// klog.D("callNode: %s, callCount:%d", nodeDst.Type, *callNodes)
	// defer catchError()

	// 1: Create a new frame and save the data and node to it
	newframe := NewCallFrame(caller, nodeDst, data, dataFormat, hint)

	// 2: Call nodeDst.Processor with new frame
	data, dataFormat, hint = nodeDst.Processor(newframe)
	newframe.Dump()

	// 3. Push to next stage.
	klog.D("PROCESSORR: Node:(%s | %d | %s), Fmt:%s, hint:%s:%d.", nodeDst.Type, nodeDst.TypeId, nodeDst.Name, NpStr(dataFormat), NpStr(hint), hint)
	nm.sendtoSubs(newframe, nodeDst, data, dataFormat, hint)
}

// XXX: Directly send data to node
//
// @caller: current(caller's) call frame
// @nodeDst: node who will process the data
// @data: data generated by node
func (nm *KNodeManager) SendtoNode(caller *KCallFrame, nodeDst *KNode, data []byte, dataFormat uint, hint uint) {
	if nodeDst.Processor != nil {
		go nm.callNode(caller, nodeDst, data, dataFormat, hint)
	}
}

func (nm *KNodeManager) getNexts(node *KNode, hint uint) []*KNode {
	utype := uint(node.TypeId)
	uhint := uint(hint)

	var hashkey uint
	var nexts []*KNode

	//
	// try: a+b a+* *+b *+*
	//

	hashkey = utype * uhint
	if subNodes, ok := nm.followMap[hashkey]; ok {
		nexts = append(nexts, subNodes...)
	}

	hashkey = uint(NT_ANY) * uhint
	if subNodes, ok := nm.followMap[hashkey]; ok {
		nexts = append(nexts, subNodes...)
	}

	hashkey = utype * uint(HT_ANY)
	if subNodes, ok := nm.followMap[hashkey]; ok {
		nexts = append(nexts, subNodes...)
	}

	hashkey = uint(NT_ANY) * uint(HT_ANY)
	if subNodes, ok := nm.followMap[hashkey]; ok {
		nexts = append(nexts, subNodes...)
	}

	klog.D("---------------------------------------------")
	klog.D("THIS IS : %s, T:%s, S:%d, IT'S NEXTS ARE:", node.Name, node.Type, uhint)
	for i, n := range nexts {
		if n == node {
			klog.D("- %d: %s (SKIPPED)", i, n.Name)
			continue
		}
		klog.D("- %d: %s", i, n.Name)
	}
	return nexts
}

// XXX: Send data to down stream node
//
// @caller: current(caller's) call frame
// @nodeDst: node who generate the data
// @data: data generated by node
func (nm *KNodeManager) sendtoSubs(caller *KCallFrame, nodeSrc *KNode, data []byte, dataFormat uint, hint uint) {
	nextNodes := nm.getNexts(nodeSrc, hint)
	if nextNodes == nil {
		// XXX: end this pipeline
		nm.endPipeline(nodeSrc, data, dataFormat, hint)
	}

	if hint == 0 {
		klog.D("xxxxxxxxx")
	}

	// klog.F("%s", string(nextNodes))
	for _, nodeNext := range nextNodes {
		// FIXME: "nodeNext != nodeSrc" OK?
		if nodeNext.Processor != nil && nodeNext != nodeSrc {
			go nm.callNode(caller, nodeNext, data, dataFormat, hint)
		}
	}
}

func (nm *KNodeManager) Run() {
	NpDump()

	klog.D(">>> followMap")
	klog.Dump(nm.followMap)
	klog.D("<<< followMap")

	fmt.Println(">>>>>>>> START nm.nodes DUMP")
	for i, n := range nm.nodes {
		fmt.Println("\n--- NODE(", i, ")")
		klog.Dump(n)
	}
	fmt.Println("\n<<<<<<<< END nm.nodes DUMP")

	//
	// buildTopo the topology
	//
	nm.startNodes()

	//
	// Dump link Map
	//
	for hashKey, nodes := range nm.followMap {
		klog.D("HASHKEY: '%d'", hashKey)
		for _, node := range nodes {
			klog.D("- %s:%s", node.Type, node.Name)
		}
	}

	for {
		select {}
	}
}

var NodeManager *KNodeManager

func init() {
	NodeManager = new(KNodeManager)
	NodeManager.followMap = make(map[uint][]*KNode)
}
