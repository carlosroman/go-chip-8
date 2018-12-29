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
	assert.Equal(t, c.pc, int16(514))
	assert.Equal(t, c.ir, uint16(0x2F0))
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
		crry   uint8
	}{
		{
			name:   "0x8XY0",
			opcode: 0x80e0,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    8, // expect vx to become vy
			crry:   0x0,
		},
		{
			name:   "0x8XY1",
			opcode: 0x80e1,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    13, // Vx=Vx|Vy
			crry:   0x0,
		},
		{
			name:   "0x8XY2",
			opcode: 0x80e2,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    1, // Vx=Vx&Vy
			crry:   0x0,
		},
		{
			name:   "0x8XY3",
			opcode: 0x80e3,
			x:      0,
			vx:     5,
			y:      14,
			vy:     9,
			exp:    12, // Vx=Vx^Vy
			crry:   0x0,
		},
		{
			name:   "0x8XY4 no carry",
			opcode: 0x80e4,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    20,
			crry:   0x0,
		},
		{
			name:   "0x8XY4 carry the one",
			opcode: 0x80e4,
			x:      0,
			vx:     200,
			y:      14,
			vy:     60,
			exp:    4, // rolls over, 200 + 60 = 260, 5 over so 0,1,2,3,*4*
			crry:   0x1,
		},
		{
			name:   "0x8XY5 no borrow",
			opcode: 0x80e5,
			x:      0,
			vx:     12,
			y:      14,
			vy:     8,
			exp:    4,
			crry:   0x1,
		},
		{
			name:   "0x8XY5 borrow",
			opcode: 0x80e5,
			x:      0,
			vx:     8,
			y:      14,
			vy:     249,
			exp:    15,
			crry:   0x0,
		},
		{
			name:   "0x8XY6",
			opcode: 0x80e6,
			x:      0,
			vx:     33,
			y:      14,
			vy:     8,
			exp:    16,
			crry:   0x1,
		},
		{
			name:   "0x8XY7 no borrow",
			opcode: 0x80e7,
			x:      0,
			vx:     8,
			y:      14,
			vy:     12,
			exp:    4,
			crry:   0x1,
		},
		{
			name:   "0x8XY7 borrow",
			opcode: 0x80e7,
			x:      0,
			vx:     249,
			y:      14,
			vy:     8,
			exp:    15,
			crry:   0x0,
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
			assert.Equal(t, tc.crry, c.v[15], "No one to carry over or borrow")
			assert.Equal(t, tc.exp, c.v[0], "X not equal expected")
		})
	}
}

func opCodeToBytes(opcode uint16) (result []byte) {
	result = make([]byte, 2)
	binary.BigEndian.PutUint16(result, opcode)
	return result
}

//func Test0x8XY4(t *testing.T) {
//	opcode := uint16(0x8fa4)
//	fmt.Printf("%#04x:%X%v\n", opcode, opcode, opcode)
//	m := opcode & 0x00F0 >> 4
//	fmt.Printf("%#04x:%X:%v\n", m, m, m)
//	ms := m >> 4
//	fmt.Printf("%#04x:%X:%v\n", ms, ms, ms)
//	fmt.Printf("%#04x:%X:%v\n", (opcode&0x0F00)>>8, (opcode&0x0F00)>>8, (opcode&0x0F00)>>8)
//}
//
//func TestSomething(t *testing.T) {
//	fmt.Printf("%v\n", 0xF)    // 15 	// 1111
//	fmt.Printf("%v\n", 0xFF)   // 255 	// 11111111
//	fmt.Printf("%v\n", 0x200)  // 512 	// 1000000000
//	fmt.Printf("%v\n", 0x0F00) // 3840 	// 000000111100000000
//	fmt.Printf("%v\n", 0x0FFF) // 4095 	// 000000111111111111
//	fmt.Printf("%v\n", 0xF000) // 61440	// 1111000000000000
//}

/*
	fmt.Printf("%v\n", 0xF000)
	i := 0xa2f0 & 0x0FFF
	fmt.Printf("%#04x:%X:%v\n", i, i, i)
	fmt.Printf("%#04x:%X:%v\n", 0xA000, 0xA000, 0xF000)
*/
// 0x0FFF == 4095
// 0xF000 == 61440
