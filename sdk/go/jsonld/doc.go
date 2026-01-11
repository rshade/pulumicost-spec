// Copyright 2026 PulumiCost/FinFocus Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jsonld provides JSON-LD 1.1 serialization for FOCUS cost data.
//
// This package transforms protobuf messages (FocusCostRecord, ContractCommitment) into
// JSON-LD format with Schema.org vocabulary mappings and custom FOCUS namespace support.
//
// # Basic Usage
//
//	serializer := jsonld.NewSerializer()
//	output, err := serializer.Serialize(record)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(string(output))
//
// # Context Customization
//
//	ctx := jsonld.NewContext().
//		WithCustomMapping("billingAccountId", "yourOrg:accountIdentifier")
//	serializer := jsonld.NewSerializer(
//		jsonld.WithContext(ctx),
//	)
//
// # Batch Serialization (Streaming)
//
//	serializer := jsonld.NewSerializer()
//	err := serializer.SerializeStream(recordChannel, writer)
//
// # Performance
//
// This package is optimized for high-throughput serialization:
//   - Single record: ~15.3µs for fully-populated FocusCostRecord
//   - Batch: ~182ms for 10,000 records
//   - Streaming: ~197ms for 10,000 records with bounded memory usage
//
// # JSON-LD 1.1 Compliance
//
// Output conforms to JSON-LD 1.1 specification:
//   - @context defines vocabulary mappings (Schema.org + FOCUS namespace)
//   - @id provides unique identifiers (user-provided or SHA256 fallback)
//   - @type declares record types
//   - Property names use compact IRIs defined in context
//
// See README.md for detailed examples and configuration options.
package jsonld
