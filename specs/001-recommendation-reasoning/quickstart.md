# Quickstart: Using Recommendation Reasons

## Overview

The `RecommendationReason` enum provides a standardized way to understand _why_ a
recommendation was made, regardless of the underlying cloud provider.

## Usage in Go SDK

### Consuming Recommendations

```go
package main

import (
    "fmt"
    pb "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
)

func processRecommendation(rec *pb.Recommendation) {
    fmt.Printf("Recommendation: %s\n", rec.Description)

    switch rec.PrimaryReason {
    case pb.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED:
        fmt.Println("Action: Downsize resource")
    case pb.RecommendationReason_RECOMMENDATION_REASON_IDLE:
        fmt.Println("Action: Terminate resource")
    case pb.RecommendationReason_RECOMMENDATION_REASON_OBSOLETE_GENERATION:
        fmt.Println("Action: Upgrade generation")
    default:
        fmt.Println("Action: Manual review needed")
    }

    if len(rec.SecondaryReasons) > 0 {
        fmt.Println("Contributing factors:")
        for _, reason := range rec.SecondaryReasons {
            fmt.Printf("- %s\n", reason)
        }
    }
}
```

### Creating Recommendations (Plugin Developer)

```go
rec := &pb.Recommendation{
    Id:          "rec-123",
    Description: "Downsize instance to t3.micro",
    PrimaryReason: pb.RecommendationReason_RECOMMENDATION_REASON_OVER_PROVISIONED,
    SecondaryReasons: []pb.RecommendationReason{
        pb.RecommendationReason_RECOMMENDATION_REASON_IDLE,
    },
    // ... other fields
}
```

```protobuf

```
