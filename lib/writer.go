package lib

import (
	"bytes"
	"net/http"
)

/*
Writer Object
*/
type Writer struct {
	buffer     bytes.Buffer
	StatusCode int
	Headers    *ObjectMap
}

func (w Writer) init() {
	w.Flush()
}

/*
Write Method
*/
func (w *Writer) Write(bytes []byte) {
	w.buffer.Write(bytes)
}

/*
WriteString Method
*/
func (w *Writer) WriteString(str string) {
	w.buffer.WriteString(str)
}

/*
Status Method
*/
func (w *Writer) Status(code int) {
	w.StatusCode = code
}

/*
GetStatus Method
*/
func (w *Writer) GetStatus() int {
	return w.StatusCode
}

/*
Header Method
*/
func (w *Writer) Header(key string, val string) {
	w.Headers.Set(key, val)
}

/*
GetHeader Method
*/
func (w *Writer) GetHeader(key string) string {
	return w.Headers.GetString(key)
}

/*
ClearBuffer Method
*/
func (w *Writer) ClearBuffer() {
	w.buffer.Reset()
}

/*
Flush Method
*/
func (w *Writer) Flush() {
	w.ClearBuffer()
	w.Headers.Clear()
}

/*
PushTo Method
*/
func (w *Writer) PushTo(iow http.ResponseWriter) {
	if w.StatusCode > 0 {
		iow.WriteHeader(w.StatusCode)
	} else {
		iow.WriteHeader(http.StatusOK)
	}
	w.buffer.WriteTo(iow)
	headers := w.Headers.StringMap()
	for hKey, hVal := range headers {
		iow.Header().Set(hKey, hVal)
	}
}

/*
Bytes Method
*/
func (w *Writer) Bytes() []byte {
	return w.buffer.Bytes()
}

/*
NewWriter Function
*/
func NewWriter() *Writer {
	writer := &Writer{Headers: &ObjectMap{}}
	writer.init()
	return writer
}
