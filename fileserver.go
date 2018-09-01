package main

import (
	"os"

	"github.com/charleszheng44/fileserver/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("fail to start fileserver: %v", err)
		os.Exit(1)
	}
}
