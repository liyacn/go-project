package cmd

import (
	"github.com/spf13/cobra"
	"project/pkg/process"
	"project/pkg/worker"
	"project/task/internal/handler"
	"project/task/internal/service"
	"time"
)

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "示例：周期性任务/常驻任务组",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.Redis(&cfg.Service.Redis)))
		task1 := worker.NewIntervalTask(2*time.Minute, h.Example1)
		task2 := worker.NewParallelTask(4, h.Example2)
		process.Notify()
		task1.Stop()
		task2.Stop()
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}
