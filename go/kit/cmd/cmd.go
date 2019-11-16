package cmd

import (
	"fmt"
	"github.com/mfslog/lab/go/kit/config"
	"github.com/spf13/cobra"
	"github.com/go-kit/kit/transport/grpc"
	"os"
)



var (
	RootCmd = &cobra.Command{
		Use:"",
		Run: func(cmd *cobra.Command, args []string) {
			run()
			return
		},
	}

	versionCmd = &cobra.Command{
		Use:"version",
		Short:"show version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %d.%d.%d.%d\n",MAJOR,MINOR, PATCH, BUILD)
			os.Exit(0)
		},
	}
)

var (
	grpcSrv  *grpc.Server
)


func Execute(){
	if err := RootCmd.Execute(); err != nil{
		os.Exit(-1)
	}
}

func init(){
	RootCmd.AddCommand(versionCmd)
}


func run(){
	var (
		err error
	)
	err  = config.Init()
	if err != nil{
		panic(err)
	}

	grpcSrv = grpc.NewServer()


}

