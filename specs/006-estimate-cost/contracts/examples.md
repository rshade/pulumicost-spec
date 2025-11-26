# EstimateCost RPC Examples

## Request/Response Examples

### AWS EC2 Instance

**Request**:

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.micro",
    "region": "us-east-1",
    "tenancy": "default"
  }
}
```

**Response**:

```json
{
  "currency": "USD",
  "cost_monthly": "7.30"
}
```

### Azure Virtual Machine

**Request**:

```json
{
  "resource_type": "azure:compute/virtualMachine:VirtualMachine",
  "attributes": {
    "vmSize": "Standard_B1s",
    "location": "eastus",
    "osType": "Linux"
  }
}
```

**Response**:

```json
{
  "currency": "USD",
  "cost_monthly": "7.59"
}
```

### GCP Compute Instance

**Request**:

```json
{
  "resource_type": "gcp:compute/instance:Instance",
  "attributes": {
    "machineType": "e2-micro",
    "zone": "us-central1-a"
  }
}
```

**Response**:

```json
{
  "currency": "USD",
  "cost_monthly": "6.11"
}
```

### AWS S3 Bucket (Free Tier)

**Request**:

```json
{
  "resource_type": "aws:s3/bucket:Bucket",
  "attributes": {
    "region": "us-east-1"
  }
}
```

**Response** (Zero cost for bucket creation):

```json
{
  "currency": "USD",
  "cost_monthly": "0.00"
}
```

## Error Examples

### Invalid Resource Type Format

**Request**:

```json
{
  "resource_type": "aws:ec2:Instance",
  "attributes": {}
}
```

**Error** (InvalidArgument):

```json
{
  "code": "INVALID_ARGUMENT",
  "message": "resource_type must follow provider:module/resource:Type format, got: aws:ec2:Instance"
}
```

### Unsupported Resource Type

**Request**:

```json
{
  "resource_type": "aws:lambda/function:Function",
  "attributes": {}
}
```

**Error** (NotFound):

```json
{
  "code": "NOT_FOUND",
  "message": "resource type aws:lambda/function:Function is not supported by this plugin"
}
```

### Missing Required Attributes

**Request**:

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {}
}
```

**Error** (InvalidArgument):

```json
{
  "code": "INVALID_ARGUMENT",
  "message": "missing required attributes for aws:ec2/instance:Instance: [instanceType, region]"
}
```

### Pricing Source Unavailable

**Request**:

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.micro",
    "region": "us-east-1"
  }
}
```

**Error** (Unavailable):

```json
{
  "code": "UNAVAILABLE",
  "message": "pricing source unavailable: AWS Pricing API timeout after 30s"
}
```

## Configuration Comparison Use Case

### Comparing Instance Sizes

**Request 1** (t3.micro):

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.micro",
    "region": "us-east-1"
  }
}
```

**Response 1**:

```json
{
  "currency": "USD",
  "cost_monthly": "7.30"
}
```

**Request 2** (t3.small):

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.small",
    "region": "us-east-1"
  }
}
```

**Response 2**:

```json
{
  "currency": "USD",
  "cost_monthly": "14.60"
}
```

**Request 3** (t3.large):

```json
{
  "resource_type": "aws:ec2/instance:Instance",
  "attributes": {
    "instanceType": "t3.large",
    "region": "us-east-1"
  }
}
```

**Response 3**:

```json
{
  "currency": "USD",
  "cost_monthly": "58.40"
}
```

**Analysis**: Developer can compare costs before selecting instance size, enabling informed
infrastructure decisions based on budget constraints.
