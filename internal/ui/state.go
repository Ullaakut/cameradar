package ui

import (
	"context"
	"strings"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState struct {
	steps           []cameradar.Step
	status          map[cameradar.Step]state
	logs            []logMsg
	summaryStreams  []cameradar.Stream
	summaryFinal    bool
	buildInfo       BuildInfo
	cancel          context.CancelFunc
	debug           bool
	spinner         spinner.Model
	progress        progress.Model
	width           int
	height          int
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
}

func (m *modelState) handleLogMsg(msg logMsg) {
	m.logs = append(m.logs, msg)
}

func (m *modelState) handleSummaryMsg(msg summaryMsg) {
	m.summaryStreams = msg.streams
	m.summaryFinal = msg.final
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
	m.height = msg.Height
	m.progress.Width = progressWidth(msg.Width)
}

func (m *modelState) View() string {
	var builder strings.Builder
	header := sectionStyle.Render(m.buildInfo.TUIHeader())
	headerLines := splitLines(header)
	builder.WriteString(strings.Join(headerLines, "\n"))
	builder.WriteString("\n\n")

	stepsLines := m.renderSteps()
	builder.WriteString(strings.Join(stepsLines, "\n"))
	builder.WriteString("\n\n")

	summaryHeight, logsHeight := m.layoutHeights(len(headerLines), len(stepsLines))
	logsLines := m.renderLogs(logsHeight)
	builder.WriteString(sectionStyle.Render("Logs"))
	builder.WriteString("\n")
	builder.WriteString(strings.Join(logsLines, "\n"))
	builder.WriteString("\n\n")

	rowsToShow := max(1, summaryHeight-2)
	summaryTitle := renderSummaryTitle(m.summaryStreams)
	summaryTables := buildSummaryTables(m.summaryStreams, m.width, m.status, m.summaryFinal, rowsToShow)
	builder.WriteString(sectionStyle.Render(summaryTitle))
	builder.WriteString("\n")
	for i, summary := range summaryTables {
		if summary.emptyMessage != "" {
			builder.WriteString(dimStyle.Render(summary.emptyMessage))
			builder.WriteString("\n")
			continue
		}
		builder.WriteString(summaryTableStyle.Render(summary.table.View()))
		if i < len(summaryTables)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (m *modelState) FinalView() string {
	var builder strings.Builder
	header := sectionStyle.Render(m.buildInfo.TUIHeader())
	headerLines := splitLines(header)
	builder.WriteString(strings.Join(headerLines, "\n"))
	builder.WriteString("\n\n")

	stepsLines := m.renderSteps()
	builder.WriteString(strings.Join(stepsLines, "\n"))
	builder.WriteString("\n\n")

	builder.WriteString(sectionStyle.Render("Logs"))
	builder.WriteString("\n")
	logLines := m.renderLogsAll()
	if len(logLines) == 0 {
		builder.WriteString(dimStyle.Render("No events yet."))
	} else {
		builder.WriteString(strings.Join(logLines, "\n"))
	}
	builder.WriteString("\n\n")

	summaryTitle := renderSummaryTitle(m.summaryStreams)
	visibility := summaryVisibility(summaryStatusAllDone())
	accessible, others := partitionStreams(m.summaryStreams)
	rows := append(buildSummaryRows(accessible, visibility), buildSummaryRows(others, visibility)...)
	if len(rows) == 0 {
		rows = []table.Row{emptySummaryRow()}
	}
	columns := summaryColumns(m.width, rows)
	builder.WriteString(sectionStyle.Render(summaryTitle))
	builder.WriteString("\n")
	builder.WriteString(renderSummaryTablePlain(columns, rows))
	return builder.String()
}

func (m *modelState) renderSteps() []string {
	lines := []string{sectionStyle.Render("Steps"), renderProgress(m)}
	spinnerView := m.spinner.View()
	for _, step := range m.steps {
		lines = append(lines, renderStep(step, m.status[step], spinnerView))
	}
	return lines
}

func (m *modelState) renderLogs(height int) []string {
	if height <= 0 {
		return nil
	}
	if len(m.logs) == 0 {
		lines := []string{dimStyle.Render("No events yet.")}
		return padLines(lines, height)
	}

	start := 0
	if len(m.logs) > height {
		start = len(m.logs) - height
	}
	lines := make([]string, 0, min(height, len(m.logs)))
	for _, entry := range m.logs[start:] {
		lines = append(lines, renderLog(entry))
	}
	return padLines(lines, height)
}

func (m *modelState) renderLogsAll() []string {
	if len(m.logs) == 0 {
		return nil
	}
	lines := make([]string, 0, len(m.logs))
	for _, entry := range m.logs {
		lines = append(lines, renderLog(entry))
	}
	return lines
}

func (m *modelState) layoutHeights(headerLines, stepsLines int) (summaryHeight, logsHeight int) {
	if m.height <= 0 {
		return summaryMinHeight, len(m.logs)
	}

	reserved := headerLines + 1 + stepsLines + 1 + 1 + 1
	remaining := m.height - reserved
	remaining = max(0, remaining)

	switch {
	case remaining < summaryMinHeight:
		summaryHeight = max(3, remaining)
	case remaining > summaryMaxHeight:
		summaryHeight = summaryMaxHeight
	default:
		summaryHeight = remaining
	}

	logsHeight = max(0, remaining-summaryHeight)

	return summaryHeight, logsHeight
}

func padLines(lines []string, height int) []string {
	if height <= 0 {
		return lines
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return lines
}

func splitLines(value string) []string {
	return strings.Split(value, "\n")
}
