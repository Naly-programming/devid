package review

import (
	"fmt"
	"strings"

	devsync "github.com/Naly-programming/devid/internal/sync"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Decision represents the user's choice for a candidate.
type Decision int

const (
	DecisionPending Decision = iota
	DecisionApprove
	DecisionReject
	DecisionSkip
)

// Result contains the decisions made during review.
type Result struct {
	Decisions []Decision
}

var (
	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	addStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	removeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	hintStyle    = lipgloss.NewStyle().Faint(true)
	sourceStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	candidates []devsync.Candidate
	index      int
	decisions  []Decision
	quitting   bool
}

func initialModel(candidates []devsync.Candidate) model {
	decisions := make([]Decision, len(candidates))
	return model{
		candidates: candidates,
		decisions:  decisions,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.decisions[m.index] = DecisionApprove
			return m.advance()
		case "r":
			m.decisions[m.index] = DecisionReject
			return m.advance()
		case "s":
			m.decisions[m.index] = DecisionSkip
			return m.advance()
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) advance() (tea.Model, tea.Cmd) {
	if m.index >= len(m.candidates)-1 {
		return m, tea.Quit
	}
	m.index++
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	c := m.candidates[m.index]
	header := fmt.Sprintf("devid review - candidate %d of %d", m.index+1, len(m.candidates))
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(sourceStyle.Render(fmt.Sprintf("Source: %s  Time: %s", c.Source, c.Timestamp.Format("2006-01-02 15:04"))))
	b.WriteString("\n\n")

	// Render diff with colours
	if c.Diff != "" {
		for _, line := range strings.Split(c.Diff, "\n") {
			if strings.HasPrefix(line, "+ ") {
				b.WriteString(addStyle.Render(line))
			} else if strings.HasPrefix(line, "- ") {
				b.WriteString(removeStyle.Render(line))
			} else {
				b.WriteString(line)
			}
			b.WriteString("\n")
		}
	} else {
		b.WriteString("(no diff available)\n")
	}

	b.WriteString("\n")
	b.WriteString(hintStyle.Render("[a]pprove  [r]eject  [s]kip  [q]uit"))
	b.WriteString("\n")

	return b.String()
}

// Run starts the interactive review TUI and returns the decisions.
func Run(candidates []devsync.Candidate) (Result, error) {
	if len(candidates) == 0 {
		return Result{}, nil
	}

	m := initialModel(candidates)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return Result{}, err
	}

	fm := final.(model)
	return Result{Decisions: fm.decisions}, nil
}
