package isuperagent

import (
	"context"
)

type Context interface {
	Set(key, value interface{}) Context
	Get(key interface{}) interface{}

	SetReq(Request) Context
	GetReq() Request
	SetRes(Response) Context
	GetRes() Response
}

type icontext struct {
	context.Context

	Req Request
	Res Response
}

func NewContext(ctx context.Context, req Request, res Response) Context {
	return &icontext{
		Context: ctx,
		Req:     req,
		Res:     res,
	}
}

func (ctx *icontext) Set(key, value interface{}) Context {
	ctx.Context = context.WithValue(ctx.Context, key, value)

	return ctx
}

func (ctx *icontext) Get(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

func (ctx *icontext) SetReq(req Request) Context {
	ctx.Req = req

	return ctx
}

func (ctx *icontext) GetReq() Request {
	return ctx.Req
}

func (ctx *icontext) SetRes(res Response) Context {
	ctx.Res = res

	return ctx
}

func (ctx *icontext) GetRes() Response {
	return ctx.Res
}
