import { createPromiseClient, PromiseClient, Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { CostSourceService } from "../generated/finfocus/v1/costsource_connect.js";
import { ValidationError } from "../errors/validation-error.js";
import {
  GetActualCostRequest,
  GetActualCostResponse,
  GetProjectedCostRequest,
  GetProjectedCostResponse,
  GetRecommendationsRequest,
  GetRecommendationsResponse,
  GetPricingSpecRequest,
  GetPricingSpecResponse,
  EstimateCostRequest,
  EstimateCostResponse,
  DismissRecommendationRequest,
  DismissRecommendationResponse,
  GetPluginInfoRequest,
  GetPluginInfoResponse,
  DryRunRequest,
  DryRunResponse,
  NameRequest,
  NameResponse,
  SupportsRequest,
  SupportsResponse
} from "../generated/finfocus/v1/costsource_pb.js";
import {
  GetBudgetsRequest,
  GetBudgetsResponse
} from "../generated/finfocus/v1/budget_pb.js";

export interface CostSourceClientConfig {
  baseUrl: string;
  transport?: Transport;
}

export class CostSourceClient {
  private client: PromiseClient<typeof CostSourceService>;

  constructor(config: CostSourceClientConfig) {
    const transport = config.transport || createConnectTransport({
      baseUrl: config.baseUrl,
      useBinaryFormat: false,
    });
    this.client = createPromiseClient(CostSourceService, transport);
  }

  async name(req: NameRequest = new NameRequest()): Promise<NameResponse> {
    return this.client.name(req);
  }

  async supports(req: SupportsRequest = new SupportsRequest()): Promise<SupportsResponse> {
    return this.client.supports(req);
  }

  async getActualCost(req: GetActualCostRequest): Promise<GetActualCostResponse> {
    if (!req.resourceId && !req.arn) {
        throw new ValidationError("Resource ID or ARN is required", "resourceId");
    }
    return this.client.getActualCost(req);
  }

  async getProjectedCost(req: GetProjectedCostRequest): Promise<GetProjectedCostResponse> {
    if (!req.resource) throw new ValidationError("Resource is required", "resource");
    return this.client.getProjectedCost(req);
  }

  async getPricingSpec(req: GetPricingSpecRequest = new GetPricingSpecRequest()): Promise<GetPricingSpecResponse> {
    return this.client.getPricingSpec(req);
  }

  async estimateCost(req: EstimateCostRequest): Promise<EstimateCostResponse> {
    return this.client.estimateCost(req);
  }

  async getRecommendations(req: GetRecommendationsRequest = new GetRecommendationsRequest()): Promise<GetRecommendationsResponse> {
    return this.client.getRecommendations(req);
  }

  async dismissRecommendation(req: DismissRecommendationRequest): Promise<DismissRecommendationResponse> {
    if (!req.recommendationId) throw new ValidationError("Recommendation ID is required", "recommendationId");
    return this.client.dismissRecommendation(req);
  }

  async getBudgets(req: GetBudgetsRequest = new GetBudgetsRequest()): Promise<GetBudgetsResponse> {
    return this.client.getBudgets(req);
  }

  async getPluginInfo(req: GetPluginInfoRequest = new GetPluginInfoRequest()): Promise<GetPluginInfoResponse> {
    return this.client.getPluginInfo(req);
  }

  async dryRun(req: DryRunRequest = new DryRunRequest()): Promise<DryRunResponse> {
    return this.client.dryRun(req);
  }
}