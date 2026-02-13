package auditor

type RiskLevel int

const (
	Safe RiskLevel = iota
	Warning
	Danger
)

func Audit(command string) RiskLevel {
	// TODO: Implement Tier 1, 2, 3 checks
	return Safe
}
