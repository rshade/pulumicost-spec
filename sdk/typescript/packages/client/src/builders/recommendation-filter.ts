import { RecommendationFilter } from "../generated/finfocus/v1/costsource_pb.js";
import { RecommendationActionType, RecommendationPriority } from "../generated/finfocus/v1/enums_pb.js";

export class RecommendationFilterBuilder {
  private filter: RecommendationFilter;

  constructor() {
    this.filter = new RecommendationFilter();
  }

  forProvider(provider: string): this {
    this.filter.provider = provider;
    return this;
  }

  forRegion(region: string): this {
    this.filter.region = region;
    return this;
  }

  forResourceType(resourceType: string): this {
    this.filter.resourceType = resourceType;
    return this;
  }

  withActionType(actionType: RecommendationActionType): this {
    this.filter.actionType = actionType;
    return this;
  }

  withPriority(priority: RecommendationPriority): this {
    this.filter.priority = priority;
    return this;
  }

  minEstimatedSavings(amount: number): this {
    this.filter.minEstimatedSavings = amount;
    return this;
  }

  build(): RecommendationFilter {
    return this.filter;
  }
}
