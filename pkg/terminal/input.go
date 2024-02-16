package terminal

import (
	"bytes"
	"fmt"
	"wakelesstuna/pkg/ui/list"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

// We need to have a wrapper for our bubbles as they don't currently implement the tea.Model interface
// textinput, textarea

type shortAnswerField struct {
	textinput textinput.Model
}

func NewShortAnswerField() *shortAnswerField {
	a := shortAnswerField{}

	model := textinput.New()
	model.Placeholder = "Your answer here"
	model.Focus()

	a.textinput = model
	return &a
}

func (a *shortAnswerField) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *shortAnswerField) Init() tea.Cmd {
	return nil
}

func (a *shortAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *shortAnswerField) View() string {
	return a.textinput.View()
}

func (a *shortAnswerField) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *shortAnswerField) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *shortAnswerField) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *shortAnswerField) Value() string {
	return a.textinput.Value()
}

type longAnswerField struct {
	textarea textarea.Model
}

func NewLongAnswerField() *longAnswerField {
	a := longAnswerField{}

	model := textarea.New()
	model.Placeholder = "Your answer here"
	model.Focus()

	a.textarea = model
	return &a
}

func (a *longAnswerField) Blink() tea.Msg {
	return textarea.Blink()
}

func (a *longAnswerField) Init() tea.Cmd {
	return nil
}

func (a *longAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textarea, cmd = a.textarea.Update(msg)
	return a, cmd
}

func (a *longAnswerField) View() string {
	return a.textarea.View()
}

func (a *longAnswerField) Focus() tea.Cmd {
	return a.textarea.Focus()
}

func (a *longAnswerField) SetValue(s string) {
	a.textarea.SetValue(s)
}

func (a *longAnswerField) Blur() tea.Msg {
	return a.textarea.Blur
}

func (a *longAnswerField) Value() string {
	return a.textarea.Value()
}

type listAnswerField struct {
	list list.Model
}

func NewListAnswerField() *listAnswerField {
	a := listAnswerField{}

	model := list.InitialModel()
	a.list = model
	return &a
}

func (a *listAnswerField) Blink() tea.Msg {
	return textarea.Blink()
}

func (a *listAnswerField) Init() tea.Cmd {
	return nil
}

func (a *listAnswerField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model
	model, cmd = a.list.Update(msg)
	a.list, _ = model.(list.Model) // Perform type assertion here
	return a, cmd
}

func (a *listAnswerField) View() string {
	return a.list.View()
}

func (a *listAnswerField) Focus() tea.Cmd {
	return nil
	// return a.list.Focus()
}

func (a *listAnswerField) SetValue(s string) {
	// do nothing
	// a.list.SetValue(s)
}

func (a *listAnswerField) Blur() tea.Msg {
	return nil
	//return a.list.Blur()
}

func (a *listAnswerField) Value() string {
	fmt.Println()
	return createValuePairs(a.list.Selected())
	//return a.list.Value()
}

func createValuePairs(m map[int]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%d=\"%s\"\n", key, value)
	}
	return b.String()
}
