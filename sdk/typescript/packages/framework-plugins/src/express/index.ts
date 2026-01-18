import type { Request, Response, NextFunction } from 'express';
import { RESTGateway, RESTGatewayConfig } from 'finfocus-middleware';

/**
 * Creates Express middleware for FinFocus REST Gateway.
 * Integrates the FinFocus REST gateway with Express applications.
 *
 * @param config - RESTGateway configuration with client instances
 * @returns Express middleware function
 *
 * @example
 * ```typescript
 * import express from 'express';
 * import { createExpressMiddleware } from 'finfocus-framework-plugins';
 * import { CostSourceClient } from 'finfocus-client';
 *
 * const app = express();
 * const client = new CostSourceClient({ baseUrl: 'https://plugin.example.com' });
 * app.post('*', createExpressMiddleware({ costSourceClient: client }));
 * app.listen(3000);
 * ```
 */
export function createExpressMiddleware(config: RESTGatewayConfig) {
  const gateway = new RESTGateway(config);

  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      // Only handle POST requests matching the FinFocus REST API pattern
      if (req.method !== 'POST') {
        next();
        return;
      }

      const path = req.path;
      if (!path.match(/^\/finfocus\.v1\./)) {
        next();
        return;
      }

      // Create a mock IncomingMessage and ServerResponse for the gateway
      // since Express req/res are compatible with Node.js http objects
      await gateway.handleRequest(req as any, res as any);
    } catch (error) {
      next(error);
    }
  };
}

/**
 * Express Router configuration helper.
 * Sets up a dedicated router for FinFocus REST endpoints.
 *
 * @example
 * ```typescript
 * import express from 'express';
 * import { createExpressRouter } from 'finfocus-framework-plugins';
 *
 * const app = express();
 * const router = createExpressRouter(config);
 * app.use('/', router);
 * ```
 */
export function createExpressRouter(config: RESTGatewayConfig) {
  const express = require('express');
  const router = express.Router();

  const middleware = createExpressMiddleware(config);
  router.post('*', middleware);

  return router;
}
