package lib

/*
Patterns Object
*/
type Patterns StringMap

/*
Set Method
*/
func (p *Patterns) Set(key string, val string) {
	(*p)[key] = val
}

/*
Get Method
*/
func (p *Patterns) Get(key string, fallback ...string) string {
	val, found := (*p)[key]
	if found {
		return val
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

/*
Has Method
*/
func (p *Patterns) Has(key string) bool {
	_, found := (*p)[key]
	return found
}

/*
Copy Method
*/
func (p *Patterns) Copy(p2 Patterns) {
	for k, v := range p2 {
		p.Set(k, v)
	}
}

/*
PatternMixin Object
*/
type PatternMixin struct {
	patterns *Patterns
}

func (pm *PatternMixin) initPatterns() {
	pm.patterns = &Patterns{}
}

/*
SetPattern Method
*/
func (pm *PatternMixin) SetPattern(name string, pattern string) *PatternMixin {
	pm.patterns.Set(name, pattern)
	return pm
}

/*
GetPattern Method
*/
func (pm *PatternMixin) GetPattern(name string, fallback ...string) string {
	return pm.patterns.Get(name, fallback...)
}

/*
HasPattern Method
*/
func (pm *PatternMixin) HasPattern(name string) bool {
	return pm.patterns.Has(name)
}

/*
Patterns Method
*/
func (pm *PatternMixin) Patterns() *Patterns {
	return pm.patterns
}

/*
ClearPatterns Method
*/
func (pm *PatternMixin) ClearPatterns() *PatternMixin {
	pm.patterns = &Patterns{}
	return pm
}

/*
CopyPatterns Method
*/
func (pm *PatternMixin) CopyPatterns(pm2 PatternMixin) {
	pm.patterns.Copy(*pm2.Patterns())
}
