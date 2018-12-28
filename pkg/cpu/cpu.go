package cpu

import (
	"encoding/binary"
	"fmt"
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
)

type cpu struct {
	m     state.Memory // CPU Memory
	pc    int16        // Program counter
	ir    uint16       // Index register - 16bit register (For memory address) (Similar to void pointer)
	sp    int16        // Stack pointer
	stack *state.Stack // Stack
	v     []byte       // CPU registers
}

func (c *cpu) Tick() (err error) {
	opcode := binary.BigEndian.Uint16([]byte{c.m[c.pc], c.m[c.pc+1]})
	//fmt.Printf("%#04x:%X:[%#02x %#02x]:%v\n", opcode, opcode, c.m[c.pc], c.m[c.pc+1], opcode)
	//opCodeA := (uint16(c.m[c.pc]) << 8) | uint16(c.m[c.pc+1])
	//fmt.Printf("%#04x:%X:[%#02x %#02x]:%v\n", opCodeA, opCodeA, c.m[c.pc], c.m[c.pc+1], opCodeA)

	switch val := opcode & 0xF000; val {
	case 0x0000: // 0x00
		switch sub := opcode & 0x000F; sub {
		case 0x0000:
			// 0x00E0, Display, disp_clear(), Clears the screen.
			log.Debugf("Clear screen")
		case 0x00EE:
			// 0x00EE, Flow, return;, Returns from a subroutine.
			log.Debugf("Return from a subroutine")
		default:
			fmt.Printf("Unknown opcode [0x0000]: %#04x:%X\n", val, val)
		}
	case 0xA000:
		// 0xANNN, MEM, I = NNN, Sets I to the address NNN.
		log.Debugf("Opcode: 0xANNN")
		c.ir = opcode & 0x0FFF
		c.pc += 2
	case 0x8000: // 0x8
		switch sub := opcode & 0x000F; sub {
		case 0x0004:
			// 0x8XY4, Math, Vx += Vy , Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			log.Debugf("Opcode: 0x0004")
			vy := (opcode & 0x00F0) >> 4
			vx := (opcode & 0x0F00) >> 8
			log.
				WithField("vy", vy).
				WithField("y", c.v[vy]).
				WithField("vx", vx).
				WithField("x", c.v[vx]).
				Debug("Got vy vx")
			if c.v[vy] > (0xFF - c.v[vx]) {
				log.Debug("carrying the one")
				c.v[0xF] = 1 // carry
			} else {
				c.v[0xF] = 0
			}
			c.v[vx] += c.v[vy]
			c.pc += 2
		}
	default:
		log.Debugf("Unknown opcode: %#04x:%X\n", val, val)
	}
	return err
}

func NewCPU(memory state.Memory) *cpu {
	return &cpu{
		m:     memory,
		pc:    0x200,            // Program counter starts at 0x200 (512)
		v:     make([]byte, 16), // The Chip 8 has 15 8-bit general purpose registers and the 16th register is used  for the ‘carry flag’.
		stack: state.InitStack(),
	}
}
