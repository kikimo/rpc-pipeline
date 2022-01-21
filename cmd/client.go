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

	// "github.com/gogf/gf/os/glog"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type ClientOpts struct {
	usePipeline bool
	loops       int
}

var clientOpts ClientOpts

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		glog.Info("client called")
		if !clientOpts.usePipeline {
			glog.Infof("runing sync client...")
			runSyncClient(clientOpts.loops)
		} else {
			glog.Infof("runing pipeline client...")
			runPiplineClient(clientOpts.loops)
		}
		// piplineClient(loops)
	},
}

func runSyncClient(loop int) {
	client, err := rpc.Dial("tcp", "localhost:8848")
	if err != nil {
		glog.Fatalf("dialing: %+v", err)
	}

	for i := 0; i < loop; i++ {
		var reply string
		err = client.Call("HelloService.Hello", fmt.Sprintf("client %d", i), &reply)
		if err != nil {
			glog.Fatal(err)
		}

		glog.V(1).Infof("reply: %+v", reply)
	}
}

// in our test, pipeline is three times as good as sync client
// for 524288 rpc invoke, pipe take:
//
// real    0m22.410s
// user    0m10.519s
// sys     0m44.574s
//
// for sync clients:
// real    1m4.834s
// user    0m12.472s
// sys     0m47.499s
//
// and for concurrent clients:
// real    0m28.981s
// user    0m14.819s
// sys     1m17.591s
// and sync client takes:

func runPiplineClient(loop int) {
	client, err := rpc.Dial("tcp", "localhost:8848")
	if err != nil {
		glog.Fatalf("dialing: %+v", err)
	}

	var wg sync.WaitGroup
	wg.Add(loop)

	for i := 0; i < loop; i++ {
		var reply string
		helloCall := client.Go("HelloService.Hello", fmt.Sprintf("client %d", i), &reply, nil)
		go func(call *rpc.Call, replyp *string) {
			<-call.Done
			glog.V(1).Infof("reply: %+v", reply)
			wg.Done()
		}(helloCall, &reply)
	}
	wg.Wait()
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.
	clientCmd.Flags().BoolVarP(&clientOpts.usePipeline, "usePipeline", "", false, "use pipeline, defualt false")
	clientCmd.Flags().IntVarP(&clientOpts.loops, "loops", "", 65536, "loops")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
