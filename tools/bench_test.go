package tools

import (
	"log/slog"
	"testing"
)

func BenchmarkBuildCRUDTools(b *testing.B) {
	log := slog.New(slog.DiscardHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct := NewCRUDTools(nil, log)
		_, _ = ct.GetListResourcesTool()
		_, _ = ct.GetGetResourceTool()
		_, _ = ct.GetCreateResourceTool()
		_, _ = ct.GetUpdateResourceTool()
		_, _ = ct.GetDeleteResourceTool()
	}
}

func BenchmarkNewRegistry(b *testing.B) {
	log := slog.New(slog.DiscardHandler)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRegistry(nil, log)
	}
}
