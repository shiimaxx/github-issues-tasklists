package tasklists_test

import (
	"fmt"
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
				Tasks: []tasklists.Task{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: false},
					{Text: "Draft", Checked: false},
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
				Tasks: []tasklists.Task{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: false},
					{Text: "Draft", Checked: false},
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
				Tasks: []tasklists.Task{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: false},
					{Text: "Draft", Checked: false},
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

			if got, want := len(got.Tasks), len(c.want.Tasks); got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			for i, task := range got.Tasks {
				if got, want := task.Text, c.want.Tasks[i].Text; got != want {
					t.Errorf("got %s, want %s", got, want)
				}
				if got, want := task.Checked, c.want.Tasks[i].Checked; got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		name       string
		inBody     string
		inTasklist tasklists.Tasklist
		want       string
	}{
		{
			name: "toggle checked",
			inBody: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
			inTasklist: tasklists.Tasklist{
				Title: "",
				Tasks: []tasklists.Task{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: true},
					{Text: "Draft", Checked: false},
				},
			},
			want: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [x] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
		},
		{
			name: "add task",
			inBody: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
%s`,
			inTasklist: tasklists.Tasklist{
				Title: "",
				Tasks: []tasklists.Task{
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
					{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: false},
					{Text: "Draft", Checked: false},
					{Text: "New task", Checked: false},
				},
			},
			want: `example description
%s
- [x] https://github.com/shiimaxx/github-tasklist/issues/123
- [ ] https://github.com/shiimaxx/github-tasklist/issues/124
- [ ] Draft
- [ ] New task
%s`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inBody := fmt.Sprintf(c.inBody, tasklists.TasklistPrefix, tasklists.TasklistSuffix)
			if got, want := tasklists.Replace(inBody, c.inTasklist), fmt.Sprintf(c.want, tasklists.TasklistPrefix, tasklists.TasklistSuffix); got != want {
				t.Errorf("got %s, want %s", got, want)
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
		Tasks: []tasklists.Task{
			{Text: "https://github.com/shiimaxx/github-tasklist/issues/123", Checked: true},
			{Text: "https://github.com/shiimaxx/github-tasklist/issues/124", Checked: false},
			{Text: "Draft", Checked: false},
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
