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
	"go.uber.org/zap"
)

func main() {
	workRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workRoot))

	gcm := gitgo.New(workRoot).WithDebug()

	rootCmd := cobra.Command{
		Use:   "tago",
		Short: "tago",
		Long:  "tago",
		Run: func(cmd *cobra.Command, args []string) {
			eroticgo.BLUE.ShowMessage(rese.V1(gcm.SortedGitTags()))
		},
	}

	rootCmd.AddCommand(newGitTagBumpCmd(gcm))

	must.Done(rootCmd.Execute())
}

func newGitTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	var versionBase = 0
	tagBumpCmd := &cobra.Command{
		Use:   "bump",
		Short: "tago bump",
		Long:  "tago bump",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				eroticgo.PINK.ShowMessage("UNKNOWN")
				zaplog.LOG.Warn("unknown-subcommand-param-args", zap.Strings("args", args))
				os.Exit(1)
			}

			success := rese.V1(tagbump.BumpGitTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100")

	tagBumpCmd.AddCommand(newMainTagBumpCmd(gcm))
	tagBumpCmd.AddCommand(newSubModuleTagBumpCmd(gcm))
	return tagBumpCmd
}

func newMainTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	var versionBase = 0
	tagBumpCmd := &cobra.Command{
		Use:   "main",
		Short: "tago bump main",
		Long:  "tago bump main",
		Run: func(cmd *cobra.Command, args []string) {
			success := rese.V1(tagbump.BumpMainTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100")
	return tagBumpCmd
}

func newSubModuleTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	var versionBase = 0
	tagBumpCmd := &cobra.Command{
		Use:   "sub-module",
		Short: "tago bump sub-module",
		Long:  "tago bump sub-module",
		Run: func(cmd *cobra.Command, args []string) {
			must.Different(rese.C1(os.Getwd()), rese.C1(gcm.GetTopPath()))

			success := rese.V1(tagbump.BumpSubModuleTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100")
	return tagBumpCmd
}
