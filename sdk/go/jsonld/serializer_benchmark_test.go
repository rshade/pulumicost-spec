package jsonld_test

import (
	"bytes"
	"testing"

	"github.com/rshade/pulumicost-spec/sdk/go/jsonld"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// createTestRecord creates a typical FocusCostRecord for benchmarking.
func createTestRecord(i int) *pbc.FocusCostRecord {
	return &pbc.FocusCostRecord{
		BillingAccountId:   "123456789012",
		BillingAccountName: "Production Account",
		ChargePeriodStart: &timestamppb.Timestamp{
			Seconds: int64(1735689600 + i*86400),
		},
		ChargePeriodEnd: &timestamppb.Timestamp{
			Seconds: int64(1735776000 + i*86400),
		},
		ServiceCategory: pbc.FocusServiceCategory_FOCUS_SERVICE_CATEGORY_COMPUTE,
		ServiceName:     "Amazon EC2",
		ResourceId:      "i-1234567890abcdef0",
		ResourceName:    "production-web-server",
		ResourceType:    "m5.large",
		RegionId:        "us-east-1",
		RegionName:      "US East (N. Virginia)",
		BilledCost:      125.50,
		ListCost:        150.00,
		EffectiveCost:   125.50,
		BillingCurrency: "USD",
		Tags: map[string]string{
			"environment": "production",
			"team":        "engineering",
			"cost-center": "CC-12345",
		},
		ServiceProviderName: "Amazon Web Services",
		HostProviderName:    "Amazon Web Services",
	}
}

// BenchmarkSerialize_SingleRecord benchmarks serialization of a single record.
func BenchmarkSerialize_SingleRecord(b *testing.B) {
	serializer := jsonld.NewSerializer()
	record := createTestRecord(0)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := serializer.Serialize(record)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerialize_SingleRecord_PrettyPrint benchmarks with pretty printing enabled.
func BenchmarkSerialize_SingleRecord_PrettyPrint(b *testing.B) {
	serializer := jsonld.NewSerializer(jsonld.WithPrettyPrint(true))
	record := createTestRecord(0)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, err := serializer.Serialize(record)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerializeBatch_100 benchmarks batch serialization of 100 records.
func BenchmarkSerializeBatch_100(b *testing.B) {
	serializer := jsonld.NewSerializer()
	records := make([]*pbc.FocusCostRecord, 100)
	for i := range 100 {
		records[i] = createTestRecord(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _, err := serializer.SerializeBatch(records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerializeBatch_1000 benchmarks batch serialization of 1000 records.
func BenchmarkSerializeBatch_1000(b *testing.B) {
	serializer := jsonld.NewSerializer()
	records := make([]*pbc.FocusCostRecord, 1000)
	for i := range 1000 {
		records[i] = createTestRecord(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _, err := serializer.SerializeBatch(records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerializeBatch_10000 benchmarks batch serialization of 10000 records.
func BenchmarkSerializeBatch_10000(b *testing.B) {
	serializer := jsonld.NewSerializer()
	records := make([]*pbc.FocusCostRecord, 10000)
	for i := range 10000 {
		records[i] = createTestRecord(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _, err := serializer.SerializeBatch(records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerializeStream_10000 benchmarks streaming serialization of 10000 records.
func BenchmarkSerializeStream_10000(b *testing.B) {
	serializer := jsonld.NewSerializer()
	records := make([]*pbc.FocusCostRecord, 10000)
	for i := range 10000 {
		records[i] = createTestRecord(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		var buf bytes.Buffer
		_, err := serializer.SerializeSlice(records, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSerializeStream_Memory measures memory usage during streaming.
func BenchmarkSerializeStream_Memory(b *testing.B) {
	serializer := jsonld.NewSerializer()
	records := make([]*pbc.FocusCostRecord, 10000)
	for i := range 10000 {
		records[i] = createTestRecord(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		var buf bytes.Buffer
		buf.Grow(1024 * 1024) // Pre-allocate 1MB to reduce reallocation
		_, err := serializer.SerializeSlice(records, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkIDGenerator benchmarks ID generation.
func BenchmarkIDGenerator(b *testing.B) {
	gen := jsonld.NewIDGenerator()
	record := createTestRecord(0)

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = gen.Generate(record)
	}
}

// BenchmarkContextBuild benchmarks context building.
func BenchmarkContextBuild(b *testing.B) {
	ctx := jsonld.NewContext()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = ctx.Build()
	}
}
