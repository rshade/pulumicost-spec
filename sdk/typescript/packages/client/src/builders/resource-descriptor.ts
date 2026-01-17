import { ResourceDescriptor } from "../generated/finfocus/v1/costsource_pb.js";
import { Provider } from "../generated/finfocus/v1/enums_pb.js";

export class ResourceDescriptorBuilder {
  private descriptor: ResourceDescriptor;

  constructor() {
    this.descriptor = new ResourceDescriptor();
  }

  withProvider(provider: string): this {
    this.descriptor.provider = provider;
    return this;
  }

  withResourceType(type: string): this {
    this.descriptor.resourceType = type;
    return this;
  }

  withRegion(region: string): this {
    this.descriptor.region = region;
    return this;
  }

  withSku(sku: string): this {
    this.descriptor.sku = sku;
    return this;
  }

  withTags(tags: { [key: string]: string }): this {
    this.descriptor.tags = tags;
    return this;
  }

  withArn(arn: string): this {
    this.descriptor.resourceId = arn;
    return this;
  }

  withId(id: string): this {
      this.descriptor.resourceId = id;
      return this;
  }

  build(): ResourceDescriptor {
    return this.descriptor;
  }
}
