package gen

import (
	"fmt"
)

func makeWindows(space *Space) (string, error) {
	out := ""

	for i := 0; i < len(space.Windows); i++ {
		window, err := makeWindow(space, &space.Windows[i])

		if err != nil {
			return "", err
		}

		out += window
	}

	return out, nil
}

func makeWindow(space *Space, window *Window) (string, error) {
	out := ""

	for i := 0; i < len(window.Cmds); i++ {
		cmd := &window.Cmds[i]
		madeCmd, err := makeCmd(space, cmd, &FullTarget{
			session: space.Space,
			window: window.index,
			pane: -1,
		})

		if err != nil {
			return "", err
		}

		out += madeCmd
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
		pane, err := makePanes(space, window)

		if err != nil {
			return "", err
		}

		out += pane
	} else if window.Init != "" {
		init, err := tmuxExecuteInit(
			&HalfTarget{session: space.Space, window: window.index},
			window.Init,
		) 

		if err != nil {
			return "", err
		}

		out +=  init
	}

	return out, nil
}

func makePanes(space *Space, window *Window) (string, error) {
	out := ""

	for i := 0; i < len(window.Panes); i++ {
		pane, err := makePane(space, window.index, &window.Panes[i])

		if err != nil {
			return "", err
		}

		out += pane
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

	return out, nil 
}

func makePane(space *Space, windowIndex int, pane *Pane) (string, error) {
	out := ""

	for i := 0; i < len(pane.Cmds); i++ {
		cmd := &pane.Cmds[i]

		madeCmd, err := makeCmd(space, cmd, &FullTarget{
			session: space.Space,
			window: windowIndex,
			pane: pane.rel_index,
		})

		if err != nil {
			return "", err
		}

		out += madeCmd
	}
	
	// If this is pane 0 in the window, then we only need to execute
	// this pane's init script
	if pane.rel_index == 0 {
		init, err := tmuxExecuteInit(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index,
			},
			pane.Init,
		)

		if err != nil {
			return "", err
		}

		out += init
	} else {
		// If this pane is not pane 0 in window, then we need to split 
		out += tmuxSplitWindow(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index - 1,
			},
		)

		init, err := tmuxExecuteInit(
			&FullTarget{
				session: space.Space,
				window: windowIndex,
				pane: pane.rel_index,
			},
			pane.Init,
		)

		if err != nil {
			return "", err
		}

		out += init
	}

	return out, nil
}

func makeCmd(space *Space, command *Cmd, source *FullTarget) (string, error) {
	script, err := resolveSpaceScript(command.Cmd) 

	if err != nil {
		return "", err
	}

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
	), nil
}
