package webhandle

import (
	"bytes"
	"fmt"
	"net/http"
)

var globalStringCache map[string]string

type FakeResponseWriter struct {
	Buf bytes.Buffer // initialized by default
}

// Create a new fake ResponseWriter.
// Useful for testing the output from http.HandlerFunc functions.
func NewFakeResponseWriter() *FakeResponseWriter {
	return &FakeResponseWriter{}
}

func (f *FakeResponseWriter) Header() http.Header {
	return http.Header{}
}

func (f *FakeResponseWriter) Write(b []byte) (int, error) {
	return f.Buf.Write(b)
}

func (f *FakeResponseWriter) WriteHeader(i int) {
}

func (f *FakeResponseWriter) String() string {
	return f.Buf.String()
}

// Wrap a http.HandlerFunc so that the output is cached (with an id)
// Do not cache functions with side-effects! (that sets the mimetype for instance)
func NewCacheWrapper(id string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Initialize the map if it isn't already initialized
		if globalStringCache == nil {
			globalStringCache = make(map[string]string)
		}
		if _, ok := globalStringCache[id]; !ok {
			fake := &FakeResponseWriter{}
			// Make the handler write to the buffer
			handler(fake, req)
			globalStringCache[id] = fake.String()
		}
		fmt.Fprint(w, globalStringCache[id])
	}
}

// Fetch the entire cache map
func GetCache() map[string]string {
	return globalStringCache
}
