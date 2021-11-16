package cmd

import (
	"bytes"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/vetcher/go-astra"
)

// TempData 表示生成template所需要的数据结构
type TempData struct {
	Module        string
	ModulePath    string
	FileName      string
	PackageName   string
	StructName    string
	TableName     string
	CacheData     string //数据缓存时长
	CacheList     string //list缓存时长
	CacheLimit    string //list缓存长度
	Columns       map[string][]string
	PrimaryKey    []string
	HasPrimaryKey bool
	HasTime       bool
	HasCache      bool
}

//handleFile 处理model文件
func handleFile(module, modulePath, filename string) error {
	tempData := new(TempData)
	tempData.Module = module
	tempData.ModulePath = modulePath

	fset := token.NewFileSet()
	var src interface{}
	_, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		showError(err)
		panic(err)
	}

	////////////
	file, err := astra.ParseFile(filename,
		astra.IgnoreVariables|astra.IgnoreConstants|astra.IgnoreFunctions|
			astra.IgnoreInterfaces|astra.IgnoreTypes|astra.IgnoreMethods)
	if err != nil {
		showError(err)
	}
	tempData.PackageName = file.Name

	for _, stru := range file.Structures {
		tempData.TableName = ""
		tempData.HasCache = false
		tempData.HasPrimaryKey = false
		tempData.CacheData = ""
		tempData.CacheList = ""
		tempData.CacheLimit = ""
		tempData.PrimaryKey = nil
		tempData.Columns = make(map[string][]string)
		tempData.FileName = filename
		tempData.StructName = stru.Name
		//解析struct文档
		parseDocs(tempData, stru.Docs)
		if tempData.TableName == "" {
			continue
			//tempData.TableName = tempData.StructName
		}

		for _, field := range stru.Fields {
			var pk string
			var _namejson = make([]string, 3)
			for k, v := range field.Tags {
				if k == "json" {
					_namejson[1] = v[0] //json_name
				} else if k == XORM_TAG {
					_namejson[0], pk = parseTagsForXORM(v) //column_name
				} else if k == GORM_TAG {
					_namejson[0] = parseTagsForGORM(v) //column_name
				}
			}
			_namejson[2] = field.Type.String()
			if _namejson[2] == "time.Time" {
				tempData.HasTime = true
			}

			if _namejson[1] == "" {
				if _namejson[0] == "" {
					_namejson[1] = getFieldName(field.Name)
				} else {
					_namejson[1] = _namejson[0]
				}
			}
			if _namejson[0] == "" {
				if _namejson[1] == "" {
					_namejson[0] = getFieldName(field.Name)
				} else {
					_namejson[0] = _namejson[1]
				}
			}
			tempData.Columns[field.Name] = _namejson
			if pk != "" {
				tempData.PrimaryKey = _namejson
				tempData.HasPrimaryKey = true
			}
		}
		//如果struct名称为空,或者是一个私有struct,或者field为空,返回
		if len(tempData.StructName) == 0 ||
			tempData.StructName[:1] == strings.ToLower(tempData.StructName[:1]) ||
			len(tempData.Columns) == 0 {
			return nil
		}

		if debug {
			return tempData.writeTo(os.Stdout)
		}
		//写model文件
		//if !functions["*"+stru.Name] && !functions[stru.Name] {
		err = tempData.writeToModel(filename)
		if err != nil {
			showError(err)
			return err
		}
		//}

		//写table文件
		err := tempData.writeToTable()
		if err != nil {
			showError(err.Error())
			return err
		}

		//写dao文件
		if tempData.HasPrimaryKey {
			err = tempData.writeToDao(filename)
			if err != nil {
				showError(err)
				return err
			}
		}
	}

	return err
}

func parseTagsForXORM(matchs []string) (columnName string, key string) {
	s := strings.Split(matchs[0], " ")
	if len(s) == 1 {
		columnName = strings.Replace(s[0], "'", "", -1)
		return
	}
	col := &columnName
	k := new(string)
	for _, v := range s {
		if v == "" {
			continue
		}
		if v[:1] == "'" {
			*col = strings.Replace(v, "'", "", -1)
			continue
		}
		if strings.ToLower(v) == "pk" {
			k = col
			continue
		}
	}
	key = *k
	return

	//if len(matchs) >= 1 {
	//	_matchs := regexp.MustCompile(`'(.*?)'`).FindStringSubmatch(matchs[0])
	//	if len(_matchs) >= 1 {
	//		return _matchs[1]
	//	}
	//}
	//return ""
}

func parseTagsForGORM(matchs []string) string {
	if len(matchs) >= 1 {
		_matchs := regexp.MustCompile(`(?i:column):(.*?)(?:;|$)`).FindStringSubmatch(matchs[0])
		if len(_matchs) >= 1 {
			return _matchs[1]
		}
	}
	return ""
}

func parseDocs(tmp *TempData, docs []string) {
	for _, doc := range docs {
		doc = strings.TrimLeft(doc, " /")
		if strings.Contains(doc, "tablename") {
			tmp.TableName = strings.TrimSpace(strings.TrimLeft(doc, "tablename"))
			continue
		}
		if strings.HasPrefix(doc, "cache ") {
			tmp.HasCache = true
			cache := strings.Replace(strings.TrimSpace(strings.TrimLeft(doc, "cache")), "  ", " ", -1)
			caches := strings.Split(cache, " ")
			if len(caches) >= 3 {
				tmp.CacheData = caches[0]
				tmp.CacheList = caches[1]
				tmp.CacheLimit = caches[2]
			}
			continue
		} else {
			if strings.HasPrefix(doc, "cachedata") {
				tmp.HasCache = true
				tmp.CacheData = strings.TrimSpace(strings.TrimLeft(doc, "cachedata"))
				continue
			}
			if strings.HasPrefix(doc, "cachelist") {
				tmp.HasCache = true
				tmp.CacheList = strings.TrimSpace(strings.TrimLeft(doc, "cachelist"))
				continue
			}
			if strings.HasPrefix(doc, "cachelimit") {
				tmp.HasCache = true
				tmp.CacheLimit = strings.TrimSpace(strings.TrimLeft(doc, "cachelimit"))
				continue
			}
		}
	}
}

func getFilepath(filename string) string {
	absPath, _ := filepath.Abs(filename)
	return filepath.Join(filepath.Dir(absPath), "table")
}
func getBaseFilename(filename string) string {
	f := filepath.Base(filename)
	pos := strings.LastIndex(f, ".")
	if pos == -1 {
		return f
	}
	return f[:pos]
}

func (d *TempData) tableFilename() string {
	return filepath.Join(getFilepath(d.FileName), getBaseFilename(d.FileName)+"_"+d.StructName+"_table.go")
}

func (d *TempData) writeTo(w io.Writer) error {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}
	return template.Must(template.New("tableTpl").Funcs(funcMap).Parse(tableTpl)).Execute(w, d)
}

// writeToTable 将生成好的模块文件写到本地
func (d *TempData) writeToTable() error {
	file, err := os.Create(d.tableFilename())
	if err != nil {
		showError(err.Error())
		return err
	}
	defer file.Close()
	var buf bytes.Buffer
	err = d.writeTo(&buf)
	if err != nil {
		showError(err.Error())
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		showError(err.Error())
		return err
	}
	_, err = file.Write(formatted)
	if err != nil {
		showError(err.Error())
		return err
	}
	return err
}

//
func getFieldName(name string) string {
	bs := bytes.NewBuffer([]byte{})

	pre_lower := true //前一个字母是小写
	for i, s := range name {
		//如果是大写字母
		if s >= 65 && s <= 90 {
			s += 32 //转成小写
			if i == 0 {
				bs.WriteByte(byte(s))
			} else {
				if pre_lower {
					bs.WriteByte(byte(95)) //写下划线
				}
				bs.WriteByte(byte(s))
			}
			pre_lower = false
			continue
		}
		pre_lower = true
		bs.WriteByte(byte(s))
	}
	return bs.String()
}
