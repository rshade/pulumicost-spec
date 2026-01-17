import { CostSourceClient } from "../clients/cost-source.js";
import { GetRecommendationsRequest, Recommendation } from "../generated/finfocus/v1/costsource_pb.js";

export async function* recommendationsIterator(
  client: CostSourceClient,
  baseRequest: GetRecommendationsRequest
): AsyncGenerator<Recommendation, void, unknown> {
  // Clone request to avoid mutating the original
  const request = new GetRecommendationsRequest(baseRequest);
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
