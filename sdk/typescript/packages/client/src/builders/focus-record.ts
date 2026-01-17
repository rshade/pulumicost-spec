import { FocusCostRecord } from "../generated/finfocus/v1/focus_pb.js";
import { Timestamp } from "@bufbuild/protobuf";
import { ValidationError } from "../errors/validation-error.js";

export class FocusRecordBuilder {
  private record: FocusCostRecord;

  constructor() {
    this.record = new FocusCostRecord();
  }

  withBilledCost(amount: number, currencyCode: string): this {
    if (!currencyCode || currencyCode.trim() === '') {
      throw new ValidationError("Currency code cannot be empty", "currencyCode");
    }
    this.record.billedCost = amount;
    this.record.billingCurrency = currencyCode;
    return this;
  }

  withBillingPeriod(start: Date, end: Date): this {
    if (end < start) {
      throw new ValidationError("Billing period end cannot be before start", "billingPeriod");
    }
    this.record.billingPeriodStart = Timestamp.fromDate(start);
    this.record.billingPeriodEnd = Timestamp.fromDate(end);
    return this;
  }

  withResourceId(id: string): this {
    if (!id || id.trim() === '') {
      throw new ValidationError("Resource ID cannot be empty", "resourceId");
    }
    this.record.resourceId = id;
    return this;
  }

  withProvider(provider: string): this {
    if (!provider || provider.trim() === '') {
      throw new ValidationError("Provider cannot be empty", "provider");
    }
    this.record.providerName = provider;
    return this;
  }

  build(): FocusCostRecord {
    // Create a copy to prevent mutation issues when reusing the builder
    const copy = new FocusCostRecord();
    copy.billedCost = this.record.billedCost;
    copy.billingCurrency = this.record.billingCurrency;
    copy.billingPeriodStart = this.record.billingPeriodStart;
    copy.billingPeriodEnd = this.record.billingPeriodEnd;
    copy.resourceId = this.record.resourceId;
    copy.providerName = this.record.providerName;
    return copy;
  }
}
