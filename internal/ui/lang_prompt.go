package ui

import (
	"fmt"
	"strings"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func SelectLanguage(ls []core.Language) (core.Language, error) {
	searcher := func(input string, index int) bool {
		l := ls[index]
		name := strings.Replace(strings.ToLower(l.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}:",
		Active:   fmt.Sprintf("%s {{ .Name | cyan | underline }}", promptui.IconSelect),
		Inactive: "  {{ .Name }}",
		Selected: `{{ "Select language:" | faint}} {{ .Name }}`,
	}

	langPrompt := promptui.Select{
		Label:     "Select language",
		Items:     ls,
		Searcher:  searcher,
		Templates: templates,
	}

	i, _, err := langPrompt.Run()
	cobra.CheckErr(err)

	return ls[i], nil
}
