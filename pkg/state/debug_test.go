//+build debug

package state

import (
	"encoding/binary"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestPrintDump(t *testing.T) {
	m := loadMemoryTest(t)

	for i := range m {
		if i < 0x200 { // 0x200 == 512
			continue
		}

		if i%2 != 0 {
			continue
		}

		opCode := binary.BigEndian.Uint16([]byte{m[i], m[i+1]})
		fmt.Printf("%#04x:%X:%#04x:[%#02x %#02x]\n", opCode, opCode, opCode&0x0FFF, m[i], m[i+1])
	}
}
