# Data Model: FinOps FOCUS 1.2 Integration

## Enums (`proto/pulumicost/v1/enums.proto`)

### FocusServiceCategory
| Name | Value | Description |
| :--- | :--- | :--- |
| `SERVICE_CATEGORY_UNSPECIFIED` | 0 | Default |
| `SERVICE_CATEGORY_COMPUTE` | 1 | Virtual Machines, Containers, Lambda |
| `SERVICE_CATEGORY_STORAGE` | 2 | Object Storage, Block Storage |
| `SERVICE_CATEGORY_NETWORK` | 3 | VPC, Load Balancers, Data Transfer |
| `SERVICE_CATEGORY_DATABASE` | 4 | RDS, DynamoDB, SQL |
| `SERVICE_CATEGORY_ANALYTICS` | 5 | BigQuery, Redshift |
| `SERVICE_CATEGORY_MACHINE_LEARNING` | 6 | SageMaker, Vertex AI |
| `SERVICE_CATEGORY_MANAGEMENT` | 7 | CloudWatch, Stackdriver |
| `SERVICE_CATEGORY_SECURITY` | 8 | IAM, KMS |
| `SERVICE_CATEGORY_DEVELOPER_TOOLS` | 9 | CodeBuild, DevOps |
| `SERVICE_CATEGORY_OTHER` | 10 | Fallback |

### FocusChargeCategory
| Name | Value | Description |
| :--- | :--- | :--- |
| `CHARGE_CATEGORY_UNSPECIFIED` | 0 | Default |
| `CHARGE_CATEGORY_USAGE` | 1 | Consumption-based |
| `CHARGE_CATEGORY_PURCHASE` | 2 | Upfront fees |
| `CHARGE_CATEGORY_CREDIT` | 3 | Discounts/Vouchers |
| `CHARGE_CATEGORY_TAX` | 4 | Taxes |
| `CHARGE_CATEGORY_REFUND` | 5 | Reimbursements |
| `CHARGE_CATEGORY_ADJUSTMENT` | 6 | Other adjustments (Superset of Refund sometimes) |

### FocusPricingCategory
| Name | Value | Description |
| :--- | :--- | :--- |
| `PRICING_CATEGORY_UNSPECIFIED` | 0 | Default |
| `PRICING_CATEGORY_STANDARD` | 1 | On-Demand / List |
| `PRICING_CATEGORY_COMMITTED` | 2 | Reserved / Savings Plan |
| `PRICING_CATEGORY_DYNAMIC` | 3 | Spot / Preemptible |
| `PRICING_CATEGORY_OTHER` | 4 | Fallback |

## Messages (`proto/pulumicost/v1/focus.proto`)

### FocusCostRecord
| Field Name | Type | ID | Description |
| :--- | :--- | :--- | :--- |
| `provider_name` | string | 1 | e.g., "AWS" |
| `billing_account_id` | string | 2 | Mandatory |
| `billing_account_name` | string | 3 | |
| `charge_period_start` | google.protobuf.Timestamp | 4 | Mandatory |
| `charge_period_end` | google.protobuf.Timestamp | 5 | Mandatory |
| `service_category` | FocusServiceCategory | 6 | Mandatory |
| `service_name` | string | 7 | Mandatory |
| `charge_category` | FocusChargeCategory | 8 | Mandatory |
| `pricing_category` | FocusPricingCategory | 9 | Mandatory |
| `region_id` | string | 10 | |
| `region_name` | string | 11 | |
| `resource_id` | string | 12 | |
| `resource_name` | string | 13 | |
| `sku_id` | string | 14 | |
| `billed_cost` | double | 15 | Mandatory |
| `list_cost` | double | 16 | |
| `effective_cost` | double | 17 | |
| `currency` | string | 18 | Mandatory (if cost > 0) |
| `invoice_id` | string | 19 | |
| `usage_quantity` | double | 20 | |
| `usage_unit` | string | 21 | |
| `tags` | map<string, string> | 22 | Resource Tags |
| `extended_columns` | map<string, string> | 23 | Backpack |
