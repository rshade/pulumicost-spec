import { RecommendationFilter } from "../generated/finfocus/v1/costsource_pb.js";
import { RecommendationActionType, RecommendationPriority } from "../generated/finfocus/v1/enums_pb.js";
import { ValidationError } from "../errors/validation-error.js";

export class RecommendationFilterBuilder {
  private filter: RecommendationFilter;

  constructor() {
    this.filter = new RecommendationFilter();
  }

  forProvider(provider: string): this {
    if (!provider || provider.trim() === '') {
      throw new ValidationError("Provider cannot be empty", "provider");
    }
    this.filter.provider = provider;
    return this;
  }

  forRegion(region: string): this {
    if (!region || region.trim() === '') {
      throw new ValidationError("Region cannot be empty", "region");
    }
    this.filter.region = region;
    return this;
  }

  forResourceType(resourceType: string): this {
    if (!resourceType || resourceType.trim() === '') {
      throw new ValidationError("Resource type cannot be empty", "resourceType");
    }
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
    if (amount < 0) {
      throw new ValidationError("Minimum estimated savings cannot be negative", "minEstimatedSavings");
    }
    this.filter.minEstimatedSavings = amount;
    return this;
  }

  build(): RecommendationFilter {
    // Create a copy to prevent mutation issues when reusing the builder
    const copy = new RecommendationFilter();
    copy.provider = this.filter.provider;
    copy.region = this.filter.region;
    copy.resourceType = this.filter.resourceType;
    copy.actionType = this.filter.actionType;
    copy.priority = this.filter.priority;
    copy.minEstimatedSavings = this.filter.minEstimatedSavings;
    return copy;
  }
}
