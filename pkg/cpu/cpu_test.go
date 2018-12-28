package cpu

import (
	"encoding/binary"
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestNewCPU(t *testing.T) {
	m := make(state.Memory, 10)
	m[5] = 8
	c := NewCPU(m)
	assert.NotNil(t, c)
	assert.Equal(t, m, c.m)
}

func TestCpu_Tick_0xANNN(t *testing.T) {
	bs := opCodeToBytes(0xa2F0)
	m := state.InitMemory()
	m[512] = bs[0]
	m[513] = bs[1]
	c := NewCPU(m)
	err := c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, c.pc, int16(514))
	assert.Equal(t, c.ir, uint16(0x2F0))
}

func TestCpu_Tick_0x8XY4(t *testing.T) {
	bs := opCodeToBytes(0x8004)
	m := state.InitMemory()
	m[512] = bs[0]
	m[513] = bs[1]
	c := NewCPU(m)
	err := c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, c.pc, int16(514))
	assert.Equal(t, c.ir, uint16(0x2F0))
}

func opCodeToBytes(opCode uint16) (result []byte) {
	result = make([]byte, 2)
	binary.BigEndian.PutUint16(result, opCode)
	return result
}

/*
	fmt.Printf("%v\n", 0xF000)
	i := 0xa2f0 & 0x0FFF
	fmt.Printf("%#04x:%X:%v\n", i, i, i)
	fmt.Printf("%#04x:%X:%v\n", 0xA000, 0xA000, 0xF000)
*/
// 0x0FFF == 4095
// 0xF000 == 61440
