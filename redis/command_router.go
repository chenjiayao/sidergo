package redis

type ExecFunc func(args [][]byte)

var (
	commandTables = make(map[string]Command)
)

type Command struct {
	cmdName  string
	execFunc ExecFunc
}

func registerCommand(cmdName string, execFunc ExecFunc) {
	commandTables[cmdName] = Command{
		cmdName:  cmdName,
		execFunc: execFunc,
	}
}
