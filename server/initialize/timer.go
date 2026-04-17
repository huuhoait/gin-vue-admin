package initialize

import (
	"fmt"
	"github.com/huuhoait/gin-vue-admin/server/task"

	"github.com/robfig/cron/v3"

	"github.com/huuhoait/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// CleanDBscheduled task
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // scheduled taskmethodSetAttaskFilePackageIn
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "SetWhenCleanDatabase[Log, Blacklist]content", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// Othersscheduled taskSetAtHere ReferenceUpperSideUsemethod

		//_, err := global.GVA_Timer.AddTaskByFunc("scheduled taskIdentifier", "cornExpression", func() {
		//	ToolBodyExecutecontent...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}
	}()
}
