package services

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Metadata fs.FileInfo
	Children map[string]*Node
}

func ReadFolder() *Node {
	fileSystem := os.DirFS("mnt")
	rootFileInfo, err := fs.Stat(fileSystem, ".")
	if err != nil {
		log.Fatal(err)
	}

	rootNode := &Node{
		Name:     ".",
		Path:     ".",
		IsDir:    true,
		Metadata: rootFileInfo,
		Children: make(map[string]*Node),
	}

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == rootNode.Path || d.IsDir() {
			return nil
		}

		parts := strings.Split(filepath.ToSlash(path), "/")
		currentNode := rootNode
		for i, part := range parts {
			child, exists := currentNode.Children[part]
			if !exists {
				metadata, err := fs.Stat(fileSystem, strings.Join(parts[:i+1], "/"))
				if err != nil {
					return err
				}
				child = &Node{
					Name:     part,
					Path:     strings.Join(parts[:i+1], "/"),
					IsDir:    i < len(parts)-1,
					Metadata: metadata,
					Children: make(map[string]*Node),
				}
				currentNode.Children[part] = child
			}
			currentNode = child
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return rootNode
}
