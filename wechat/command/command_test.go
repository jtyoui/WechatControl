package command_test

import (
	"fmt"
	"testing"

	"github.com/jtyoui/WechatControl/wechat/command"
)

func TestAI(t *testing.T) {
	result, err := command.OpenAI("1+1等于几？")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(result)
}
