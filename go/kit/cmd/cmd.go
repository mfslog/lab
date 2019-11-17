package cmd

import (
	"fmt"
	"github.com/mfslog/lab/go/kit/application"
	"github.com/mfslog/lab/go/kit/config"
	"github.com/mfslog/lab/go/kit/domain/service"
	transportGRPC "github.com/mfslog/lab/go/kit/interfaces/grpc"
	"github.com/mfslog/lab/go/kit/idl/account/acctsvc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	mgoDB "github.com/mfslog/lab/go/kit/infrastructure/account/mgo"
	"github.com/mfslog/lab/go/kit/logger"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"net"
	"os"
	"time"
)

var (
	RootCmd = &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			run()
			return
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %d.%d.%d.%d\n", MAJOR, MINOR, PATCH, BUILD)
			os.Exit(0)
		},
	}
)

var (
	grpcSrv *grpc.Server
	defaultMgo *mgo.Session
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func run() {
	var (
		err error
	)
	err = config.Init()
	if err != nil {
		panic(err)
	}

	dialInfo := &mgo.DialInfo{
		Addrs:     viper.GetStringSlice("mgoAddr"),
		Direct:    false,
		Timeout:   time.Second * 3,
		PoolLimit: viper.GetInt("mgoLimit"),
		Username:  viper.GetString("mgoUser"),
		Password:  viper.GetString(" mgoPassword"),
	}

	defaultMgo, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		os.Exit(-1)
	}

	repo := mgoDB.NewRepository(defaultMgo)
	svc := service.NewService(logger.Logger,repo)
	endPoint := application.New(svc, logger.Logger)

	srv := transportGRPC.NewGRPCServer(endPoint)

	baseSrv := grpc.NewServer()
	acctsvc.RegisterAccountSvcServer(baseSrv, srv)
	l,err := net.Listen("tcp",fmt.Sprintf(":%d",viper.GetInt("grpcPort")))
	baseSrv.Serve(l)
}
