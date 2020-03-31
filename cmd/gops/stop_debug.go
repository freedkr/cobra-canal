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
	port, err := helps.GetPort(cmd)
	if err != nil {
		panic(err)
	}
	host, _ := cmd.Flags().GetString("host")
	pretty, _ := cmd.Flags().GetBool("pretty")
	var Addr string
	if pretty {
		Addr = fmt.Sprintf("http://%s:%s/gops/debug/stop?pretty", host, port)
	} else {
		Addr = fmt.Sprintf("http://%s:%s/gops/debug/stop", host, port)
	}
	req, _ := http.NewRequest("GET", Addr, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	helps.CmdPrint(resp)
}
