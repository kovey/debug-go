package debug

import "testing"

func TestLogDebug(t *testing.T) {
	log := LogWith("traceId", "spanId")
	SetFileLine(File_Line_Show)
	log.Info("test [%s] ...", Debug_Info)
	log.Dbug("test [%s] ...", Debug_Dbug)
	log.Warn("test [%s] ...", Debug_Warn)
	log.Erro("test [%s] ...", Debug_Erro)
	log.Test("test [%s] ...", Debug_Test)
	SetFileLine(File_Line_Off)
	log.Info("test [%s] ...", Debug_Info)
	log.Dbug("test [%s] ...", Debug_Dbug)
	log.Warn("test [%s] ...", Debug_Warn)
	log.Erro("test [%s] ...", Debug_Erro)
	log.Test("test [%s] ...", Debug_Test)
}
