// Package z80 provides the implementation of the Zilog Z80 CPU.
package z80

func init() {
	initTables()
	initMisc()
	initLD()
	initADD()
	initSUB()
	initLogic()
	initIncDec()
	initPushPop()
	initJump()
	initExchange()
	initRot()
	initBCD()
	initBIT()
	initIO()
	initBlock()
	PopulateIndexTables()
}
