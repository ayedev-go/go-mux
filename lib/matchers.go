package lib

import (
	"regexp"
	"strings"
)

/*
MatchCaseEmpty Constant
*/
const MatchCaseEmpty = ""

/*
MatcherItem Object
*/
type MatcherItem struct {
	value   string
	compare string
	reg     *regexp.Regexp
}

/*
Valid Method
*/
func (mi *MatcherItem) Valid(val string) bool {
	matches := false
	if val != "" {
		if len(mi.value) > 0 {
			switch mi.compare {
			case "equal":
				matches = val == mi.value
			case "regexp":
				if mi.reg == nil {
					reg, regErr := regexp.Compile(mi.value)
					if regErr == nil {
						mi.reg = reg
					}
				}
				if mi.reg != nil {
					matches = mi.reg.MatchString(val)
				}
			case "list":
				list := strings.Split(mi.value, "|")
				if list != nil && len(list) > 0 {
					matches = InStringSlice(list, val)
				}
			default:
				matches = false
			}
		} else {
			matches = true
		}
	}
	return matches
}

/*
SetRegexp Method
*/
func (mi *MatcherItem) SetRegexp(reg *regexp.Regexp) {
	mi.reg = reg
	mi.SetCompare("regexp")
}

/*
SetRegexpCompile Method
*/
func (mi *MatcherItem) SetRegexpCompile(str string) {
	reg, regErr := regexp.Compile(str)
	if regErr == nil {
		mi.SetRegexp(reg)
	}
}

/*
GetRegexp Method
*/
func (mi *MatcherItem) GetRegexp() *regexp.Regexp {
	return mi.reg
}

/*
SetCompare Method
*/
func (mi *MatcherItem) SetCompare(compare string) {
	mi.compare = compare
}

/*
HeaderMatcher Object
*/
type HeaderMatcher struct {
	MatcherItem
	header string
}

/*
Match Method
*/
func (hm *HeaderMatcher) Match(ctx *Context) bool {
	return hm.Valid(ctx.Request.Header.Get(hm.header))
}

/*
HeaderMatch Function
*/
func HeaderMatch(key string, val string) *HeaderMatcher {
	return &HeaderMatcher{MatcherItem: MatcherItem{value: val, compare: "equal"}, header: key}
}

/*
HeaderMatchRegexp Function
*/
func HeaderMatchRegexp(key string, val string) *HeaderMatcher {
	matcher := HeaderMatch(key, val)
	matcher.MatcherItem.SetRegexpCompile(val)
	return matcher
}

/*
SchemaMatcher Object
*/
type SchemaMatcher struct {
	MatcherItem
}

/*
Match Method
*/
func (sm *SchemaMatcher) Match(ctx *Context) bool {
	return sm.Valid(ctx.Schema)
}

/*
SchemaMatch Function
*/
func SchemaMatch(schema string) *SchemaMatcher {
	return &SchemaMatcher{MatcherItem{value: schema, compare: "equal"}}
}

/*
SchemaMatchRegexp Function
*/
func SchemaMatchRegexp(schema string) *SchemaMatcher {
	matcher := SchemaMatch(schema)
	matcher.MatcherItem.SetRegexpCompile(schema)
	return matcher
}

/*
HostMatcher Object
*/
type HostMatcher struct {
	MatcherItem
}

/*
Match Method
*/
func (hm *HostMatcher) Match(ctx *Context) bool {
	return hm.Valid(ctx.Host)
}

/*
HostMatch Function
*/
func HostMatch(host string) *HostMatcher {
	return &HostMatcher{MatcherItem{value: host, compare: "equal"}}
}

/*
HostMatchRegexp Function
*/
func HostMatchRegexp(host string) *HostMatcher {
	matcher := HostMatch(host)
	matcher.MatcherItem.SetRegexpCompile(host)
	return matcher
}

/*
PortMatcher Object
*/
type PortMatcher struct {
	MatcherItem
}

/*
Match Method
*/
func (hm *PortMatcher) Match(ctx *Context) bool {
	return hm.Valid(ctx.Port)
}

/*
PortMatch Function
*/
func PortMatch(port string) *PortMatcher {
	return &PortMatcher{MatcherItem{value: port, compare: "equal"}}
}

/*
PortMatchRegexp Function
*/
func PortMatchRegexp(port string) *PortMatcher {
	matcher := PortMatch(port)
	matcher.MatcherItem.SetRegexpCompile(port)
	return matcher
}

/*
QueryMatcher Object
*/
type QueryMatcher struct {
	MatcherItem
	key string
}

/*
Match Method
*/
func (qm *QueryMatcher) Match(ctx *Context) bool {
	return qm.Valid(ctx.QueryValue(qm.key))
}

/*
QueryMatch Function
*/
func QueryMatch(key string, val string) *QueryMatcher {
	return &QueryMatcher{MatcherItem: MatcherItem{value: val, compare: "equal"}, key: key}
}

/*
QueryMatchRegexp Function
*/
func QueryMatchRegexp(key string, val string) *QueryMatcher {
	matcher := QueryMatch(key, val)
	matcher.MatcherItem.SetRegexpCompile(val)
	return matcher
}

/*
MethodMatcher Object
*/
type MethodMatcher struct {
	MatcherItem
}

/*
Match Method
*/
func (mm *MethodMatcher) Match(ctx *Context) bool {
	return mm.Valid(ctx.Method)
}

/*
Add Method
*/
func (mm *MethodMatcher) Add(method string) {
	methods := mm.Methods()
	methods = append(methods, method)
	mm.MatcherItem.value = strings.Join(methods, "|")
}

/*
Set Method
*/
func (mm *MethodMatcher) Set(methods ...string) {
	mm.MatcherItem.value = strings.Join(methods, "|")
}

/*
Set Method
*/
func (mm *MethodMatcher) Methods() []string {
	return strings.Split(mm.MatcherItem.value, "|")
}

/*
MethodMatch Function
*/
func MethodMatch(methods ...string) *MethodMatcher {
	methodsStr := strings.Join(methods, "|")
	return &MethodMatcher{MatcherItem{value: methodsStr, compare: "list"}}
}

/*
PathMatcher Object
*/
type PathMatcher struct {
	MatcherItem
	formats []string
}

/*
Match Method
*/
func (pm *PathMatcher) Match(ctx *Context) bool {
	return pm.Valid(ctx.Path)
}

/*
SetRegexp Method
*/
func (pm *PathMatcher) SetRegexp(reg *regexp.Regexp) {
	pm.MatcherItem.reg = reg
	pm.MatcherItem.compare = "regexp"
}

/*
SetCompare Method
*/
func (pm *PathMatcher) SetCompare() {
	pm.MatcherItem.compare = "equal"
}

/*
GetFormats Method
*/
func (pm *PathMatcher) GetFormats() []string {
	return pm.formats
}

/*
PathMatch Function
*/
func PathMatch(path string, formats ...string) *PathMatcher {
	if path == "" || path == "/" {
		formats = []string{}
	}
	if len(path) > 0 && path[0:1] != "/" {
		path = "/" + path
	}
	return &PathMatcher{MatcherItem: MatcherItem{value: path}, formats: formats}
}
