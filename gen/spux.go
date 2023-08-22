package gen 

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYml(filename string) *Space {
	space := &Space{}

	fileContents, err := os.ReadFile(filename)

	if (err != nil) {
		log.Fatalf("Error reading file: %v", err)
	}

	err = yaml.Unmarshal(fileContents, space)

	if (err != nil) {
		log.Fatalf("Error reading file: %v", err)
	}

	return space
}

func (space *Space) GenerateScript() (string, error) {
	space.preprocess()

	out := "#!/bin/bash\n"
	out += fmt.Sprintf(`cd %s`, space.Root) + "\n"
	out += tmuxNewSession(space.Space)
	windows, err := makeWindows(space)

	if err != nil {
		return "", err
	}

	out += windows

	return out, nil
}
