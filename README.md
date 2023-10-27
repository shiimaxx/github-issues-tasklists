# github-issues-tasklists

Simple library to parse [GitHub Issues tasklists](https://docs.github.com/en/issues/managing-your-tasks-with-tasklists).

## Example

```go
package main

import (
    "fmt"
    "os"

    tasklists "github.com/shiimaxx/github-issues-tasklists"
)

func main() {
	body := "example description\n"
	body += "```[tasklist]"
	body += `
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
`

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

	tl.Items[0].Done = false
	fmt.Println(tl.Render())
	// Output:
	// ```[tasklist]
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// ```

	tl.Items = append(tl.Items, tasklists.Item{Text: "New item", Done: false})
	fmt.Println(tl.Render())
	// Output:
	// ```[tasklist]
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/123
	// - [ ] https://github.com/shiimaxx/github-tasklist/issues/124
	// - [ ] Draft
	// - [ ] New item
	// ```
}
```
