//go:build ignore

package main

import (
	"fmt"
	"os"

	tasklists "github.com/shiimaxx/github-issues-tasklists"
)

func main() {
	body := "example description\n"
	body += "\n"
	body += "```[tasklist]]\n"
	body += "- [x] https://github.com/shiimaxx/github-tasklist/issues/123\n"
	body += "- [ ] https://github.com/shiimaxx/github-tasklist/issues/124\n"
	body += "- [ ] Draft\n"
	body += "```"

	fmt.Println(body)
	// Output:
	// example description
	//
	// ```[tasklist]]
	// - [x] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// ```

	tl, err := tasklists.Extract(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(tl.Render())
	// Output:
	// ```[tasklist]
	// - [x] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// ```

	tl.Tasks[0].Checked = false
	fmt.Println(tl.Render())
	// Output:
	// ```[tasklist]
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// ```

	tl.Tasks = append(tl.Tasks, tasklists.Task{Text: "New item", Checked: false})
	fmt.Println(tl.Render())
	// Output:
	// ```[tasklist]
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// - [ ] New item
	// ```

	replacedBody := tasklists.Replace(body, tl)
	fmt.Println(replacedBody)
	// Output:
	// example description
	//
	// ```[tasklist]
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// - [ ] New item
	// ```
}
