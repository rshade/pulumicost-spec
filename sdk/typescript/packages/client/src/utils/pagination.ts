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

export async function* recommendationsIterator(
  client: CostSourceClient,
  baseRequest: GetRecommendationsRequest
): AsyncGenerator<Recommendation, void, unknown> {
  // Clone request to avoid mutating the original
  const request = clone(GetRecommendationsRequestSchema, baseRequest);
  // Preserve caller-supplied pageToken for resuming pagination
  let nextPageToken = request.pageToken ?? "";

  do {
    request.pageToken = nextPageToken;
    const response = await client.getRecommendations(request);

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
  // Preserve caller-supplied pageToken for resuming pagination
  let nextPageToken = request.pageToken ?? "";

  do {
    request.pageToken = nextPageToken;
    const response = await client.getActualCost(request);

    for (const result of response.results) {
      yield result;
    }

    nextPageToken = response.nextPageToken;
  } while (nextPageToken);
}
