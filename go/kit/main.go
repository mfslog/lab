package  main

import (
	"github.com/mfslog/lab/go/kit/cmd"
	"os"
	"os/signal"
)


func main(){
	cmd.Execute()
	wait()
}

func wait(){
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}