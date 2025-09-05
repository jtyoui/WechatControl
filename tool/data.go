package tool

import "errors"

// 记录所有的变量数据

type Status struct {
	Result string // 返回的正常结果
	OK     bool   // 执行是否成功
	Error  error  // 异常信息
	retry  int    // 重试次数
}

// Join 加入错误信息
func (s *Status) Join(err error) {
	s.Error = errors.Join(s.Error, err)
}

// CanRetry 是否还能重试
func (s *Status) CanRetry() bool {
	return s.retry < 3
}

// Stop 是否要停止
func (s *Status) Stop(err error) (ok bool) {
	s.Join(err)
	if s.CanRetry() {
		s.AddRetry()
		return
	}
	ok = true
	return
}

// AddRetry 重试一次
func (s *Status) AddRetry() {
	s.retry += 1
}

type WxInfo struct {
	Data    string // OCR识别的数据
	Pusher  Status // 推送
	Command Status // 执行
	Over    bool   // 全部执行完
}

var (
	WxNormal = make(chan WxInfo, 100) // ocr识别的数据
	WxError  = make(chan WxInfo, 100)
)
