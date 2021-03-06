package fastcgi

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/TeaWeb/code/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddAction actions.Action

// 添加
func (this *AddAction) Run(params struct {
	From       string
	Server     string
	LocationId string
}) {
	this.Data["from"] = params.From
	this.Data["server"] = maps.Map{
		"filename": params.Server,
	}
	this.Data["locationId"] = params.LocationId

	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	Server      string
	LocationId  string
	On          bool
	Pass        string
	ReadTimeout int
	ParamNames  []string
	ParamValues []string
	PoolSize    int
	Must        *actions.Must
}) {
	params.Must.
		Field("pass", params.Pass).
		Require("请输入Fastcgi地址").
		Field("poolSize", params.PoolSize).
		Gte(0, "连接池尺寸不能小于0")

	paramsMap := map[string]string{}
	for index, paramName := range params.ParamNames {
		if index < len(params.ParamValues) {
			paramsMap[paramName] = params.ParamValues[index]
		}
	}

	server, err := teaconfigs.NewServerConfigFromFile(params.Server)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgiList, err := server.FindFastcgiList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	fastcgi := teaconfigs.NewFastcgiConfig()
	fastcgi.On = params.On
	fastcgi.Pass = params.Pass
	fastcgi.ReadTimeout = fmt.Sprintf("%ds", params.ReadTimeout)
	fastcgi.Params = paramsMap
	fastcgi.PoolSize = params.PoolSize
	fastcgiList.AddFastcgi(fastcgi)
	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
