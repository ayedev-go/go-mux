package lib

import "fmt"

/*
StringMap Object
*/
type StringMap map[string]string

/*
ObjectMap Object
*/
type ObjectMap map[string]interface{}

/*
Set Method
*/
func (d *ObjectMap) Set(key string, val interface{}) {
	(*d)[key] = val
}

/*
Add Method
*/
func (d *ObjectMap) Add(key string, val interface{}) {
	d.Set(key, val)
}

/*
Get Method
*/
func (d *ObjectMap) Get(key string, def ...interface{}) interface{} {
	val, found := (*d)[key]
	if found {
		return val
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

/*
GetString Method
*/
func (d *ObjectMap) GetString(key string) string {
	val := d.Get(key)
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}

/*
Has Method
*/
func (d *ObjectMap) Has(key string) bool {
	_, found := (*d)[key]
	return found
}

/*
Clear Method
*/
func (d *ObjectMap) Clear() {
	newMap := &ObjectMap{}
	*d = *newMap
}

/*
Copy Method
*/
func (d *ObjectMap) Copy(d2 ObjectMap) {
	for key, val := range d2 {
		d.Set(key, val)
	}
}

/*
StringMap Method
*/
func (d *ObjectMap) StringMap() StringMap {
	data := StringMap{}
	for key, val := range *d {
		data[key] = fmt.Sprintf("%v", val)
	}
	return data
}
