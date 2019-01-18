package cpu

import (
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	bcChip8TestPath = "../../test/roms/BC_test.ch8"
)

func BenchmarkStack(b *testing.B) {
	m := state.InitMemory()
	log.SetLevel(log.WarnLevel)
	addToMemory(m, 0x8004, 512)
	addToMemory(m, 0x80e2, 514)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := getCPU(b, m)
		// add a value for Y
		c.v[14] = uint8(8)
		// add a value for X
		c.v[0] = uint8(12)
		err := c.Tick()
		assert.NoError(b, err)
		err = c.Tick()
		assert.NoError(b, err)
	}
}

func Benchmark_BC_Chip8Test(b *testing.B) {
	log.SetLevel(log.WarnLevel)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err := getMemory(b)
		c := getCPU(b, m)
		for tc := 0; tc < 250; tc++ { // only need 250 cycles to process all of BC_test.ch8
			err = c.Tick()
			assert.NoError(b, err)
		}
	}
}

func getMemory(b *testing.B) (state.Memory, error) {
	b.StopTimer()
	defer b.StartTimer()
	m := state.InitMemory()
	f, err := os.Open(bcChip8TestPath)
	assert.NoError(b, err)
	err = m.LoadMemory(f)
	assert.NoError(b, err)
	return m, err
}

func addToMemory(m state.Memory, opcode uint16, loc int) {
	bs := opCodeToBytes(opcode)
	m[loc] = bs[0]
	m[loc+1] = bs[1]
}

func getCPU(b *testing.B, m state.Memory) *cpu {
	b.StopTimer()
	defer b.StartTimer()
	ti, sc := setupTimer()
	go func(s <-chan byte) {
		for {
			select {
			case <-s:
				// noop
			}
		}
	}(sc)
	return getNewCPU(m, NewKeyboard(), ti, &noopScreen{})
}

type noopScreen struct {
}

func (s *noopScreen) Draw(frameBuffer []byte) {

}
