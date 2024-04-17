package main

import (
	"context"
	"fmt"
	"github.com/dstgo/maxwell/internal/app"
	"github.com/dstgo/maxwell/internal/app/conf"
	"github.com/dstgo/maxwell/pkg/cfgx"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var (
	Version    string
	BuildTime  string
	ConfigFile string
)

var rootCmd = &cobra.Command{
	Use:          "maxwell [commands] [-flags]",
	Short:        "maxwell is a opensource personal don‘t starve together web panel",
	Long:         "maxwell is a opensource personal don‘t starve together web panel",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// load config file
		var appconf conf.AppConf
		if err := cfgx.LoadConfigAndMapTo(ConfigFile, &appconf); err != nil {
			return err
		}
		appconf.Version = Version
		appconf.BuildTime = BuildTime

		// print banner
		if err := app.PrintBanner(os.Stderr); err != nil {
			return err
		}

		appconf.Log.Prompt = "[maxwell]"
		// initialize app logger
		logger, err := app.NewLogger(appconf.Log)
		if err != nil {
			return err
		}
		defer logger.Close()

		// set it to the default logger
		slog.SetDefault(logger.Slog())
		slog.Info(fmt.Sprintf("logging in level: %s", appconf.Log.Level))

		// this is the root context for the whole program
		rootCtx := context.Background()

		// initialize app
		server, err := app.NewApp(rootCtx, &appconf)
		if err != nil {
			return err
		}

		// run the server
		return server.Spin()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "f", "conf.yaml", "app configuration file")
}

func main() {
	rootCmd.Execute()
}
