package table

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/CaiJinKen/gostruct/handler"
)

const tableNameStr = `
func (*%s) TableName() string {
	return "%s"
}
`

type Config struct {
	UseGormTag bool
	UseJsonTag bool
	SortField  bool   // sort field
	PkgName    string // output file package name
}

func (c *Config) Build() *Table {
	table := newTable()
	table.Config = *c
	return table
}

type defaultValue struct {
	value string
	valid bool
}

type Field struct {
	table *Table

	Name          string
	RawName       string
	Default       defaultValue
	TypeName      string
	RawTypeName   string
	Comment       string
	AutoIncrement bool
	Unsigned      bool
	NotNull       bool

	Type *Type

	indexes []*Index
}

type Type struct {
	size        uint
	decimalSize uint
	name        string
}

func (t *Type) String() string {
	if t.name == "" {
		return "interface{}"
	}
	return t.name
}

type Table struct {
	Config

	Name         string
	RawName      string
	Fields       []*Field
	PrimaryKeys  []string
	Indexes      []*Index
	nameFiledMap map[string]*Field
	Imports      map[string]string
	Comment      string

	model *model
}

type Index struct {
	Fields  []string
	Type    int
	RawName string
	Comment string
}

func newTable() *Table {
	return &Table{
		Name:         "",
		RawName:      "",
		Fields:       make([]*Field, 0),
		Indexes:      make([]*Index, 0),
		PrimaryKeys:  make([]string, 0),
		nameFiledMap: make(map[string]*Field),
		Imports:      make(map[string]string),

		model: newModel(),
	}
}

func (f *Field) parseType() {
	if f.TypeName == "" {
		return
	}
	tp := &Type{}
	f.Type = tp
	typeNameSlice := strings.Split(f.TypeName, "(")
	f.TypeName = typeNameSlice[0]
	if len(typeNameSlice) > 1 {
		lengths := strings.Split(typeNameSlice[1], "")
		str := strings.Join(lengths[:len(lengths)-1], "")
		sizes := strings.Split(str, ",")

		size, _ := strconv.Atoi(sizes[0])
		tp.size = uint(size)
		if len(sizes) > 1 {
			size, _ = strconv.Atoi(sizes[1])
			tp.decimalSize = uint(size)
		}
	}

	f.getType()
}

func (f *Field) getType() {
	switch f.TypeName {
	case "tinyint":
		f.Type.name = reflect.Int.String()
		if f.Type.size == 1 {
			f.Type.name = reflect.Bool.String()
		}
	case "smallint":
		f.Type.name = reflect.Int16.String()
		if f.Unsigned {
			f.Type.name = reflect.Uint16.String()
		}
	case "int", "integer":
		f.Type.name = reflect.Int.String()
		if f.Unsigned {
			f.Type.name = reflect.Uint.String()
		}
	case "bigint":
		f.Type.name = reflect.Int64.String()
		if f.Unsigned {
			f.Type.name = reflect.Uint64.String()
		}
	case "decimal", "float":
		f.Type.name = reflect.Float64.String()
	case "char", "varchar", "text", "longtext":
		f.Type.name = reflect.String.String()
	case "date", "datetime", "timestamp", "time":
		f.Type.name = "time.Time"
		f.table.Imports[" "] = "time"
	case "json":
		f.Type.name = "interface{}"
	}
	return
}

func (t *Table) parseField(line []byte) {
	contents := bytes.Split(line, []byte{' '})
	if len(contents) < 2 {
		return
	}
	f := &Field{table: t}
	t.Fields = append(t.Fields, f)

	name := trim(contents[0])
	f.RawName = string(name)
	f.Name = string(title(name))
	f.TypeName = string(trim(contents[1]))
	f.RawTypeName = f.TypeName

	t.nameFiledMap[f.RawName] = f

	for i := 2; i < len(contents); {
		value := contents[i]
		switch string(value) {
		case "NOT":
			if string(contents[i+1]) == "NULL" {
				f.NotNull = true
			}
			i += 2
			continue

		case "DEFAULT":
			f.Default = defaultValue{
				value: string(bytes.Trim(contents[i+1], "'")),
				valid: true,
			}
			if len(contents[i+1]) == 2 {
				f.Default.value = "''"
			}

			i += 1
			continue

		case "COMMENT":
			f.Comment = string(bytes.Trim(contents[i+1], "'"))
			i += 2
			continue

		case "AUTO_INCREMENT":
			f.AutoIncrement = true
		case "unsigned":
			f.Unsigned = true
		}
		i++
	}
	f.parseType()

}

func (t *Table) parseComment(line []byte) {
	idx := bytes.Index(line, []byte("COMMENT"))
	if idx < 0 {
		return
	}
	line = line[idx:]
	contents := bytes.Split(line, []byte{' '})
	line = contents[0]
	contents = bytes.Split(line, []byte{'='})
	if len(contents) < 2 {
		return
	}
	line = contents[1]
	line = bytes.Trim(line, ";")
	t.Comment = string(line)
}

func (t *Table) parseKey(line []byte) {
	contents := bytes.Split(line, []byte{' '})
	if len(contents) < 3 {
		return
	}

	for i := 2; i < len(contents); i++ {
		content := contents[i]
		content = bytes.TrimSpace(content)
		content = bytes.TrimPrefix(content, []byte{'('})
		content = bytes.TrimSuffix(content, []byte{')'})
		content = bytes.TrimSuffix(content, []byte{','})
		content = bytes.Trim(content, "`")
		t.PrimaryKeys = append(t.PrimaryKeys, string(content))
	}
}

func (t *Table) parseUniqueIndex(line []byte) {
	contents := bytes.Split(line, []byte{' '})
	if len(contents) < 3 {
		return
	}
	idx := &Index{
		Fields:  nil,
		Type:    2,
		RawName: string(trim(contents[2])),
		Comment: "",
	}
	size := len(contents[3])
	fields := contents[3][1 : size-1]
	for _, v := range bytes.Split(fields, []byte{','}) {
		idx.Fields = append(idx.Fields, string(v[1:len(v)-1]))
	}
	for i := 4; i < len(contents); i++ {
		v := contents[i]
		if string(v) != "COMMENT" {
			continue
		}
		idx.Comment = string(contents[i+1])
	}
	for _, v := range idx.Fields {
		if filed, ok := t.nameFiledMap[v]; ok && filed != nil {
			filed.indexes = append(filed.indexes, idx)
		}
	}
}

func (t *Table) parseIndex(line []byte) {
	contents := bytes.Split(line, []byte{' '})
	if len(contents) < 2 {
		return
	}
	idx := &Index{
		Fields:  nil,
		Type:    1,
		RawName: string(trim(contents[1])),
		Comment: "",
	}
	size := len(contents[2])
	fields := contents[2][1 : size-1]
	for _, v := range bytes.Split(fields, []byte{','}) {
		idx.Fields = append(idx.Fields, string(v[1:len(v)-1]))
	}
	for i := 3; i < len(contents); i++ {
		v := contents[i]
		if string(v) != "COMMENT" {
			continue
		}
		idx.Comment = string(contents[i+1])
	}

	for _, v := range idx.Fields {
		if filed, ok := t.nameFiledMap[v]; ok && filed != nil {
			filed.indexes = append(filed.indexes, idx)
		}
	}
}

func (t *Table) Marshal() {
	buf := t.model.reader
	buf.WriteString(fmt.Sprintf("package %s\n\n", t.PkgName))
	if len(t.Imports) > 0 {
		buf.WriteString("import (\n")
		for alias, path := range t.Imports {
			alias = strings.TrimSpace(alias)
			buf.WriteString(fmt.Sprintf("\t%s \"%s\"\n", alias, path))
		}
		buf.WriteString(")\n\n")
	}

	if t.SortField {
		sort.Slice(t.Fields, func(i, j int) bool {
			return t.Fields[i].Name < t.Fields[j].Name
		})
	}

	buf.WriteString(fmt.Sprintf("type %s struct {\n", t.Name))
	for _, filed := range t.Fields {
		buf.WriteString(fmt.Sprintf("\t%s\t%s", filed.Name, filed.Type.String()))
		if t.UseJsonTag || t.UseGormTag {
			buf.WriteString("\t`")
			if t.UseJsonTag {
				str := `json:"%s"`
				if !filed.NotNull {
					str = `json:"%s,omitempty"`
				}
				buf.WriteString(fmt.Sprintf(str, filed.RawName))
			}
			if t.UseGormTag {
				if t.UseJsonTag {
					buf.WriteByte(' ')
				}
				var contents []string
				contents = append(contents,
					fmt.Sprintf("COLUMN:%s", filed.RawName),
					fmt.Sprintf("TYPE:%s", filed.RawTypeName),
				)
				if filed.Default.valid {
					contents = append(contents, fmt.Sprintf("DEFAULT:%s", filed.Default.value))
				}
				if filed.NotNull {
					contents = append(contents, "NOT NULL")
				}
				if filed.AutoIncrement {
					contents = append(contents, "AUTOINCREMENT")
				}
				if filed.Type.size > 0 {
					contents = append(contents, fmt.Sprintf("SIZE:%d", filed.Type.size))
				}
				for _, v := range t.PrimaryKeys {
					if filed.RawName == v {
						contents = append(contents, "PRIMARYKEY")
					}
				}

				for _, v := range filed.indexes {
					switch v.Type {
					case 1:
						contents = append(contents, fmt.Sprintf("INDEX:%s", v.RawName))
					case 2:
						contents = append(contents, fmt.Sprintf("UNIQUEINDEX:%s", v.RawName))
					}
				}

				buf.WriteString(fmt.Sprintf(`gorm:"%s"`, strings.Join(contents, ";")))

			}
			buf.WriteByte('`')
		}
		if filed.Comment != "" {
			buf.WriteString(fmt.Sprintf("\t// %s", filed.Comment))
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")

	buf.WriteString(fmt.Sprintf(tableNameStr, t.Name, t.RawName))

	return
}

func (t *Table) Format() {
	t.model.format()
}

func (t *Table) Data() []byte {
	t.Format()
	return t.model.result
}

func (t *Table) Parse(data []byte) {

	if len(data) == 0 {
		return
	}

	var continueComment bool

	reader := bufio.NewReader(bytes.NewReader(data))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			handler.PrintErrAndExit(err)
		}
		line = trimLine(line)

		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("/*")) {
			continueComment = true
			continue
		}

		if bytes.HasPrefix(line, []byte("*/")) {
			continueComment = false
			continue
		}

		if bytes.HasPrefix(line, []byte("--")) || continueComment || bytes.HasPrefix(line, []byte("SET")) || bytes.HasPrefix(line, []byte("DROP")) {
			continue
		}

		if bytes.HasPrefix(line, []byte("CREATE TABLE")) {
			name := trim(bytes.Split(line, []byte{' '})[2])
			t.RawName = string(name)
			t.Name = string(title(name))
			continue
		}

		line = bytes.Trim(line, ",")
		switch line[0] {
		case ')':
			t.parseComment(line)
		case '`':
			t.parseField(line)
		case 'P':
			t.parseKey(line)
		case 'I', 'K':
			t.parseIndex(line)
		case 'U':
			t.parseUniqueIndex(line)
		}
	}
	return
}
