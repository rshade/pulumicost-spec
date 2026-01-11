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

package pluginsdk

import (
	"errors"
	"fmt"
	"time"

	"github.com/rshade/pulumicost-spec/sdk/go/currency"
	pbc "github.com/rshade/pulumicost-spec/sdk/go/proto/pulumicost/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ContractCommitmentBuilder handles the construction of FOCUS 1.3 ContractCommitment records.
// This builder creates records for the Contract Commitment supplemental dataset,
// which tracks contractual obligations separately from cost line items.
//
// Reference: FOCUS 1.3 Contract Commitment Dataset.
type ContractCommitmentBuilder struct {
	record *pbc.ContractCommitment
}

// NewContractCommitmentBuilder creates a new builder instance for ContractCommitment records.
func NewContractCommitmentBuilder() *ContractCommitmentBuilder {
	return &ContractCommitmentBuilder{
		record: &pbc.ContractCommitment{},
	}
}

// WithIdentity sets the identity fields for the contract commitment.
// ContractCommitmentId is the unique identifier for this specific commitment (REQUIRED).
// ContractId is the identifier of the parent contract containing this commitment (REQUIRED).
// FOCUS 1.3 Section: Contract Commitment ID, Contract ID.
func (b *ContractCommitmentBuilder) WithIdentity(
	commitmentID, contractID string,
) *ContractCommitmentBuilder {
	b.record.ContractCommitmentId = commitmentID
	b.record.ContractId = contractID
	return b
}

// WithCategory sets the commitment category.
// Category indicates whether this is a SPEND or USAGE commitment.
// FOCUS 1.3 Section: Contract Commitment Category.
func (b *ContractCommitmentBuilder) WithCategory(
	category pbc.FocusContractCommitmentCategory,
) *ContractCommitmentBuilder {
	b.record.ContractCommitmentCategory = category
	return b
}

// WithType sets the provider-specific commitment type.
// Examples: "Reserved Instance", "Savings Plan", "Committed Use Discount", "Enterprise Agreement"
// FOCUS 1.3 Section: Contract Commitment Type.
func (b *ContractCommitmentBuilder) WithType(
	commitmentType string,
) *ContractCommitmentBuilder {
	b.record.ContractCommitmentType = commitmentType
	return b
}

// WithCommitmentPeriod sets the start and end of the commitment period.
// This is when the specific commitment obligations are active.
// FOCUS 1.3 Section: Contract Commitment Period Start/End.
func (b *ContractCommitmentBuilder) WithCommitmentPeriod(
	start, end time.Time,
) *ContractCommitmentBuilder {
	b.record.ContractCommitmentPeriodStart = timestamppb.New(start)
	b.record.ContractCommitmentPeriodEnd = timestamppb.New(end)
	return b
}

// WithContractPeriod sets the start and end of the overall contract period.
// This is when the parent contract agreement is active.
// FOCUS 1.3 Section: Contract Period Start/End.
func (b *ContractCommitmentBuilder) WithContractPeriod(
	start, end time.Time,
) *ContractCommitmentBuilder {
	b.record.ContractPeriodStart = timestamppb.New(start)
	b.record.ContractPeriodEnd = timestamppb.New(end)
	return b
}

// WithFinancials sets all financial fields for the commitment.
// - cost: Monetary amount of the commitment (for SPEND category)
// - quantity: Quantity amount of the commitment (for USAGE category)
// - unit: Unit of measure for quantity (e.g., "Hours", "GB", "vCPU-Hours")
// - currencyCode: ISO 4217 currency code for monetary values (REQUIRED)
//
// For granular control, use WithCost(), WithQuantity(), and WithCurrency() instead.
// FOCUS 1.3 Section: Contract Commitment Cost, Quantity, Unit, Billing Currency.
func (b *ContractCommitmentBuilder) WithFinancials(
	cost, quantity float64, unit, currencyCode string,
) *ContractCommitmentBuilder {
	b.record.ContractCommitmentCost = cost
	b.record.ContractCommitmentQuantity = quantity
	b.record.ContractCommitmentUnit = unit
	b.record.BillingCurrency = currencyCode
	return b
}

// WithCost sets the monetary commitment amount for SPEND category commitments.
// Use this for commitments based on spend thresholds (e.g., "$100,000/year").
// FOCUS 1.3 Section: Contract Commitment Cost.
func (b *ContractCommitmentBuilder) WithCost(cost float64) *ContractCommitmentBuilder {
	b.record.ContractCommitmentCost = cost
	return b
}

// WithQuantity sets the quantity and unit for USAGE category commitments.
// Use this for commitments based on consumption (e.g., "1000 vCPU-Hours").
// FOCUS 1.3 Section: Contract Commitment Quantity, Contract Commitment Unit.
func (b *ContractCommitmentBuilder) WithQuantity(quantity float64, unit string) *ContractCommitmentBuilder {
	b.record.ContractCommitmentQuantity = quantity
	b.record.ContractCommitmentUnit = unit
	return b
}

// WithCurrency sets the ISO 4217 currency code for the commitment.
// This field is REQUIRED for all commitments.
// FOCUS 1.3 Section: Billing Currency.
func (b *ContractCommitmentBuilder) WithCurrency(currencyCode string) *ContractCommitmentBuilder {
	b.record.BillingCurrency = currencyCode
	return b
}

// Build validates and returns the constructed ContractCommitment record.
// Returns an error if required fields are missing or validation rules are violated.
func (b *ContractCommitmentBuilder) Build() (*pbc.ContractCommitment, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}
	return b.record, nil
}

// validate checks all validation rules for the ContractCommitment record.
func (b *ContractCommitmentBuilder) validate() error {
	// Required fields validation
	if b.record.GetContractCommitmentId() == "" {
		return errors.New("contract_commitment_id is required")
	}
	if b.record.GetContractId() == "" {
		return errors.New("contract_id is required")
	}
	if b.record.GetBillingCurrency() == "" {
		return errors.New("billing_currency is required")
	}

	// Category validation - must be explicitly set to SPEND or USAGE
	category := b.record.GetContractCommitmentCategory()
	if category != pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_SPEND &&
		category != pbc.FocusContractCommitmentCategory_FOCUS_CONTRACT_COMMITMENT_CATEGORY_USAGE {
		return fmt.Errorf("contract_commitment_category must be SPEND or USAGE, got %v", category)
	}

	// Currency validation using existing ISO 4217 validator
	if !currency.IsValid(b.record.GetBillingCurrency()) {
		return fmt.Errorf(
			"billing_currency must be a valid ISO 4217 currency code, got %q",
			b.record.GetBillingCurrency(),
		)
	}

	// Period consistency validation
	if err := b.validatePeriods(); err != nil {
		return err
	}

	// Non-negative value validation
	if b.record.GetContractCommitmentCost() < 0 {
		return errors.New("contract_commitment_cost must be non-negative")
	}
	if b.record.GetContractCommitmentQuantity() < 0 {
		return errors.New("contract_commitment_quantity must be non-negative")
	}

	return nil
}

// validatePeriods ensures period end >= period start for both commitment and contract periods.
func (b *ContractCommitmentBuilder) validatePeriods() error {
	// Validate commitment period if both are set
	if b.record.GetContractCommitmentPeriodStart() != nil && b.record.GetContractCommitmentPeriodEnd() != nil {
		start := b.record.GetContractCommitmentPeriodStart().AsTime()
		end := b.record.GetContractCommitmentPeriodEnd().AsTime()
		if end.Before(start) {
			return fmt.Errorf(
				"contract_commitment_period_end (%s) must be >= contract_commitment_period_start (%s)",
				end.Format(time.RFC3339), start.Format(time.RFC3339),
			)
		}
	}

	// Validate contract period if both are set
	if b.record.GetContractPeriodStart() != nil && b.record.GetContractPeriodEnd() != nil {
		start := b.record.GetContractPeriodStart().AsTime()
		end := b.record.GetContractPeriodEnd().AsTime()
		if end.Before(start) {
			return fmt.Errorf(
				"contract_period_end (%s) must be >= contract_period_start (%s)",
				end.Format(time.RFC3339), start.Format(time.RFC3339),
			)
		}
	}

	return nil
}
