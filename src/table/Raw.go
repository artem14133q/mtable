package table

type Raw struct {
	parent *Table
	items []*RawItem
	size int
}

func NewRaw(parent *Table) *Raw {
	return &Raw{parent: parent, size: 0}
}

func (raw *Raw) GetItem(i int) *RawItem {
	return raw.items[i]
}

func (raw *Raw) AddItem(item *RawItem)  {
	raw.items = append(raw.items, item)
	raw.size += 1
}

func (raw *Raw) CreateItem(value string) *RawItem {
	rawItem := NewRawItem(raw, value)

	raw.AddItem(rawItem)

	return rawItem
}