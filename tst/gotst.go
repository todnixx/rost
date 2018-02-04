package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

///////////////////// VARs /////////////////////////

type task struct {
	parent *task
	index  int
	count  int
	name   string
	done   string
}

type tTask struct {
	TaskListName string
	task         []task
}

var traceL int = 0

const (
	EXIT = "q"
	ADD  = "a"
	DEL  = "d"
	SAVE = "s"
)

///////////////////// FUNCs /////////////////////////

func clScr() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (t *tTask) findParent(parentName string) (parentTask task, err error) {

	var nt task

	for index, tsk := range t.task {
		if tsk.name == parentName {
			return t.task[index], nil
		}
	}
	err = errors.New("There is no parent with such name!")
	return nt, err
}

func (t *tTask) add() (err error) {
	nt := task{parent: nil,
		index: 0,
		count: 0,
		name:  "",
		done:  "n"}

	fmt.Println("Adding task...")
	fmt.Println("Enter parent task name: ")

	parentName := ""
	fmt.Scan(&parentName)

	fmt.Println("Enter new task name: ")
	fmt.Scan(&nt.name)

	if parentName == "/" {
		if len(t.task) == 0 {
			t.task = append(t.task, nt)
			return nil
		}

		nt.parent = &t.task[0]
		t.task[0].count++
		nt.index = nt.parent.count - 1
		t.task = append(t.task, nt)
		return nil
	}

	var parentPtr *task = nil

	for index, tsk := range t.task {
		if tsk.name == parentName {
			parentPtr = &t.task[index]
		}
	}
	if parentPtr == nil {
		err = errors.New("There is no parent with such name!")
		return err
	}
	nt.parent = parentPtr
	parentPtr.count++
	nt.index = nt.parent.count - 1
	t.task = append(t.task, nt)

	return nil
}

func (t *tTask) save() (err error) {

}

func findChild(parentName string, taskList []task) {
	//fmt.Println(traceL)
	for i := 0; i < len(taskList); i++ {
		if taskList[i].parent != nil {
			if taskList[i].parent.name == parentName {
				//	fmt.Print(" ")
				ss := strings.Repeat(" |", traceL)
				fmt.Println(" |" + ss + "-" + taskList[i].name)
				if taskList[i].count > 0 {
					traceL++
					findChild(taskList[i].name, taskList)
				}
			}
		}
	}
	if traceL > 0 {
		traceL--
	}
}

func (t *tTask) print() {

	if len(t.task) == 0 {
		fmt.Println("There are no tasks." + "\n")

	} else {
		fmt.Println("Tasks are:" + "\n")
		fmt.Println(t.task[0].name)
		findChild(t.task[0].name, t.task)
		fmt.Println("")
	}
}

func main() {
	Task := tTask{"Main Task", nil}
	cmnd := ""

	for cmnd != EXIT {
		clScr()
		Task.print()
		fmt.Printf("Enter command %s, %s, %s: ", ADD, DEL, EXIT)
		fmt.Scan(&cmnd)

		if cmnd == ADD {
			if err := Task.add(); err != nil {
				fmt.Println(err)
				fmt.Scan(&cmnd)
			}
		}
	}
	clScr()
}
