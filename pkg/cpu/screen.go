package cpu

type Screen interface {
	Draw(frameBuffer []byte)
}
