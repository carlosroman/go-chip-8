package cpu

import (
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

func BenchmarkStack(b *testing.B) {
	bs := opCodeToBytes(0x8004)
	m := state.InitMemory()
	m[512] = bs[0]
	m[513] = bs[1]
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := getCPU(b, m)
		err := c.Tick()
		assert.NoError(b, err)
	}
}

func getCPU(b *testing.B, m state.Memory) *cpu {
	b.StopTimer()
	defer b.StartTimer()
	return getNewCPU(m, newKeyboard(), NewTimer(), &noopScreen{})
}

type noopScreen struct {
}

func (s *noopScreen) Draw(frameBuffer []byte) {

}
