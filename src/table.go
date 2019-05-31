package src

import (
	"bytes"
	"fmt"
)

type content map[string][]string
type tables struct {
	name        string //table name
	field       content
	index       content
	uniqueIndex content
	imports     map[string]string
}

func (t *tables) parseField(line []byte) {
	contents := bytes.Split(line, space)
	if len(contents) < 2 {
		return
	}
	var unsign bool
	if len(contents) >= 3 {
		unsign = bytes.Equal(contents[2], unsigned)
	}

	name := trim(contents[0])
	tagName := string(name)
	nameStr := string(title(name))
	types := contents[1]
	typesStr := toString(types)
	contents = contents[2:]
	if types[len(types)-1] == ')' {
		types = types[:len(types)-1]
	}
	tmpBytes := bytes.Split(types, []byte{'('})
	var lenBytes = []byte{0}
	if len(tmpBytes) > 1 {
		types = tmpBytes[0]
		lenBytes = tmpBytes[1]
	}

	length := int(lenBytes[0])
	types = tmpBytes[0]
	tp := string(types)
	switch tp {
	case "tinyint":
		if length == 1 {
			tp = "boole"
		} else {
			tp = "int8"
		}
	case "smallint":
		tp = "int16"
	case "integer":
		tp = "int"
	case "bigint":
		tp = "int64"
	case "decimal", "float":
		tp = "float64"
	case "char", "varchar":
		tp = "string"
	case "date", "datetime", "timestamp", "time":
		tp = "time.Time"
		if t.imports == nil {
			t.imports = make(map[string]string)
		}
		t.imports["time"] = "time"
	}

	if unsign {
		tp = "u" + tp
	}

	if t.field == nil {
		t.field = make(map[string][]string)
	}

	t.field[nameStr] = append(t.field[nameStr], tp)

	if i := bytes.Index(line, comment); i > 0 { //has comment
		var commentByts []byte
		var j int
		for k, v := range contents {
			if j > 0 && k > j {
				commentByts = append(commentByts, v...)
			}
			if bytes.Equal(v, comment) {
				j = k
			}
		}
		commentByts = commentByts[1:]
		commentByts = commentByts[:len(commentByts)-2]

		t.field[nameStr] = append(t.field[nameStr], string(commentByts))
	} else {
		t.field[nameStr] = append(t.field[nameStr], "")
	}

	t.field[nameStr] = append(t.field[nameStr], tagName)

	if *gorm { //use gorm tag
		//slice := append(t.field[nameStr], fmt.Sprintf("TYPE:%s;SIZE:%d", string(types), length))
		slice := append(t.field[nameStr], fmt.Sprintf("TYPE:%s", typesStr))
		if bytes.Contains(line, notNull) {
			slice = append(slice, toString(notNull))
		}
		if bytes.Contains(line, autoIncrement) {
			slice = append(slice, toString(autoIncrement))
		}
		if bytes.Contains(line, dft) {
			for k := range contents {
				if bytes.Equal(contents[k], dft) {
					slice = append(slice, "DEFAULT:"+string(bytes.TrimRight(contents[k+1], ",")))
				}
			}
		}
		t.field[nameStr] = slice
	}
}

func (t *tables) parseKey(line []byte) {
	if !*gorm {
		return
	}
	contents := bytes.Split(line, space)
	if len(contents) >= 3 {
		key := contents[2]
		key = key[2 : len(key)-2]
		nameStr := string(title(key))
		t.field[nameStr] = append(t.field[nameStr], "PRIMARY_KEY")
	}

}

func (t *tables) parseUniqueIndex(line []byte) {
	if !*gorm {
		return
	}

	t.parseIndex(line)
}
func (t *tables) parseIndex(line []byte) {
	if !*gorm {
		return
	}

	contents := bytes.Split(line, []byte{'`'})
	if len(contents) < 5 {
		return
	}

	if t.index == nil {
		t.index = make(map[string][]string)
	}
	contents = contents[1 : len(contents)-1]
	for i := 1; i <= len(contents)/2; i++ {
		nameStr := string(title(contents[i*2]))
		t.index[nameStr] = append(t.index[nameStr], "INDEX:"+string(contents[0]))
	}
}
