package tasklists

import (
	"fmt"
	"strings"
)

const (
	gfmCheckboxUncheckedPrefix = "- [ ] "
	gfmCheckboxCheckedSuffix   = "- [x] "

	tasklistPrefix = "```[tasklist]"
	tasklistSuffix = "```"
)

type Tasklist struct {
	Title string
	Tasks []Task
}

type Task struct {
	Text    string
	Checked bool
}

func Extract(body string) (Tasklist, error) {
	return extract(body)
}

func extract(body string) (Tasklist, error) {
	lines := strings.Split(body, "\n")

	var i int
	var title string

	// seek to beginning of tasklist
	for i < len(lines) {
		if strings.HasPrefix(lines[i], tasklistPrefix) {
			// peek next line to check if title exists
			if strings.HasPrefix(lines[i+1], "# ") {
				title = strings.TrimPrefix(lines[i+1], "# ")
				i += 1
			} else if strings.HasPrefix(lines[i+1], "## ") {
				title = strings.TrimPrefix(lines[i+1], "## ")
				i += 1
			} else if strings.HasPrefix(lines[i+1], "### ") {
				title = strings.TrimPrefix(lines[i+1], "### ")
				i += 1
			} else if strings.HasPrefix(lines[i+1], "#### ") {
				title = strings.TrimPrefix(lines[i+1], "#### ")
				i += 1
			} else if strings.HasPrefix(lines[i+1], "##### ") {
				title = strings.TrimPrefix(lines[i+1], "##### ")
				i += 1
			} else if strings.HasPrefix(lines[i+1], "###### ") {
				title = strings.TrimPrefix(lines[i+1], "###### ")
				i += 1
			}

			i += 1
			break
		}
		i += 1
	}

	var tasks []Task
	for i < len(lines) {
		if strings.HasPrefix(lines[i], gfmCheckboxUncheckedPrefix) {
			task := Task{Text: strings.TrimPrefix(lines[i], gfmCheckboxUncheckedPrefix), Checked: false}
			tasks = append(tasks, task)
			i += 1
			continue
		}

		if strings.HasPrefix(lines[i], gfmCheckboxCheckedSuffix) {
			task := Task{Text: strings.TrimPrefix(lines[i], gfmCheckboxCheckedSuffix), Checked: true}
			tasks = append(tasks, task)
			i += 1
			continue
		}

		if strings.HasPrefix(lines[i], tasklistSuffix) {
			break
		}

		if lines[i] == "" {
			i += 1
			if lines[i] == tasklistSuffix {
				break
			}
			return Tasklist{}, fmt.Errorf("blank line must be only end of tasklist")
		}

		return Tasklist{}, fmt.Errorf("invalid tasklist format")
	}

	return Tasklist{Title: title, Tasks: tasks}, nil
}

func (t *Tasklist) Render() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "%s\n", tasklistPrefix)
	if t.Title != "" {
		fmt.Fprintf(&builder, "### %s\n", t.Title)
	}
	for _, t := range t.Tasks {
		if t.Checked {
			fmt.Fprintf(&builder, "%s%s\n", gfmCheckboxCheckedSuffix, t.Text)
		} else {
			fmt.Fprintf(&builder, "%s%s\n", gfmCheckboxUncheckedPrefix, t.Text)
		}
	}
	fmt.Fprint(&builder, tasklistSuffix)

	return builder.String()
}
