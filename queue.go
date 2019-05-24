package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
	"golang.org/x/tools/imports"
)

type Queue []QueueAction

func NewQueue(actions ...QueueAction) Queue {
	return Queue(actions)
}

func (q Queue) Description() string {
	t := treeprint.New()
	t = t.AddBranch("Actions to be committed")
	for i, action := range q {
		t.AddMetaNode(fmt.Sprintf("%d", i+1), action.Comment())
	}
	return t.String()
}

func (q Queue) Exec() bool {
	qq := make(Queue, 0, len(q))
	revertPrevious := func(qq Queue) {
		for i := len(qq) - 1; i >= 0; i-- {
			if err := qq[i].Revert(); err != nil {
				log.Printf("Revert Action#%d failed: %v", i+1, err)
			} else {
				log.Printf("Reverted Action#%d", i+1)
			}
		}
	}
	for i, action := range q {
		log.Printf("Action#%d: %s", i+1, action.Comment())
		f, err := action.Run()
		if err != nil {
			log.Printf("Action#%d error: %v", i+1, err)
			revertPrevious(qq)
			return false
		}
		qq = append(qq, action)
		if err := action.Finalize(f); err != nil {
			log.Printf("Finalizer#%d error: %v", i+1, err)
			revertPrevious(qq)
			return false
		}
	}
	return true
}

type QueueAction interface {
	Run() (*os.File, error)
	Comment() string
	Finalize(f *os.File) error
	Revert() error
}

func CheckDirAction(path string) QueueAction {
	return &queueAction{
		action: func() (*os.File, error) {
			info, err := os.Stat(path)
			if err != nil {
				return nil, err
			}
			if !info.IsDir() {
				return nil, errors.New(path + " is not a dir")
			}
			return nil, nil
		},
		comment: fmt.Sprintf("dir %s must exist", projectPath(path)),
	}
}

func NewDirAction(path string) QueueAction {
	return &queueAction{
		action: func() (*os.File, error) {
			err := os.MkdirAll(path, 0755)
			return nil, err
		},
		comment: fmt.Sprintf("new dir %s if not exists", projectPath(path)),
		revert: func() error {
			return os.Remove(path)
		},
	}
}

func CreateNewFileAction(path string, contents []byte) QueueAction {
	return &queueAction{
		action: func() (f *os.File, err error) {
			return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		},
		comment: fmt.Sprintf("new file %s with %d lines of content (no overwrite)",
			projectPath(path), lineCount(contents)),
		finalize: func(f *os.File) error {
			if f == nil {
				return nil
			}
			defer f.Close()
			return flushBufferToFile(contents, f, filepath.Ext(path) == ".go")
		},
		revert: func() error {
			return os.Remove(path)
		},
	}
}

func OverwriteFileAction(path string, contents []byte) QueueAction {
	return &queueAction{
		action: func() (f *os.File, err error) {
			return os.Create(path)
		},
		comment: fmt.Sprintf("overwrite file %s with %d lines of content",
			projectPath(path), lineCount(contents)),
		finalize: func(f *os.File) error {
			if f == nil {
				return nil
			}
			defer f.Close()
			return flushBufferToFile(contents, f, filepath.Ext(path) == ".go")
		},
		revert: func() error {
			return os.Remove(path)
		},
	}
}

type queueAction struct {
	action   func() (*os.File, error)
	comment  string
	finalize func(f *os.File) error
	revert   func() error
}

func (q *queueAction) Run() (*os.File, error) {
	if q.action != nil {
		return q.action()
	}
	return nil, nil
}

func (q *queueAction) Comment() string {
	return q.comment
}

func (q *queueAction) Finalize(f *os.File) error {
	if q.finalize != nil {
		return q.finalize(f)
	}
	return nil
}

func (q *queueAction) Revert() error {
	if q.revert != nil {
		return q.revert()
	}
	return nil
}

func lineCount(contents []byte) int {
	var lines int
	s := bufio.NewScanner(bytes.NewReader(contents))
	for s.Scan() {
		lines++
	}
	return lines
}

func projectPath(path string) string {
	return filepath.Join("[project]", strings.TrimPrefix(path, *projectDir))
}

func flushBufferToFile(buf []byte, f *os.File, fmt bool) error {
	if fmt {
		if fmtBuf, err := imports.Process(f.Name(), buf, nil); err == nil {
			_, err = f.Write(fmtBuf)
			return err
		} else {
			log.Printf("Warning: cannot gofmt %s: %s\n", f.Name(), err.Error())
			f.Write(buf)
			return nil
		}
	}
	_, err := f.Write(buf)
	return err
}
