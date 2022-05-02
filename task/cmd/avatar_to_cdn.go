package cmd

import (
	"github.com/spf13/cobra"
	"project/model/queue"
	"project/pkg/gnsq"
	"project/pkg/process"
	"project/task/internal/handler"
	"project/task/internal/service"
)

var avatarToCdnCmd = &cobra.Command{
	Use:   "avatar:to:cdn",
	Short: "从临时链接获取用户头像转存到CND",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.Orm(&cfg.Service.Mysql)))
		consumer := gnsq.NewConsumer(cfg.Service.Nsq.Consumer, queue.AvatarToCdnTP, queue.DefaultCH, 4, h.AvatarToCdn)
		process.Notify()
		consumer.Stop()
	},
}

func init() {
	rootCmd.AddCommand(avatarToCdnCmd)
}
