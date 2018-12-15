package cmd

import (
	"github.com/charleszheng44/fileserver/pkg/fileserver"
	"github.com/spf13/cobra"
)

var (
	uploadPath string
	addr       string
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fileserver",
		Short: "fileserver start a basic http fileserver",
		Long:  "fileserver support basic upload and download operations",
	}

	rootCmd.PersistentFlags().StringVar(&addr, "addr", "127.0.0.1:9344",
		"network address fileserver will liseten to. (<ip:port>)")
	rootCmd.PersistentFlags().StringVar(&uploadPath, "path", "/data/workspace",
		"path to the directory that stores uploaded data.")
	fs := fileserver.NewFileServer(addr, uploadPath, uploadPath)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		fs.StartServer(addr, uploadPath)
	}

	return rootCmd
}
