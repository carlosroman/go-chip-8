package state

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

const (
	C8PIC_sha  = "4bbb2f95d64e976bb90837df5c9248e4632974c4137759d33589891ee3622977" // includes fonts
	C8PIC_path = "../../test/roms/C8PIC.ch8"
)

func TestInitMemory(t *testing.T) {
	t.Parallel()
	m := InitMemory()
	assert.NotNil(t, m)
	assert.Len(t, m, 4096)
}

func TestLoadMemory(t *testing.T) {
	t.Parallel()
	m := loadMemoryTest(t)
	hash := sha256.Sum256(m)
	assert.Equal(t,
		C8PIC_sha,
		fmt.Sprintf("%x", hash))
	expF := getFonts()
	for i := 0; i < 80; i++ {
		assert.Equal(t, expF[i], m[i], "Expected m[%v] to have value '%v'", i, expF[i])
	}
}

func loadMemoryTest(tb testing.TB) Memory {
	m := InitMemory()
	f, err := os.Open(C8PIC_path)
	assert.NoError(tb, err)
	err = m.LoadMemory(f)
	assert.NoError(tb, err)
	return m
}

func TestInitStack(t *testing.T) {
	t.Parallel()
	s := InitStack()
	assert.NotNil(t, s)
	assert.Len(t, s.s, 16)
}

func TestStack_PushThenPop(t *testing.T) {
	t.Parallel()
	s := InitStack()
	ex := int16(1337)
	for i := int16(0); i < 16; i++ {
		s.Push(ex + i)
	}
	for i := int16(15); i > -1; i-- {
		ac := s.Pop()
		assert.Equal(t, ex+i, ac)
	}
}
