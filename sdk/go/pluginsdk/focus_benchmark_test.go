package pluginsdk_test

import (
	"testing"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
)

func BenchmarkFocusRecordBuilder(b *testing.B) {
	start := time.Now()
	end := start.Add(time.Hour)

	b.ResetTimer()
	for range b.N {
		builder := pluginsdk.NewFocusRecordBuilder()
		builder.WithIdentity("aws", "acc-123", "My Account")
		builder.WithChargePeriod(start, end)
		builder.WithServiceCategory(pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE)
		builder.WithChargeDetails(
			pbc.FocusChargeCategory_FOCUS_CHARGE_CATEGORY_USAGE,
			pbc.FocusPricingCategory_FOCUS_PRICING_CATEGORY_STANDARD,
		)
		builder.WithFinancials(10.5, 12.0, 10.0, "USD", "inv-001")
		builder.WithUsage(1.0, "Hour")
		builder.WithExtension("Env", "Prod")
		_, _ = builder.Build()
	}
}
