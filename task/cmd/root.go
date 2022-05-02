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
	rootCmd.PersistentFlags().StringP("config", "c", "conf.yaml", "config file")
}

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg = config.Load[handler.Config, service.Config](cmd.Flag("config").Value.String())
		process.SetEnv(cfg.Env)
		cfg.Logger.Topic = regexp.MustCompile(`\W`).ReplaceAllString(cmd.Use, "_")
		logger.Initialize(&cfg.Logger)
		ctx := core.NewContext(process.GetHostname(), cmd.Use, process.GetIP(), process.GetEnv())
		cmd.SetContext(ctx)
		logger.FromContext(ctx).Info("start")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		process.DoCleanup()
		logger.FromContext(cmd.Context()).Info("stop")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
