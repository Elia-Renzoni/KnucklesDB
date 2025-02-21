package wal

type CrashFaultRecover struct {
}

func NewRecover() *CrashFaultRecover {
	return &CrashFaultRecover{}
}
