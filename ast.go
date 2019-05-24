package main

import (
	"errors"
	"fmt"
)

func NewMethodsCollection(ifaceName string, srcDir string) (*MethodsCollection, error) {
	path, id, err := findInterface(ifaceName, srcDir)
	if err != nil {
		return nil, err
	}
	m := &MethodsCollection{
		Path: path,
		ID:   id,
	}
	methods, srcPath, err := methodsOf(path, id, ifaceName, srcDir)
	if err != nil {
		return nil, err
	}
	m.Methods = methods
	m.SrcPath = srcPath
	return m, nil
}

type MethodsCollection struct {
	Path    string
	SrcPath string
	ID      string
	Methods []Method
}

var ErrStopRange = errors.New("stop range")

func (m *MethodsCollection) ForEachMethod(fn func(m *Method) error) error {
	for _, mm := range m.Methods {
		if err := fn(&mm); err == ErrStopRange {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (m *MethodsCollection) String() string {
	return fmt.Sprintf("(%s) %s: %d methods", m.Path, m.ID, len(m.Methods))
}
