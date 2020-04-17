package cmd

import (
	"github.com/always-waiting/cobra-canal/cmd/consumer"
	"github.com/always-waiting/cobra-canal/cmd/filter"
	"github.com/always-waiting/cobra-canal/cmd/monitor"
	"github.com/always-waiting/cobra-canal/cmd/pipeline"
	"github.com/always-waiting/cobra-canal/cmd/transfer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdb_cobra",
	Short: "mysql监控命令组",
	Long: `进行mysql的binlog监控，可以根据不同需要，开发不同的
过滤规则以及下游消费`,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("port", "", "程序监听的端口号")
	rootCmd.PersistentFlags().String("pid", "", "程序的pid号")
	rootCmd.PersistentFlags().String("host", "127.0.0.1", "监控程序运行地址")
	rootCmd.PersistentFlags().Bool("pretty", false, "返回结果是否格式化")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(monitor.RootCmd)
	rootCmd.AddCommand(filter.RootCmd)
	rootCmd.AddCommand(transfer.RootCmd)
	rootCmd.AddCommand(consumer.RootCmd)
	rootCmd.AddCommand(pipeline.RootCmd)
}

const (
	SUCCESS1       = "%s successfully\n"
	SERVICE_PREFIX = "cobra."
)
