package test

import (
	"encoding/json"
	"testing"
)

func BenchmarkMirroredQuery(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < b.N; i++ {
		query := mirroredQuery()
		m[query]++
	}
	marshalBytes, _ := json.Marshal(m)
	b.Log(string(marshalBytes))
}
