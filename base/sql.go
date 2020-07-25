package base

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"time"
)

const (
	EnPage = "SQL_CALC_FOUND_ROWS" //开启分页
)

type SqlDiff string

const (
	//相等
	Equal SqlDiff = "="
	//小于
	Lt SqlDiff = "<"
	//大于
	Gt SqlDiff = ">"
	//大于
	Like SqlDiff = "like"
	//范围区间
	Range SqlDiff = "between"
)

//sql编辑对象
type SqlModify struct {
	//表名
	table string
	//and条件列
	and []string
	//or条件列
	or []string
	//and条件值
	av []interface{}
	//or条件值
	ov []interface{}
	//数据模型对象实例
	obj interface{}
	//排序
	order string
	//分页
	page int
	//单页数量
	size int
}

func (sm *SqlModify) SetTable(table string) *SqlModify {
	sm.table = table
	return sm
}

//设置数据模型
func (sm *SqlModify) SetModel(DataModel interface{}) *SqlModify {
	sm.obj = DataModel
	return sm
}

func (sm *SqlModify) SetPage(p int) *SqlModify {
	sm.page = p
	return sm
}

func (sm *SqlModify) SetSize(s int) *SqlModify {
	sm.size = s
	return sm
}

//设置and条件
func (sm *SqlModify) And(c string, diff SqlDiff, val ...interface{}) *SqlModify {
	switch diff {
	case Range:
		if len(val) == 2 {
			sm.and = append(sm.and, fmt.Sprintf(`(%s between ? and ? )`, c))
			sm.av = append(sm.av, val...)
		}
		break
	case Like:
		if len(val) == 1 {
			sm.and = append(sm.and, "("+c+" like ?)")
			sm.av = append(sm.av, "%"+fmt.Sprintf(`%v`, val[0])+"%")
		}
		break
	default:
		if len(val) == 1 {
			sm.and = append(sm.and, fmt.Sprintf(`( %s %s ?)`, c, diff))
			sm.av = append(sm.av, val...)
		}
	}
	return sm
}

//设置or条件
func (sm *SqlModify) Or(c string, diff SqlDiff, val ...interface{}) *SqlModify {
	switch diff {
	case Range:
		if len(val) == 2 {
			sm.or = append(sm.or, fmt.Sprintf(`(%s between ? and ? )`, c))
			sm.ov = append(sm.ov, val...)
		}
		break
	case Like:
		if len(val) == 1 {
			sm.or = append(sm.or, "("+c+" like ?)")
			sm.ov = append(sm.ov, "%"+fmt.Sprintf(`%v`, val[0])+"%")
		}
		break
	default:
		if len(val) == 1 {
			sm.or = append(sm.or, fmt.Sprintf(`( %s %s ?)`, c, diff))
			sm.ov = append(sm.ov, val...)
		}
	}
	return sm
}

func (sm *SqlModify) QueryList() (string, []interface{}) {
	var sql string

	and, or := strings.Join(sm.and, " and "), strings.Join(sm.or, " or ")
	var condition []string
	if and != "" {
		condition = append(condition, and)
	}
	if or != "" {
		condition = append(condition, or)
	}
	if len(condition) > 0 {
		sql = fmt.Sprintf(`select %s * from %s where %s `, EnPage, sm.table, strings.Join(condition, " or "))
	} else {
		sql = fmt.Sprintf(`select %s * from %s `, EnPage, sm.table)
	}

	return sql, append(sm.av, sm.ov...)
}

func (sm *SqlModify) queryList(db *sql.DB, sql string, val ...interface{}) (list []map[string]interface{}, total int, err error) {

	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			fmt.Println(err, "\n", sql, val)
		}
		tx.Commit()
	}()
	rows, err := tx.Query(sql, val...)
	if err != nil {
		return
	}
	defer rows.Close()
	list, err = sm.parse(rows)
	if err != nil {
		return
	}
	rows, err = tx.Query(`SELECT FOUND_ROWS() total`)

	if err != nil {
		return
	}
	m, err := sm.parse(rows)
	if err != nil || len(m) != 1 {
		return
	}
	total, err = ValToInt(m[0]["total"])
	return
}

func (sm *SqlModify) Out(data []map[string]interface{}) (list []interface{}) {
	mv := reflect.TypeOf(sm.obj)
	for _, v := range data {
		newObj := reflect.New(mv)
		for i := 0; i < newObj.Elem().NumField(); i++ {
			one := v[mv.Field(i).Tag.Get("sql")]
			newObj.Elem().Field(i).Set(reflect.ValueOf(one))
		}
		list = append(list, newObj.Interface())
	}
	return list
}

//解析数据库返回数据
func (sm *SqlModify) parse(rows *sql.Rows) (list []map[string]interface{}, err error) {
	//返回所有列
	cols, _ := rows.ColumnTypes()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	for rows.Next() {
		//填充数据
		rows.Scan(scans...)

		//每行数据
		row := make(map[string]interface{})
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k].Name()
			bytesBuffer := bytes.NewBuffer(v)

			switch cols[k].DatabaseTypeName() {
			case "INT", "TINYINT", "BIGINT":
				row[key], _ = ValToInt(bytesBuffer.String())
				break
			case "TIMESTAMP":
				row[key], _ = time.Parse("2006-01-02 15:04:05", bytesBuffer.String())
				break
			default:
				row[key] = bytesBuffer.String()
				break
			}
		}
		//放入结果集
		list = append(list, row)
	}
	return
}
