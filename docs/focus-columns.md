# FOCUS 1.2 Column Reference

This document provides a comprehensive reference for all 57 columns defined in the
FinOps FOCUS 1.2 (FinOps Open Cost and Usage Specification) as implemented in
PulumiCost.

Reference: <https://focus.finops.org/focus-specification/v1-2/>

## Column Summary

| Level | Count | Description |
|-------|-------|-------------|
| Mandatory | 14 | Required for all cost records |
| Recommended | 1 | Strongly suggested for completeness |
| Conditional | 42 | Required when applicable conditions are met |
| **Total** | **57** | Complete FOCUS 1.2 coverage |

## Mandatory Columns (14)

These columns MUST be present in every FOCUS-compliant cost record.

### Identity & Hierarchy

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **ProviderName** | string | Cloud provider name | AWS, Azure, GCP, Kubernetes |
| **BillingAccountId** | string | Billing account identifier | AWS: Payer Account ID, Azure: Billing Account ID, GCP: Billing Account ID |
| **BillingAccountName** | string | Display name for billing account | AWS: Account Alias, Azure: Billing Account Name, GCP: Billing Account Name |

### Billing Period

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **BillingPeriodStart** | DateTime | Start of billing period | First day of billing month |
| **BillingPeriodEnd** | DateTime | End of billing period | Last day of billing month |
| **BillingCurrency** | string | ISO 4217 currency code | USD, EUR, GBP, etc. |

### Charge Period

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **ChargePeriodStart** | DateTime | When charge began | Line item start time |
| **ChargePeriodEnd** | DateTime | When charge ended | Line item end time |

### Charge Details

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **ChargeCategory** | Enum | Nature of charge | Usage, Purchase, Credit, Tax, Refund, Adjustment |
| **ChargeClass** | Enum | Classification | Regular, Correction |
| **ChargeDescription** | string | Human-readable description | AWS: lineItem/LineItemDescription |

### Service

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **ServiceName** | string | Name of the service | AWS: EC2, S3; Azure: Virtual Machines; GCP: Compute Engine |

### Financial Amounts

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **BilledCost** | Decimal | Amount billed | AWS: lineItem/BlendedCost, Azure: BilledCost |
| **ContractedCost** | Decimal | ContractedUnitPrice × PricingQuantity | Calculated from contracted rates |

## Recommended Column (1)

This column is strongly recommended but not mandatory.

| Column | Type | Description | Provider Mapping |
|--------|------|-------------|------------------|
| **ChargeFrequency** | Enum | How often charge applies | OneTime, Recurring, UsageBased |

## Conditional Columns (42)

These columns are required when specific conditions apply.

### Identity & Hierarchy - Conditional

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **SubAccountId** | string | When sub-accounts exist | AWS: Account ID, Azure: Subscription ID, GCP: Project ID |
| **SubAccountName** | string | When sub-accounts exist | Account/Subscription/Project name |
| **BillingAccountType** | string | When account types vary | Enterprise, PayAsYouGo, etc. |
| **SubAccountType** | string | When sub-account types vary | LinkedAccount, Subscription, Project |

### Location

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **RegionId** | string | When resource has region | AWS: us-east-1, Azure: eastus, GCP: us-central1 |
| **RegionName** | string | When resource has region | US East (N. Virginia), East US, etc. |
| **AvailabilityZone** | string | When AZ is applicable | AWS: us-east-1a, Azure: 1, GCP: us-central1-a |

### Resource Details

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **ResourceId** | string | When resource is identifiable | AWS: ARN, Azure: Resource ID, GCP: Resource ID |
| **ResourceName** | string | When resource has name | User-assigned resource name |
| **ResourceType** | string | When resource type is known | m5.large, Standard_D2s_v3, n1-standard-1 |

### SKU Details

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **SkuId** | string | When SKU is known | Provider-specific SKU identifier |
| **SkuPriceId** | string | When SKU price is known | Specific price list ID |
| **SkuMeter** | string | When meter exists | AWS: BoxUsage, Azure: Meter Name |
| **SkuPriceDetails** | string | When additional details exist | Operating system, license type, etc. |

### Pricing Details

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **PricingQuantity** | Decimal | When pricing is quantifiable | Number of units priced |
| **PricingUnit** | string | When pricing unit exists | Hours, GB, Requests |
| **ListUnitPrice** | Decimal | When list price is known | Public on-demand rate |
| **ContractedUnitPrice** | Decimal | When contracted price exists | Reserved/committed rate |
| **PricingCategory** | Enum | When pricing model varies | Standard, Committed, Dynamic, Other |
| **PricingCurrency** | string | When different from billing currency | ISO 4217 code |
| **PricingCurrencyContractedUnitPrice** | Decimal | When pricing currency differs | Contracted price in pricing currency |
| **PricingCurrencyEffectiveCost** | Decimal | When pricing currency differs | Effective cost in pricing currency |
| **PricingCurrencyListUnitPrice** | Decimal | When pricing currency differs | List price in pricing currency |

### Financial Amounts - Conditional

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **ListCost** | Decimal | When list price is known | On-demand cost at list price |
| **EffectiveCost** | Decimal | When discounts apply | Actual cost after all adjustments |

### Consumption/Usage

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **ConsumedQuantity** | Decimal | When usage is measurable | AWS: UsageQuantity, Azure: Quantity |
| **ConsumedUnit** | string | When usage unit exists | Hours, GB, Requests |

### Service - Conditional

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **ServiceCategory** | Enum | When categorization is possible | Compute, Storage, Network, Database, etc. |
| **ServiceSubcategory** | string | When granular categorization exists | Virtual Machine, Container, Serverless |
| **Publisher** | string | When publisher differs from provider | Amazon Web Services, Microsoft, Google, Third-party |

### Commitment Discounts

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **CommitmentDiscountCategory** | Enum | When commitment exists | Spend, Usage |
| **CommitmentDiscountId** | string | When commitment exists | RI ID, Savings Plan ID |
| **CommitmentDiscountName** | string | When commitment has name | User-assigned name |
| **CommitmentDiscountQuantity** | Decimal | When commitment has quantity | Amount purchased/consumed |
| **CommitmentDiscountStatus** | Enum | When commitment applies | Used, Unused |
| **CommitmentDiscountType** | string | When type is known | Reserved Instance, Savings Plan |
| **CommitmentDiscountUnit** | string | When unit exists | Hours, USD/Hour |

### Capacity Reservation

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **CapacityReservationId** | string | When capacity reservation exists | CR ID from provider |
| **CapacityReservationStatus** | Enum | When CR applies | Used, Unused |

### Invoice Details

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **InvoiceId** | string | When invoice exists | AWS: bill/InvoiceId, Azure: InvoiceId |
| **InvoiceIssuer** | string | When issuer is identifiable | Legal entity name |

### Metadata

| Column | Type | Condition | Provider Mapping |
|--------|------|-----------|------------------|
| **Tags** | map\<string,string\> | When resource has tags | User-defined key-value pairs |

## Provider Mapping Examples

### AWS

```text
ProviderName:           "AWS"
BillingAccountId:       "123456789012"
ServiceName:            "Amazon EC2"
ResourceId:             "arn:aws:ec2:us-east-1:123456789012:instance/i-0abc123"
RegionId:               "us-east-1"
CommitmentDiscountType: "Standard Reserved Instance"
```

### Azure

```text
ProviderName:           "Azure"
BillingAccountId:       "ea12345678"
SubAccountId:           "00000000-0000-0000-0000-000000000000"
ServiceName:            "Virtual Machines"
ResourceId:             "/subscriptions/.../resourceGroups/.../providers/Microsoft.Compute/virtualMachines/vm-name"
RegionId:               "eastus"
CommitmentDiscountType: "Reservation"
```

### GCP

```text
ProviderName:           "GCP"
BillingAccountId:       "012345-ABCDEF-678901"
SubAccountId:           "my-project-id"
ServiceName:            "Compute Engine"
ResourceId:             "projects/my-project/zones/us-central1-a/instances/instance-1"
RegionId:               "us-central1"
CommitmentDiscountType: "Committed Use Discount"
```

### Kubernetes (via Kubecost)

```text
ProviderName:           "Kubernetes"
BillingAccountId:       "cluster-name"
SubAccountId:           "namespace-name"
ServiceName:            "Kubernetes Workload"
ResourceId:             "namespace/pod-name"
ServiceCategory:        FOCUS_SERVICE_CATEGORY_COMPUTE
```

## Common Use Cases

### Cost Allocation by Team

Use Tags to allocate costs:

```go
builder.WithTags(map[string]string{
    "team":        "platform",
    "cost_center": "CC-1001",
    "environment": "production",
})
```

### Reserved Instance Analysis

Track commitment utilization:

```go
builder.WithCommitmentDiscount(
    pbc.FocusCommitmentDiscountCategory_FOCUS_COMMITMENT_DISCOUNT_CATEGORY_USAGE,
    "ri-123456789",
    "EC2 m5.large RI",
)
builder.WithCommitmentDiscountDetails(
    730.0, // hours
    pbc.FocusCommitmentDiscountStatus_FOCUS_COMMITMENT_DISCOUNT_STATUS_USED,
    "Standard Reserved Instance",
    "Hours",
)
```

### Multi-Currency Billing

Handle pricing in different currencies:

```go
builder.WithBillingPeriod(start, end, "EUR")
builder.WithPricingCurrency("USD")
builder.WithPricingCurrencyPrices(
    0.089,  // contracted unit price in USD
    65.0,   // effective cost in USD
    0.10,   // list unit price in USD
)
```

### Capacity Reservation Tracking

Monitor capacity reservation utilization:

```go
builder.WithCapacityReservation(
    "cr-0abc123456789def0",
    pbc.FocusCapacityReservationStatus_FOCUS_CAPACITY_RESERVATION_STATUS_USED,
)
```

## Query Patterns

### Total Cost by Service

```sql
SELECT ServiceName, SUM(BilledCost) as TotalCost
FROM focus_records
GROUP BY ServiceName
ORDER BY TotalCost DESC
```

### Reserved Instance Utilization

```sql
SELECT
    CommitmentDiscountId,
    CommitmentDiscountName,
    SUM(CASE WHEN CommitmentDiscountStatus = 'USED' THEN BilledCost ELSE 0 END) as UsedCost,
    SUM(CASE WHEN CommitmentDiscountStatus = 'UNUSED' THEN BilledCost ELSE 0 END) as WastedCost
FROM focus_records
WHERE CommitmentDiscountId IS NOT NULL
GROUP BY CommitmentDiscountId, CommitmentDiscountName
```

### Cost by Tag

```sql
SELECT
    Tags['team'] as Team,
    SUM(BilledCost) as TotalCost
FROM focus_records
GROUP BY Tags['team']
```

### Regional Cost Distribution

```sql
SELECT
    ProviderName,
    RegionId,
    SUM(BilledCost) as TotalCost
FROM focus_records
GROUP BY ProviderName, RegionId
ORDER BY ProviderName, TotalCost DESC
```
