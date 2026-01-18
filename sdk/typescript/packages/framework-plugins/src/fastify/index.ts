import type { FastifyRequest, FastifyReply, FastifyInstance, FastifyPluginAsync } from 'fastify';
import { RESTGateway, RESTGatewayConfig } from 'finfocus-middleware';
import * as http from 'http';

/**
 * Creates a Fastify plugin for FinFocus REST Gateway.
 * Integrates the FinFocus REST gateway with Fastify applications.
 *
 * @param config - RESTGateway configuration with client instances
 * @returns Fastify plugin async function
 *
 * @example
 * ```typescript
 * import Fastify from 'fastify';
 * import { createFastifyPlugin } from 'finfocus-framework-plugins';
 * import { CostSourceClient } from 'finfocus-client';
 *
 * const fastify = Fastify();
 * const client = new CostSourceClient({ baseUrl: 'https://plugin.example.com' });
 * await fastify.register(createFastifyPlugin({ costSourceClient: client }));
 * await fastify.listen({ port: 3000 });
 * ```
 */
export function createFastifyPlugin(config: RESTGatewayConfig): FastifyPluginAsync {
  const gateway = new RESTGateway(config);

  return async (fastify: FastifyInstance) => {
    fastify.post('/*', async (request: FastifyRequest, reply: FastifyReply) => {
      try {
        // Only handle FinFocus REST API paths
        const path = request.url;
        if (!path.match(/^\/finfocus\.v1\./)) {
          return reply.status(404).send({ error: 'Not found' });
        }

        // Convert Fastify request/reply to Node.js http objects
        // Create a mock IncomingMessage
        const mockReq = {
          method: request.method,
          url: request.url,
          headers: request.headers,
          on: (_event: string, _callback: Function) => {},
          once: (_event: string, _callback: Function) => {},
        } as any as http.IncomingMessage;

        // Create a mock ServerResponse
        const mockRes = {
          writeHead: (statusCode: number, headers?: any) => {
            reply.status(statusCode);
            if (headers) {
              Object.entries(headers).forEach(([key, value]) => {
                reply.header(key, value as string);
              });
            }
          },
          end: (data?: string | Buffer) => {
            if (data) {
              reply.send(data);
            }
          },
          write: (data: string | Buffer) => {
            // This will be called if body is split
          },
        } as any as http.ServerResponse;

        // Handle the request
        await gateway.handleRequest(mockReq, mockRes);
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Unknown error';
        return reply.status(500).send({ error: errorMessage });
      }
    });
  };
}

/**
 * Alternative route-based approach for Fastify.
 * More explicit control over the routing.
 *
 * @example
 * ```typescript
 * const fastify = Fastify();
 * await fastify.register(createFastifyRoutes, { config });
 * ```
 */
export const createFastifyRoutes: FastifyPluginAsync<{ config: RESTGatewayConfig }> = async (
  fastify: FastifyInstance,
  { config }
) => {
  const plugin = createFastifyPlugin(config);
  await plugin(fastify);
};
