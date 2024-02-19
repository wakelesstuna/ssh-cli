package terminal

import (
	"fmt"
	"log"
	"wakelesstuna/pkg/git"
	"wakelesstuna/pkg/github"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Main struct {
	styles    *Styles
	index     int
	questions []Question
	width     int
	height    int
	done      bool
	gitUrl    string
}

type Question struct {
	question string
	answer   string
	input    Input
}

func newQuestion(q string) Question {
	return Question{question: q}
}

func newShortQuestion(q string) Question {
	question := newQuestion(q)
	model := NewShortAnswerField()
	question.input = model
	return question
}

func newLongQuestion(q string) Question {
	question := newQuestion(q)
	model := NewLongAnswerField()
	question.input = model
	return question
}

func newListChoice() Question {
	question := newQuestion("")
	model := NewListAnswerField()
	question.input = model
	return question
}

func New(questions []Question) *Main {
	styles := DefaultStyles()
	return &Main{styles: styles, questions: questions}
}

func (m Main) Init() tea.Cmd {
	return m.questions[m.index].input.Blink
}

func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.questions[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
				m.gitUrl = DoStuff(m)
			}
			current.answer = current.input.Value()
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

func (m Main) View() string {
	current := m.questions[m.index]
	if m.done {
		return fmt.Sprintf("Repo created! @ --> %s\n", m.gitUrl)
	}
	if m.width == 0 {
		return "loading..."
	}
	// stack some left-aligned strings together in the center of the window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
		),
	)
}

func DoStuff(m Main) string {
	projectDir := "C://Users/oscfor/gc/fun/ssh-cli/temp-project"
	gitUrl, err := github.CreateRepository(m.questions[0].answer, m.questions[1].answer)
	if err != nil {
		log.Fatalf("Error creating git repo: %s", err)
	}

	err = git.GitAddCommitPush(gitUrl, projectDir)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return gitUrl
}

func (m *Main) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func InitTerminalWizard(s ssh.Session) *tea.Program { //(tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	_, _, active := s.Pty()

	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		tea.Quit()
	}

	questions := []Question{
		newShortQuestion("Enter the name of your new github repo?"),
		newShortQuestion("Add an description?"),
		//newListChoice(),
		//newLongQuestion("what's your favourite quote?"),
	}
	main := New(questions)

	return tea.NewProgram(*main, []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
		tea.WithAltScreen(),
	}...)
}
