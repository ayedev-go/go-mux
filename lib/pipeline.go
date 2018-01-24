package lib

import (
	"net/http"
)

/*
PipelineCallback Function
*/
type PipelineCallback func()

/*
PipelineErrorCallback Function
*/
type PipelineErrorCallback func(string, PipelineCallback, PipelineCallback)

/*
Pipeline Object
*/
type Pipeline struct {
	done     bool
	position int
	cbBefore PipelineCallback
	cbAfter  PipelineCallback
	cbError  PipelineErrorCallback
	ctx      *Context
	queue    []Middleware
}

/*
Add Method
*/
func (pl *Pipeline) Add(handlers ...Middleware) {
	for _, hn := range handlers {
		pl.queue = append(pl.queue, hn)
	}
}

/*
Size Method
*/
func (pl *Pipeline) Size() int {
	return len(pl.queue)
}

/*
Reset Method
*/
func (pl *Pipeline) Reset() {
	pl.ClearCallbacks()
	pl.queue = []Middleware{}
}

/*
ClearCallbacks Method
*/
func (pl *Pipeline) ClearCallbacks() {
	pl.cbAfter = nil
	pl.cbBefore = nil
	pl.cbError = nil
}

/*
Copy Method
*/
func (pl *Pipeline) Copy(source interface{}) {
	switch tp := source.(type) {
	case *Middlewares:
		for _, mw := range *tp {
			pl.Add(mw)
		}
	case *Route:
		for _, mw := range *tp.Middlewares() {
			pl.Add(mw)
		}
	case *Router:
		for _, mw := range *tp.Middlewares() {
			pl.Add(mw)
		}
	case *MiddlewareMixin:
		for _, mw := range *tp.Middlewares() {
			pl.Add(mw)
		}
	}
}

/*
Start Method
*/
func (pl *Pipeline) Start(cb ...PipelineCallback) {
	pl.done = false
	pl.position = -1
	if len(cb) == 1 {
		pl.cbAfter = cb[0]
	} else if len(cb) == 2 {
		pl.cbBefore = cb[0]
		pl.cbAfter = cb[1]
	} else if len(cb) == 3 {
		pl.cbBefore = cb[0]
		pl.cbAfter = cb[1]
	}
	if pl.cbBefore != nil {
		pl.cbBefore()
	}
	pl.next()
}

/*
OnError Method
*/
func (pl *Pipeline) OnError(cb PipelineErrorCallback) {
	pl.cbError = cb
}

func (pl *Pipeline) next() {
	newPosition := pl.position + 1
	if (len(pl.queue) - 1) >= newPosition {
		pl.position = newPosition
		defer func() {
			if r := recover(); r != nil {
				if pl.cbError != nil {
					pl.cbError(r.(string), pl.next, pl.stop)
				} else {
					pl.ctx.Error(http.StatusInternalServerError, r.(string))
					pl.stopOnError(r.(string))
				}
			}
		}()
		pl.queue[newPosition](pl.ctx, pl.next)
	} else {
		pl.done = true
		if pl.cbAfter != nil {
			pl.cbAfter()
		}
	}
}

func (pl *Pipeline) stopOnError(err string) {
	pl.stop()
}

func (pl *Pipeline) stop() {
	if pl.cbAfter != nil {
		pl.cbAfter()
	}
}

/*
NewPipeline Function
*/
func NewPipeline(ctx *Context) *Pipeline {
	pipeline := &Pipeline{ctx: ctx}
	pipeline.Reset()
	return pipeline
}
