package tasklists_test

import (
	"fmt"
	"os"
	"testing"

	tasklists "github.com/shiimaxx/github-issues-tasklists"
)

func TestExtract(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want tasklists.Tasklist
	}{
		{
			name: "tasklist",
			in: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
			want: tasklists.Tasklist{
				Title: "",
				Items: []tasklists.Item{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Done: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Done: false},
					{Text: "Draft", Done: false},
				},
			},
		},
		{
			name: "tasklist with title",
			in: `example description
%s
### asdf
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
			want: tasklists.Tasklist{
				Title: "asdf",
				Items: []tasklists.Item{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Done: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Done: false},
					{Text: "Draft", Done: false},
				},
			},
		},
		{
			name: "tasklist with blank line",
			in: `example description
%s
### asdf
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft

%s`,
			want: tasklists.Tasklist{
				Title: "asdf",
				Items: []tasklists.Item{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Done: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Done: false},
					{Text: "Draft", Done: false},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := tasklists.Extract(fmt.Sprintf(c.in, tasklists.TasklistPrefix, tasklists.TasklistSuffix))
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if got, want := got.Title, c.want.Title; got != want {
				t.Errorf("got %s, want %s", got, want)
			}

			if got, want := len(got.Items), len(c.want.Items); got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			for i, item := range got.Items {
				if got, want := item.Text, c.want.Items[i].Text; got != want {
					t.Errorf("got %s, want %s", got, want)
				}
				if got, want := item.Done, c.want.Items[i].Done; got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestExtract_invalidFormat(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{
			name: "invalide tasklist",
			in: `example description
%s
- x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
		},
		{
			name: "blank line except at the end",
			in: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123

- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := tasklists.Extract(fmt.Sprintf(c.in, tasklists.TasklistPrefix, tasklists.TasklistSuffix))
			if err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestRender(t *testing.T) {
	tl := tasklists.Tasklist{
		Title: "asdf",
		Items: []tasklists.Item{
			{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Done: true},
			{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Done: false},
			{Text: "Draft", Done: false},
		},
	}

	want := fmt.Sprintf(`%s
### asdf
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`, tasklists.TasklistPrefix, tasklists.TasklistSuffix)

	if got := tl.Render(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestReadme(t *testing.T) {
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

	tl.Items[0].Done = false

	fmt.Println(tl.Render())

	tl.Items = append(tl.Items, tasklists.Item{Text: "New item", Done: false})

	fmt.Println(tl.Render())
}
