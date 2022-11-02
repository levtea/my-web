package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 自定义 Context
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	// 是否超时标记位
	hasTimeout bool
	// 写保护机制
	writeMux *sync.Mutex
}

// region base function
// 对外暴露锁
func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writeMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

// end region

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// region implement context.Context
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// end region

// region query url
func (ctx *Context) QueryInt(key string, def int) int {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

// end region

// region from post
func (ctx *Context) FromInt(key string, def int) int {
	params := ctx.FromAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			intval, err := strconv.Atoi(vals[len-1])
			if err != nil {
				return def
			}
			return intval
		}
	}
	return def
}

func (ctx *Context) FromString(key string, def string) string {
	params := ctx.FromAll()
	if vals, ok := params[key]; ok {
		len := len(vals)
		if len > 0 {
			return vals[len-1]
		}
	}
	return def
}

func (ctx *Context) FromArray(key string, def []string) []string {
	params := ctx.FromAll()
	if vals, ok := params[key]; ok {
		return vals
	}
	return def
}

func (ctx *Context) FromAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

// end region

// region application/json post
func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// end region

// region response
func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

// end region
