package main

import (
	"fmt"
	"os"
	"path/filepath"
	"rost/tasklist"
	"strings"
	"tdlib"
)

///////////////////// VARs /////////////////////////

var userTasks tasklist.TaskList

var traceL int = 0

const (
	EXIT     = "q"
	ADD      = "a"
	DEL      = "del"
	SAVE     = "s"
	NAME     = "n"
	MUP      = "u"
	MDOWN    = "d"
	RENAME   = "rn"
	REPLACE  = "rp"
	PRINTRAW = "praw"

	//Default file name with tasks
	defaultTaskFile = "default.tsk"
)

///////////////////// FUNCs /////////////////////////

func printChild(parentName string, list *tasklist.TaskList) {

	if t := list.List[parentName]; t.Subtask != nil {
		for _, p := range t.Subtask {
			d := list.List[p]
			fmt.Println(strings.Repeat(" |", traceL) + "-" + p + " " + d.Done)
			traceL++
			printChild(p, list)
		}
	}
	if traceL > 0 {
		traceL--
	}
}

func printTasks(t *tasklist.TaskList) {

	if len(t.List) < 2 {
		fmt.Println("...no tasks" + "\n")
	} else {
		//fmt.Println("Tasks are:" + "\n")
		fmt.Println()
		printChild("/", t)
		fmt.Println()
	}
}

func printHeader() {
	fmt.Println("====== This is ROST ======")
	fmt.Println("file: " + userTasks.File + " " + userTasks.IsModify)
	fmt.Println("name: " + userTasks.Name)
}

func fileArg() (taskFileName string, err error) {

	// if programm arguments are -h or --help, print usage info and exit
	if len(os.Args) > 1 && (os.Args[1] == "--h" || os.Args[1] == "--help") {
		err = fmt.Errorf("usage: %s userfile.tsk",
			filepath.Base(os.Args[0]))
		return "", err
	}

	//nt := new(tTask)

	//if there is an argument
	if len(os.Args) > 1 {
		taskFileName = os.Args[1]
	}

	//this code is quite clear )
	if taskFileName == "" {
		taskFileName = defaultTaskFile
	}

	return taskFileName, nil
}

func main() {
	cmnd := ""
	exitFlag := ""

	tdlib.ClearScr()

	// Getting filename from command arguments
	filename, err := fileArg()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	userTasks.File = filename
	userTasks.Name = "Testing"

	//Loading tasks data from file
	if err := userTasks.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Main cycle
	for exitFlag != EXIT {
		tdlib.ClearScr()
		printHeader()
		printTasks(&userTasks)
		fmt.Println("Command menu:")
		fmt.Println(NAME + " - rename main task")
		fmt.Println(ADD + " - add task")
		fmt.Println(DEL + " - delete task")
		fmt.Println(SAVE + " - save task list to file")
		fmt.Println(MUP + " - move task up")
		fmt.Println(MDOWN + " - move task down")
		fmt.Println(RENAME + " - rename task")
		fmt.Println(REPLACE + " - replace task")
		fmt.Println(PRINTRAW + " - print raw List")
		fmt.Println(EXIT + " - quit")
		fmt.Println()
		fmt.Print("Enter command: ")
		fmt.Scan(&cmnd)

		if cmnd == ADD {
			fmt.Println("Adding task...")
			fmt.Println("Enter parent task name: ")
			ptn := ""
			fmt.Scan(&ptn)

			if !userTasks.IsParent(ptn) {
				tdlib.Msg("There is no parent with such name!")
				continue
			}

			fmt.Println("Enter new task name: ")
			ntn := ""
			fmt.Scan(&ntn)

			if err := userTasks.Add(ptn, ntn); err != nil {
				tdlib.Msg(err.Error())
			}
		}

		if cmnd == NAME {
			fmt.Print("Enter task list name: ")
			fmt.Scan(&userTasks.Name)
			userTasks.IsModify = "*"
		}

		if cmnd == SAVE {
			if tdlib.YesNo("SAVING... Are you shure?") {
				if err := userTasks.Save(); err != nil {
					tdlib.Msg(err.Error())
				}
			}
		}

		if cmnd == DEL {
			fmt.Println("What task you want to delete?: ")
			fmt.Scan(&cmnd)
			if !userTasks.IsTask(cmnd) {
				tdlib.Msg("There is no such task.")
				continue
			}
			if tdlib.YesNo("DELETING...Are you shure?") {
				if err := userTasks.Del(cmnd); err != nil {
					tdlib.Msg(err.Error())
				}
			}
		}

		if cmnd == MUP {
			fmt.Print("What task to move up: ")
			fmt.Scan(&cmnd)
			if !userTasks.IsTask(cmnd) {
				tdlib.Msg("There is no such task.")
				continue
			}
			if err := userTasks.MoveUp(cmnd); err != nil {
				tdlib.Msg(err.Error())
			}
		}

		if cmnd == MDOWN {
			fmt.Print("What task to move down: ")
			fmt.Scan(&cmnd)
			if !userTasks.IsTask(cmnd) {
				tdlib.Msg("There is no such task.")
				continue
			}
			if err := userTasks.MoveDown(cmnd); err != nil {
				tdlib.Msg(err.Error())
			}
		}

		if cmnd == RENAME {
			fmt.Print("What task to rename: ")
			fmt.Scan(&cmnd)
			if !userTasks.IsTask(cmnd) {
				tdlib.Msg("There is no such task.")
				continue
			}
			nn := ""
			fmt.Print("New name: ")
			fmt.Scan(&nn)

			if err := userTasks.Rename(cmnd, nn); err != nil {
				tdlib.Msg(err.Error())
			}
		}

		if cmnd == REPLACE {
			fmt.Print("What task to replace: ")
			fmt.Scan(&cmnd)
			if !userTasks.IsTask(cmnd) {
				tdlib.Msg("There is no such task.")
				continue
			}
			pn := ""
			fmt.Print("Where: ")
			fmt.Scan(&pn)

			if err := userTasks.Replace(cmnd, pn); err != nil {
				tdlib.Msg(err.Error())
			}
		}

		if cmnd == PRINTRAW {
			for k, t := range userTasks.List {
				tdlib.ClearScr()
				fmt.Print(k + ": ")
				fmt.Println(t)
			}
			tdlib.Msg("")
		}

		if cmnd == EXIT {
			if userTasks.IsModify == "*" {
				if tdlib.YesNo("Save changes?") {
					userTasks.Save()
				}
			}
			exitFlag = EXIT
		}
	}
	tdlib.ClearScr()
}
