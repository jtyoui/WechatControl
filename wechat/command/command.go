package command

// Func 执行命令 输入一个文本 返回结果
type Func func(text string) (result string, err error)
