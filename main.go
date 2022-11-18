package main

/*
For the Vue.js .env handling:
 see: https://cli.vuejs.org/guide/mode-and-env.html#example-staging-mode
*/

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type replacement struct {
	odjVariable string
	replacement string
}

var (
	odjVariables = []replacement{}
	root         = "."
	env          = ""
)

// For the beginning we just brake at any error
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	getRootFolder()
	getEnvFile()
	getVariables()
	iterate()
}

// the root folder of the distribution must be tolled as first argument
func getRootFolder() {
	//check input parameter
	if len(os.Args) < 2 || os.Args[1] == "" {
		panic("Please define path to your distribution folder as first argument")
	}
	root = os.Args[1]

	//check if root is a folder
	if fileInfo, ok := os.Stat(root); ok != nil {
		fmt.Println("It was not possible to check the root folder of your distribution")
		panic(ok)
	} else if !fileInfo.IsDir() {
		panic("The first argument must point to a directory with your distribution in it but it isn't")
	}

}

// the location of the .env.odj file must be tolled as second argument
func getEnvFile() {
	//check input parameter
	if len(os.Args) < 3 || os.Args[2] == "" {
		panic("Please define the path to your .env.odj file as second argument")
	}
	env = os.Args[2]

	//check if env is a file
	if fileInfo, ok := os.Stat(env); ok != nil {
		fmt.Println("It was not possible to check the .env.odj file")
		panic(ok)
	} else if fileInfo.IsDir() {
		panic("The second argument must point to your .env.odj file but it isn't")
	}

}

// will read .env.odj and extract all replacements
func getVariables() {
	//check if .env.odj exists
	if fileInfo, ok := os.Stat(env); ok != nil {
		fmt.Println("It was not possible to check the .env.odj file")
		panic(ok)
	} else if fileInfo.IsDir() {
		panic("There must be a file '.env.odj' in this folder but isn't")
	}

	file, err := os.Open(env)
	handleError(err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isVar := strings.Contains(line, "ODJ_"); isVar {
			elements := strings.Split(line, "=")
			odjVariables = append(odjVariables, replacement{strings.TrimSpace(elements[1]), os.Getenv(strings.TrimSpace(elements[1]))})
		}

	}

	for _, rep := range odjVariables {
		fmt.Printf("Will replace '%s' with '%s'\n", rep.odjVariable, rep.replacement)
	}

	handleError(scanner.Err())

	defer handleError(file.Close())
}

// iterate over the file system recursively
func iterate() {
	fmt.Printf("Start to iterate over the file system. Starting at: %s\n", root)
	err := filepath.WalkDir(root,
		func(path string, info os.DirEntry, err error) error {
			handleError(err)

			if info.IsDir() {
				fmt.Printf("+ Dir: %s\n", info.Name())
			} else {
				fileInfo, _ := info.Info()
				selectRelevantFiles(path, fileInfo)
			}

			return nil
		})
	handleError(err)
}

// Filter JS files and ignore the rest
func selectRelevantFiles(path string, info os.FileInfo) {
	if isJs := strings.HasSuffix(strings.ToLower(info.Name()), ".js"); isJs {
		replaceOdjVariables(path, info)
	} else {
		fmt.Printf("\t--> Ignoring no JS file %s\n", info.Name())
	}

}

// will replace any ODJ variable if there is any
func replaceOdjVariables(path string, info os.FileInfo) {
	text := getFileContent(path)
	found := false
	for _, replace := range odjVariables {
		// it is not necessary to check if the dedicated ODJ variable exists
		// but we wont to save the file just once and just if there was something changed
		if strings.Contains(*text, replace.odjVariable) {
			fmt.Printf("\t--> File %s contains ODJ variable '%s' and will be replaced with '%s'\n", info.Name(), replace.odjVariable, replace.replacement)
			*text = strings.ReplaceAll(*text, replace.odjVariable, replace.replacement)
			saveFile(path, text)
			found = true
		}
	}
	if !found {
		fmt.Printf("\t--> File %s doesn't contain any ODJ variable \n", info.Name())
	}
}

func getFileContent(path string) *string {
	dat, ok := os.ReadFile(path)
	handleError(ok)
	text := string(dat)
	return &text

}

func saveFile(path string, text *string) {
	file, err := os.Create(path)
	handleError(err)
	_, err = file.WriteString(*text)
	handleError((err))
	defer handleError(file.Close())
}
