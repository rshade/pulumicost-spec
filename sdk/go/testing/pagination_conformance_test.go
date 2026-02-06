// Copyright 2024 The FinFocus Authors
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

package testing_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rshade/finfocus-spec/sdk/go/pluginsdk"
	pbc "github.com/rshade/finfocus-spec/sdk/go/proto/finfocus/v1"
	plugintesting "github.com/rshade/finfocus-spec/sdk/go/testing"
)

// TestPaginationConformance_FirstPage verifies the first page returns the
// correct record count and a non-empty next token when more records exist.
//
// IMPORTANT: Do NOT use t.Parallel() in subtests below.
// They share a single gRPC harness that will be closed on function exit.
func TestPaginationConformance_FirstPage(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(200)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(200)
	resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
	})
	require.NoError(t, err)
	require.Len(t, resp.GetResults(), 50, "first page should have exactly page_size records")
	require.NotEmpty(t, resp.GetNextPageToken(), "first page should have a next page token")
	require.Equal(t, int32(200), resp.GetTotalCount(), "total_count should match dataset size")
}

// TestPaginationConformance_MiddlePage verifies a middle page returns the
// correct records using a continuation token from a previous page.
func TestPaginationConformance_MiddlePage(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(200)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(200)

	// Get first page to obtain token
	resp1, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.GetNextPageToken())

	// Get middle page using token
	resp2, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
		PageToken:  resp1.GetNextPageToken(),
	})
	require.NoError(t, err)
	require.Len(t, resp2.GetResults(), 50, "middle page should have page_size records")
	require.NotEmpty(t, resp2.GetNextPageToken(), "middle page should have next token")
	require.Equal(t, resp1.GetTotalCount(), resp2.GetTotalCount(), "total_count should be consistent")
}

// TestPaginationConformance_LastPage verifies the last page returns
// remaining records and an empty next_page_token.
func TestPaginationConformance_LastPage(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(120)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(120)

	// Iterate to last page
	const maxPages = 100
	var lastResp *pbc.GetActualCostResponse
	pageToken := ""
	pageCount := 0
	for {
		resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
			ResourceId: "i-abc123",
			Start:      start,
			End:        end,
			PageSize:   50,
			PageToken:  pageToken,
		})
		require.NoError(t, err)
		pageCount++
		require.LessOrEqual(t, pageCount, maxPages, "pagination loop exceeded maxPages safety cap")
		lastResp = resp
		if resp.GetNextPageToken() == "" {
			break
		}
		pageToken = resp.GetNextPageToken()
	}

	require.Equal(t, 3, pageCount, "120 records / 50 per page = 3 pages (50+50+20)")
	require.Len(t, lastResp.GetResults(), 20, "last page should have remaining 20 records")
	require.Empty(t, lastResp.GetNextPageToken(), "last page should have empty next token")
}

// TestPaginationConformance_EmptyResult verifies that an offset beyond the
// data returns empty results with an empty next token.
func TestPaginationConformance_EmptyResult(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(50)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(50)

	// Get first (and only) page
	resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
	})
	require.NoError(t, err)
	require.Empty(t, resp.GetNextPageToken(), "single page should have empty next token")

	// Use a token pointing beyond the data
	beyondToken := pluginsdk.EncodePageToken(999)
	resp2, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
		PageToken:  beyondToken,
	})
	require.NoError(t, err)
	require.Empty(t, resp2.GetResults(), "beyond-range token should return empty results")
	require.Empty(t, resp2.GetNextPageToken(), "beyond-range should have empty next token")
}

// TestPaginationConformance_InvalidToken verifies that a malformed token
// returns a gRPC InvalidArgument error.
func TestPaginationConformance_InvalidToken(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(50)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(50)

	_, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
		PageToken:  "not-a-valid-base64-token!!!",
	})
	require.Error(t, err, "invalid token should return an error")
	require.Equal(t, codes.InvalidArgument, status.Code(err),
		"invalid token should return InvalidArgument")
}

// TestPaginationConformance_OversizedPage verifies the response never contains
// more records than the effective page size.
func TestPaginationConformance_OversizedPage(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(500)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(500)

	pageSizes := []int32{10, 50, 100, 200}

	for _, ps := range pageSizes {
		pageSize := ps // capture for closure
		t.Run(fmt.Sprintf("pageSize=%d", pageSize), func(t *testing.T) {
			resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
				ResourceId: "i-abc123",
				Start:      start,
				End:        end,
				PageSize:   pageSize,
			})
			require.NoError(t, err)
			require.LessOrEqual(t, int32(len(resp.GetResults())), pageSize,
				"page should not exceed requested page_size=%d", pageSize)
		})
	}
}

// TestPaginationConformance_DefaultPageSize verifies that page_size=0 on a
// continuation request (with a non-empty page_token) uses DefaultPageSize (50).
// Note: When both page_size=0 AND page_token="" (proto3 defaults), the mock
// returns all results for backward compatibility. The default page size only
// takes effect when pagination is explicitly started.
func TestPaginationConformance_DefaultPageSize(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(200)
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(200)

	// First, get page 1 with explicit page_size to start pagination
	resp1, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   50,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.GetNextPageToken())

	// Now request page 2 with page_size=0 and a valid page_token.
	// The mock should use its default page size (50) for this request.
	resp2, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   0,
		PageToken:  resp1.GetNextPageToken(),
	})
	require.NoError(t, err)
	require.Len(t, resp2.GetResults(), pluginsdk.DefaultPageSize,
		"page_size=0 with page_token should use DefaultPageSize=%d", pluginsdk.DefaultPageSize)
}

// TestPaginationConformance_BackwardCompat verifies that a legacy plugin
// (no pagination awareness) passes basic conformance with pagination fields
// present in the request.
func TestPaginationConformance_BackwardCompat(t *testing.T) {
	plugin := plugintesting.NewMockPlugin()
	plugin.SetActualCostDataPoints(24) // Small dataset like legacy plugins
	harness := plugintesting.NewTestHarness(plugin)
	harness.Start(t)
	defer harness.Stop()

	client := harness.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start, end := plugintesting.CreateTimeRange(24)

	// Request with pagination params against a plugin with small dataset
	// Plugin should return all results in one page
	resp, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   100, // Larger than dataset
	})
	require.NoError(t, err)
	require.Len(t, resp.GetResults(), 24, "should return all records when dataset < page_size")
	require.Empty(t, resp.GetNextPageToken(), "single page should have empty next token")
	require.Equal(t, int32(24), resp.GetTotalCount())

	// Also verify with zero page_size (backward compat mode)
	resp2, err := client.GetActualCost(ctx, &pbc.GetActualCostRequest{
		ResourceId: "i-abc123",
		Start:      start,
		End:        end,
		PageSize:   0,
		PageToken:  "",
	})
	require.NoError(t, err)
	require.Len(t, resp2.GetResults(), 24, "backward compat: should return all records")
}
