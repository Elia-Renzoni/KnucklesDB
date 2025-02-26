package wal

type CrashFaultRecover struct {
	wal *WAL
}

func NewRecover(wal *WAL) *CrashFaultRecover {
	return &CrashFaultRecover{
		wal: wal,
	}
}

func (c *CrashFaultRecover) StartRecoveryProcedure() {

}
