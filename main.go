package main

import (
	"BaselineCheck/client/baselinelinux"
	"BaselineCheck/server"
	"BaselineCheck/server/compliance"
	"BaselineCheck/server/config"
	"BaselineCheck/server/repository"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config.yaml", "配置文件路径")
	rootCmd.PersistentFlags().StringVarP(&pushUrl, "push_url", "p", "http://127.0.0.1:9527/check", "push 地址")
	rootCmd.AddCommand(serverCmd, checkCmd)
	rootCmd.Execute()
}

var (
	configFilePath string
	pushUrl        string

	rootCmd = &cobra.Command{
		Use:   "BaselineCheck",
		Short: "BaselineCheck is a tool for baseline check",
		Long:  "BaselineCheck is a tool for baseline check",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	serverCmd = &cobra.Command{
		Use:   "start",
		Short: "start the server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化配置
			conf, err := config.InitConfig(configFilePath)
			if err != nil {
				panic(err)
			}
			db, err := repository.NewDb(conf.Repo)
			if err != nil {
				panic(err)
			}
			db.AutoMigrate(&compliance.ComplianceResult{}, &compliance.ComplianceDetails{})
			comparableRepo := compliance.NewRepo(db)
			comparableHandler := compliance.NewHandler(comparableRepo)
			server.NewServer(conf, comparableHandler).Start()

		},
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "check the baseline",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			baselinelinux.Run(pushUrl)
		},
	}
)
