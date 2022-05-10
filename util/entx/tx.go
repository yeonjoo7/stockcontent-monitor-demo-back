package entx

import "stockcontent-monitor-demo-back/ent"

type Tx interface {
	GetTx() *ent.Tx
}
