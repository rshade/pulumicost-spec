import { createPromiseClient, PromiseClient, Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { PluginRegistryService } from "../generated/finfocus/v1/registry_connect.js";
import { ObservabilityService } from "../generated/finfocus/v1/costsource_connect.js";
import { 
    DiscoverPluginsRequest, DiscoverPluginsResponse,
    GetPluginManifestRequest, GetPluginManifestResponse,
    ValidatePluginRequest, ValidatePluginResponse,
    InstallPluginRequest, InstallPluginResponse,
    UpdatePluginRequest, UpdatePluginResponse,
    RemovePluginRequest, RemovePluginResponse,
    ListInstalledPluginsRequest, ListInstalledPluginsResponse,
    CheckPluginHealthRequest, CheckPluginHealthResponse
} from "../generated/finfocus/v1/registry_pb.js";
import {
    HealthCheckRequest, HealthCheckResponse,
    GetMetricsRequest, GetMetricsResponse,
    GetServiceLevelIndicatorsRequest, GetServiceLevelIndicatorsResponse
} from "../generated/finfocus/v1/costsource_pb.js";

export interface ClientConfig {
  baseUrl: string;
  transport?: Transport;
}

export class RegistryClient {
  private client: PromiseClient<typeof PluginRegistryService>;

  constructor(config: ClientConfig) {
    const transport = config.transport || createConnectTransport({
      baseUrl: config.baseUrl,
      useBinaryFormat: false,
    });
    this.client = createPromiseClient(PluginRegistryService, transport);
  }

  async discoverPlugins(req: DiscoverPluginsRequest = new DiscoverPluginsRequest()): Promise<DiscoverPluginsResponse> {
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

  async listInstalledPlugins(req: ListInstalledPluginsRequest = new ListInstalledPluginsRequest()): Promise<ListInstalledPluginsResponse> {
      return this.client.listInstalledPlugins(req);
  }

  async checkPluginHealth(req: CheckPluginHealthRequest): Promise<CheckPluginHealthResponse> {
      return this.client.checkPluginHealth(req);
  }
}

export class ObservabilityClient {
  private client: PromiseClient<typeof ObservabilityService>;

  constructor(config: ClientConfig) {
    const transport = config.transport || createConnectTransport({
      baseUrl: config.baseUrl,
      useBinaryFormat: false,
    });
    this.client = createPromiseClient(ObservabilityService, transport);
  }

  async healthCheck(req: HealthCheckRequest = new HealthCheckRequest()): Promise<HealthCheckResponse> {
      return this.client.healthCheck(req);
  }

  async getMetrics(req: GetMetricsRequest = new GetMetricsRequest()): Promise<GetMetricsResponse> {
      return this.client.getMetrics(req);
  }

  async getServiceLevelIndicators(req: GetServiceLevelIndicatorsRequest = new GetServiceLevelIndicatorsRequest()): Promise<GetServiceLevelIndicatorsResponse> {
      return this.client.getServiceLevelIndicators(req);
  }
}
