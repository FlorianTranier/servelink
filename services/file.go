package services

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Node represents a single node in a hierarchical file system structure.
// Name is the name of the file or directory for this node.
// Path is the full path of the file or directory represented by this node.
// IsDir indicates whether this node represents a directory.
// Metadata provides file information for this node.
// Children hold the child nodes, keyed by their respective names.
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Metadata fs.FileInfo
	Children map[string]*Node
}

func ReadFolderV2(baseDir string, parentNode *Node, introspectionLevel int, maxIntrospectionLevel int) *Node {
	if baseDir == "" {
		baseDir = "mnt"
	}
	fileSystem := os.DirFS(baseDir)
	var rootNode *Node
	if parentNode != nil {
		rootNode = parentNode
	} else {
		rootNode = &Node{
			Name:     ".",
			Path:     ".",
			IsDir:    true,
			Metadata: nil,
			Children: make(map[string]*Node),
		}
	}

	rootDir, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, 10)

	for _, entry := range rootDir {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		node := &Node{
			Name:     entry.Name(),
			Path:     filepath.ToSlash(baseDir + "/" + entry.Name()),
			IsDir:    entry.IsDir(),
			Metadata: info,
			Children: make(map[string]*Node),
		}

		mu.Lock()
		rootNode.Children[entry.Name()] = node
		mu.Unlock()

		if entry.IsDir() {
			wg.Add(1)
			go func(parentNode *Node, dirPath string, introspectionLevel int, maxIntrospectionLevel int) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				log.Println("Running with introspection level: " + strconv.Itoa(introspectionLevel))
				log.Println("Max Introspection level: " + strconv.Itoa(maxIntrospectionLevel))
				log.Println("Processing directory: " + dirPath)

				if introspectionLevel <= maxIntrospectionLevel {
					ReadFolderV2(dirPath, parentNode, introspectionLevel, maxIntrospectionLevel)
				}

				//readSubFolder(fileSystem, dirPath, parentNode)
			}(node, node.Path, introspectionLevel+1, maxIntrospectionLevel)
		}
	}

	wg.Wait()
	return rootNode
}
