package canceler_test

import (
	"testing"

	"github.com/mway/pkg/x/sync/canceler"
)

func BenchmarkCancelerC(b *testing.B) {
	c := canceler.New(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.C()
	}
}
