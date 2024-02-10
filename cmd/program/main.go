package program

import (
	"wakelesstuna/pkg/views"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
)

// You can wire any Bubble Tea model up to the middleware with a function that
// handles the incoming ssh.Session. Here we just grab the terminal info and
// pass it to the new model. You can also return tea.ProgramOptions (such as
// tea.WithAltScreen) on a session by session basis.
func MainHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	_, _, active := s.Pty()

	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	// When running a Bubble Tea app over SSH, you shouldn't use the default
	// lipgloss.NewStyle function.
	// That function will use the color profile from the os.Stdin, which is the
	// server, not the client.
	// We provide a MakeRenderer function in the bubbletea middleware package,
	// so you can easily get the correct renderer for the current session, and
	// use it to create the styles.
	// The recommended way to use these styles is to then pass them down to
	// your Bubble Tea model.
	/*renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8")) */

	/* m := model{
		term:      pty.Term,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		txtStyle:  txtStyle,
		quitStyle: quitStyle,
	} */

	items := []list.Item{
		views.Item("Ramen"),
		views.Item("Tomato Soup"),
		views.Item("Hamburgers"),
		views.Item("Cheeseburgers"),
		views.Item("Currywurst"),
		views.Item("Okonomiyaki"),
		views.Item("Pasta"),
		views.Item("Fillet Mignon"),
		views.Item("Caviar"),
		views.Item("Just Wine"),
	}

	m := views.InitModel("What do you want for dinner?", items, s)

	return m, []tea.ProgramOption{
		tea.WithInput(s),
		tea.WithOutput(s),
	}
}
