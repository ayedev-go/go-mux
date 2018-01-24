package lib

/*
Param Object
*/
type Param struct {
	key         string
	placeholder string
	value       string
}

/*
Set Method
*/
func (param *Param) Set(value string) {
	param.value = value
}

/*
Value Method
*/
func (param *Param) Value() string {
	return param.value
}

/*
Params Object
*/
type Params []*Param

/*
Add Method
*/
func (params *Params) Add(key string, placeholder string) {
	params.AddParam(&Param{key: key, placeholder: placeholder})
}

/*
AddParam Method
*/
func (params *Params) AddParam(param *Param) {
	*params = append(*params, param)
}

/*
FindParam Method
*/
func (params *Params) Find(key string) (*Param, bool) {
	for _, param := range *params {
		if param.key == key {
			return param, true
		}
	}
	return nil, false
}

/*
Param Function
*/
func (params *Params) Param(key string) *Param {
	param, _ := params.Find(key)
	return param
}

/*
Set Method
*/
func (params *Params) Set(key string, val string) {
	param := params.Param(key)
	if param != nil {
		param.Set(val)
	}
}

/*
Get Method
*/
func (params *Params) Get(key string, fallback ...string) string {
	param := params.Param(key)
	if param != nil {
		return param.Value()
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

/*
Has Method
*/
func (params *Params) Has(key string) bool {
	_, found := params.Find(key)
	return found
}

/*
Copy Method
*/
func (params *Params) Copy(p2 Params) {
	for _, p := range p2 {
		params.AddParam(p)
	}
}

/*
Size Method
*/
func (params *Params) Size() int {
	return len(*params)
}

/*
StringMap Method
*/
func (params *Params) StringMap() StringMap {
	data := StringMap{}
	for _, param := range *params {
		data[param.key] = param.value
	}
	return data
}

/*
ParamMixin Object
*/
type ParamMixin struct {
	params *Params
}

func (pm *ParamMixin) initParams() {
	pm.params = &Params{}
}

/*
CopyParams Method
*/
func (pm *ParamMixin) CopyParams(params Params) {
	for _, param := range params {
		pm.params.AddParam(param)
	}
}

/*
Params Method
*/
func (pm *ParamMixin) Params() *Params {
	return pm.params
}

/*
HasParam Method
*/
func (pm *ParamMixin) HasParam(key string) bool {
	return pm.params.Has(key)
}

/*
Param Method
*/
func (pm *ParamMixin) Param(key string) Param {
	return *(pm.params.Param(key))
}

/*
StringMap Method
*/
func (pm *ParamMixin) StringMap() StringMap {
	return pm.params.StringMap()
}
