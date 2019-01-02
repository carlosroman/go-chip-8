package cpu

import (
	"encoding/binary"
	"github.com/carlosroman/go-chip-8/pkg/state"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type cpu struct {
	m     state.Memory // CPU Memory
	pc    int16        // Program counter
	ir    uint16       // Index register - 16bit register (For memory address) (Similar to void pointer)
	sp    int16        // Stack pointer
	stack *state.Stack // Stack
	v     []byte       // CPU registers
	r     *rand.Rand   // Random number generator
}

func (c *cpu) Tick() (err error) {
	opcode := binary.BigEndian.Uint16([]byte{c.m[c.pc], c.m[c.pc+1]})
	//fmt.Printf("%#04x:%X:[%#02x %#02x]:%v\n", opcode, opcode, c.m[c.pc], c.m[c.pc+1], opcode)
	//opCodeA := (uint16(c.m[c.pc]) << 8) | uint16(c.m[c.pc+1])
	//fmt.Printf("%#04x:%X:[%#02x %#02x]:%v\n", opCodeA, opCodeA, c.m[c.pc], c.m[c.pc+1], opCodeA)

	switch val := opcode & 0xF000; val {
	case 0x0000: // 0x00
		switch sub := opcode & 0x00FF; sub {
		case 0x00E0:
			log.Info("Opcode: 00E0")
			// 0x00E0, Display, disp_clear(), Clears the screen.
			log.Debug("Clear screen")
		case 0x00EE:
			// 0x00EE, Flow, return;, Returns from a subroutine.
			log.Info("Opcode: 00EE")
			c.pc = c.stack.Pop()
		default:
			log.Warnf("Unknown opcode [0x0000]: %#04x:%#04x\n", val, sub)
		}
	case 0xA000:
		// 0xANNN, MEM, I = NNN, Sets I to the address NNN.
		log.Info("Opcode: 0xANNN")
		c.ir = opcode & 0x0FFF
		c.pc += 2
	case 0xB000:
		// 0xBNNN, Flow, PC=V0+NNN , Jumps to the address NNN plus V0.
		log.Info("Opcode: 0xBNNN")
		nnn := opcode & 0x0FFF
		c.pc = int16(c.v[0]) + int16(nnn)
	case 0xC000:
		// 0xCXNN, Rand, Vx=rand()&NN, Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		log.Info("Opcode: 0xCXNN")
		r := c.r.Intn(256)
		x := getX(opcode)
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, opcode&0x00FF)
		nn := bs[1]
		c.v[x] = byte(r) & nn
		c.pc += 2
	case 0x1000:
		// 0x1NNN, Flow, goto NNN;, Jumps to address NNN.
		log.Info("Opcode: 1NNN")
		nnn := opcode & 0x0FFF
		log.Debugf("nnn:%v", nnn)
		c.pc = int16(nnn)
	case 0x2000:
		// 0x2NNN, Flow, *(0xNNN)(), Calls subroutine at NNN.
		log.Info("Opcode: 2NNN")
		nnn := opcode & 0x0FFF
		log.Debugf("nnn:%v", nnn)
		c.stack.Push(c.pc)
		c.pc = int16(nnn)
	case 0x3000:
		// 0x3XNN, Cond, if(Vx==NN), Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
		log.Info("Opcode: 3XNN")
		x := getX(opcode)
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, opcode&0x00FF)
		c.pc += 2
		if c.v[x] == bs[1] {
			c.pc += 2
		}
	case 0x4000:
		// 0x4XNN, Cond, if(Vx!=NN), Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		log.Info("Opcode: 4XNN")
		x := getX(opcode)
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, opcode&0x00FF)
		c.pc += 2
		if c.v[x] != bs[1] {
			c.pc += 2
		}
	case 0x5000:
		// 0x5XY0, Cond, if(Vx==Vy) 	Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
		log.Info("Opcode: 5XY0")
		x, y := getXY(opcode, c)
		c.pc += 2
		if c.v[x] == c.v[y] {
			c.pc += 2
		}
	case 0x6000:
		// 0x6XNN, Const, Vx = NN, Sets VX to NN.
		log.Info("Opcode: 6XNN")
		x := getX(opcode)
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, opcode&0x00FF)
		c.v[x] = bs[1]
		c.pc += 2
	case 0x7000:
		// 0x7XNN, Const, Vx += NN, Adds NN to VX. (Carry flag is not changed)
		log.Info("Opcode: 7XNN")
		x := getX(opcode)
		bs := make([]byte, 2)
		binary.BigEndian.PutUint16(bs, opcode&0x00FF)
		c.v[x] += bs[1]
		c.pc += 2
	case 0x8000: // 0x8
		switch sub := opcode & 0x000F; sub {
		case 0x0000:
			// 0x8XY0, Assign, Vx=Vy, Sets VX to the value of VY.
			log.Info("Opcode: 0x8XY0")
			x, y := getXY(opcode, c)
			c.v[x] = c.v[y]
		case 0x0001:
			// 0x8XY1, BitOp, Vx=Vx|Vy, Sets VX to VX or VY. (Bitwise OR operation)
			log.Info("Opcode: 0x8XY1")
			x, y := getXY(opcode, c)
			c.v[x] |= c.v[y]
		case 0x0002:
			// 0x8XY2, BitOp, Vx=Vx&Vy, Sets VX to VX and VY. (Bitwise AND operation)
			log.Info("Opcode: 0x8XY2")
			x, y := getXY(opcode, c)
			c.v[x] &= c.v[y]
		case 0x0003:
			// 0x8XY3, BitOp, Vx=Vx^Vy, Sets VX to VX xor VY.
			log.Info("Opcode: 0x8XY3")
			x, y := getXY(opcode, c)
			c.v[x] ^= c.v[y]
		case 0x0004:
			// 0x8XY4, Math, Vx += Vy , Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			log.Info("Opcode: 0x8XY4")
			x, y := getXY(opcode, c)
			if c.v[y] > (0xFF - c.v[x]) {
				log.Debug("carrying the one")
				c.v[0xF] = 1 // carry
			} else {
				c.v[0xF] = 0
			}
			c.v[x] += c.v[y]
		case 0x0005:
			// 0x8XY5, Math, Vx -= Vy, VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			log.Info("Opcode: 0x8XY5")
			x, y := getXY(opcode, c)
			if c.v[y] > c.v[x] {
				log.Debug("borrowing")
				c.v[0xF] = 0 // borrow
			} else {
				c.v[0xF] = 1
			}
			c.v[x] -= c.v[y]
		case 0x0006:
			// 0x8XY6, BitOp, Vx>>=1, Stores the least significant bit of VX in VF and then shifts VX to the right by 1.
			log.Info("Opcode: 0x8XY6")
			x, _ := getXY(opcode, c)
			c.v[0xF] = c.v[x] & 0x1
			c.v[x] >>= 1
		case 0x0007:
			// 0x8XY7, Math, Vx=Vy-Vx, Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			log.Info("Opcode: 0x8XY7")
			x, y := getXY(opcode, c)
			if c.v[x] > c.v[y] {
				log.Debug("borrowing")
				c.v[0xF] = 0 // borrow
			} else {
				c.v[0xF] = 1
			}
			c.v[x] = c.v[y] - c.v[x]
		case 0x000E:
			// 0x8XYE, BitOp, Vx<<=1, Stores the most significant bit of VX in VF and then shifts VX to the left by 1.
			log.Info("Opcode: 0x8XYE")
			x, _ := getXY(opcode, c)
			c.v[0xF] = c.v[x] >> 7
			c.v[x] <<= 1
		}
		c.pc += 2
	case 0x9000:
		// 0x9XY0, Cond, if(Vx!=Vy), Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
		log.Info("Opcode: 9XY0")
		x, y := getXY(opcode, c)
		c.pc += 2
		if c.v[x] != c.v[y] {
			c.pc += 2
		}
	case 0xF000:
		switch sub := opcode & 0x00FF; sub {
		case 0x001E:
			// 0xFX1E, MEM, I +=Vx 	Adds VX to I.
			log.Info("Opcode: FX1E")
			x := getX(opcode)
			c.ir += uint16(c.v[x])
		case 0x0029:
			// 0xFX29, MEM, I=sprite_addr[Vx], Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			log.Info("Opcode: FX29")
			x := getX(opcode)
			c.ir = uint16(c.v[x] * 0x5)
		case 0x55:
			// 0xFX55, MEM, reg_dump(Vx,&I), Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.
			log.Info("Opcode: FX55")
			x := getX(opcode)
			for i := uint16(0); i <= x; i++ {
				c.m[c.ir+i] = c.v[i]
			}
			c.ir += x + 1
			log.Infof("c.ir=%v", c.ir)
		default:
			log.Warnf("Unknown opcode [0xF000]: %#04x:%#04x\n", val, sub)
		}
		c.pc += 2
	default:
		log.Debugf("Unknown opcode: %#04x:%X\n", val, val)
	}
	return err
}

func getX(opcode uint16) (x uint16) {
	x = (opcode & 0x0F00) >> 8
	if log.IsLevelEnabled(log.DebugLevel) {
		log.
			WithField("x", x).
			Debug("Got x")
	}
	return x
}

func getXY(opcode uint16, c *cpu) (x uint16, y uint16) {
	x = (opcode & 0x0F00) >> 8
	y = (opcode & 0x00F0) >> 4
	if log.IsLevelEnabled(log.DebugLevel) {
		log.
			WithField("y", y).
			WithField("vy", c.v[y]).
			WithField("x", x).
			WithField("vx", c.v[x]).
			Debug("Got vy vx")
	}
	return x, y
}

func NewCPU(memory state.Memory, rgen *rand.Rand) *cpu {
	return &cpu{
		m:     memory,
		pc:    0x200,            // Program counter starts at 0x200 (512)
		v:     make([]byte, 16), // The Chip 8 has 15 8-bit general purpose registers and the 16th register is used  for the ‘carry flag’.
		stack: state.InitStack(),
		r:     rgen,
	}
}
