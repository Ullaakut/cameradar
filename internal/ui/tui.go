package ui

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Ullaakut/cameradar/v6"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	statePending state = iota
	stateActive
	stateDone
	stateError
)

type logLevel int

const (
	logInfo logLevel = iota
	logDebug
	logError
)

type stepMsg struct {
	step    cameradar.Step
	state   state
	message string
}

type logMsg struct {
	level   logLevel
	step    cameradar.Step
	message string
}

type progressMsg struct {
	step      cameradar.Step
	total     int
	increment int
}

type closeMsg struct{}

type summaryMsg struct {
	streams []cameradar.Stream
	final   bool
}

type summaryTable struct {
	table table.Model
}

const (
	summaryMinHeight   = 8
	summaryMaxHeight   = 10
	summaryColumnCount = 8
)

// TUIReporter renders a Bubble Tea based UI.
type TUIReporter struct {
	program *tea.Program
	debug   bool
	once    sync.Once
	closed  chan struct{}
	mu      sync.Mutex
	last    []cameradar.Stream
}

// NewTUIReporter creates a new Bubble Tea reporter.
func NewTUIReporter(debug bool, out io.Writer, buildInfo BuildInfo, cancel context.CancelFunc) (*TUIReporter, error) {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithFillCharacters('━', '·'),
		progress.WithoutPercentage(),
		progress.WithWidth(28),
	)

	initial := &modelState{
		steps:          cameradar.Steps(),
		status:         make(map[cameradar.Step]state),
		debug:          debug,
		buildInfo:      buildInfo,
		cancel:         cancel,
		spinner:        spin,
		progress:       prog,
		progressTotals: make(map[cameradar.Step]int),
		progressCounts: make(map[cameradar.Step]int),
	}

	p := tea.NewProgram(initial, tea.WithInputTTY(), tea.WithOutput(out), tea.WithAltScreen())
	reporter := &TUIReporter{program: p, debug: debug, closed: make(chan struct{})}

	go func() {
		model, err := p.Run()
		if err != nil {
			_, _ = fmt.Fprintf(out, "Error running TUI: %v\n", err)
			close(reporter.closed)
			return
		}

		if rendered, ok := model.(*modelState); ok {
			output := rendered.FinalView()
			if len(rendered.summaryStreams) == 0 {
				fallback := reporter.snapshotSummary()
				if len(fallback) > 0 {
					tmp := &modelState{
						summaryStreams: fallback,
						width:          rendered.width,
						status:         summaryStatusAllDone(),
					}
					output = tmp.FinalView()
				}
			}
			_, _ = fmt.Fprintln(out, output)
		}
		close(reporter.closed)
	}()

	return reporter, nil
}

// Start implements Reporter.
func (r *TUIReporter) Start(step cameradar.Step, message string) {
	r.send(stepMsg{step: step, state: stateActive, message: message})
}

// Done implements Reporter.
func (r *TUIReporter) Done(step cameradar.Step, message string) {
	r.send(stepMsg{step: step, state: stateDone, message: message})
}

// Progress implements Reporter.
func (r *TUIReporter) Progress(step cameradar.Step, message string) {
	if kind, value, ok := cameradar.ParseProgressMessage(message); ok {
		msg := progressMsg{step: step}
		if kind == "total" {
			msg.total = value
		}
		if kind == "tick" {
			msg.increment = value
		}
		r.send(msg)
		return
	}

	r.send(logMsg{level: logInfo, step: step, message: message})
}

// Debug implements Reporter.
func (r *TUIReporter) Debug(step cameradar.Step, message string) {
	if !r.debug {
		return
	}

	r.send(logMsg{level: logDebug, step: step, message: message})
}

// Error implements Reporter.
func (r *TUIReporter) Error(step cameradar.Step, err error) {
	if err == nil {
		return
	}

	r.send(stepMsg{step: step, state: stateError, message: err.Error()})
}

// Summary implements Reporter.
func (r *TUIReporter) Summary(streams []cameradar.Stream, _ error) {
	cloned := copyStreams(streams)
	r.recordSummary(cloned)
	r.send(summaryMsg{streams: cloned, final: true})
}

// UpdateSummary updates the summary section with partial results.
func (r *TUIReporter) UpdateSummary(streams []cameradar.Stream) {
	cloned := copyStreams(streams)
	r.recordSummary(cloned)
	r.send(summaryMsg{streams: cloned, final: false})
}

func (r *TUIReporter) recordSummary(streams []cameradar.Stream) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last = streams
}

func (r *TUIReporter) snapshotSummary() []cameradar.Stream {
	r.mu.Lock()
	defer r.mu.Unlock()
	return copyStreams(r.last)
}

// Close implements Reporter.
func (r *TUIReporter) Close() {
	r.once.Do(func() {
		r.send(closeMsg{})
	})

	// Timeout after 2 seconds to avoid hanging forever.
	select {
	case <-r.closed:
	case <-time.After(2 * time.Second):
	}
}

func (r *TUIReporter) send(msg tea.Msg) {
	if r.program == nil {
		return
	}

	r.program.Send(msg)
}

func renderStep(step cameradar.Step, state state, spinnerView string) string {
	label := cameradar.StepLabel(step)
	symbol := "·"
	style := dimStyle
	switch state {
	case stateActive:
		symbol = spinnerView
		style = activeStyle
	case stateDone:
		symbol = "✓"
		style = successStyle
	case stateError:
		symbol = "✗"
		style = errorStyle
	}
	return style.Render(fmt.Sprintf("%s %s", symbol, label))
}

func renderLog(entry logMsg) string {
	prefix := "INFO"
	style := infoStyle
	if entry.level == logDebug {
		prefix = "DEBUG"
		style = debugStyle
	}
	if entry.level == logError {
		prefix = "ERROR"
		style = errorStyle
	}
	return style.Render(fmt.Sprintf("[%s] %s: %s", prefix, cameradar.StepLabel(entry.step), entry.message))
}

func renderProgress(m *modelState) string {
	completed, total := progressCounts(m.steps, m.status)
	percent := progressPercent(m.steps, m.status, m.progressTotals, m.progressCounts)
	countLabel := dimStyle.Render(fmt.Sprintf("%3.0f%% %d/%d complete", percent*100, completed, total))
	return fmt.Sprintf("%s %s", m.progress.ViewAs(m.progressVisible), countLabel)
}

func progressCounts(steps []cameradar.Step, status map[cameradar.Step]state) (int, int) {
	if len(steps) == 0 {
		return 0, 0
	}

	completed := 0
	for _, step := range steps {
		switch status[step] {
		case stateDone, stateError:
			completed++
		}
	}

	return completed, len(steps)
}

func progressPercent(steps []cameradar.Step, status map[cameradar.Step]state, totals, counts map[cameradar.Step]int) float64 {
	weights := stepWeights()
	percent := 0.0
	for _, step := range steps {
		weight := weights[step]
		if weight <= 0 {
			continue
		}
		percent += weight * stepProgress(step, status, totals, counts)
	}
	if percent > 1 {
		return 1
	}
	return percent
}

func stepWeights() map[cameradar.Step]float64 {
	return map[cameradar.Step]float64{
		cameradar.StepScan:              0.15,
		cameradar.StepAttackRoutes:      0.25,
		cameradar.StepDetectAuth:        0.05,
		cameradar.StepAttackCredentials: 0.35,
		cameradar.StepValidateStreams:   0.2,
		cameradar.StepSummary:           0.0,
	}
}

func stepProgress(step cameradar.Step, status map[cameradar.Step]state, totals, counts map[cameradar.Step]int) float64 {
	if total := totals[step]; total > 0 {
		count := counts[step]
		if count >= total {
			return 1
		}
		return float64(count) / float64(total)
	}

	switch status[step] {
	case stateDone, stateError:
		return 1
	default:
		return 0
	}
}

func queueProgressUpdate(m *modelState) {
	desired := progressPercent(m.steps, m.status, m.progressTotals, m.progressCounts)
	if desired <= m.progressTarget {
		return
	}
	m.progressTarget = desired
}

func advanceProgress(m *modelState) {
	if m.progressVisible >= m.progressTarget {
		return
	}
	remaining := m.progressTarget - m.progressVisible
	step := remaining * 0.2
	if step < 0.02 {
		step = 0.02
	}
	if m.quitting && step < 0.08 {
		step = 0.08
	}
	if remaining < step {
		m.progressVisible = m.progressTarget
		return
	}
	m.progressVisible += step
}

func progressComplete(m modelState) bool {
	return m.progressVisible >= m.progressTarget
}

func markStepComplete(m *modelState, step cameradar.Step) {
	if m.progressTotals[step] == 0 {
		m.progressTotals[step] = 1
	}
	if m.progressCounts[step] < m.progressTotals[step] {
		m.progressCounts[step] = m.progressTotals[step]
	}
}

func progressWidth(width int) int {
	if width <= 0 {
		return 28
	}
	if width < 60 {
		return 20
	}
	if width < 100 {
		return 28
	}
	return 36
}

func buildSummaryTables(streams []cameradar.Stream, width int, status map[cameradar.Step]state, maxRows int) []summaryTable {
	visibility := summaryVisibility(status)
	accessible, others := partitionStreams(streams)
	rows := append(buildSummaryRows(accessible, visibility), buildSummaryRows(others, visibility)...)
	if len(rows) == 0 {
		rows = []table.Row{emptySummaryRow()}
	}

	if maxRows > 0 {
		switch {
		case len(rows) > maxRows:
			if maxRows == 1 {
				rows = []table.Row{summaryOverflowRow(len(rows))}
			} else {
				visibleRows := maxRows - 1
				hidden := len(rows) - visibleRows
				rows = append(rows[:visibleRows], summaryOverflowRow(hidden))
			}
		case len(rows) < maxRows:
			rows = padSummaryRows(rows, maxRows)
		}
	}

	columns := summaryColumns(width, rows)
	model := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)
	model.SetStyles(summaryTableStyles())

	return []summaryTable{{table: model}}
}

func renderSummaryTitle(streams []cameradar.Stream) string {
	accessible, _ := partitionStreams(streams)
	return fmt.Sprintf("Summary - Streams (%d accessible / %d total)", len(accessible), len(streams))
}

func summaryStatusAllDone() map[cameradar.Step]state {
	status := make(map[cameradar.Step]state)
	for _, step := range cameradar.Steps() {
		status[step] = stateDone
	}
	return status
}

const emptyEntry = "—"

func emptySummaryRow() table.Row {
	row := make(table.Row, summaryColumnCount)
	for i := range row {
		row[i] = emptyEntry
	}
	return row
}

func padSummaryRows(rows []table.Row, maxRows int) []table.Row {
	for len(rows) < maxRows {
		rows = append(rows, emptySummaryRow())
	}
	return rows
}

func summaryOverflowRow(hidden int) table.Row {
	row := emptySummaryRow()
	if hidden <= 0 {
		return row
	}
	label := "\u2026 1 more stream"
	if hidden > 1 {
		label = fmt.Sprintf("\u2026 %d more streams", hidden)
	}
	row[0] = label
	return row
}

func renderSummaryTablePlain(columns []table.Column, rows []table.Row) string {
	colWidths := make([]int, len(columns))
	for i, col := range columns {
		colWidths[i] = max(col.Width, len([]rune(col.Title)))
	}

	var builder strings.Builder
	builder.WriteString(renderSummaryBorder("┌", "┬", "┐", colWidths))
	builder.WriteString("\n")
	builder.WriteString(renderSummaryRow(columnTitles(columns), colWidths))
	builder.WriteString("\n")
	builder.WriteString(renderSummaryBorder("├", "┼", "┤", colWidths))
	for _, row := range rows {
		builder.WriteString("\n")
		builder.WriteString(renderSummaryRow(row, colWidths))
	}
	builder.WriteString("\n")
	builder.WriteString(renderSummaryBorder("└", "┴", "┘", colWidths))
	return builder.String()
}

func renderSummaryBorder(left, middle, right string, widths []int) string {
	parts := make([]string, 0, len(widths))
	for _, width := range widths {
		parts = append(parts, strings.Repeat("─", width+2))
	}
	return left + strings.Join(parts, middle) + right
}

func renderSummaryRow(cells []string, widths []int) string {
	var builder strings.Builder
	builder.WriteString("│")
	for i, width := range widths {
		value := ""
		if i < len(cells) {
			value = cells[i]
		}
		builder.WriteString(" ")
		builder.WriteString(padAndTrim(value, width))
		builder.WriteString(" │")
	}
	return builder.String()
}

func padAndTrim(value string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(value)
	if len(runes) > width {
		return string(runes[:width])
	}
	if len(runes) < width {
		return string(runes) + strings.Repeat(" ", width-len(runes))
	}
	return value
}

func columnTitles(columns []table.Column) []string {
	if len(columns) == 0 {
		return nil
	}
	titles := make([]string, len(columns))
	for i, col := range columns {
		titles[i] = col.Title
	}
	return titles
}

func buildSummaryRows(streams []cameradar.Stream, visibility summaryVisibilityState) []table.Row {
	rows := make([]table.Row, 0, len(streams))
	for _, stream := range streams {
		target := fmt.Sprintf("%s:%d", stream.Address.String(), stream.Port)
		device := emptyEntry
		if visibility.showDevice && stream.Device != "" {
			device = stream.Device
		}

		routes := emptyEntry
		if visibility.showRoutes && len(stream.Routes) > 0 {
			routes = strings.Join(stream.Routes, ", ")
		}

		credentials := emptyEntry
		if visibility.showCredentials && stream.CredentialsFound {
			credentials = fmt.Sprintf("%s:%s", stream.Username, stream.Password)
		}

		available := emptyEntry
		if visibility.showAvailable {
			available = "no"
			if stream.Available {
				available = "yes"
			}
		}

		rtspURL := emptyEntry
		if visibility.showCredentials && stream.RouteFound && stream.CredentialsFound {
			rtspURL = formatRTSPURL(stream)
		}

		authType := emptyEntry
		if visibility.showAuth {
			authType = authTypeLabel(stream.AuthenticationType)
		}

		rows = append(rows, table.Row{
			target,
			device,
			authType,
			routes,
			credentials,
			available,
			rtspURL,
			adminPanelLabel(stream, visibility),
		})
	}

	return rows
}

func summaryColumns(width int, rows []table.Row) []table.Column {
	columns := []table.Column{
		{Title: "Target", Width: 18},
		{Title: "Device", Width: 14},
		{Title: "Auth", Width: 8},
		{Title: "Routes", Width: 18},
		{Title: "Credentials", Width: 16},
		{Title: "Available", Width: 9},
		{Title: "RTSP URL", Width: 30},
		{Title: "Admin", Width: 24},
	}
	columns[6].Width = maxColumnWidth(columns[6].Title, rows, 6, columns[6].Width)
	columns[7].Width = maxColumnWidth(columns[7].Title, rows, 7, columns[7].Width)

	if width <= 0 {
		return columns
	}

	columns = clampColumns(columns, max(width-2, 60))

	return columns
}

func clampColumns(columns []table.Column, maxWidth int) []table.Column {
	padding := 2 * len(columns)
	contentWidth := 0
	for _, col := range columns {
		contentWidth += col.Width
	}
	contentWidth += padding
	if contentWidth <= maxWidth {
		return columns
	}

	over := contentWidth - maxWidth
	shrinkOrder := []int{7, 3, 4, 1}
	minWidths := map[int]int{
		7: 10,
		3: 10,
		4: 10,
		1: 10,
	}
	for over > 0 {
		changed := false
		for _, idx := range shrinkOrder {
			minWidth := minWidths[idx]
			if columns[idx].Width > minWidth {
				columns[idx].Width--
				over--
				changed = true
				if over == 0 {
					break
				}
			}
		}
		if !changed {
			break
		}
	}

	return columns
}

func summaryTableStyles() table.Styles {
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	styles.Selected = lipgloss.NewStyle()
	styles.Cell = styles.Cell.Padding(0, 1)
	return styles
}

func maxColumnWidth(title string, rows []table.Row, idx, minWidth int) int {
	width := max(len(title), minWidth)
	for _, row := range rows {
		if idx >= len(row) {
			continue
		}
		if len(row[idx]) > width {
			width = len(row[idx])
		}
	}
	return width
}

func adminPanelLabel(stream cameradar.Stream, visibility summaryVisibilityState) string {
	if !visibility.showCredentials || !stream.CredentialsFound {
		return emptyEntry
	}
	return formatAdminPanelURL(stream)
}

type summaryVisibilityState struct {
	showDevice      bool
	showRoutes      bool
	showAuth        bool
	showCredentials bool
	showAvailable   bool
}

func summaryVisibility(status map[cameradar.Step]state) summaryVisibilityState {
	return summaryVisibilityState{
		showDevice:      stepComplete(status, cameradar.StepScan),
		showRoutes:      stepComplete(status, cameradar.StepAttackRoutes),
		showAuth:        stepComplete(status, cameradar.StepDetectAuth),
		showCredentials: stepComplete(status, cameradar.StepAttackCredentials),
		showAvailable:   stepComplete(status, cameradar.StepValidateStreams),
	}
}

func stepComplete(status map[cameradar.Step]state, step cameradar.Step) bool {
	if status == nil {
		return false
	}
	switch status[step] {
	case stateDone, stateError:
		return true
	default:
		return false
	}
}

func copyStreams(streams []cameradar.Stream) []cameradar.Stream {
	if len(streams) == 0 {
		return nil
	}

	cloned := make([]cameradar.Stream, len(streams))
	copy(cloned, streams)
	return cloned
}
