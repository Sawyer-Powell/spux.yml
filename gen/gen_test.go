package gen

import (
	"testing"
) 

func TestTestData(t *testing.T) {
	space := ParseYml("../test_data")
	t.Logf("%s successfully interpreted\n", space.Space)

	script, err := space.GenerateScript()

	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Logf("script:\n%s\n", script)
}
