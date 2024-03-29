package main

import (
	"context"
	"fmt"
	"github.com/dstgo/contrib/util/cfgx"
	"github.com/dstgo/maxwell/app"
	"github.com/dstgo/maxwell/app/conf"
	"github.com/spf13/cobra"
	"log/slog"
	"os/signal"
	"syscall"
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

		// register os signal listener
		ctx := context.Background()
		notifyContext, cancelFunc := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
		defer cancelFunc()

		// initialize app
		server, err := app.NewApp(notifyContext, &appconf)
		if err != nil {
			return err
		}

		done := make(chan struct{})
		go func() {
			slog.Info(fmt.Sprintf("server is listening at %s", appconf.Server.Address))
			server.ListenAndServe()
			done <- struct{}{}
			close(done)
		}()

		select {
		case <-notifyContext.Done():
			slog.InfoContext(ctx, "received os signal, ready to shutdown")
		case <-done:
			slog.InfoContext(ctx, "shutdown")
		}
		server.Shutdown(ctx)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "f", "conf.yaml", "app configuration file")
}

func main() {
	rootCmd.Execute()
}
