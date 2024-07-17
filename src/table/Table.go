package table

type Table struct {
	columns []*Column
	raws []*Raw
	boxSettings *BoxSettings
	pure bool
}

func New() *Table {
	return &Table{boxSettings: NewBoxSettings()}
}

func (table *Table) SetBoxSettings(boxSettings *BoxSettings) {
	table.boxSettings = boxSettings
}

func (table *Table) GetBoxSettings() *BoxSettings {
	return table.boxSettings
}

func (table *Table) GetColumnSizes() []int {
	var list []int

	for _, column := range table.columns {
		list = append(list, column.width)
	}

	return list
}

func (table *Table) AppendColumn(column *Column)  {
	table.columns = append(table.columns, column)
}

func (table *Table) CreateColumn(name string) *Column {
	column := NewColumn(table, name)

	table.AppendColumn(column)

	return column
}

func (table *Table) CreateRaw() *Raw {
	raw := NewRaw(table)

	table.raws = append(table.raws, raw)

	return raw
}

func (table *Table) CreateRawWithValues(values []string) *Raw {
	raw := NewRaw(table)

	for i := 0; i < len(values); i++ {
		raw.CreateItem(values[i])
	}

	table.raws = append(table.raws, raw)

	return raw
}

func (table *Table) AppendRaw(raw *Raw)  {
	table.raws = append(table.raws, raw)
}

func (table *Table) SetPure(pure bool)  {
	table.pure = pure
}