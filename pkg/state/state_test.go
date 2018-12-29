package state

import (
	"crypto/sha256"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

const (
	C8PIC_sha  = "2f1dece0b277f953c84e0f48ae39a32f53b1796aae67ce43e873d9dd2abb4d0a"
	C8PIC_path = "../../test/roms/C8PIC.ch8"
)

func TestInitMemory(t *testing.T) {
	m := InitMemory()
	assert.NotNil(t, m)
	assert.Len(t, m, 4096)
}

func TestLoadMemory(t *testing.T) {
	m := loadMemoryTest(t)
	hash := sha256.Sum256(m)
	assert.Equal(t,
		C8PIC_sha,
		fmt.Sprintf("%x", hash))
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
	s := InitStack()
	assert.NotNil(t, s)
	assert.Len(t, s.s, 16)
}

func TestStack_PushThenPop(t *testing.T) {
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
