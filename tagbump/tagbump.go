package tagbump

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-xlan/gitgo"
	"github.com/yyle88/done"
	"github.com/yyle88/erero"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func BumpGitTag(gcm *gitgo.Gcm, versionBase int) (bool, error) {
	zaplog.SUG.Debugln("version-base", versionBase)
	tagName, err := gcm.LatestGitTag()
	if err != nil {
		return false, erero.Wro(err)
	}
	if tagName == "" {
		zaplog.SUG.Debugln("no tag")
		return false, erero.New("no tag")
	}
	tagHash, err := gcm.GitCommitHash(tagName)
	if err != nil {
		return false, erero.Wro(err)
	}
	if tagHash == "" {
		return false, erero.New("impossible")
	}
	mainHash, err := gcm.GitCommitHash("main")
	if err != nil {
		return false, erero.Wro(err)
	}
	if mainHash == "" {
		return false, erero.New("impossible")
	}
	if tagHash == mainHash {
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
		return false, nil
	}
	zaplog.SUG.Debugln(tagName)
	matches := regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)$`).FindStringSubmatch(tagName)
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
	newTagName := fmt.Sprintf("v%d.%d.%d", vAx, vBx, vCx)
	zaplog.SUG.Debugln(newTagName)
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
