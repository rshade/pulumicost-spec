import { RecommendationFilter, RecommendationCategory, RecommendationSortBy, SortOrder } from "../generated/finfocus/v1/costsource_pb.js";
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

  withCategory(category: RecommendationCategory): this {
    this.filter.category = category;
    return this;
  }

  withSku(sku: string): this {
    if (!sku || sku.trim() === '') {
      throw new ValidationError("SKU cannot be empty", "sku");
    }
    this.filter.sku = sku;
    return this;
  }

  withTags(tags: Record<string, string>): this {
    if (!tags) {
      throw new ValidationError("Tags cannot be null or undefined", "tags");
    }
    for (const [key, value] of Object.entries(tags)) {
      this.filter.tags[key] = value;
    }
    return this;
  }

  withSource(source: string): this {
    if (!source || source.trim() === '') {
      throw new ValidationError("Source cannot be empty", "source");
    }
    this.filter.source = source;
    return this;
  }

  withAccountId(accountId: string): this {
    if (!accountId || accountId.trim() === '') {
      throw new ValidationError("Account ID cannot be empty", "accountId");
    }
    this.filter.accountId = accountId;
    return this;
  }

  withSortBy(sortBy: RecommendationSortBy): this {
    this.filter.sortBy = sortBy;
    return this;
  }

  withSortOrder(sortOrder: SortOrder): this {
    this.filter.sortOrder = sortOrder;
    return this;
  }

  withMinConfidenceScore(score: number): this {
    if (score < 0 || score > 1) {
      throw new ValidationError("Confidence score must be between 0.0 and 1.0", "minConfidenceScore");
    }
    this.filter.minConfidenceScore = score;
    return this;
  }

  withMaxAgeDays(days: number): this {
    if (days < 0 || !Number.isInteger(days)) {
      throw new ValidationError("Max age days must be a non-negative integer", "maxAgeDays");
    }
    this.filter.maxAgeDays = days;
    return this;
  }

  withResourceId(resourceId: string): this {
    if (!resourceId || resourceId.trim() === '') {
      throw new ValidationError("Resource ID cannot be empty", "resourceId");
    }
    this.filter.resourceId = resourceId;
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
    copy.category = this.filter.category;
    copy.sku = this.filter.sku;
    for (const [key, value] of Object.entries(this.filter.tags)) {
      copy.tags[key] = value;
    }
    copy.source = this.filter.source;
    copy.accountId = this.filter.accountId;
    copy.sortBy = this.filter.sortBy;
    copy.sortOrder = this.filter.sortOrder;
    copy.minConfidenceScore = this.filter.minConfidenceScore;
    copy.maxAgeDays = this.filter.maxAgeDays;
    copy.resourceId = this.filter.resourceId;
    return copy;
  }
}
