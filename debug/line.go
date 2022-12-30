package debug

type FileLine int

const (
	File_Line_Off  FileLine = 0
	File_Line_Show FileLine = 1
)

var fileLineSwitch FileLine = File_Line_Off

func SetFileLine(fl FileLine) {
	fileLineSwitch = fl
}
