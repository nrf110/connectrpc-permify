package connectpermify

type stubCheckable struct {
	Checkable
	checks CheckConfig
}

func (r *stubCheckable) GetChecks() CheckConfig {
	return r.checks
}

func alwaysEnabled() bool { return true }
