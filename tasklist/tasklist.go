package tasklist

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"tdlib"
)

type Task struct {
	Parent  string
	Subtask []string
	Done    string
}

type TaskList struct {
	Name     string          // Name of the Main Task
	List     map[string]Task // List of taskas(key: task name; answer: task data)
	File     string          // filename.tsk
	IsModify string          // flag of editing
}

// Func to ADD task in List
func (t *TaskList) Add(tParent, newTaskName string) (err error) {
	nt := Task{
		Parent:  tParent,
		Subtask: nil,
		Done:    "n"}

	//Checking the parent
	if st, foundparent := t.List[nt.Parent]; foundparent {
		//check overlaping task name in list
		if _, founddouble := t.List[newTaskName]; !founddouble {
			t.List[newTaskName] = nt
			st.Subtask = append(st.Subtask, newTaskName)
			t.List[nt.Parent] = st
		} else {
			err = errors.New("This task name is alredy in list!")
			return err
		}
	} else {
		err = errors.New("There is no parent with such name!")
		return err
	}

	t.IsModify = "*"
	return nil
}

// Method to save task data in file
func (t *TaskList) Save() (err error) {
	fmt.Println("saving...")
	//Preapreing buffered output to file
	outFile := os.Stdout
	if outFile, err = os.Create(t.File); err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	//FLushing buffer if error
	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	//Writing name of tasks in 1st line
	writeLine := "N: " + t.Name + "\n"
	if _, err = writer.WriteString(writeLine); err != nil {
		return err
	}
	//creating string with task data and write it into the file
	for taskname, tsk := range t.List {
		writeLine = "T: " + taskname + " " + tsk.Parent + " " + tsk.Done
		if tsk.Subtask != nil {
			for _, subname := range tsk.Subtask {
				writeLine += " " + subname
			}
		} //end of subtask
		writeLine += "\n"
		// Writing generated line to the file
		if _, err = writer.WriteString(writeLine); err != nil {
			return err
		} //end of tasks in List
	}

	t.IsModify = ""
	return nil
}

// Method to load task data from file
func (t *TaskList) Load() (err error) {

	t.List = make(map[string]Task)

	//Prepareing buffered reading from file
	inFile := os.Stdin
	if inFile, err = os.Open(t.File); err != nil {
		inFile, err = os.Create((t.File))
		// create root task for new list
		rootTask := Task{
			Parent:  "root",
			Subtask: nil,
			Done:    "n"}
		t.List["/"] = rootTask
		return nil
	}
	defer inFile.Close()
	reader := bufio.NewReader(inFile)

	eof := false

	for !eof {
		readln := ""

		//reading string line from file
		readln, err = reader.ReadString('\n')

		if err == io.EOF { // handling err
			err = nil //its not real error
			eof = true
		} else if err != nil {
			return err // if there is real error, then exit
		}

		sfields := strings.Fields(readln)

		if len(sfields) > 0 {
			//Reading Name of Tasks
			if sfields[0] == "N:" {
				t.Name = sfields[1]
			}
			// Reading tasks
			// line format is:
			// sfields[0] - line descriptor(N: - name of tasks, T: - task data line
			// sfields[1] - task name
			// sfields[2] - task parent
			// sfields[3] - task status, Done "y"es or "n"o
			// sfields[4...n] - list of subtasks

			if sfields[0] == "T:" {
				nt := Task{
					Subtask: nil,
					Parent:  sfields[2],
					Done:    sfields[3]}
				if len(sfields) > 4 {
					for i := 4; i < len(sfields); i++ {
						nt.Subtask = append(nt.Subtask, sfields[i])
					}
				}
				t.List[sfields[1]] = nt
				//		fmt.Println(nt)
			}
		}
	}
	//	s := ""
	//	fmt.Scan(&s)
	return nil
}

// REcursive function for deleting parent task with all subtasks
func deleteSubs(taskName string, list *TaskList) {
	//If Task have subtask
	if t := list.List[taskName]; t.Subtask != nil {
		//Run subs
		for _, p := range t.Subtask {
			deleteSubs(p, list)
		}
	}
	parentName := list.List[taskName].Parent
	delete(list.List, taskName) //Deleting from List
	//Deleting form parent subtask list
	par := list.List[parentName]
	par.Subtask = tdlib.DelStringFromSlice(taskName, par.Subtask)
	list.List[parentName] = par
}

func (t *TaskList) Del(tasktodelete string) (err error) {
	var found bool

	// Checking existence of the Task
	if _, found = t.List[tasktodelete]; !found {
		err := errors.New("Task not found.")
		return err
	}

	//Delete Task with all Subs
	deleteSubs(tasktodelete, t)

	t.IsModify = "*"
	//tdlib.Msg(t.IsModify)
	return nil
}

// function for rename task
func (t *TaskList) Rename(oldName, newName string) (err error) {
	if _, found := t.List[oldName]; !found {
		err = errors.New("Task not found.")
		return err
	}
	if _, found := t.List[newName]; found {
		err = errors.New("This task name is alredy in list!")
		return err
	}
	// make copy of the required oldName task
	// delete oldName task and its record from the parent subtask list
	// create new Task in the t.List with newName (key)
	// add newName into parents subtask list !!!SAVE ORDER SUBTASK IN THE LIST
	nt := t.List[oldName]
	delete(t.List, oldName)
	t.List[newName] = nt

	p := t.List[nt.Parent]
	for i, s := range p.Subtask {
		if s == oldName {
			p.Subtask[i] = newName
		}
	}
	t.List[nt.Parent] = p

	var tsk Task
	for _, s := range nt.Subtask {
		tsk = t.List[s]
		tsk.Parent = newName
	}

	t.IsModify = "*"
	return nil
}

// func to replace current task to the different task(parent)
func (t *TaskList) Replace(whatTask, toParent string) (err error) {
	if !IsTask(whatTask, t) {
		err := errors.New("Task not found")
		return err
	}
	if !IsParent(toParent, t) {
		err := errors.New("Parent Task not found")
		return err
	}

	p := t.List[t.List[whatTask].Parent]
	p.Subtask = tdlib.DelStringFromSlice(whatTask, p.Subtask)
	t.List[t.List[whatTask].Parent] = p

	p = t.List[toParent]
	p.Subtask = append(p.Subtask, whatTask)
	t.List[toParent] = p

	p = t.List[whatTask]
	p.Parent = toParent
	t.List[whatTask] = p

	t.IsModify = "*"
	return nil
}

// function to change task position in list (UP)
func (t *TaskList) MoveUp(taskName string) (err error) {

	if _, found := t.List[taskName]; !found {
		err := errors.New("Task not found")
		return err
	}

	p := t.List[t.List[taskName].Parent]

	for i, s := range p.Subtask {
		if s == taskName && i > 0 {
			p.Subtask[i], p.Subtask[i-1] = p.Subtask[i-1], p.Subtask[i]
		}
	}
	t.List[t.List[taskName].Parent] = p
	t.IsModify = "*"
	return nil
}

// function to change task position in list (DOWN)
func (t *TaskList) MoveDown(taskName string) (err error) {

	if _, found := t.List[taskName]; !found {
		err := errors.New("Task not found")
		return err
	}

	p := t.List[t.List[taskName].Parent]

	for i, s := range p.Subtask {
		if s == taskName && i != (len(p.Subtask)-1) {
			p.Subtask[i], p.Subtask[i+1] = p.Subtask[i+1], p.Subtask[i]
		}
	}
	t.List[t.List[taskName].Parent] = p
	t.IsModify = "*"

	return nil
}

func (t *TaskList) IsParent(parentName string) bool {
	if _, found := t.List[parentName]; found {
		return true
	}
	return false
}

func (t *TaskList) IsTask(taskName string) bool {
	if _, found := t.List[taskName]; found {
		return true
	}
	return false
}

func IsParent(parentName string, t *TaskList) bool {
	if _, found := t.List[parentName]; found {
		return true
	}
	return false
}

func IsTask(taskName string, t *TaskList) bool {
	if _, found := t.List[taskName]; found {
		return true
	}
	return false
}
