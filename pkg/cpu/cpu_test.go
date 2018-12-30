package cpu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestNewCPU(t *testing.T) {
	m := make(state.Memory, 512+10)
	bf := bytes.NewBuffer([]byte{8})
	fmt.Println(bf.Len())
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	assert.NotNil(t, c)
	assert.Equal(t, m, c.m)
}

func TestCpu_Tick_0xANNN(t *testing.T) {
	bs := opCodeToBytes(0xa2F0)
	m := state.InitMemory()
	bf := bytes.NewBuffer(bs)
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	err = c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, int16(514), c.pc)
	assert.Equal(t, c.ir, uint16(0x2F0))
}

func TestCpu_Tick_0x1NNN(t *testing.T) {
	bs := opCodeToBytes(0x14ef)
	m := state.InitMemory()
	bf := bytes.NewBuffer(bs)
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	err = c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, int16(1263), c.pc)
}

func TestCpu_Tick_0x2NNN(t *testing.T) {
	bs := opCodeToBytes(0x24ef)
	m := state.InitMemory()
	bf := bytes.NewBuffer(bs)
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	err = c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, int16(1263), c.pc)
	assert.Equal(t, int8(1), c.stack.Len())
	assert.Equal(t, int16(512), c.stack.Pop())
}

func TestCpu_Tick_0x6XNN(t *testing.T) {
	bs := opCodeToBytes(0x64ee)
	m := state.InitMemory()
	bf := bytes.NewBuffer(bs)
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	err = c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, int16(514), c.pc)
	assert.Equal(t, uint8(238), c.v[4])
}

func TestCpu_Tick_0x7XNN(t *testing.T) {
	bs := opCodeToBytes(0x741f)
	m := state.InitMemory()
	bf := bytes.NewBuffer(bs)
	err := m.LoadMemory(bf)
	assert.NoError(t, err)
	c := NewCPU(m)
	c.v[4] = uint8(0x0b) // 11
	err = c.Tick()
	assert.NoError(t, err)
	assert.Equal(t, int16(514), c.pc)
	assert.Equal(t, uint8(0x2a), c.v[4])
}

func TestCpu_Tick_0x8(t *testing.T) {
	var testCases = []struct {
		name   string
		opcode uint16
		x      int
		vx     uint8
		y      int
		vy     uint8
		exp    uint8
		cry    uint8
	}{
		{
			name:   "0x8XY0",
			opcode: 0x80e0,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    8, // expect vx to become vy
			cry:    0x0,
		},
		{
			name:   "0x8XY1",
			opcode: 0x80e1,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    13, // Vx=Vx|Vy
			cry:    0x0,
		},
		{
			name:   "0x8XY2",
			opcode: 0x80e2,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    1, // Vx=Vx&Vy
			cry:    0x0,
		},
		{
			name:   "0x8XY3",
			opcode: 0x80e3,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    12, // Vx=Vx^Vy
			cry:    0x0,
		},
		{
			name:   "0x8XY4 no carry",
			opcode: 0x80e4,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    20,
			cry:    0x0,
		},
		{
			name:   "0x8XY4 carry the one",
			opcode: 0x80e4,
			x:      0,
			vx:     200,
			y:      14,
			vy:     60,
			exp:    4, // rolls over, 200 + 60 = 260, 5 over so 0,1,2,3,*4*
			cry:    0x1,
		},
		{
			name:   "0x8XY5 no borrow",
			opcode: 0x80e5,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    4,
			cry:    0x1,
		},
		{
			name:   "0x8XY5 borrow",
			opcode: 0x80e5,
			x:      0,
			vx:     8,
			y:      14,
			vy:     249,
			exp:    15,
			cry:    0x0,
		},
		{
			name:   "0x8XY6",
			opcode: 0x80e6,
			x:      0,
			vx:     33,
			y:      14,
			vy:     8,
			exp:    16,
			cry:    0x1,
		},
		{
			name:   "0x8XY7 no borrow",
			opcode: 0x80e7,
			x:      0,
			vx:     8,
			y:      14,
			vy:     12,
			exp:    4,
			cry:    0x1,
		},
		{
			name:   "0x8XY7 borrow",
			opcode: 0x80e7,
			x:      0,
			vx:     249,
			y:      14,
			vy:     8,
			exp:    15,
			cry:    0x0,
		},
		{
			name:   "0x8XYE",
			opcode: 0x80ee,
			x:      0,
			vx:     251,
			y:      14,
			vy:     8,
			exp:    246,
			cry:    0x1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs := opCodeToBytes(tc.opcode)
			m := state.InitMemory()
			bf := bytes.NewBuffer(bs)
			err := m.LoadMemory(bf)
			assert.NoError(t, err)
			c := NewCPU(m)
			// add a value for Y
			c.v[tc.y] = tc.vy
			// add a value for X
			c.v[tc.x] = tc.vx
			err = c.Tick()
			assert.NoError(t, err)
			assert.Equal(t, int16(514), c.pc, "should have moved program counter on two")
			assert.Equal(t, uint16(0x0), c.ir, "No index register to change")
			assert.Equal(t, tc.cry, c.v[15], "No one to carry over or borrow")
			assert.Equal(t, tc.exp, c.v[0], "X not equal expected")
		})
	}
}

func opCodeToBytes(opcode uint16) (result []byte) {
	result = make([]byte, 2)
	binary.BigEndian.PutUint16(result, opcode)
	return result
}
