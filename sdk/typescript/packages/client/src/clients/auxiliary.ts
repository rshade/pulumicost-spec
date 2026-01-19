import { createClient, Client, Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { create } from "@bufbuild/protobuf";
import {
  PluginRegistryService,
  DiscoverPluginsRequest,
  DiscoverPluginsRequestSchema,
  DiscoverPluginsResponse,
  GetPluginManifestRequest,
  GetPluginManifestResponse,
  ValidatePluginRequest,
  ValidatePluginResponse,
  InstallPluginRequest,
  InstallPluginResponse,
  UpdatePluginRequest,
  UpdatePluginResponse,
  RemovePluginRequest,
  RemovePluginResponse,
  ListInstalledPluginsRequest,
  ListInstalledPluginsRequestSchema,
  ListInstalledPluginsResponse,
  CheckPluginHealthRequest,
  CheckPluginHealthResponse
} from "../generated/finfocus/v1/registry_pb.js";
import {
  ObservabilityService,
  HealthCheckRequest,
  HealthCheckRequestSchema,
  HealthCheckResponse,
  GetMetricsRequest,
  GetMetricsRequestSchema,
  GetMetricsResponse,
  GetServiceLevelIndicatorsRequest,
  GetServiceLevelIndicatorsRequestSchema,
  GetServiceLevelIndicatorsResponse
} from "../generated/finfocus/v1/costsource_pb.js";

export interface ClientConfig {
  baseUrl: string;
  transport?: Transport;
}

export class RegistryClient {
  private client: Client<typeof PluginRegistryService>;

  constructor(config: ClientConfig) {
    const transport = config.transport || createConnectTransport({
      baseUrl: config.baseUrl,
      useBinaryFormat: false,
    });
    this.client = createClient(PluginRegistryService, transport);
  }

  async discoverPlugins(req: DiscoverPluginsRequest = create(DiscoverPluginsRequestSchema)): Promise<DiscoverPluginsResponse> {
    return this.client.discoverPlugins(req);
  }

  async getPluginManifest(req: GetPluginManifestRequest): Promise<GetPluginManifestResponse> {
      return this.client.getPluginManifest(req);
  }

  async validatePlugin(req: ValidatePluginRequest): Promise<ValidatePluginResponse> {
      return this.client.validatePlugin(req);
  }

  async installPlugin(req: InstallPluginRequest): Promise<InstallPluginResponse> {
      return this.client.installPlugin(req);
  }

  async updatePlugin(req: UpdatePluginRequest): Promise<UpdatePluginResponse> {
      return this.client.updatePlugin(req);
  }

  async removePlugin(req: RemovePluginRequest): Promise<RemovePluginResponse> {
      return this.client.removePlugin(req);
  }

  async listInstalledPlugins(req: ListInstalledPluginsRequest = create(ListInstalledPluginsRequestSchema)): Promise<ListInstalledPluginsResponse> {
      return this.client.listInstalledPlugins(req);
  }

  async checkPluginHealth(req: CheckPluginHealthRequest): Promise<CheckPluginHealthResponse> {
      return this.client.checkPluginHealth(req);
  }
}

export class ObservabilityClient {
  private client: Client<typeof ObservabilityService>;

  constructor(config: ClientConfig) {
    const transport = config.transport || createConnectTransport({
      baseUrl: config.baseUrl,
      useBinaryFormat: false,
    });
    this.client = createClient(ObservabilityService, transport);
  }

  async healthCheck(req: HealthCheckRequest = create(HealthCheckRequestSchema)): Promise<HealthCheckResponse> {
      return this.client.healthCheck(req);
  }

  async getMetrics(req: GetMetricsRequest = create(GetMetricsRequestSchema)): Promise<GetMetricsResponse> {
      return this.client.getMetrics(req);
  }

  async getServiceLevelIndicators(req: GetServiceLevelIndicatorsRequest = create(GetServiceLevelIndicatorsRequestSchema)): Promise<GetServiceLevelIndicatorsResponse> {
      return this.client.getServiceLevelIndicators(req);
  }
}
