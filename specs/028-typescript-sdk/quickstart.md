# Quick Start

## Installation

```bash
npm install finfocus-client
```

## Basic Usage (Browser)

```typescript
import { CostSourceClient, ResourceDescriptorBuilder } from 'finfocus-client';

// 1. Initialize client with plugin URL
const client = new CostSourceClient({
  baseUrl: 'https://plugin-aws.example.com',
});

// 2. Define resource
const resource = new ResourceDescriptorBuilder()
  .withProvider('AWS')
  .withResourceType('EC2')
  .withRegion('us-east-1')
  .build();

// 3. Fetch cost
try {
  const cost = await client.getActualCost({ resource });
  console.log(`Cost: ${cost.total.amount} ${cost.total.currency}`);
} catch (err) {
  console.error('Failed to fetch cost:', err);
}
```

## REST Middleware (Express)

```typescript
import express from 'express';
import { createExpressMiddleware } from 'finfocus-client/express';

const app = express();

// Mounts all 22 RPC endpoints at /api/finfocus
app.use('/api/finfocus', createExpressMiddleware({
  targetPluginUrl: 'https://plugin-aws.example.com'
}));

app.listen(3000);
```
