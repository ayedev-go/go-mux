package lib

/*
Handler Interface
*/
type Handler interface {
	Handle(*Context)
	Match(*Context) bool
}

/*
Handlers Object
*/
type Handlers []Handler

/*
Add Method
*/
func (hs *Handlers) Add(h Handler) {
	*hs = append(*hs, h)
}

/*
Copy Method
*/
func (hs *Handlers) Copy(hs2 Handlers) {
	for _, h := range hs2 {
		hs.Add(h)
	}
}

/*
HandlerMixin Object
*/
type HandlerMixin struct {
	handlers *Handlers
}

func (hm *HandlerMixin) initHandlers() {
	hm.handlers = &Handlers{}
}

/*
AddHandler Method
*/
func (hm *HandlerMixin) AddHandler(h Handler) *HandlerMixin {
	hm.handlers.Add(h)
	return hm
}

/*
Handlers Method
*/
func (hm *HandlerMixin) Handlers() *Handlers {
	return hm.handlers
}

/*
CopyHandlers Method
*/
func (hm *HandlerMixin) CopyHandlers(hm2 HandlerMixin) {
	hm.handlers.Copy(*hm2.Handlers())
}
