package gen 

import "strings"

func cleanScriptStr(ss string) string {
	if (len(ss) > 0 && ss[len(ss) - 1] == ';') {
		return ss[:len(ss) - 1]
	}

	return ss
}

func resolveSpaceScript(ss string) string {
	out := ""
	snippet := ss;
	start := strings.IndexRune(snippet, '{')

	for start != -1 {
		end := strings.IndexRune(snippet, '}')
		if end != -1 && start > 0 {
			s := strings.TrimSpace(snippet[:start])
			//fmt.Printf("This is snippet: %s\n", snippet)
			//fmt.Printf("This is s: %s\n", s)

			if s != "" {
				out += `'` + cleanScriptStr(s) + `' `
			}

			snippet = snippet[start:]
		} else if end != -1 {
			switch snippet[start+1:end] {
			case "interrupt":
				out += "C-c "
				snippet = snippet[end+1:]
			case "clear":
				out += "C-l "
				snippet = snippet[end+1:]
			case "enter":
				out += "C-m "
				snippet = snippet[end+1:]
			}
		}
		start = strings.IndexRune(snippet, '{')

	}

	s := cleanScriptStr(strings.TrimSpace(snippet))
	if s != "" {
		out += `'` + s + `'`
	}

	return out
}
