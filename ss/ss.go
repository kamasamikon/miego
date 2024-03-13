package ss

import (
	"miego/xmap"
)

// ///////////////////////////////////////////////////////////////////////
// SS: Query
type Input struct {
	// Var: VueVariables
	// Val: Value
	// Hint: Placeholder
	// Label: Label for checkbox
	Kind string

	Class string

	VarName string
	Value   string

	Hint  string
	Label string

	// Address
	VarName_Province string
	Value_Province   string

	VarName_City string
	Value_City   string

	VarName_District string
	Value_District   string

	VarName_Address string
	Value_Address   string
}

// Unit of masonry
type Query struct {
	Name      string
	Hint      string
	Class     string
	InputList []*Input
}

func QueryNew(Name string, Hint string, Class string) *Query {
	return &Query{
		Name:  Name,
		Hint:  Hint,
		Class: Class,
	}
}

func (q *Query) AddItem(i *Input) *Query {
	q.InputList = append(q.InputList, i)
	return q
}

func (q *Query) AddMap(mp xmap.Map) *Query {
	return q.AddItem(&Input{

		Kind:  mp.S("Kind"),
		Class: mp.S("Class"),

		VarName: mp.S("VarName"),
		Value:   mp.S("Value"),

		Hint:  mp.S("Hint"),
		Label: mp.S("Label"),

		VarName_Province: mp.S("VarName_Province"),
		Value_Province:   mp.S("Value_Province"),

		VarName_City: mp.S("VarName_City"),
		Value_City:   mp.S("Value_City"),

		VarName_District: mp.S("VarName_District"),
		Value_District:   mp.S("Value_District"),

		VarName_Address: mp.S("VarName_Address"),
		Value_Address:   mp.S("Value_Address"),
	})
}

func (q *Query) Add(Kind string, args ...interface{}) *Query {
	mp := xmap.Make(args...)
	mp["Kind"] = Kind
	return q.AddMap(mp)
}

// QueryGroup > Query > Input
type QueryGroup struct {
	QueryList []*Query
}

func QueryGroupNew() *QueryGroup {
	return &QueryGroup{}
}

func (g *QueryGroup) Add(q *Query) *QueryGroup {
	g.QueryList = append(g.QueryList, q)
	return g
}

func (g *QueryGroup) Items() []*Query {
	return g.QueryList
}

/////////////////////////////////////////////////////////////////////////
// : Columns

type ColumnItem struct {
	Show  string
	Label string

	Hidden     bool
	Sortable   bool
	Filterable bool
	DataType   string
}

type Columns struct {
	ColumnItemList []*ColumnItem
}

func ColumnsNew() *Columns {
	return &Columns{}
}

func (c *Columns) AddItem(item *ColumnItem) *Columns {
	if item.Label == "" {
		item.Label = item.Show
	}
	c.ColumnItemList = append(c.ColumnItemList, item)
	return c
}

func (c *Columns) Add(Show string, Label string, Hidden bool, Sortable bool, Filterable bool, Numeric bool) *Columns {
	DataType := "string"
	if Numeric {
		DataType = "numeric"
	}
	return c.AddItem(&ColumnItem{
		Show:       Show,
		Label:      Label,
		Hidden:     Hidden,
		Sortable:   Sortable,
		Filterable: Filterable,
		DataType:   DataType,
	})
}

func (c *Columns) AddMap(mp xmap.Map) *Columns {
	return c.AddItem(&ColumnItem{
		Show:       mp.S("Show"),
		Label:      mp.S("Label"),
		Hidden:     mp.Bool("Hidden", false),
		Sortable:   mp.Bool("Sortable", false),
		Filterable: mp.Bool("Filterable", false),
	})
}

func (c *Columns) Items() []*ColumnItem {
	return c.ColumnItemList
}
