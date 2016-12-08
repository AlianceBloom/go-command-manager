package command_manager

import (
	"log"
	"time"
	"os/exec"
	"os"
	"strconv"
)

const (
	// EXECUTE_COMMAND - execute command
	EXECUTE_COMMAND = 1
	// TERMINATE_COMMAND - kill command
	TERMINATE_COMMAND = -1
	// FINISH_COMMAND - complete command
	FINISH_COMMAND = 0
)

type CommandDescriptor struct {
	command *exec.Cmd
	commandName string
	commandUID  string
	errorCode  int8
	statusChange chan int
}

func CreateCommand(commandName string, commandUID string) *CommandDescriptor {

	command := CommandDescriptor{commandName: commandName, commandUID: commandUID, statusChange: make(chan int)}

	if !existInRedis(command.commandUID) {
		go command.handle()
		command.statusChange <- EXECUTE_COMMAND
	} else {
		return nil
	}

	return &command
}

func GetCommandInfo(commandUID string) string {
	var pidInt int
	var err error

	storablePid := getFromRedisByKey(commandUID)
	if len(storablePid) == 0 {
		return "Process not found"
	}

	pidInt, err = strconv.Atoi(storablePid)
	if err != nil {
		return "Wrong PID passed"
	}

	var process *os.Process
	process, err = os.FindProcess(pidInt)
	if err != nil {
		return "Process not found by PID"
	}

	return "Process working, pid: " + strconv.Itoa(process.Pid)
}

func AbortCommand(commandUID string) string {
	storablePid := getFromRedisByKey(commandUID)
	if len(storablePid) == 0 {
		return "Process not found"
	}
	pidInt, err := strconv.Atoi(storablePid)
	if err != nil {
		return "Wrong PID passed"
	}
	var process *os.Process
	process, err = os.FindProcess(pidInt)
	if err != nil {
		return "Process not found by PID"
	}
	err = process.Signal(os.Kill)

	if err != nil {
		return "Error while process aborting"
	}

	return "Abort process, pid: " + strconv.Itoa(process.Pid)
}

func (cm *CommandDescriptor) handle() {
	defer close(cm.statusChange)
	for {
		select {
			case command := <- cm.statusChange:
				if command == EXECUTE_COMMAND {
					go cm.execute()
				}
				if command == TERMINATE_COMMAND {
					go cm.terminateExecution()
				}
				if command == FINISH_COMMAND {
					cm.finish()
					return
				}
		}
	}
}

func (cm *CommandDescriptor) enableTimeout() {
	timer := time.NewTimer(20 * time.Second)
	go func() {
		<- timer.C
		if cm.statusChange != nil {
			_, ok := <- cm.statusChange

			if !ok {
				log.Println("Chanel of", cm.commandUID, "CLOSED")
			} else {
				cm.statusChange <- TERMINATE_COMMAND
			}
		}

		timer.Stop()
	}()
}

func (cm *CommandDescriptor) terminateExecution() {
	if cm.command.ProcessState != nil && !cm.command.ProcessState.Exited() {
		err := cm.command.Process.Signal(os.Kill)
		log.Println("Command", cm.commandUID, " interrupted")
		deleteFromRedisByKey(cm.commandUID)
		log.Println(err)
	} else {
		log.Println("Command", cm.commandUID, "already interrupted")

	}
}

func (cm *CommandDescriptor) finish()  {
	deleteFromRedisByKey(cm.commandUID)
	log.Println("Command", cm.commandUID, "finished")
}

func (cm *CommandDescriptor) execute() {
	var err error
	log.Println("Command", cm.commandName, "ID", cm.commandUID, "starting")
	cm.command = exec.Command("/bin/sleep", "60")

	err = cm.command.Start()
	PanicError(err)

	cm.enableTimeout()
	storeInRedis(cm.commandUID, cm.command.Process.Pid)

	err = cm.command.Wait()
	if err != nil {
		cm.statusChange <- TERMINATE_COMMAND
		return
	}
	cm.statusChange <- FINISH_COMMAND
}


