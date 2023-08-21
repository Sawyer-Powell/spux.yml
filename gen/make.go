package gen 

import "fmt"

func makeWindows(space *Space) string {
	out := ""

	for i := 0; i < len(space.Windows); i++ {
		out += makeWindow(space, &space.Windows[i])
	}

	return out
}

func makeWindow(space *Space, window *Window) string {
	out := ""

	for i := 0; i < len(window.Cmds); i++ {
		cmd := &window.Cmds[i]
		out += makeCmd(space, cmd, &FullTarget{
			session: space.Space,
			window: window.index,
			pane: -1,
		})
	}

	// Name/Make the window
	if window.index == 0 {
		out += tmuxRenameWindow(
			&HalfTarget{
				session: space.Space,
				window: window.index,
			},
			window.Name,
		)
	} else {
		out += tmuxNewWindow(
			&HalfTarget{
				session: space.Space,
				window: window.index,
			},
			window.Name,
		)
	}
	
	if len(window.Panes) > 0 {
		out += makePanes(space, window)
	} else if window.Init != "" {
		out += tmuxExecuteInit(
			&HalfTarget{session: space.Space, window: window.index},
			window.Init,
		)
	}

	return out
}

func makePanes(space *Space, window *Window) string {
	out := ""

	for i := 0; i < len(window.Panes); i++ {
		out += makePane(space, window.index, &window.Panes[i])
	}

	if (window.Layout != "") {
		out += tmuxSelectLayout(
			&HalfTarget{
				session: space.Space,
				window: window.index,
			},
			window.Layout,
		)
	}

	return out
}

func makePane(space *Space, windowIndex int, pane *Pane) string {
	out := ""

	for i := 0; i < len(pane.Cmds); i++ {
		cmd := &pane.Cmds[i]
		out += makeCmd(space, cmd, &FullTarget{
			session: space.Space,
			window: windowIndex,
			pane: pane.rel_index,
		})
	}
	
	// If this is pane 0 in the window, then we only need to execute
	// this pane's init script
	if pane.rel_index == 0 {
		out += tmuxExecuteInit(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index,
			},
			pane.Init,
		)
	} else {
		// If this pane is not pane 0 in window, then we need to split 
		out += tmuxSplitWindow(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index - 1,
			},
		)

		out += tmuxExecuteInit(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index,
			},
			pane.Init,
		)
	}

	return out
}

func makeCmd(space *Space, command *Cmd, source *FullTarget) string {
	script := resolveSpaceScript(command.Cmd) 
	target := space.targetMap[command.Tgt]

	spaceCmd := fmt.Sprintf(
		`tmux send-keys -t %s %s`,
		target.String(), cleanScriptStr(script),
	) + ";"

	ht := &HalfTarget{
		session: target.session,
		window: target.window,
	}
	spaceCmd += fmt.Sprintf(
		`tmux select-window -t %s;`,
		ht.String(),
	)

	if (target.pane >= 0) {
		spaceCmd += fmt.Sprintf(
			`tmux select-pane -t %s`,
			target.String(),
		)
	}

	return tmuxSetEnv(
		source,
		"SPUX_CMD_" + fmt.Sprint(command.index), 
		`"` + spaceCmd + `"`,
	)
}
