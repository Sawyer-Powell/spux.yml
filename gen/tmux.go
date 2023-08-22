package gen 

import "fmt"

// For converting the init parameters of panes and windows into tmux
// commands
func tmuxExecuteInit(target Target, init string) (string, error) {
	safe_init, err := resolveSpaceScript(init)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`tmux send-keys -t %s C-l %s`,
		target.String(),
		safe_init,
	) + "\n", nil
}

func tmuxSplitWindow(target Target) string {
	return fmt.Sprintf(`tmux split-window -t %s`, target.String()) + "\n"
}

func tmuxRenameWindow(target *HalfTarget, name string) string {
	return fmt.Sprintf(`tmux rename-window -t %s %s`,
		target.String(), name) + "\n"
}

func tmuxNewWindow(target *HalfTarget, name string) string {
	return fmt.Sprintf(`tmux new-window -t %s -n %s`,
		target.String(), name) + "\n"
}

func tmuxSelectLayout(target *HalfTarget, layout string) string {
	return fmt.Sprintf(`tmux select-layout -t %s %s`,
		target.String(), layout) + "\n"
}

func tmuxNewSession(session string) string {
	return fmt.Sprintf(`tmux new-session -d -s %s`, session) + "\n"
}

func tmuxAttachSession(session string) string {
	return fmt.Sprintf(`tmux attach-session -t %s`, session) + "\n"
}

func tmuxSetEnv(source *FullTarget, name string, value string) string {
	return fmt.Sprintf(`tmux setenv -t %s %s %s`,
		source.String(), name, value) + "\n"
}
