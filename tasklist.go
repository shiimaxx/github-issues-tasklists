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
	var title string
	var tasks []Task

	var blankInTasklist bool
	var posInTasklist int

	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, tasklistPrefix) {
			posInTasklist = 1
			continue
		}

		if posInTasklist == 0 {
			continue
		}

		if line == "" {
			blankInTasklist = true
			posInTasklist += 1
			continue
		}

		if blankInTasklist && (line != "" && line != tasklistSuffix) {
			return Tasklist{}, fmt.Errorf("blank line must be only end of tasklist")
		}

		if strings.HasPrefix(line, "### ") {
			if posInTasklist != 1 {
				return Tasklist{}, fmt.Errorf("invalid tasklist format")
			}

			title = strings.TrimPrefix(line, "### ")
			posInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, gfmCheckboxUncheckedPrefix) {
			task := Task{Text: strings.TrimPrefix(line, gfmCheckboxUncheckedPrefix), Checked: false}
			tasks = append(tasks, task)
			posInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, gfmCheckboxCheckedSuffix) {
			task := Task{Text: strings.TrimPrefix(line, gfmCheckboxCheckedSuffix), Checked: true}
			tasks = append(tasks, task)
			posInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, tasklistSuffix) {
			break
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
