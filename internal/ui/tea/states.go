package tea

// TUIState represents the current operational state of the TUI.
type TUIState int

const (
	StateNormal TUIState = iota
	StateCommandMenu
	StateInput
	StatePrompt
	StateRunning
	StatePostRun
	StateDeleteConfirm
	StateDeleteDepsConfirm
)

type InputModeType string

const (
	InputModePath InputModeType = "Path"
	InputModeName InputModeType = "Name"
)
