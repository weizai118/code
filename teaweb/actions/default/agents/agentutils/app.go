package agentutils

import (
	"context"
	"github.com/TeaWeb/code/teaconfigs/agents"
	"github.com/TeaWeb/code/teamongo"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
)

// App菜单
func InitAppData(actionWrapper actions.ActionWrapper, agentId string, appId string, tabbar string) *agents.AppConfig {
	agent := agents.NewAgentConfigFromId(agentId)
	action := actionWrapper.Object()
	if agent == nil {
		action.Fail("找不到Agent")
	}

	app := FindAgentApp(agent, appId)
	if app == nil {
		action.Fail("找不到App")
	}

	action.Data["agentId"] = agentId
	action.Data["app"] = maps.Map{
		"id":                 app.Id,
		"name":               app.Name,
		"on":                 app.On,
		"countItems":         len(app.Items),
		"countBootTasks":     len(app.FindBootingTasks()),
		"countScheduleTasks": len(app.FindSchedulingTasks()),
		"countManualTasks":   len(app.FindManualTasks()),
		"isSystem":           app.IsSystem,
	}
	action.Data["selectedTabbar"] = tabbar

	return app
}

// 格式化任务信息
func FormatTask(task *agents.TaskConfig, agentId string) maps.Map {
	// 最近执行
	ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
	cursor, err := teamongo.FindCollection("logs.agent." + agentId).Find(ctx, map[string]interface{}{
		"taskId": task.Id,
	}, findopt.Sort(map[string]interface{}{
		"_id": -1,
	}), findopt.Limit(1))
	runTime := ""
	if err == nil {
		if cursor.Next(context.Background()) {
			log := &ProcessLog{}
			err = cursor.Decode(log)
			if err == nil {
				runTime = timeutil.Format("Y-m-d H:i:s", time.Unix(log.Timestamp, 0))
			}
		}
		cursor.Close(context.Background())
	}

	return maps.Map{
		"id":        task.Id,
		"on":        task.On,
		"name":      task.Name,
		"script":    task.Script,
		"isBooting": task.IsBooting,
		"isManual":  task.IsManual,
		"schedules": lists.Map(task.Schedule, func(k int, v interface{}) interface{} {
			return v.(*agents.ScheduleConfig).Summary()
		}),
		"runTime": runTime,
	}
}
