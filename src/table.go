package src

import "bytes"

type content map[string][]string
type tables struct {
	name        string
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
		unsign = bytes.Equal(contents[2], unsigns)
	}

	name := trim(contents[0])
	tagName := string(name)
	nameStr := string(title(name))
	types := contents[1]
	contents = contents[2:]
	types = types[:len(types)-1]
	tmpBytes := bytes.Split(types, []byte{'('})
	var lenBytes []byte
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
	t.field[nameStr] = append(t.field[nameStr], tagName)

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
	}

}

func (t *tables) parseKey(line []byte) {

}

func (t *tables) parseUniqueIndex(line []byte) {

}
func (t *tables) parseIndex(line []byte) {

}
