import { clone } from "@bufbuild/protobuf";
import { CostSourceClient } from "../clients/cost-source.js";
import type {
  ActualCostResult,
  GetActualCostRequest,
  GetRecommendationsRequest,
  Recommendation,
} from "../generated/finfocus/v1/costsource_pb.js";
import {
  GetActualCostRequestSchema,
  GetRecommendationsRequestSchema,
} from "../generated/finfocus/v1/costsource_pb.js";

/**
 * Maximum number of consecutive empty pages before the iterator aborts.
 * Prevents infinite loops when a buggy server returns empty pages with
 * continuation tokens. Matches the Go SDK maxEmptyPages constant.
 */
const MAX_EMPTY_PAGES = 10;

/** Default page size used when the caller does not specify one. */
const DEFAULT_PAGE_SIZE = 50;

export async function* recommendationsIterator(
  client: CostSourceClient,
  baseRequest: GetRecommendationsRequest
): AsyncGenerator<Recommendation, void, unknown> {
  // Clone request to avoid mutating the original
  const request = clone(GetRecommendationsRequestSchema, baseRequest);
  // Preserve caller-supplied pageToken for resuming pagination
  let nextPageToken = request.pageToken ?? "";
  let emptyPageCount = 0;

  do {
    request.pageToken = nextPageToken;
    const response = await client.getRecommendations(request);

    if (response.recommendations.length === 0 && response.nextPageToken) {
      emptyPageCount++;
      if (emptyPageCount >= MAX_EMPTY_PAGES) {
        throw new Error(
          `Pagination safety: exceeded ${MAX_EMPTY_PAGES} consecutive empty pages`
        );
      }
    } else {
      emptyPageCount = 0;
    }

    for (const rec of response.recommendations) {
      yield rec;
    }

    nextPageToken = response.nextPageToken;
  } while (nextPageToken);
}

export async function* actualCostIterator(
  client: CostSourceClient,
  baseRequest: GetActualCostRequest
): AsyncGenerator<ActualCostResult, void, unknown> {
  // Clone request to avoid mutating the original
  const request = clone(GetActualCostRequestSchema, baseRequest);
  // Default pageSize to avoid triggering legacy "return all" server behaviour
  if (!request.pageSize) {
    request.pageSize = DEFAULT_PAGE_SIZE;
  }
  // Preserve caller-supplied pageToken for resuming pagination
  let nextPageToken = request.pageToken ?? "";
  let emptyPageCount = 0;

  do {
    request.pageToken = nextPageToken;
    const response = await client.getActualCost(request);

    if (response.results.length === 0 && response.nextPageToken) {
      emptyPageCount++;
      if (emptyPageCount >= MAX_EMPTY_PAGES) {
        throw new Error(
          `Pagination safety: exceeded ${MAX_EMPTY_PAGES} consecutive empty pages`
        );
      }
    } else {
      emptyPageCount = 0;
    }

    for (const result of response.results) {
      yield result;
    }

    nextPageToken = response.nextPageToken;
  } while (nextPageToken);
}
