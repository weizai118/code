package apps

import (
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改设置
func (this *UpdateAction) Run(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	this.Data["agent"] = agent

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	AgentId    string
	Name       string
	Host       string
	AllowAllIP bool
	IPs        []string `alias:"ips"`
	On         bool
	Key        string
	Must       *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入主机名").
		Field("host", params.Host).
		Require("请输入主机地址").
		Field("key", params.Key).
		Require("请输入密钥")

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}
	agent.On = params.On
	agent.Name = params.Name
	agent.Host = params.Host
	agent.AllowAll = params.AllowAllIP
	agent.Allow = params.IPs
	agent.Key = params.Key
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_AGENT", maps.Map{}))

	this.Success()
}
