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

	fb lbfb.FrameBuffer
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
		fb:          lbfb.NewFrameBuffer(0, 0),
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
		m.fb.Resize(msg.Width, msg.Height)
	}

	return m, cmd
}

func (m model) View() string {
	var (
		vsplit = l.Vercital().
			Constrains(
				lbl.NewConstrain(lbl.Percent, 25),
				lbl.NewConstrain(lbl.Percent, 50),
				lbl.NewConstrain(lbl.Percent, 25),
			).Split(m.fb.Size())

		timeRect = vsplit[1]

		time = m.stopwatchGetElapsed()
	)

	m.fb.Clear()

	m.fb.RenderString(time, timeRect, lb.Center, lb.Center)

	return m.fb.View()
}

func (m model) Init() tea.Cmd {
	return nil
}
