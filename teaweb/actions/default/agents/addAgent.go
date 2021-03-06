package agents

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/string"
)

type AddAgentAction actions.Action

// 添加代理
func (this *AddAgentAction) Run(params struct{}) {
	this.Show()
}

// 提价保存
func (this *AddAgentAction) RunPost(params struct {
	Name       string
	Host       string
	AllowAllIP bool
	IPs        []string `alias:"ips"`
	On         bool
	Must       *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入主机名").
		Field("host", params.Host).
		Require("请输入主机地址")

	agentList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	agent := agents.NewAgentConfig()
	agent.On = params.On
	agent.Name = params.Name
	agent.Host = params.Host
	agent.AllowAll = params.AllowAllIP
	agent.Allow = params.IPs
	agent.Key = stringutil.Rand(32)
	err = agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	agentList.AddAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Data["agentId"] = agent.Id

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("ADD_AGENT", maps.Map{}))

	this.Success()
}
