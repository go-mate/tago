// Package main: tago CLI tool for Git tag version management
// Provides smart tag creation, bumping, and Git repository versioning operations
// Supports main project tags, submodule tags, and version base configuration
//
// main: tago CLI 工具，用于 Git 标签版本管理
// 提供智能标签创建、升级和 Git 仓库版本管理操作
// 支持主项目标签、子模块标签和版本基数配置
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

// main initializes and executes the tago command with Git tag management
// Sets up Git repository interface and provides tag listing and bumping capabilities
// Displays sorted tags when run without subcommands
//
// main 初始化并执行 tago 命令，进行 Git 标签管理
// 设置 Git 仓库接口并提供标签列表和升级功能
// 当不带子命令运行时显示排序后的标签
func main() {
	// Get current working DIR as project root
	// 获取当前工作目录作为项目根目录
	workRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(workRoot))

	// Initialize Git command manager with debug mode
	// 初始化带调试模式的 Git 命令管理器
	gcm := gitgo.New(workRoot).WithDebug()

	// Create root command for tago CLI
	// 为 tago CLI 创建根命令
	rootCmd := cobra.Command{
		Use:   "tago",
		Short: "Git tag version management tool",
		Long:  "tago provides smart Git tag creation, bumping, and version management operations",
		Run: func(cmd *cobra.Command, args []string) {
			// Display sorted Git tags when no subcommand is provided
			// 当没有提供子命令时显示排序的 Git 标签
			eroticgo.BLUE.ShowMessage(rese.V1(gcm.SortedGitTags()))
		},
	}

	// Add tag bump command with all subcommands
	// 添加带所有子命令的标签升级命令
	rootCmd.AddCommand(newGitTagBumpCmd(gcm))

	// Execute CLI application
	// 执行 CLI 应用程序
	must.Done(rootCmd.Execute())
}

// newGitTagBumpCmd creates the main tag bump command with version base support
// Provides automatic tag version bumping with configurable version base numbers
// Supports main project and submodule tag management subcommands
//
// newGitTagBumpCmd 创建主要标签升级命令，支持版本基数
// 提供可配置版本基数的自动标签版本升级
// 支持主项目和子模块标签管理子命令
func newGitTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	// Version base configuration for automatic carry-over
	// 用于自动进位的版本基数配置
	var versionBase = 0

	// Create main bump command
	// 创建主要的 bump 命令
	tagBumpCmd := &cobra.Command{
		Use:   "bump",
		Short: "Bump Git tag version with version base support",
		Long:  "Automatically increment Git tag version with configurable version base (1/10/100) for version control",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate that no unexpected arguments are provided
			// 验证没有提供意外的参数
			if len(args) > 0 {
				eroticgo.PINK.ShowMessage("UNKNOWN")
				zaplog.LOG.Warn("unknown-subcommand-param-args", zap.Strings("args", args))
				os.Exit(1)
			}

			// Execute tag bump operation and display result
			// 执行标签升级操作并显示结果
			success := rese.V1(tagbump.BumpGitTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}
	// Configure version base flag for tag bump command
	// 为标签升级命令配置版本基数标志
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100 for automatic version carry-over")

	// Add main project and submodule subcommands
	// 添加主项目和子模块子命令
	tagBumpCmd.AddCommand(newMainTagBumpCmd(gcm))
	tagBumpCmd.AddCommand(newSubModuleTagBumpCmd(gcm))
	return tagBumpCmd
}

// newMainTagBumpCmd creates command for main project tag version bumping
// Handles main project tag operations with version base configuration
// Used when working in the root DIR of the main project
//
// newMainTagBumpCmd 创建主项目标签版本升级命令
// 处理带版本基数配置的主项目标签操作
// 在主项目根目录中使用
func newMainTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	// Version base configuration for main project tags
	// 主项目标签的版本基数配置
	var versionBase = 0

	// Create main project tag bump command
	// 创建主项目标签升级命令
	tagBumpCmd := &cobra.Command{
		Use:   "main",
		Short: "Bump main project Git tag version",
		Long:  "Bump version tag for the main project with configurable version base system",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute main project tag bump and display result
			// 执行主项目标签升级并显示结果
			success := rese.V1(tagbump.BumpMainTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}

	// Configure version base flag for main command
	// 为 main 命令配置版本基数标志
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100 for automatic version carry-over")
	return tagBumpCmd
}

// newSubModuleTagBumpCmd creates command for submodule tag version bumping
// Handles submodule-specific tag operations with path prefix support
// Requires execution from within a submodule DIR, not the main project root
//
// newSubModuleTagBumpCmd 创建子模块标签版本升级命令
// 处理带路径前缀支持的子模块特定标签操作
// 需要从子模块目录内执行，而非主项目根目录
func newSubModuleTagBumpCmd(gcm *gitgo.Gcm) *cobra.Command {
	// Version base configuration for submodule tags
	// 子模块标签的版本基数配置
	var versionBase = 0

	// Create submodule tag bump command
	// 创建子模块标签升级命令
	tagBumpCmd := &cobra.Command{
		Use:   "sub-module",
		Short: "Bump submodule Git tag version",
		Long:  "Bump version tag for submodule with path prefix, must be run from within submodule DIR",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate we are inside a submodule, not at project root
			// 验证我们在子模块内部，而非项目根目录
			must.Different(rese.C1(os.Getwd()), rese.C1(gcm.GetTopPath()))

			// Execute submodule tag bump and display result
			// 执行子模块标签升级并显示结果
			success := rese.V1(tagbump.BumpSubModuleTag(gcm, versionBase))
			if success {
				eroticgo.BLUE.ShowMessage("SUCCESS")
			} else {
				eroticgo.PINK.ShowMessage("FAILURE")
			}
		},
	}

	// Configure version base flag for submodule command
	// 为子模块命令配置版本基数标志
	tagBumpCmd.Flags().IntVarP(&versionBase, "vb", "b", 0, "version-base-num: 1/10/100 for automatic version carry-over")
	return tagBumpCmd
}
