/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/rpc"
	"sync"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type ConcurrentClientOpts struct {
	clients int
}

var concurrentClientOpts ConcurrentClientOpts

// concurrentClientCmd represents the concurrentClient command
var concurrentClientCmd = &cobra.Command{
	Use:   "concurrentClient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("concurrentClient called")
		runConncurrentClient()
	},
}

// not as good as pipeline and for rpcs like apache Thrift concurrent
// write to client might not work beacause theire clients are not thread
// safe
func runConncurrentClient() {
	client, err := rpc.Dial("tcp", "localhost:8848")
	if err != nil {
		glog.Fatalf("dialing: %+v", err)
	}

	var wg sync.WaitGroup
	wg.Add(concurrentClientOpts.clients)
	for i := 0; i < concurrentClientOpts.clients; i++ {
		go func(id int) {
			var reply string
			err = client.Call("HelloService.Hello", fmt.Sprintf("client %d", id), &reply)
			if err != nil {
				glog.Fatal(err)
			}

			glog.V(1).Infof("reply: %+v", reply)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func init() {
	rootCmd.AddCommand(concurrentClientCmd)
	concurrentClientCmd.Flags().IntVarP(&concurrentClientOpts.clients, "clients", "", 65536, "number of clients")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// concurrentClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// concurrentClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
