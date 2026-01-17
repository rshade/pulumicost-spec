import { PromiseClient, MethodInfo } from "@connectrpc/connect";
import * as http from "http";
import { CostSourceClient, ObservabilityClient, RegistryClient } from "finfocus-client";

export interface RESTGatewayConfig {
  costSourceClient: CostSourceClient;
  observabilityClient?: ObservabilityClient;
  registryClient?: RegistryClient;
}

/**
 * Generic REST Gateway that translates JSON HTTP requests to gRPC/Connect RPC calls.
 * Provides a simple HTTP handler that can be integrated with Node.js web frameworks.
 *
 * Maps HTTP POST requests to RPC methods:
 * - POST /finfocus.v1.CostSourceService/MethodName -> RPC call
 * - POST /finfocus.v1.ObservabilityService/MethodName -> RPC call
 * - POST /finfocus.v1.PluginRegistryService/MethodName -> RPC call
 *
 * @example
 * ```typescript
 * const gateway = new RESTGateway({
 *   costSourceClient: new CostSourceClient({ baseUrl: 'https://plugin.example.com' })
 * });
 *
 * const server = http.createServer((req, res) => {
 *   gateway.handleRequest(req, res);
 * });
 * ```
 */
export class RESTGateway {
  private costSourceClient: CostSourceClient;
  private observabilityClient?: ObservabilityClient;
  private registryClient?: RegistryClient;

  constructor(config: RESTGatewayConfig) {
    this.costSourceClient = config.costSourceClient;
    this.observabilityClient = config.observabilityClient;
    this.registryClient = config.registryClient;
  }

  /**
   * Main HTTP request handler.
   * Parses the request path and dispatches to appropriate RPC method.
   */
  async handleRequest(req: http.IncomingMessage, res: http.ServerResponse): Promise<void> {
    // Only handle POST requests
    if (req.method !== 'POST') {
      res.writeHead(405, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Method not allowed' }));
      return;
    }

    // Parse path: /finfocus.v1.ServiceName/MethodName
    const path = req.url || '';
    const pathMatch = path.match(/^\/finfocus\.v1\.([A-Za-z]+)\/([A-Za-z]+)$/);

    if (!pathMatch) {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Not found' }));
      return;
    }

    const [, serviceName, methodName] = pathMatch;

    // Read request body
    let body = '';
    req.on('data', chunk => {
      body += chunk.toString();
    });

    req.on('end', async () => {
      try {
        let requestData: any = {};
        if (body) {
          requestData = JSON.parse(body);
        }

        const response = await this.dispatchRequest(serviceName, methodName, requestData);

        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(response));
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Unknown error';
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: errorMessage }));
      }
    });
  }

  /**
   * Dispatches RPC calls to appropriate client based on service and method names.
   */
  private async dispatchRequest(serviceName: string, methodName: string, requestData: any): Promise<any> {
    const client = this.getClient(serviceName);
    if (!client) {
      throw new Error(`Service ${serviceName} not found`);
    }

    const method = methodName.charAt(0).toLowerCase() + methodName.slice(1);
    if (typeof (client as any)[method] !== 'function') {
      throw new Error(`Method ${methodName} not found in ${serviceName}`);
    }

    return await (client as any)[method](requestData);
  }

  /**
   * Gets the appropriate client based on service name.
   */
  private getClient(serviceName: string): any {
    switch (serviceName) {
      case 'CostSourceService':
        return this.costSourceClient;
      case 'ObservabilityService':
        return this.observabilityClient;
      case 'PluginRegistryService':
        return this.registryClient;
      default:
        return null;
    }
  }
}

/**
 * Express-style middleware wrapper for RESTGateway.
 *
 * @example
 * ```typescript
 * const gateway = new RESTGateway({...});
 * app.post('*', createRESTMiddleware(gateway));
 * ```
 */
export function createRESTMiddleware(gateway: RESTGateway) {
  return (req: http.IncomingMessage, res: http.ServerResponse, next?: () => void) => {
    gateway.handleRequest(req, res).catch(error => {
      if (next) next();
    });
  };
}
