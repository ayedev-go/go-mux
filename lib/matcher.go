package lib

import (
	"fmt"
	"reflect"
)

/*
MatchHandler Object
*/
type MatchHandler func(*Context) bool

/*
Matcher Object
*/
type Matcher interface {
	Match(*Context) bool
}

/*
Matchers Object
*/
type Matchers []Matcher

/*
Add Method
*/
func (ms *Matchers) Add(m Matcher) {
	*ms = append(*ms, m)
}

/*
Copy Method
*/
func (ms *Matchers) Copy(ms2 Matchers) {
	for _, m := range ms2 {
		ms.Add(m)
	}
}

/*
Size Method
*/
func (ms *Matchers) Size() int {
	return len(*ms)
}

/*
Matches Method
*/
func (ms *Matchers) Matches(ctx *Context) bool {
	matches := true
	for _, matcher := range *ms {
		if !matcher.Match(ctx) {
			matches = false
			break
		}
	}
	return matches
}

/*
MatcherMixin Object
*/
type MatcherMixin struct {
	matchers *Matchers
}

func (mm *MatcherMixin) initMatchers() {
	mm.matchers = &Matchers{}
}

/*
AddMatcher Method
*/
func (mm *MatcherMixin) AddMatcher(ms ...interface{}) *MatcherMixin {
	for _, m := range ms {
		switch tp := m.(type) {
		case Matcher:
			mm.matchers.Add(tp)
		case func(*Context) bool:
			mm.matchers.Add(MakeMatcher(tp))
		default:
			panic(fmt.Sprintf("Unsupported matcher type: %s\n", reflect.TypeOf(m)))
		}
	}
	return mm
}

/*
Matchers Method
*/
func (mm *MatcherMixin) Matchers() *Matchers {
	return mm.matchers
}

/*
Matches Method
*/
func (mm *MatcherMixin) Matches(ctx *Context) bool {
	matches := true
	for _, matcher := range *mm.matchers {
		if !matcher.Match(ctx) {
			matches = false
			break
		}
	}
	return matches
}

/*
CopyMatchers Method
*/
func (mm *MatcherMixin) CopyMatchers(mm2 MatcherMixin) {
	mm.matchers.Copy(*mm2.Matchers())
}

/*
MatcherWrapper Object
*/
type MatcherWrapper struct {
	matcher MatchHandler
}

/*
Match Method
*/
func (nm *MatcherWrapper) Match(ctx *Context) bool {
	return nm.matcher(ctx)
}

/*
MakeMatcher Function
*/
func MakeMatcher(handler MatchHandler) Matcher {
	return &MatcherWrapper{matcher: handler}
}
