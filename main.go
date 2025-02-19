package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lb "github.com/crolbar/lipbalm"
	lbfb "github.com/crolbar/lipbalm/framebuffer"
	lbl "github.com/crolbar/lipbalm/layout"
)

var l lbl.Layout = lbl.DefaultLayout()

type model struct {
	swStartTime time.Time
	swStopTime  string

	width  int
	height int
}

type TickMsg struct{}

func Tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(50 * time.Millisecond)
		return TickMsg{}
	}
}

func (m model) stopwatchGetElapsed() string {
	if !m.stopwatchIsRunning() {
		if m.swStopTime != "" {
			return m.swStopTime
		}
		return digits[0]
	}

	elapsed := fmt.Sprintf("%.2f", time.Since(m.swStartTime).Seconds())
	return getAsAscii(elapsed)
}

func (m model) stopwatchIsRunning() bool {
	return m.swStartTime.Second() != 0
}

func (m *model) stopwatchStopTime() {
	m.swStopTime = m.stopwatchGetElapsed()
	m.swStartTime = time.Time{}
}

func (m *model) stopwatchResetTime() {
	m.swStopTime = ""
}

func (m *model) stopwatchStart() {
	m.swStartTime = time.Now()
}

func new_model() model {
	return model{
		swStartTime: time.Time{},
		swStopTime:  "",
	}
}

func main() {
	p := tea.NewProgram(model(new_model()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case " ":
			if m.stopwatchIsRunning() && m.swStopTime == "" {
				m.stopwatchStopTime()
				return m, nil
			}

			if m.swStopTime != "" {
				m.stopwatchResetTime()
				return m, nil
			}

			m.stopwatchStart()
			return m, Tick()
		}

	case TickMsg:
		if m.stopwatchIsRunning() {
			return m, Tick()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	var (
		fb = lbfb.NewFrameBuffer(uint16(m.width), uint16(m.height))

		vsplit = l.Vercital().
			Constrains(
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 50),
				lbl.NewConstrain(lbl.Percent, 25),
			).Split(fb.Size())

		timeRect = vsplit[1]

		time = lb.MarginLeft(max(0, (int(timeRect.Width)/2)-(7*4)/2),
			lb.ExpandVertical(int(timeRect.Height), lb.Center,
				m.stopwatchGetElapsed(),
			),
		)
	)

	fb.RenderString(time, timeRect)

	return fb.View()
}

func (m model) Init() tea.Cmd {
	return nil
}
