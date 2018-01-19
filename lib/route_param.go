package lib

import "fmt"

/*
RouteParam Object
*/
type RouteParam struct {
	key   string
	dummy string
	value string
}

/*
Set Function
*/
func (param *RouteParam) Set(value string) {
	param.value = value
}

/*
Value Function
*/
func (param *RouteParam) Value() string {
	return param.value
}

/*
RouteParams Object
*/
type RouteParams []*RouteParam

/*
Add Function
*/
func (params *RouteParams) Add(key string, dummy string) {
	params.AddParam(&RouteParam{key: key, dummy: dummy})
}

/*
AddParam Function
*/
func (params *RouteParams) AddParam(param *RouteParam) {
	*params = append(*params, param)
}

/*
Param Function
*/
func (params *RouteParams) Param(key string) (*RouteParam, error) {
	for _, param := range *params {
		if param.key == key {
			return param, nil
		}
	}
	return nil, fmt.Errorf("Param not found: %s", key)
}

/*
FindParam Function
*/
func (params *RouteParams) FindParam(key string) *RouteParam {
	param, _ := params.Param(key)
	return param
}
