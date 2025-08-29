package pusher_test

import (
	"testing"

	"github.com/jtyoui/WechatControl/wechat/pusher"
)

func TestServerChan(t *testing.T) {
	if err := pusher.ServerChan("测试"); err != nil {
		t.Error(err)
		return
	}
}
