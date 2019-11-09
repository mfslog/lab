package server


type Server interface {
	Init()
	Start()
	Stop()
	Handle()
}
