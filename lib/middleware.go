package lib

/*
Middleware Object
*/
type Middleware func(*Context, PipelineCallback)

/*
Middlewares Object
*/
type Middlewares []Middleware

/*
Add Method
*/
func (mw *Middlewares) Add(m Middleware) {
	*mw = append(*mw, m)
}

/*
Copy Method
*/
func (mw *Middlewares) Copy(mw2 Middlewares) {
	for _, m := range mw2 {
		mw.Add(m)
	}
}

/*
Size Method
*/
func (mw *Middlewares) Size() int {
	return len(*mw)
}

/*
MiddlewareMixin Object
*/
type MiddlewareMixin struct {
	middlewares *Middlewares
}

func (mwm *MiddlewareMixin) initMiddlewares() {
	mwm.middlewares = &Middlewares{}
}

/*
AddMiddleware Method
*/
func (mwm *MiddlewareMixin) AddMiddleware(mws ...Middleware) *MiddlewareMixin {
	mwm.middlewares.Copy(mws)
	return mwm
}

/*
Middlewares Method
*/
func (mwm *MiddlewareMixin) Middlewares() *Middlewares {
	return mwm.middlewares
}

/*
CopyMiddlewares Method
*/
func (mwm *MiddlewareMixin) CopyMiddlewares(mwm2 MiddlewareMixin) {
	mwm.middlewares.Copy(*mwm2.Middlewares())
}
