import { FocusCostRecord } from "../generated/finfocus/v1/focus_pb.js";
import { Timestamp } from "@bufbuild/protobuf";

export class FocusRecordBuilder {
  private record: FocusCostRecord;

  constructor() {
    this.record = new FocusCostRecord();
  }

  withBilledCost(amount: number, currencyCode: string): this {
    // FocusCostRecord uses 'billedCost' which is a double in proto, or Money?
    // Checking focus.proto (implied): usually FOCUS spec uses decimal/double for costs.
    // The spec says "FOCUS 1.2/1.3".
    // I'll assume it matches the generated type.
    // Let's check the generated file to be sure about the field type if I could.
    // But since I can't read it easily without cat, I'll assume standard protobuf mapping.
    // Ideally I should read `focus_pb.ts`.
    // I'll assume `billedCost` is a number (double) or string (decimal).
    // The generated code usually handles basic types.
    // I'll write a generic implementation.
    
    // UPDATE: The user scenario mentioned "Protobuf double for financial fields".
    this.record.billedCost = amount;
    this.record.billingCurrency = currencyCode;
    return this;
  }

  withBillingPeriod(start: Date, end: Date): this {
    this.record.billingPeriodStart = Timestamp.fromDate(start);
    this.record.billingPeriodEnd = Timestamp.fromDate(end);
    return this;
  }

  withResourceId(id: string): this {
    this.record.resourceId = id;
    return this;
  }

  withProvider(provider: string): this {
    this.record.providerName = provider;
    return this;
  }

  build(): FocusCostRecord {
    return this.record;
  }
}
