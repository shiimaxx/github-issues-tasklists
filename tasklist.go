package tasklists

import (
	"fmt"
	"strings"
)

const (
	gfmCheckboxUnchecked = "- [ ] "
	gfmCheckboxChecked   = "- [x] "
	tasklistPrefix       = "```[tasklist]"
	tasklistSuffix       = "```"
)

type Tasklist struct {
	Title string
	Items []Item
}

type Item struct {
	Text string
	Done bool
}

func Extract(body string) (Tasklist, error) {
	var title string
	var items []Item

	var blankInTasklist bool
	var whatLineInTasklist int

	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, tasklistPrefix) {
			whatLineInTasklist = 1
			continue
		}

		if whatLineInTasklist == 0 {
			continue
		}

		if line == "" {
			blankInTasklist = true
			whatLineInTasklist += 1
			continue
		}

		if blankInTasklist && (line != "" && line != tasklistSuffix) {
			return Tasklist{}, fmt.Errorf("blank line must be only end of tasklist")
		}

		if strings.HasPrefix(line, "### ") {
			if whatLineInTasklist != 1 {
				return Tasklist{}, fmt.Errorf("invalid tasklist format")
			}

			title = strings.TrimPrefix(line, "### ")
			whatLineInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, gfmCheckboxUnchecked) {
			item := Item{Text: strings.TrimPrefix(line, gfmCheckboxUnchecked), Done: false}
			items = append(items, item)
			whatLineInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, gfmCheckboxChecked) {
			item := Item{Text: strings.TrimPrefix(line, gfmCheckboxChecked), Done: true}
			items = append(items, item)
			whatLineInTasklist += 1
			continue
		}

		if strings.HasPrefix(line, tasklistSuffix) {
			break
		}

		return Tasklist{}, fmt.Errorf("invalid tasklist format")
	}

	return Tasklist{Title: title, Items: items}, nil
}

func (t *Tasklist) Render() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "%s\n", tasklistPrefix)
	if t.Title != "" {
		fmt.Fprintf(&builder, "### %s\n", t.Title)
	}
	for _, item := range t.Items {
		if item.Done {
			fmt.Fprintf(&builder, "%s%s\n", gfmCheckboxChecked, item.Text)
		} else {
			fmt.Fprintf(&builder, "%s%s\n", gfmCheckboxUnchecked, item.Text)
		}
	}
	fmt.Fprint(&builder, tasklistSuffix)

	return builder.String()
}
