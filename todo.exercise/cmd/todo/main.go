package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bernylinville/command/todo.exercise"
)

// Default file name for the todo list
var todoFileName = ".todo.json"

func main() {
	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	list := flag.Bool("list", false, "list all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	add := flag.Bool("add", false, "Add task to the ToDo list")
	del := flag.Int("del", 0, "Delete task from the ToDo list")
	verbose := flag.Bool("v", false, "Show verbose output with additional information like date/time of the task")
	hideCompleted := flag.Bool("hideCompleted", false, "Hide completed items")

	flag.Parse()

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read todo items from the file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on provided flags
	switch {
	// For no extra arguments, print the list
	case *list:
		if *verbose {
			// Show verbose output with additional information like date/time of the task
			fmt.Print(l.Verbose())
			break
		}

		if *hideCompleted {
			// Hide completed items
			fmt.Print(l.HideCompleted())
			break
		}

		// List current todo items
		fmt.Print(l)

	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		// When any arguments (excluding flags) are provided, use them as the new task
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Add tasks to the list
		for _, task := range tasks {
			l.Add(task)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *del > 0:
		// Delete the given item
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}

// getTask function decides where to get the task from: arguments or stdin
func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return args, nil
	}

	var tasks []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		if len(s.Text()) == 0 {
			return nil, fmt.Errorf("task cannot be blank")
		}

		tasks = append(tasks, s.Text())
	}

	return tasks, nil
}
