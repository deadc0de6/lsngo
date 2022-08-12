/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2022, deadc0de6
*/

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	log "github.com/sirupsen/logrus"

	"github.com/rivo/tview"
)

const (
	version = "0.1.5"
)

var (
	editorEnv  = "EDITOR"
	shellEnv   = "SHELL"
	colorReset = "[-]"
	units      = []string{"", "K", "M", "G", "T", "P"}

	help = `
		j: down
		k: up
		h: go to parent directory
		l: open file/directory
		q: exit
		esc: exit
		H: toggle hidden files
		L: toggle long format
		enter: open file/directory
		?: show this help.
	`
)

type nav struct {
	app        *tview.Application
	list       *tview.List
	textarea   *tview.TextView
	editor     string
	exit       bool
	goBack     bool
	selected   bool
	showHidden bool
	extended   bool
	help       bool
}

// returns files in directory
func getFiles(path string, showHidden bool) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var infos []fs.FileInfo
	for _, entry := range entries {
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			log.Error(err)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// handle user input
func (n *nav) eventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Rune() == 'q' {
		// exit
		n.exit = true
		n.app.Stop()
		return nil
	} else if eventKey.Key() == tcell.KeyEscape {
		// exit
		n.exit = true
		n.app.Stop()
		return nil
	} else if eventKey.Rune() == '?' {
		// help
		n.help = true
		n.app.Stop()
		return nil
	} else if eventKey.Key() == tcell.KeyEnter {
		// open
		n.selected = true
		n.app.Stop()
		return nil
	} else if eventKey.Rune() == 'j' {
		// down
		idx := (n.list.GetCurrentItem() + 1) % n.list.GetItemCount()
		n.list.SetCurrentItem(idx)
		return nil
	} else if eventKey.Rune() == 'k' {
		// up
		idx := n.list.GetCurrentItem() - 1
		n.list.SetCurrentItem(idx)
		return nil
	} else if eventKey.Rune() == 'l' || eventKey.Key() == tcell.KeyRight {
		// open
		n.selected = true
		n.app.Stop()
		return nil
	} else if eventKey.Rune() == 'h' || eventKey.Key() == tcell.KeyLeft {
		// open parent directory
		n.goBack = true
		n.app.Stop()
		return nil
	} else if eventKey.Rune() == 'H' {
		// toggle hidden files
		n.showHidden = !n.showHidden
		n.app.Stop()
		return nil
	} else if eventKey.Rune() == 'L' {
		// toggle extended mode
		n.extended = !n.extended
		n.app.Stop()
		return nil
	}
	return eventKey
}

// returns true if p is a directory
func isDir(p string) bool {
	fileInfo, err := os.Stat(p)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// return size in human readable format
func sizeToHuman(bytes int64) string {
	size := bytes
	for _, unit := range units {
		if size < 1024 {
			return fmt.Sprintf("%d%s", size, unit)
		}
		size = size / 1024
	}
	return "??"
}

// show ls long view
func extendedView(info fs.FileInfo) string {
	mode := info.Mode().String()
	size := sizeToHuman(info.Size())
	name := info.Name()
	dt := info.ModTime()
	dtstr := dt.Format("Jan 02 2006")
	if dt.Year() == time.Now().Year() {
		dtstr = dt.Format("Jan 02 15:04")
	}
	line := fmt.Sprintf("%s %4s %12s %s", mode, size, dtstr, name)
	return line
}

// file list with file infos
func (n *nav) fillList(finfos []fs.FileInfo) {
	// insert ".."
	n.list.InsertItem(-1, "..", "", 0, nil)

	// insert entries
	for _, info := range finfos {
		line := info.Name()

		if n.extended {
			line = extendedView(info)
		}

		preColor := ""
		if info.IsDir() {
			// style directory
			preColor = "[blue]"
		}
		if info.Mode()&fs.ModeSymlink > 0 {
			// style symlink
			preColor = "[cyan]"
		}
		if info.Mode()&fs.ModeSetuid > 0 {
			// style setuid
			preColor = "[red]"
		}

		line = preColor + line + colorReset
		n.list.InsertItem(-1, line, "", 0, nil)
	}
}

// update the list with new content
func (n *nav) updateList(path string, finfos []fs.FileInfo) {
	pwd, err := filepath.Abs(path)
	if err != nil {
		pwd = path
	}

	n.list.Clear()
	n.fillList(finfos)
	n.textarea.SetText(pwd)
}

// create the view
func (n *nav) createList(path string, finfos []fs.FileInfo) {
	pwd, err := filepath.Abs(path)
	if err != nil {
		pwd = path
	}

	// create list
	n.list = tview.NewList()
	n.list.ShowSecondaryText(false)
	n.list.SetWrapAround(true)
	n.fillList(finfos)

	// current working directory
	n.textarea = tview.NewTextView()
	n.textarea.SetText(pwd)
	n.textarea.SetTextColor(tcell.ColorSlateGray)

	// create layout
	content := tview.NewGrid()
	content.SetRows(1, 0)
	content.SetBorders(false)
	content.AddItem(n.textarea, 0, 0, 1, 1, 0, 0, false)
	content.AddItem(n.list, 1, 0, 1, 1, 0, 0, true)

	// add to app
	n.app.SetRoot(content, true)
	n.app.SetFocus(n.list)
}

// show modal help
func (n *nav) showHelp() {
	old := n.app.GetFocus()
	modal := tview.NewModal()
	modal.SetText(help)
	modal.AddButtons([]string{"Quit"})
	modal.SetDoneFunc(func(idx int, label string) {
		if label == "Quit" {
			n.app.Stop()
		}
	})

	// show modal
	n.app.SetRoot(modal, false)
	n.app.SetFocus(modal)
	n.app.Run()

	// reset to old primitive
	n.app.SetRoot(old, true)
	n.app.SetFocus(old)
}

// edit a file
func (n *nav) editFile(path string) {
	cmd := exec.Command(n.editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// start a shell in path
func openShell(path string) {
	if !isDir(path) {
		return
	}
	os.Chdir(path)
	shell, err := exec.LookPath(os.Getenv(shellEnv))
	if err != nil {
		log.Error(err)
		return
	}
	err = syscall.Exec(shell, []string{shell}, os.Environ())
	if err != nil {
		log.Error(err)
	}
}

// open a file/dir
func (n *nav) openPath(base string, sub string) string {
	p := filepath.Join(base, sub)
	if !isDir(p) {
		n.editFile(p)
		return ""
	}
	return p
}

// run app
func (n *nav) runApp(path string) {
	var err error

	// user inputs
	n.app.SetInputCapture(n.eventHandler)
	n.createList(path, nil)

	// navigate
	for {
		entries, err := getFiles(path, n.showHidden)
		if err != nil {
			break
		}
		n.updateList(path, entries)

		err = n.app.Run()
		if err != nil {
			break
		}

		// exit
		if n.exit {
			break
		}

		// show help
		if n.help {
			n.showHelp()
			continue
		}

		// open
		if n.selected {
			n.selected = false
			idx := n.list.GetCurrentItem()
			if idx == 0 {
				// ..
				n.goBack = true
			} else {
				// item selected
				idx--
				if idx < 0 || idx >= len(entries) {
					break
				}
				sub := entries[idx]
				newPath := n.openPath(path, sub.Name())
				if len(newPath) > 0 {
					path = newPath
				}
				continue
			}
		}

		// parent directory
		if n.goBack {
			n.goBack = false
			path, err = filepath.Abs(filepath.Join(path, ".."))
			if err != nil {
				break
			}
			continue
		}
	}

	if err != nil {
		log.Error(err)
	}

	openShell(path)
}

// print usage
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [<options>] [<path>]\n", os.Args[0])
	flag.PrintDefaults()
}

// create the app
func (n *nav) createApp(path string) {
	n.app = tview.NewApplication()
	n.runApp(path)
}

// entry point
func main() {
	var path string
	editorArg := flag.String("editor", "", "File editor")
	hiddenArg := flag.Bool("a", false, "Show hidden files")
	extendedArg := flag.Bool("l", false, "Long format")
	helpArg := flag.Bool("help", false, "Show this help")
	versArg := flag.Bool("version", false, "Show version")
	debugArg := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	if *debugArg {
		log.SetLevel(log.DebugLevel)
		log.Info("debug mode enabled")
	}

	// show help
	if *helpArg {
		usage()
		os.Exit(0)
	}

	// path
	if len(flag.Args()) < 1 {
		path = "."
	} else {
		path = flag.Args()[0]
	}

	// version
	if *versArg {
		fmt.Printf("%s v%s\n", os.Args[0], version)
		os.Exit(0)
	}

	// editor
	var editor string
	if len(*editorArg) < 1 {
		editor = os.Getenv(editorEnv)
	}

	// create struct
	n := nav{
		editor:     editor,
		showHidden: *hiddenArg,
		extended:   *extendedArg,
	}
	n.createApp(path)
}
