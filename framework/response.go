package framework

import (
	"encoding/json"
	"html/template"
)

// IResponse
type IResponse interface {
	// // Json 输出
	// Json(obj interface{}) IResponse

	// Jsonp 输出
	Jsonp(obj interface{}) IResponse

	// //xml 输出
	// Xml(obj interface{}) IResponse

	// html 输出
	Html(template string, obj interface{}) IResponse

	// // string
	// Text(format string, values ...interface{}) IResponse

	// // 重定向
	// Redirect(path string) IResponse

	// header
	SetHeader(key string, val string) IResponse

	// // Cookie
	// SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// // 设置状态码
	// SetStatus(code int) IResponse

	// // 设置 200 状态
	// SetOkStatus() IResponse
}

// Jsonp 输出
func (ctx *Context) Jsonp(obj interface{}) IResponse {
	// 获取请求 callback
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下字符过滤，否则有可能造成XSS攻击
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		return ctx
	}
	// 函数数据参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.responseWriter.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

// html 输出
func (ctx *Context) Html(file string, obj interface{}) IResponse {
	// 读取模版文件，创建 template 实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 执行 Execute 方法将 obj 和模板进行结合
	if err := t.Execute(ctx.responseWriter, obj); err != nil {
		return ctx
	}

	ctx.SetHeader("Context-Type", "application/html")
	return ctx
}

// header
func (ctx *Context) SetHeader(key string, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}
