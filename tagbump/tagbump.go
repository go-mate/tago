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
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func BumpGitTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	zaplog.LOG.Debug("bump-git-tag", zap.Int("version-base", versionBase))
	tagName, err := gcm.LatestGitTag()
	if err != nil {
		return false, erero.Wro(err)
	}
	if tagName == "" {
		return false, erero.New("no tag")
	}
	return BumpTagNum(gcm, tagName, "v", versionBase)
}

func BumpSubModuleTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	zaplog.LOG.Debug("bump-sub-module-tag", zap.Int("version-base", versionBase))
	subPath, err := gcm.GetSubPath()
	if err != nil {
		return false, erero.Wro(err)
	}
	if subPath == "" {
		return false, erero.New("not in sub-module path")
	}
	tagPrefix := filepath.Join(subPath, "v")
	tagRegexp := tagPrefix + "[0-9]*.[0-9]*.[0-9]*"
	return BumpTagMatchRegexp(gcm, tagPrefix, tagRegexp, versionBase)
}

func BumpMainTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	zaplog.LOG.Debug("bump-main-tag", zap.Int("version-base", versionBase))
	tagPrefix := "v"
	tagRegexp := tagPrefix + "[0-9]*.[0-9]*.[0-9]*"
	return BumpTagMatchRegexp(gcm, tagPrefix, tagRegexp, versionBase)
}

func BumpTagMatchRegexp(gcm *gitgo.Gcm, tagPrefix string, tagRegexp string, versionBase int) (bool, error) {
	zaplog.LOG.Debug("bump-match-regexp-tag", zap.String("tag-prefix", tagPrefix), zap.String("tag-regexp", tagRegexp))
	tagName, err := gcm.LatestGitTagMatchRegexp(tagRegexp)
	if err != nil {
		return false, erero.Wro(err)
	}
	if tagName == "" {
		return false, erero.Errorf("not match tag name with tag-prefix=((%s)) tag-regexp=((%s))", tagPrefix, tagRegexp)
	}
	return BumpTagNum(gcm, tagName, tagPrefix, versionBase)
}

func BumpTagNum(gcm *gitgo.Gcm, tagName string, tagPrefix string, versionBase int) (bool, error) {
	if rese.C1(gcm.GitCommitHash(tagName)) == rese.C1(gcm.GitCommitHash("main")) {
		if versionBase <= 1 {
			if !chooseConfirm("do you want to push the old tag? " + tagName) {
				return false, nil
			}
		}
		result, err := gcm.PushTag(tagName).ShowDebugMessage().Result()
		if err != nil {
			zaplog.SUG.Debugln(string(result))
			return false, erero.Wro(err)
		}
		return true, nil
	}
	zaplog.LOG.Info("old-tag-name", zap.String("tag", tagName))
	tagRegexp := `^` + regexp.QuoteMeta(tagPrefix) + `(\d+)\.(\d+)\.(\d+)$`
	zaplog.LOG.Info("check-tag-name-format-with-regexp", zap.String("regexp", tagRegexp))
	matches := regexp.MustCompile(tagRegexp).FindStringSubmatch(tagName)
	if len(matches) != 4 {
		return false, erero.New("no match")
	}
	vAx := done.VCE(strconv.Atoi(matches[1])).Done()
	vBx := done.VCE(strconv.Atoi(matches[2])).Done()
	vCx := done.VCE(strconv.Atoi(matches[3])).Done()
	if versionBase >= 2 {
		mustLessThan(vBx, versionBase)
		mustLessThan(vCx, versionBase)
	}
	vCx++
	if versionBase >= 2 { //就是当是0或者1的时候，就表示标签是不自动进位的，只有大于等于2时才自动进位标签
		if vCx >= versionBase {
			vCx = 0
			vBx++
		}
		if vBx >= versionBase {
			vBx = 0
			vAx++
		}
	}
	newTagName := fmt.Sprintf("%s%d.%d.%d", tagPrefix, vAx, vBx, vCx)
	zaplog.LOG.Info("new-tag-name", zap.String("tag", newTagName))
	if versionBase <= 1 {
		if !chooseConfirm("do you want to set this new tag? " + newTagName) {
			return false, nil
		}
	}
	if true {
		result, err := gcm.Tag(newTagName).ShowDebugMessage().Result()
		if err != nil {
			zaplog.SUG.Debugln(string(result))
			return false, erero.Wro(err)
		}
	}
	if versionBase <= 1 {
		if !chooseConfirm("do you want to push the new tag? " + newTagName) {
			return false, nil
		}
	}
	if true {
		result, err := gcm.PushTag(newTagName).ShowDebugMessage().Result()
		if err != nil {
			zaplog.SUG.Debugln(string(result))
			return false, erero.Wro(err)
		}
	}
	return true, nil
}

func chooseConfirm(msg string) bool {
	// 用于存储用户的回答
	var input bool

	// 定义确认问题
	prompt := &survey.Confirm{
		Message: msg,
		Default: true, // 默认值，如果用户直接按回车
	}

	// 运行提示并捕获用户输入的内容
	done.Done(survey.AskOne(prompt, &input))

	// 输出用户的回答
	if input {
		fmt.Println("You chose Yes")
		return true
	} else {
		fmt.Println("You chose Not")
		return false
	}
}

func mustLessThan(v0, v1 int) {
	if v0 >= v1 {
		zaplog.ZAPS.Skip1.LOG.Panic("V0 NOT LESS THAN V1", zap.Int("v0", v0), zap.Int("v1", v1))
	}
}
