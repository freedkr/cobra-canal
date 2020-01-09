package gops

import (
	"fmt"
	"net/http"

	"github.com/always-waiting/cobra-canal/helps"
	"github.com/spf13/cobra"
)

var stopDebugCmd = &cobra.Command{
	Use:     "stopDebug",
	Short:   "开启debug模式",
	Version: "2.0.0",
	Run:     stopDebugCmdRun,
}

func stopDebugCmdRun(cmd *cobra.Command, args []string) {
	pid, _ := cmd.Flags().GetString("pid")
	port, err := helps.GetPortByPid(pid)
	if err != nil {
		panic(err)
	}
	Addr := fmt.Sprintf("http://127.0.0.1:%s/gops/debug/stop", port)
	req, _ := http.NewRequest("GET", Addr, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf(SUCCESS1, "stopDebug")
}
