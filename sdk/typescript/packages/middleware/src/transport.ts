import { Transport, ConnectError, Code } from "@connectrpc/connect";
import { createNodeHttpTransport } from "@connectrpc/connect-node";
import * as http from "http";
import * as https from "https";

export interface NodeTransportConfig {
  baseUrl: string;
  httpClient?: http.Agent;
  httpsClient?: https.Agent;
  timeout?: number;
}

/**
 * Creates a Node.js HTTP/HTTPS transport for Connect RPC.
 * Supports both standard HTTP and secure HTTPS connections.
 *
 * @example
 * ```typescript
 * const transport = createNodeTransport({
 *   baseUrl: 'https://plugin.example.com',
 *   timeout: 30000
 * });
 * ```
 */
export function createNodeTransport(config: NodeTransportConfig): Transport {
  return createNodeHttpTransport({
    baseUrl: config.baseUrl,
    httpClient: config.httpClient,
    httpsClient: config.httpsClient,
    interceptors: config.timeout ? [
      (next) => async (request) => {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), config.timeout!);
        try {
          const response = await next(request);
          clearTimeout(timeoutId);
          return response;
        } catch (error) {
          clearTimeout(timeoutId);
          if (error instanceof Error && error.name === 'AbortError') {
            throw new ConnectError('Request timeout', Code.DeadlineExceeded);
          }
          throw error;
        }
      }
    ] : []
  });
}
