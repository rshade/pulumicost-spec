package types

type BillingMode string

const (
	PerHour    BillingMode = "per_hour"
	PerGBMonth BillingMode = "per_gb_month"
	PerRequest BillingMode = "per_request"
	FlatRate   BillingMode = "flat"
	PerDay     BillingMode = "per_day"
	PerCPUHour BillingMode = "per_cpu_hour"
)

func (b BillingMode) String() string { return string(b) }

func ValidBillingMode(s string) bool {
	switch BillingMode(s) {
	case PerHour, PerGBMonth, PerRequest, FlatRate, PerDay, PerCPUHour:
		return true
	default:
		return false
	}
}
