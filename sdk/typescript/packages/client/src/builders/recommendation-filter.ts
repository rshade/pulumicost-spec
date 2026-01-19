import { create, clone } from "@bufbuild/protobuf";
import {
  RecommendationFilter,
  RecommendationFilterSchema,
  RecommendationCategory,
  RecommendationSortBy,
  SortOrder,
  RecommendationActionType,
  RecommendationPriority
} from "../generated/finfocus/v1/costsource_pb.js";
import { ValidationError } from "../errors/validation-error.js";

export class RecommendationFilterBuilder {
  private filter: RecommendationFilter;

  constructor() {
    this.filter = create(RecommendationFilterSchema);
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

  /**
   * Sets tags, replacing any existing tags.
   * Use addTags() to merge additional tags without clearing existing ones.
   */
  withTags(tags: Record<string, string>): this {
    if (!tags) {
      throw new ValidationError("Tags cannot be null or undefined", "tags");
    }
    // Clear existing tags and set new ones (replacement semantics)
    const existingKeys = Object.keys(this.filter.tags);
    for (const key of existingKeys) {
      delete this.filter.tags[key];
    }
    for (const [key, value] of Object.entries(tags)) {
      this.filter.tags[key] = value;
    }
    return this;
  }

  /**
   * Adds tags to the existing tags without clearing them.
   * Use withTags() to replace all tags.
   */
  addTags(tags: Record<string, string>): this {
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
    if (!Number.isFinite(score) || score < 0 || score > 1) {
      throw new ValidationError("Confidence score must be between 0.0 and 1.0", "minConfidenceScore");
    }
    this.filter.minConfidenceScore = score;
    return this;
  }

  withMaxAgeDays(days: number): this {
    // Check integer first to reject NaN, Infinity, and floats before range check
    if (!Number.isInteger(days) || days < 0) {
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
    return clone(RecommendationFilterSchema, this.filter);
  }
}
