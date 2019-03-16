package cmd

const CMDSuccess  = 1
const CMDError  = 2

const (
	CMDRegister = 1 // 注册命令
	CMDRegisterAck = 2 // 注册返回命令
)

const (
	RegisterError = 2 // 注册错误！
	RegisterClientExist = 3 // 客户端已存在！
)
