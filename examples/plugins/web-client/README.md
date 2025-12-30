# PulumiCost Web Client Example

This example demonstrates how to call PulumiCost plugins from a web browser using the Connect
protocol.

## Overview

The Connect protocol enables browser-based applications to communicate with PulumiCost plugins using
simple JSON over HTTP - no special client libraries required.

## Key Benefits

1. **No Special Client**: Uses standard `fetch()` API
2. **Human-Readable**: JSON payloads are easy to debug
3. **HTTP/1.1 Compatible**: Works without HTTP/2
4. **CORS Support**: Built-in cross-origin support when configured

## Running the Example

### 1. Start a Plugin Server

First, start a PulumiCost plugin server with web support enabled:

```go
package main

import (
    "context"

    "github.com/rshade/pulumicost-spec/sdk/go/pluginsdk"
)

func main() {
    plugin := &YourPlugin{}

    err := pluginsdk.Serve(context.Background(), pluginsdk.ServeConfig{
        Plugin: plugin,
        Port:   8080,
        Web: pluginsdk.WebConfig{
            Enabled:              true,
            AllowedOrigins:       []string{"*"}, // Allow all origins for demo
            EnableHealthEndpoint: true,
        },
    })
    if err != nil {
        panic(err)
    }
}
```

### 2. Open the HTML File

Open `index.html` in a web browser. You can serve it locally with:

```bash
# Using Python
python3 -m http.server 3000

# Using Node.js
npx serve .
```

Then navigate to `http://localhost:3000`.

### 3. Configure Server URL

Enter your plugin server URL (e.g., `http://localhost:8080`) in the Server Configuration section.

## How It Works

### Connect Protocol

The Connect protocol sends JSON POST requests to predictable URL paths:

```text
POST /pulumicost.v1.CostSourceService/Name
Content-Type: application/json

{}
```

Response:

```json
{
  "name": "my-cost-plugin"
}
```

### Example: Estimate Cost

```javascript
const response = await fetch('http://localhost:8080/pulumicost.v1.CostSourceService/EstimateCost', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        resource_type: 'aws:ec2/instance:Instance',
        attributes: {
            instance_type: 't3.micro',
            region: 'us-east-1',
        },
    }),
});

const data = await response.json();
console.log(`Monthly cost: ${data.currency} ${data.cost_monthly}`);
```

### Example: Using curl

```bash
# Get plugin name
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/Name \
  -H "Content-Type: application/json" \
  -d '{}'

# Estimate cost
curl -X POST http://localhost:8080/pulumicost.v1.CostSourceService/EstimateCost \
  -H "Content-Type: application/json" \
  -d '{
    "resource_type": "aws:ec2/instance:Instance",
    "attributes": {
      "instance_type": "t3.micro"
    }
  }'
```

## CORS Configuration

For production use, configure specific allowed origins:

```go
pluginsdk.WebConfig{
    Enabled:          true,
    AllowedOrigins:   []string{"https://your-app.com"},
    AllowCredentials: true,
}
```

## Error Handling

Connect protocol errors include a `code` and `message`:

```json
{
  "code": "not_found",
  "message": "resource type not supported"
}
```

Map these to your application's error handling:

```javascript
if (!response.ok) {
    const error = await response.json();
    throw new Error(`${error.code}: ${error.message}`);
}
```

## TypeScript Support

For TypeScript applications, you can use the generated Connect-ES types (optional):

```bash
npm install @connectrpc/connect @connectrpc/connect-web
```

Or simply define your own types based on the protobuf definitions.

## See Also

- [Connect Protocol Documentation](https://connectrpc.com/docs/protocol)
- [PulumiCost SDK Documentation](../../../sdk/go/pluginsdk/README.md)
