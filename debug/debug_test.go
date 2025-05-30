package debug

import "testing"

func TestDebug(t *testing.T) {
	SetFileLine(File_Line_Show)
	Info("test [%s] ...", Debug_Info)
	Dbug("test [%s] ...", Debug_Dbug)
	Warn("test [%s] ...", Debug_Warn)
	Erro("test [%s] ...", Debug_Erro)
	Test("test [%s] ...", Debug_Test)
	SetFileLine(File_Line_Off)
	Info("test [%s] ...", Debug_Info)
	Dbug("test [%s] ...", Debug_Dbug)
	Warn("test [%s] ...", Debug_Warn)
	Erro("test [%s] ...", Debug_Erro)
	Test("test [%s] ...", Debug_Test)
}
