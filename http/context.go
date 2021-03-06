package jhttp

import (
	"bytes"
	"context"
	"github.com/gorilla/mux"
	ogjson "github.com/og/json"
	"github.com/pkg/errors"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
	router *Router
	resolvedParam bool
	param map[string]string
}
func NewContext (w http.ResponseWriter, r *http.Request, router *Router) *Context {
	return &Context{
		W: w,
		R: r,
		router: router,
	}
}
func (c *Context) Context() context.Context {
	return c.R.Context()
}
func (c *Context) Param(name string) (param string, err error) {
	data := map[string]string{}
	if c.resolvedParam {
		data = c.param
	} else {
		data = mux.Vars(c.R)
	}
	param, has := data[name]
	if !has {
		return "", errors.New(`not found param (` + name + `)`)
	}
	return param, nil
}
func (c *Context) Bytes(b []byte) error {
	_, err := c.W.Write(b)
	if err != nil {
		_, err := c.W.Write([]byte("c.Bytes(b) error"))
		if err != nil {panic(err)}
	}
	return nil
}
func (c *Context) Render(render func(buffer *bytes.Buffer)) error {
	buffer := bytes.NewBuffer(nil)
	render(buffer)
	c.W.Header().Set("Content-Type", "text/html; charset=UTF-8")
	return c.Bytes(buffer.Bytes())
}
func (c *Context) JSON(v interface{}) error {
	jsonb, err := ogjson.BytesWithErr(v)
	if err != nil {
		return err
	}
	return c.Bytes(jsonb)
}
func (c *Context) BindRequest(ptr interface{}) error {
	return BindRequest(ptr, c.R)
}

func (c *Context) CheckError(errInterface interface{}) {
	err := c.router.OnCatchError(c, errInterface)
	if err != nil {
		panic(err)
	}
}