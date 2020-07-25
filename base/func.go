//功能性方法文件
package base

import (
	"fmt"
	"strconv"
)

//数据转int
func ValToInt(v interface{}) (int, error) {
	vs := fmt.Sprintf("%v", v)
	return strconv.Atoi(vs)
}
