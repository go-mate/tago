// Package tagbump: Git tag version bumping engine with intelligent version management
// Provides comprehensive Git tag operations including main project and submodule support
// Implements semantic versioning with configurable version base systems and interactive confirmation
//
// tagbump: Git 标签版本升级引擎，具备智能版本管理
// 提供全面的 Git 标签操作，包括主项目和子模块支持
// 实现语义化版本控制，支持可配置版本基数系统和交互式确认
package tagbump

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-xlan/gitgo"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/must/mustnum"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// BumpGitTag bumps the latest Git tag version with version base support
// Retrieves the most recent tag and increments its version using specified base system
// Returns success status and handles cases where no tags exist
//
// BumpGitTag 使用版本基数支持升级最新的 Git 标签版本
// 获取最新标签并使用指定基数系统递增其版本
// 返回成功状态并处理不存在标签的情况
func BumpGitTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	// Log operation parameters for debugging
	// 记录操作参数用于调试
	zaplog.LOG.Debug("BUMP-GIT-TAG", zap.Int("version-base", versionBase))

	// Retrieve the latest Git tag from repository
	// 从仓库获取最新的 Git 标签
	tagName, err := gcm.LatestGitTag()
	if err != nil {
		return false, erero.Wro(err)
	}

	// Validate that at least one tag exists
	// 验证至少存在一个标签
	if tagName == "" {
		return false, erero.New("no tag")
	}

	// Delegate to core version bumping logic
	// 委托给核心版本升级逻辑
	return BumpTagNum(gcm, tagName, "v", versionBase)
}

// BumpSubModuleTag bumps Git tag version for submodule with path prefix
// Constructs submodule-specific tag prefix and applies version bumping logic
// Ensures operation is performed within a valid submodule context
//
// BumpSubModuleTag 使用路径前缀升级子模块的 Git 标签版本
// 构造子模块特定的标签前缀并应用版本升级逻辑
// 确保操作在有效的子模块上下文中执行
func BumpSubModuleTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	// Log submodule tag operation parameters
	// 记录子模块标签操作参数
	zaplog.LOG.Debug("BUMP-SUB-MODULE-TAG", zap.Int("version-base", versionBase))

	// Get current submodule path relative to main project
	// 获取相对于主项目的当前子模块路径
	subPath, err := gcm.GetSubPath()
	if err != nil {
		return false, erero.Wro(err)
	}

	// Ensure we are inside a submodule DIR
	// 确保我们在子模块目录内
	if subPath == "" {
		return false, erero.New("not in sub-module path")
	}

	// Construct submodule-specific tag prefix with path
	// 构建带路径的子模块特定标签前缀
	tagPrefix := filepath.Join(subPath, "v")
	tagRegexp := tagPrefix + "[0-9]*.[0-9]*.[0-9]*"

	// Apply regexp-based tag matching and bumping
	// 应用基于正则表达式的标签匹配和升级
	return BumpTagMatchRegexp(gcm, tagPrefix, tagRegexp, versionBase)
}

// BumpMainTag bumps Git tag version for main project repository
// Uses standard 'v' prefix for main project tags and applies semantic versioning
// Designed for main project root DIR operations
//
// BumpMainTag 升级主项目仓库的 Git 标签版本
// 使用标准的 'v' 前缀用于主项目标签并应用语义版本控制
// 设计用于主项目根目录操作
func BumpMainTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	// Log main project tag operation parameters
	// 记录主项目标签操作参数
	zaplog.LOG.Debug("BUMP-MAIN-TAG", zap.Int("version-base", versionBase))

	// Use standard 'v' prefix for main project tags
	// 主项目标签使用标准的 'v' 前缀
	tagPrefix := "v"
	tagRegexp := tagPrefix + "[0-9]*.[0-9]*.[0-9]*"

	// Apply regexp-based tag matching and bumping for main project
	// 为主项目应用基于正则表达式的标签匹配和升级
	return BumpTagMatchRegexp(gcm, tagPrefix, tagRegexp, versionBase)
}

// BumpTagMatchRegexp bumps Git tag version matching specified regular expression pattern
// Finds latest tag matching regexp pattern and applies version bumping logic
// Used for both main project and submodule tag operations with custom patterns
//
// BumpTagMatchRegexp 升级匹配指定正则表达式模式的 Git 标签版本
// 查找匹配正则表达式模式的最新标签并应用版本升级逻辑
// 用于主项目和子模块标签操作，支持自定义模式
func BumpTagMatchRegexp(gcm *gitgo.Gcm, tagPrefix string, tagRegexp string, versionBase int) (bool, error) {
	// Log regexp matching parameters for debugging
	// 记录正则匹配参数用于调试
	zaplog.LOG.Debug("BUMP-MATCH-REGEXP-TAG", zap.String("tag-prefix", tagPrefix), zap.String("tag-regexp", tagRegexp))

	// Find latest tag matching the specified regexp pattern
	// 查找匹配指定正则模式的最新标签
	tagName, err := gcm.LatestGitTagMatchRegexp(tagRegexp)
	if err != nil {
		return false, erero.Wro(err)
	}

	// Validate that a matching tag was found
	// 验证找到了匹配的标签
	if tagName == "" {
		return false, erero.Errorf("not match tag name with tag-prefix=((%s)) tag-regexp=((%s))", tagPrefix, tagRegexp)
	}

	// Delegate to core version bumping with found tag
	// 使用找到的标签委托给核心版本升级
	return BumpTagNum(gcm, tagName, tagPrefix, versionBase)
}

// BumpTagNum performs core semantic version incrementing with configurable version base
// Handles commit hash comparison, version parsing, increment logic, and tag creation/pushing
// Supports interactive confirmation for version base <= 1, auto mode for higher bases
//
// BumpTagNum 执行带可配置版本基数的核心语义版本递增
// 处理提交哈希比较、版本解析、递增逻辑和标签创建/推送
// 对版本基数 <= 1 支持交互式确认，更高基数支持自动模式
func BumpTagNum(gcm *gitgo.Gcm, tagName string, tagPrefix string, versionBase int) (bool, error) {
	// Create configuration and delegate to config-based implementation
	// 创建配置并委托给基于配置的实现
	config := &BumpConfig{
		TagName:     tagName,
		TagPrefix:   tagPrefix,
		VersionBase: versionBase,
		AutoConfirm: false,
		SkipGitPush: false,
	}
	return BumpTag(gcm, config)
}

// BumpConfig contains configuration options for Git tag version bumping
// Encapsulates tag prefix, version base, and flow control parameters
//
// BumpConfig 包含 Git 标签版本升级的配置选项
// 封装标签前缀、版本基数和流程控制参数
type BumpConfig struct {
	// Basic configuration
	// 基础配置
	TagName     string // Current tag name to bump from // 要升级的当前标签名
	TagPrefix   string // Tag prefix (e.g., "v", "release-") // 标签前缀（如 "v", "release-"）
	VersionBase int    // Version base for carry-over (0/1 = interactive, >=2 = auto) // 进位的版本基数（0/1 = 交互式，>=2 = 自动）

	// Testing and automation options
	// 测试和自动化选项
	AutoConfirm bool // Auto confirm operation // 自动确认操作
	SkipGitPush bool // Skip pushing to remote // 跳过推送远程
}

// BumpTag performs core semantic version incrementing with flexible configuration
// Handles commit hash comparison, version parsing, increment logic, and tag creation/pushing
// Uses BumpConfig structure for enhanced testability and future extensibility
//
// BumpTag 执行带灵活配置的核心语义版本递增
// 处理提交哈希比较、版本解析、递增逻辑和标签创建/推送
// 使用 BumpConfig 结构提供增强的可测试性和未来扩展性
func BumpTag(gcm *gitgo.Gcm, config *BumpConfig) (bool, error) {
	zaplog.SUG.Infoln("STARTING-BUMP-TAG", neatjsons.S(config))

	// Compare commit hashes to check if tag is already at HEAD
	// 比较提交哈希检查标签是否已在 HEAD 位置
	tagCommitHash := rese.C1(gcm.GitCommitHash(config.TagName))
	topCommitHash := rese.C1(gcm.GitCommitHash("main"))

	zaplog.LOG.Debug("COMMIT-HASH-COMPARISON",
		zap.String("tag-commit", tagCommitHash),
		zap.String("top-commit", topCommitHash),
	)

	if tagCommitHash == topCommitHash {
		// Tag is already at current commit, just push existing tag
		// 标签已在当前提交，只需推送现有标签
		zaplog.LOG.Info("TAG-ALREADY-AT-HEAD", zap.String("tag", config.TagName))

		// Check if we should proceed with pushing existing tag
		// 检查是否应该继续推送现有标签
		if !shouldConfirm(config, "do you want to push the old tag? "+config.TagName) {
			zaplog.LOG.Info("USER-DECLINED-PUSH-EXISTING-TAG")
			return false, nil
		}

		// Skip push if configured
		// 如果配置了则跳过推送
		if config.SkipGitPush {
			zaplog.LOG.Info("SKIPPING-PUSH-EXISTING-TAG", zap.String("tag", config.TagName))
			return true, nil
		}
		// Push existing tag to remote repository
		// 推送现有标签到远程仓库
		zaplog.LOG.Info("PUSHING-EXISTING-TAG", zap.String("tag", config.TagName))
		result, err := gcm.PushTag(config.TagName).ShowDebugMessage().Result()
		if err != nil {
			zaplog.SUG.Debugln(string(result))
			zaplog.LOG.Error("PUSH-EXISTING-TAG-FAILED", zap.Error(err))
			return false, erero.Wro(err)
		}
		zaplog.LOG.Info("SUCCESSFULLY-PUSHED-EXISTING-TAG", zap.String("tag", config.TagName))
		return true, nil
	}
	// Log current tag name for version bumping
	// 记录当前标签名用于版本升级
	zaplog.LOG.Info("OLD-TAG-NAME", zap.String("tag", config.TagName))

	// Construct regexp to parse semantic version format
	// 构造正则表达式来解析语义版本格式
	tagRegexp := `^` + regexp.QuoteMeta(config.TagPrefix) + `(\d+)\.(\d+)\.(\d+)$`
	zaplog.LOG.Info("CHECK-TAG-NAME-FORMAT-WITH-REGEXP", zap.String("regexp", tagRegexp))

	// Parse version components from tag name
	// 从标签名解析版本组件
	matches := regexp.MustCompile(tagRegexp).FindStringSubmatch(config.TagName)
	if len(matches) != 4 {
		zaplog.LOG.Error("TAG-FORMAT-MISMATCH",
			zap.String("tag", config.TagName),
			zap.String("regexp", tagRegexp),
		)
		return false, erero.New("no match")
	}
	// Extract major, minor, and patch version numbers
	// 提取主版本、次版本和补丁版本号
	vAx := done.VCE(strconv.Atoi(matches[1])).Done() // major version // 主版本号
	vBx := done.VCE(strconv.Atoi(matches[2])).Done() // minor version // 次版本号
	vCx := done.VCE(strconv.Atoi(matches[3])).Done() // patch version // 补丁版本

	zaplog.LOG.Debug("PARSED-VERSION-COMPONENTS",
		zap.Int("major", vAx),
		zap.Int("minor", vBx),
		zap.Int("patch", vCx))
	// Validate version components against version base for carry-over logic
	// 验证版本组件与版本基数的进位逻辑
	if config.VersionBase >= 2 {
		mustnum.Less(vBx, config.VersionBase)
		mustnum.Less(vCx, config.VersionBase)
	}

	// Increment patch version by default
	// 默认递增补丁版本
	vCx++
	zaplog.LOG.Debug("INCREMENTING-VERSION", zap.Int("new-patch", vCx))

	// Apply version carry-over logic for automatic mode
	// 为自动模式应用版本进位逻辑
	if config.VersionBase >= 2 { // When 0 or 1, no automatic version carry-over; >= 2 enables it // 当是0或者1时，标签不自动进位；>=2时启用自动进位
		// Check if patch version needs to carry over to minor
		// 检查补丁版本是否需要进位到次版本
		if vCx >= config.VersionBase {
			vCx = 0
			vBx++
		}
		// Check if minor version needs to carry over to major
		// 检查次版本是否需要进位到主版本
		if vBx >= config.VersionBase {
			vBx = 0
			vAx++
		}
	}
	// Construct new tag name with incremented version
	// 构造带递增版本的新标签名
	newTagName := fmt.Sprintf("%s%d.%d.%d", config.TagPrefix, vAx, vBx, vCx)
	zaplog.LOG.Info("NEW-TAG-NAME", zap.String("tag", newTagName))

	// Check if we should proceed with creating new tag
	// 检查是否应该继续创建新标签
	if !shouldConfirm(config, "do you want to set this new tag? "+newTagName) {
		zaplog.LOG.Info("USER-DECLINED-CREATE-TAG", zap.String("tag", newTagName))
		return false, nil
	}

	// Create new tag in local repository
	// 在本地仓库创建新标签
	zaplog.LOG.Info("CREATING-NEW-TAG", zap.String("tag", newTagName))
	result, err := gcm.Tag(newTagName).ShowDebugMessage().Result()
	if err != nil {
		zaplog.SUG.Debugln(string(result))
		zaplog.LOG.Error("TAG-CREATION-FAILED", zap.String("tag", newTagName), zap.Error(err))
		return false, erero.Wro(err)
	}
	zaplog.LOG.Info("SUCCESSFULLY-CREATED-TAG", zap.String("tag", newTagName))
	// Check if we should proceed with pushing new tag
	// 检查是否应该继续推送新标签
	if !shouldConfirm(config, "do you want to push the new tag? "+newTagName) {
		zaplog.LOG.Info("USER-DECLINED-PUSH-NEW-TAG", zap.String("tag", newTagName))
		return true, nil // Tag created but not pushed
	}

	// Skip push if configured
	// 如果配置了则跳过推送
	if config.SkipGitPush {
		zaplog.LOG.Info("SKIPPING-TAG-PUSH", zap.String("tag", newTagName))
		return true, nil
	}
	// Push new tag to remote repository
	// 推送新标签到远程仓库
	zaplog.LOG.Info("PUSHING-NEW-TAG", zap.String("tag", newTagName))
	result, err = gcm.PushTag(newTagName).ShowDebugMessage().Result()
	if err != nil {
		zaplog.SUG.Debugln(string(result))
		zaplog.LOG.Error("PUSH-NEW-TAG-FAILED", zap.String("tag", newTagName), zap.Error(err))
		return false, erero.Wro(err)
	}
	zaplog.LOG.Info("SUCCESSFULLY-PUSHED-NEW-TAG", zap.String("tag", newTagName))
	return true, nil
}

// shouldConfirm determines whether to proceed with an operation based on config
// Uses config flags to control confirmation flow
// Returns true to proceed, false to skip the operation
//
// shouldConfirm 根据配置确定是否继续操作
// 使用配置标志控制确认流程
// 返回 true 继续，false 跳过操作
func shouldConfirm(config *BumpConfig, message string) bool {
	// Auto-confirm if explicitly configured
	// 如果明确配置则自动确认
	if config.AutoConfirm {
		return true
	}

	// Interactive confirmation for version base <= 1
	// 对于版本基数 <= 1 使用交互式确认
	if config.VersionBase <= 1 {
		return chooseConfirm(message)
	}

	// Auto-proceed for higher version bases
	// 对于更高版本基数自动继续
	return true
}

// chooseConfirm displays interactive confirmation prompt with yes/no options
// Uses survey package to present user-friendly confirmation dialog
// Returns true for yes, false for no, with default value of true
//
// chooseConfirm 显示交互式确认提示，带有是/否选项
// 使用 survey 包呈现用户友好的确认对话框
// 返回 true 表示是，false 表示否，默认值为 true
func chooseConfirm(msg string) bool {
	// Variable to store user response
	// 用于存储用户的回答
	var input bool

	// Define confirmation prompt
	// 定义确认问题
	prompt := &survey.Confirm{
		Message: msg,
		Default: true, // Default value if user presses enter // 默认值，如果用户直接按回车
	}

	// Run prompt and capture user input
	// 运行提示并捕获用户输入的内容
	done.Done(survey.AskOne(prompt, &input))

	// Output user response
	// 输出用户的回答
	if input {
		fmt.Println("You chose Yes")
		return true
	}
	fmt.Println("You chose Not")
	return false
}
