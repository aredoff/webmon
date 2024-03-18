package httptracer

import (
	"fmt"
	"testing"
	"time"
)

func TestTracer(t *testing.T) {
	trace := New()
	trace.SetTimeout(5 * time.Second)
	result := trace.Trace("https://google.com", "GET")
	fmt.Println("Error: ", result.Error)
	fmt.Println("Name Lookup: ", result.NameLookup)
	fmt.Println("Connect: ", result.Connect)
	fmt.Println("TLS Handshake: ", result.TLSHandshake)
	fmt.Println("First Byte: ", result.FirstByte)
	fmt.Println("Full Response: ", result.FullResponse)
	fmt.Println("Body Size (byte): ", result.BodySize)

	jsonData, _ := result.ToJSON()
	fmt.Println(string(jsonData))
}
