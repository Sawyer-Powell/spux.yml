package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"spux/gen"
	"strings"

	gotmux "github.com/jubnzv/go-tmux"
	"github.com/spf13/cobra"
)

var (
	scriptFlag = false // indicating whether to send script to stdout
	recreateFlag = false // if not in tmux session, and session exists, force restart and attach to session
	listFlag = false // whether to list available spaces
	rootFlag = false // whether to send root dir of space to stdout
	killFlag = false // whether to kill the space tmux session

	ps = string(os.PathSeparator)

	Root = cobra.Command{
		Use: "spux [path/to/space.yml]",
		Short: "Spaces for tmux",
		Long: `A friendly and powerful companion for organizing, defining, 
and configuring tmux based development environments`,
		Args: cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
		Run: spux,
	}
)

func init() {
	Root.Flags().BoolVarP(
		&scriptFlag,
		"script",
		"s",
		false,
		"Sends generated tmux script to stdout, blocks activating the space",
	)
	Root.Flags().BoolVarP(
		&recreateFlag,
		"recreate",
		"r",
		false,
		"Kills any tmux session matching the space name (inferred or provided) and then recreates the space",
	)
	Root.Flags().BoolVarP(
		&listFlag,
		"list",
		"l",
		false,
		"Lists the spaces already saved in spux, blocks creating the space.",
	)
	Root.Flags().BoolVarP(
		&killFlag,
		"kill",
		"k",
		false,
		"Kills the space",
	)
}

func spux(cmd *cobra.Command, args []string) {
	if (listFlag) {
		fmt.Printf("Spaces saved in spux:\n%s", listSavedSpaces())
		return
	}

	if len(args) > 0 {
		first_arg := args[0]
		isYml := strings.Contains(first_arg, ".yml") ||
			strings.Contains(first_arg, ".yaml")
		if isYml {
			err := handleYmlArg(first_arg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			handleSpaceArg(first_arg)
		}
	} else {
		err := handleSpuxYml()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func getBinPath() string {
	home := os.Getenv("HOME")

	if home == "" {
		fmt.Printf("$HOME is either unset or empty\n")
		return "" 
	}

	spux := home + ps + ".config" + ps + "spux"
	bin := spux + ps + "bin"

	err := os.MkdirAll(bin, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to find/create directory %s: err=%v\n", bin, err)
		return ""
	}

	return bin
}

func createScript(space *gen.Space) (string, error) {
	script, err := space.GenerateScript();

	if err != nil {
		return "", err
	}

	scriptPath := getBinPath() + ps + space.Space

	err = os.WriteFile(
		scriptPath,
		[]byte(script),
		0755,
	)

	if err != nil {
		fmt.Printf("failed to write to %s: err=%v\n", scriptPath, err)
		return "", err
	}

	return scriptPath, nil
}

func executeScript(scriptPath string) {
	scriptExec := exec.Command("bash", scriptPath)
	_, err := scriptExec.CombinedOutput()

	if err != nil {
		fmt.Printf("execution of %s failed: err=%v\n", scriptPath, err)
	}
}

func handleYmlArg(filename string) error {
	space := gen.ParseYml(filename)
	script, err := createScript(space)

	if err != nil {
		return err
	}

	safeExecuteScript(script, space.Space)
	return nil
}

func getFilesFromDir(dir string) []string{
	var files []string

	inDir, err := os.ReadDir(dir)

	if err != nil {
		fmt.Printf("could not read files from %s, err=%v\n", dir, err)
		return []string{}
	}

	for _, entry := range inDir {
		if !entry.IsDir() {
			file := dir + ps + entry.Name()
			files = append(files, file)
		}
	}
	
	return files
}

func listSavedSpaces() string {
	out := ""

	bin := getBinPath()
	files := getFilesFromDir(bin)

	for _, file := range files {
		out += "- '" + filepath.Base(file) + "'\n"
	}

	return out
}

func handleSpaceArg(spaceName string) {
	bin := getBinPath()

	files := getFilesFromDir(bin)

	var scriptPath string = ""

	for _, file := range files {
		if filepath.Base(file) == spaceName {
			scriptPath = file
		}
	}

	if scriptPath == "" {
		fmt.Printf("could not find file matching %s in %s\n", spaceName, bin)
		return
	}

	safeExecuteScript(scriptPath, spaceName)
}

func handleSpuxYml() error {
	cwd, err := os.Getwd()

	if err != nil {
		return errors.New(fmt.Sprintf("could not get current directory err=%v\n", err))
	}

	files := getFilesFromDir(cwd)
	spuxYml := "spux.yml"
	ymlPath := ""

	for _, file := range files {
		if filepath.Base(file) == spuxYml {
			ymlPath = file
		}
	}

	if ymlPath == "" {
		fmt.Printf("could not find file matching \"%s\" in %s\n", spuxYml, cwd)
		return nil
	}

	space := gen.ParseYml(ymlPath)
	script, err := createScript(space)

	if err != nil {
		return err
	}

	safeExecuteScript(script, space.Space)
	return nil
}

func safeExecuteScript(script string, spaceName string) {
	if rootFlag {
		return
	}

	if scriptFlag {
		fmt.Println(getFileContents(script))
		return
	}

	if killFlag {
		killSession(spaceName)
		return
	}

	if !gotmux.IsInsideTmux() {
		if !spaceAlreadyRunning(spaceName) {
			executeScript(script)
		} else if(recreateFlag) {
			killSession(spaceName)
			executeScript(script)
		}

		session := gotmux.Session{
			Name: spaceName,
		}
		session.AttachSession()
	} else {
		fmt.Printf("already inside a tmux session, please detach and run spux again\n")
	}
}

func killSession(session string) {
	cmd := exec.Command("tmux", "kill-session", "-t", session)
	_, err := cmd.Output()
	
	if err != nil {
		fmt.Printf("error running \"tmux kill-session -t %s\" err=%v\n", session, err)
	}
}

func getFileContents(file string) string {
	cmd := exec.Command("cat", file)
	fileContents, err := cmd.Output()

	if err != nil {
		fmt.Printf("error running \"cat %s\" err=%v\n", file, err)
	}

	return string(fileContents)
}

func spaceAlreadyRunning(spaceName string) bool {
	tmuxLs := exec.Command("tmux", "ls", "-F", "#S")

	outBytes, err := tmuxLs.Output()

	if err != nil {
		//fmt.Printf("error running \"tmux ls\" err=%v\n", err)
		return false
	}

	lines := strings.Split(string(outBytes), "\n")

	for _, line := range lines {
		if line == spaceName {
			return true
		}
	}
	
	return false
}

