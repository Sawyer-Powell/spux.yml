package gen 

import (
	"fmt"
)

type Space struct {
	Space string
	Root string
	Windows []Window
	targetMap map[string]FullTarget
}

func (space *Space) preprocess() {
	space.targetMap = make(map[string]FullTarget)
	// Assign indices to windows and panes and generate targetMap

	for i := 0; i < len(space.Windows); i++ {
		window := &space.Windows[i]
		window.index = i
		space.targetMap[window.Name] = FullTarget{
			session: space.Space,
			window: window.index,
			pane: -1,
		}


		for k := 0; k < len(window.Cmds); k++ {
			cmd := &window.Cmds[k]
			cmd.index = k
		}

		for j := 0; j < len(window.Panes); j++ {
			pane := &window.Panes[j]
			pane.global_index = i + j
			pane.rel_index = j
			space.targetMap[pane.Name] = FullTarget{
				session: space.Space,
				window: window.index,
				pane: pane.rel_index,
			}

			for l := 0; l < len(pane.Cmds); l++ {
				fmt.Println(fmt.Sprint(l))
				cmd := &pane.Cmds[l]
				cmd.index = l
			}
		}
	}
}

type Pane struct {
	Name string
	Init string
	global_index int // Panes index related to the whole session
	rel_index int // Panes index related to the window
	Cmds []Cmd
}

type Window struct {
	Name string
	Panes []Pane
	Layout string
	Cmds []Cmd
	Init string
	index int
}

type Cmd struct {
	Cmd string
	Tgt string
	index int
}

type Target interface {
	String() string
}

type FullTarget struct {
	session string
	window int
	pane int
}
func (t *FullTarget) String() string {
	if t.pane >= 0 {
		return fmt.Sprintf("%s:%d.%d",
			t.session, t.window, t.pane)
	} else {
		ht :=HalfTarget{
			session: t.session,
			window: t.window,
		}
		return ht.String()
	}
}

type HalfTarget struct {
	session string
	window int
}
func (t *HalfTarget) String() string {
	return fmt.Sprintf("%s:%d",
		t.session, t.window)
}
