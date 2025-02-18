package main

import (
	"os"

	"github.com/go-mate/tago/tagbump"
	"github.com/go-xlan/gitgo"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	projectPath := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(projectPath))

	gcm := gitgo.New(projectPath).WithDebug()

	rootCmd := cobra.Command{
		Use:   "tags",
		Short: "tags",
		Long:  "tags",
		Run: func(cmd *cobra.Command, args []string) {
			eroticgo.BLUE.ShowMessage(rese.V1(gcm.SortedGitTags()))
		},
	}

	rootCmd.AddCommand(newBumpTagCmd(gcm))

	must.Done(rootCmd.Execute())
}

func newBumpTagCmd(gcm *gitgo.Gcm) *cobra.Command {
	var versionBase = 0
	bumpTagCmd := &cobra.Command{
		Use:   "bump",
		Short: "bump",
		Long:  "bump",
		Run: func(cmd *cobra.Command, args []string) {
			success := rese.V1(tagbump.BumpGitTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}
	bumpTagCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100")
	return bumpTagCmd
}
