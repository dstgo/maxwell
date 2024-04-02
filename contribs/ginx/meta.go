package ginx

import "github.com/gin-gonic/gin"

// Elem is element of Metadata
type Elem struct {
	Key string
	Val any
}

func (e *Elem) Bool() bool {
	if b, ok := e.Val.(bool); ok {
		return b
	}
	return false
}

func (e *Elem) String() string {
	if b, ok := e.Val.(string); ok {
		return b
	}
	return ""
}

func (e *Elem) Int() int {
	if b, ok := e.Val.(int); ok {
		return b
	}
	return 0
}

func (e *Elem) Float() float64 {
	if b, ok := e.Val.(float64); ok {
		return b
	}
	return 0
}

func Meta(e ...Elem) MetaData {
	if len(e) == 0 {
		return nil
	}
	meta := make(MetaData)
	for _, E := range e {
		meta[E.Key] = E.Val
	}
	return meta
}

// MetaData
// route meta info
type MetaData map[string]any

func (m MetaData) Get(key string) (Elem, bool) {
	v, e := m[key]
	if !e {
		return Elem{}, false
	}
	return Elem{Key: key, Val: v}, true
}

func (m MetaData) Has(key string) bool {
	_, b := m.Get(key)
	return b
}

var MetaKey = "gin.metadata"

func MetaDataHandler(meta MetaData) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if meta == nil {
			meta = make(MetaData)
		}
		ctx.Set(MetaKey, meta)
	}
}
func MetaDataFromCtx(ctx *gin.Context) MetaData {
	value, e := ctx.Get(MetaKey)
	if !e {
		value = make(MetaData)
	}
	return value.(MetaData)
}
