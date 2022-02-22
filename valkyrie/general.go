package valkyrie

// Valkyrie is used to general used of
type Valkyrie struct {
	*core
}

// New returns Valkyrie object
func New(option *Option) *Valkyrie {
	return &Valkyrie{construct(option)}
}

// GetProjectID is used to convert  identifier to project_id
func (v *Valkyrie) GetProjectID(identifier string) (pid int64, ok bool) {
	return v.geti2p(identifier)
}
