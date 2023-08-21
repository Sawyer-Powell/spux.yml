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

func (space *Space) GenerateScript() string {
	space.preprocess()

	out := "#!/bin/bash\n"
	out += fmt.Sprintf(`cd %s`, space.Root) + "\n"
	out += tmuxNewSession(space.Space)
	out += makeWindows(space)
	//out += tmuxAttachSession(space.Space)

	return out
}
