package loader

import (
	"fmt"
	"graph"
	"os"
	"path/filepath"
	"plugin"
)

func LookupPlugin(filename string) (string, bool) {
	candidates := make([]string, 0)
	execpath, err := os.Executable()
	execdir := filepath.Dir(execpath)
	if err == nil {
		candidates = append(candidates, filepath.Join(execdir, filename))
		candidates = append(candidates, filepath.Join(execdir, "plugin", filename))
	}
	candidates = append(candidates, filepath.Join(".", filename))
	candidates = append(candidates, filepath.Join(".", "plugin", filename))
	candidates = append(candidates, filepath.Join(os.Getenv("GOPATH"), "bin", "plugin", filename))

	for _, candidate := range candidates {
		_, err = os.Stat(candidate)
		if err == nil {
			return candidate, true
		}
	}

	return "", false
}

func LoadGraph(config string) (graph.Graph, bool) {
	// TU
	filename_T, exist_T := LookupPlugin("../bin/plugin.so")
	filename, exist := LookupPlugin("plugin.so")

	if exist_T {
		filename = filename_T
	}

	if !exist && !exist_T {
		fmt.Println("Failed to find plugin: plugin.so or ../bin/plugin.so")
		return nil, false
	}
	// TU

	// filename, exist := LookupPlugin("plugin.so")
	// if !exist {
	// 	fmt.Println("Failed to find plugin: plugin.so")
	// 	return nil, false
	// }

	p, err := plugin.Open(filename)
	if err != nil {
		fmt.Printf("Failed to load plugin: %s\n", filename)
		return nil, false
	}

	loadGraph, err := p.Lookup("LoadGraph")
	if err != nil {
		fmt.Println("Failed to lookup symbol: LoadGraph()")
		return nil, false
	}

	loadGraphFunc, ok := loadGraph.(func(string) graph.Graph)
	if !ok {
		fmt.Println("Failed for type assertion: graph.Graph")
		return nil, false
	}

	accessor := loadGraphFunc(config)

	return accessor, true
}
