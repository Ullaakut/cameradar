package cameradar

import (
	"fmt"
	"strings"
)

// Mode defines which UI renderer to use.
type Mode string

// Supported rendering modes.
const (
	ModeAuto  Mode = "auto"
	ModeTUI   Mode = "tui"
	ModePlain Mode = "plain"
)

// Step identifies a stage in the workflow.
type Step string

// Supported steps.
const (
	StepScan              Step = "scan"
	StepAttackRoutes      Step = "attack-routes"
	StepDetectAuth        Step = "detect-auth"
	StepAttackCredentials Step = "attack-credentials"
	StepValidateStreams   Step = "validate-streams"
	StepSummary           Step = "summary"
)

// StepLabel returns the human-readable label for a step.
func StepLabel(step Step) string {
	switch step {
	case StepScan:
		return "Scan targets"
	case StepAttackRoutes:
		return "Attack routes"
	case StepDetectAuth:
		return "Detect authentication"
	case StepAttackCredentials:
		return "Attack credentials"
	case StepValidateStreams:
		return "Validate streams"
	case StepSummary:
		return "Summary"
	default:
		return string(step)
	}
}

// Steps returns the ordered list of steps.
func Steps() []Step {
	return []Step{
		StepScan,
		StepAttackRoutes,
		StepDetectAuth,
		StepAttackCredentials,
		StepValidateStreams,
		StepSummary,
	}
}

// ParseMode parses a user-provided UI mode.
func ParseMode(value string) (Mode, error) {
	mode := Mode(strings.ToLower(strings.TrimSpace(value)))
	switch mode {
	case ModeAuto, ModeTUI, ModePlain:
		return mode, nil
	case "":
		return ModeAuto, nil
	default:
		return ModeAuto, fmt.Errorf("invalid ui mode %q", value)
	}
}
