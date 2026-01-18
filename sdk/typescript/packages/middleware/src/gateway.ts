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
/** Maximum request body size in bytes (1MB) */
const MAX_BODY_SIZE = 1024 * 1024;

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
   * Returns a Promise that resolves after the response is fully sent.
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

    // Wrap body reading in a Promise so handleRequest resolves after response is sent
    return new Promise<void>((resolve, reject) => {
      let body = '';
      let bodySize = 0;
      let aborted = false;

      req.on('data', chunk => {
        if (aborted) return;
        bodySize += chunk.length;
        if (bodySize > MAX_BODY_SIZE) {
          aborted = true;
          req.destroy();
          res.writeHead(413, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Request body too large' }));
          resolve();
          return;
        }
        body += chunk.toString();
      });

      req.on('error', (error) => {
        if (aborted) return;
        aborted = true;
        req.destroy();
        const errorMessage = error instanceof Error ? error.message : 'Request error';
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: errorMessage }));
        resolve();
      });

      req.on('end', async () => {
        if (aborted) return;

        // Parse JSON body separately to return 400 for malformed JSON
        let requestData: any = {};
        if (body) {
          try {
            requestData = JSON.parse(body);
          } catch (parseError) {
            // Treat all JSON parse errors as 400 Bad Request
            res.writeHead(400, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Malformed JSON' }));
            resolve();
            return;
          }
        }

        try {
          const response = await this.dispatchRequest(serviceName, methodName, requestData);

          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify(response));
          resolve();
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Unknown error';
          res.writeHead(500, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: errorMessage }));
          resolve();
        }
      });
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
  return (req: http.IncomingMessage, res: http.ServerResponse, next?: (err?: Error) => void) => {
    gateway.handleRequest(req, res).catch(error => {
      if (next) next(error instanceof Error ? error : new Error(String(error)));
    });
  };
}
