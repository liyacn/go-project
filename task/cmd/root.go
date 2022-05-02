package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"project/pkg/config"
	"project/pkg/core"
	"project/pkg/logger"
	"project/pkg/process"
	"project/task/internal/handler"
	"project/task/internal/service"
	"regexp"
)

var cfg *config.Config[handler.Config, service.Config]

func init() {
	cobra.OnInitialize(func() {
		cfg = config.Load[handler.Config, service.Config]()
	})
	rootCmd.AddCommand(
		cronjobCmd,
		avatarToCdnCmd,
	)
}

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		process.SetEnv(cfg.App.Env)
		cfg.App.Logger.Topic = regexp.MustCompile(`\W`).ReplaceAllString(cmd.Use, "_")
		logger.Setup(&cfg.App.Logger)
		ctx := core.NewContext(process.GetHostname(), cmd.Use, process.GetIP(), process.GetEnv())
		cmd.SetContext(ctx)
		logger.FromContext(ctx).Info("start", nil, nil)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logger.FromContext(cmd.Context()).Info("stop", nil, nil)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
