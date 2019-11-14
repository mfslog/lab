package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use:"",
		Run: func(cmd *cobra.Command, args []string) {
			
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









func init(){
	RootCmd.AddCommand(versionCmd)
}

func Execute()error{
	return RootCmd.Execute()
}
