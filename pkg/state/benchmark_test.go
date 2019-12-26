package state

import (
	"crypto/sha256"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkLoadMemory(b *testing.B) {
	log.SetLevel(log.WarnLevel)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := loadMemoryTest(b)
		b.StopTimer()
		hash := sha256.Sum256(m)
		assert.Equal(b,
			C8PIC_sha,
			fmt.Sprintf("%x", hash))
		b.StartTimer()
	}
}
