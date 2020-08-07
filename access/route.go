package access

//系统路由
type Route struct {
	Id int `json:"id"`
	//级别
	Level int `json:"level"`
	//接口地址,只有末梢级别才有接口地址
	Url string `json:"url"`
	//客户端显示名字
	DisplayName string `json:"display_name"`
	//上级ID
	Fid int `json:"fid"`
}
