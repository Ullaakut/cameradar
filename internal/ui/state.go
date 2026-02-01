package ui

import (
	"context"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState struct {
	steps           []cameradar.Step
	status          map[cameradar.Step]state
	logs            []logMsg
	summary         []summaryTable
	summaryStreams  []cameradar.Stream
	summaryFinal    bool
	buildInfo       BuildInfo
	cancel          context.CancelFunc
	debug           bool
	spinner         spinner.Model
	progress        progress.Model
	width           int
	quitting        bool
	progressTotals  map[cameradar.Step]int
	progressCounts  map[cameradar.Step]int
	progressTarget  float64
	progressVisible float64
}

func (m *modelState) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *modelState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch typed := msg.(type) {
	case stepMsg:
		m.handleStepMsg(typed)
	case logMsg:
		m.handleLogMsg(typed)
	case summaryMsg:
		m.handleSummaryMsg(typed)
	case progressMsg:
		m.handleProgressMsg(typed)
	case closeMsg:
		m.quitting = true
	case tea.KeyMsg:
		if typed.Type == tea.KeyCtrlC {
			if m.cancel != nil {
				m.cancel()
			}
			m.quitting = true
			return m, tea.Quit
		}
	case spinner.TickMsg:
		cmds = m.handleSpinnerMsg(typed)
	case tea.WindowSizeMsg:
		m.handleWindowSizeMsg(typed)
	case progress.FrameMsg:
	}

	if len(cmds) == 0 {
		return m, nil
	}
	return m, tea.Batch(cmds...)
}

func (m *modelState) handleStepMsg(msg stepMsg) {
	m.status[msg.step] = msg.state
	if msg.message != "" {
		level := logInfo
		if msg.state == stateError {
			level = logError
		}
		m.logs = append(m.logs, logMsg{level: level, step: msg.step, message: msg.message})
	}
	if msg.state == stateDone || msg.state == stateError {
		markStepComplete(m, msg.step)
		queueProgressUpdate(m)
	}
	m.summary = buildSummaryTables(m.summaryStreams, m.width, m.status, m.summaryFinal)
}

func (m *modelState) handleLogMsg(msg logMsg) {
	m.logs = append(m.logs, msg)
}

func (m *modelState) handleSummaryMsg(msg summaryMsg) {
	m.summaryStreams = msg.streams
	m.summaryFinal = msg.final
	m.summary = buildSummaryTables(msg.streams, m.width, m.status, msg.final)
	if msg.final {
		m.status[cameradar.StepSummary] = stateDone
		markStepComplete(m, cameradar.StepSummary)
		queueProgressUpdate(m)
		m.quitting = true
	}
}

func (m *modelState) handleProgressMsg(msg progressMsg) {
	if msg.total > 0 {
		m.progressTotals[msg.step] = msg.total
		if m.progressCounts[msg.step] > msg.total {
			m.progressCounts[msg.step] = msg.total
		}
	}

	if msg.increment > 0 {
		m.progressCounts[msg.step] += msg.increment
		total := m.progressTotals[msg.step]
		if total > 0 && m.progressCounts[msg.step] > total {
			m.progressCounts[msg.step] = total
		}
	}

	queueProgressUpdate(m)
}

func (m *modelState) handleSpinnerMsg(msg spinner.TickMsg) []tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	advanceProgress(m)
	if m.quitting && progressComplete(*m) {
		cmds = append(cmds, tea.Quit)
	}
	return cmds
}

func (m *modelState) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.progress.Width = progressWidth(msg.Width)
	m.summary = buildSummaryTables(m.summaryStreams, m.width, m.status, m.summaryFinal)
}

func (m *modelState) View() string {
	var builder strings.Builder
	builder.WriteString(sectionStyle.Render(m.buildInfo.TUIHeader()))
	builder.WriteString("\n")
	builder.WriteString(renderProgress(m))
	builder.WriteString("\n")

	spinnerView := m.spinner.View()
	for _, step := range m.steps {
		builder.WriteString(renderStep(step, m.status[step], spinnerView))
		builder.WriteString("\n")
	}

	builder.WriteString("\n")
	builder.WriteString(sectionStyle.Render("Logs"))
	builder.WriteString("\n")
	if len(m.logs) == 0 {
		builder.WriteString(dimStyle.Render("No events yet."))
		builder.WriteString("\n")
	} else {
		for _, entry := range m.logs {
			builder.WriteString(renderLog(entry))
			builder.WriteString("\n")
		}
	}

	builder.WriteString("\n")
	builder.WriteString(sectionStyle.Render("Summary"))
	builder.WriteString("\n")
	for i, summary := range m.summary {
		if summary.title != "" {
			builder.WriteString(subsectionStyle.Render(summary.title))
			builder.WriteString("\n")
		}
		if summary.emptyMessage != "" {
			builder.WriteString(dimStyle.Render(summary.emptyMessage))
			builder.WriteString("\n")
		} else {
			builder.WriteString(summaryTableStyle.Render(summary.table.View()))
			builder.WriteString("\n")
		}
		if i < len(m.summary)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
