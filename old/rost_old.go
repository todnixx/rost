// rost.go
package main

import (
	"bufio"
	"fmt"
	//	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"taskpak/task"
)

type task struct {
	name string
	done int
	task []task
}

type tTask struct {
	task     []task
	ntask    int
	taskfile string
}

//Default file name with tasks
const defaultTaskFile = "default.tsk"

//Method for loading file with tasks
//If there is no file with specified file name, this method create empty file
func (t *tTask) Load() bool {

	ufio := os.Stdin
	var err error

	// Try to open file
	ufio, err = os.Open(t.taskfile)
	defer ufio.Close()

	if err != nil {
		// If it doesn`t exist, create it
		ufio, err = os.Create(t.taskfile)
		if err != nil {
			log.Fatal(err)
		} else {
			defer ufio.Close()
			fmt.Println("- Creating file")
			t.ntask = 1
			nt := task{"No task", 0, nil}
			t.task = append(t.task, nt)
		}
	} else if err == nil {
		scanner := bufio.NewScanner(ufio)
		if scanner.Scan() {
			rl := scanner.Text()
			t.ntask, _ = strconv.Atoi(rl)
		}

		for n := 0; n < t.ntask; n++ {
			nt := task{"", 0, nil}
			if scanner.Scan() {
				nt.name = scanner.Text()
			}
			if scanner.Scan() {
				rl := scanner.Text()
				nt.done, _ = strconv.Atoi(rl)
			}
			t.task = append(t.task, nt)

		}
	}

	return true

} //Load

//Method for saving file with tasks
func (t *tTask) Save() {

	var err error

	ufio, err := os.Create(t.taskfile)
	if err != nil {
		log.Fatal(err)
	}
	defer ufio.Close()

	ufio.WriteString(strconv.Itoa(t.ntask) + "\n")

	for i := 0; i < t.ntask; i++ {
		ufio.WriteString(t.task[i].name + "\n")
		ufio.WriteString(strconv.Itoa(t.task[i].done) + "\n")
	}

	fmt.Println("Tasks are saved!")

} //Save

func (t *tTask) Add() {
	nt := task{"", 0, nil}

	fmt.Println("Adding task...")
	fmt.Print("Name: ")
	fmt.Scan(&nt.name)

	t.ntask++
	t.task = append(t.task, nt)

}

//Method for printing tasks tree
func (t *tTask) Print() {
	fmt.Println(t.ntask)
	for i := 0; i < t.ntask; i++ {
		fmt.Print(t.task[i].name + " ")
		fmt.Println(t.task[i].done)
	}
	fmt.Println("")
}

func clScr() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printMenu(uf string) {
	fmt.Println("+------------------------------------+")
	fmt.Println("| This is ROST in file " + uf + "   |")
	fmt.Println("+------------------------------------+")
	fmt.Println("")
	fmt.Println("============== MENU ==============")
	fmt.Println("= (A)dd === (D)elete ==== (Q)uit =")
	fmt.Println("")
}

func fileArg() (userFile string, err error) {

	// if programm arguments are -h or --help, print usage info and exit
	if len(os.Args) > 1 && (os.Args[1] == "--h" || os.Args[1] == "--help") {
		err = fmt.Errorf("usage: %s userfile.tsk",
			filepath.Base(os.Args[0]))
		return "", err
	}

	//if there is an argument
	if len(os.Args) > 1 {
		userFile = os.Args[1]
	}

	//this code is quite clear )
	if userFile == "" {
		userFile = defaultTaskFile
	}

	return userFile, nil
}

func main() {
	var Task tTask
	cmnd := ""
	var err error

	Task.taskfile, err = fileArg()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if Task.Load() {

		Exit := false

		// main cycle
		for !Exit {
			clScr()
			printMenu(Task.taskfile)
			Task.Print()
			fmt.Print("Enter command ((A) or (D) or (Q)): ")
			fmt.Scan(&cmnd)
			if cmnd == "q" || cmnd == "Q" {
				Exit = true
			}
			if cmnd == "a" || cmnd == "A" {
				Task.Add()

			}
		}
	}

	SAVE := false
	for !SAVE {
		fmt.Print("Save changes to " + Task.taskfile + "? (Y/n):")
		fmt.Scan(&cmnd)
		if cmnd == "Y" || cmnd == "y" || cmnd == "N" || cmnd == "n" {
			SAVE = true
		}
	}

	if cmnd == "Y" || cmnd == "y" {
		Task.Save()
	}
} // main
