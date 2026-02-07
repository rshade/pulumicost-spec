import { describe, it, expect, beforeAll, afterAll, afterEach } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import { create } from '@bufbuild/protobuf';
import { CostSourceClient } from '../src/clients/cost-source.js';
import { actualCostIterator, recommendationsIterator } from '../src/utils/pagination.js';
import {
  GetActualCostRequestSchema,
  GetRecommendationsRequestSchema,
} from '../src/generated/finfocus/v1/costsource_pb.js';

// Helper to create mock actual cost results
function createMockResults(count: number, startIndex: number = 0) {
  return Array.from({ length: count }, (_, i) => ({
    cost: startIndex + i + 1,
    usageAmount: 1,
    usageUnit: "hours",
    source: "test-plugin",
  }));
}

// Paginated GetActualCost handler that returns pages of results
function paginatedActualCostHandler(totalRecords: number, pageSize: number) {
  return http.post(
    'https://plugin-test.example.com/finfocus.v1.CostSourceService/GetActualCost',
    async ({ request }) => {
      const body = await request.json() as Record<string, unknown>;
      const requestedPageSize = (body.pageSize as number) || pageSize;
      const pageToken = (body.pageToken as string) || "";

      let offset = 0;
      if (pageToken) {
        offset = parseInt(Buffer.from(pageToken, 'base64').toString(), 10);
      }

      const end = Math.min(offset + requestedPageSize, totalRecords);
      const results = createMockResults(end - offset, offset);

      let nextPageToken = "";
      if (end < totalRecords) {
        nextPageToken = Buffer.from(end.toString()).toString('base64');
      }

      return HttpResponse.json({
        results,
        nextPageToken,
        totalCount: totalRecords,
      });
    }
  );
}

// Handler that returns an error on a specific page
function errorOnPageHandler(errorOnPage: number) {
  let callCount = 0;
  return http.post(
    'https://plugin-test.example.com/finfocus.v1.CostSourceService/GetActualCost',
    async () => {
      callCount++;
      if (callCount === errorOnPage) {
        return new HttpResponse(null, { status: 500 });
      }
      const results = createMockResults(50, (callCount - 1) * 50);
      const nextPageToken = Buffer.from((callCount * 50).toString()).toString('base64');
      return HttpResponse.json({
        results,
        nextPageToken,
        totalCount: 200,
      });
    }
  );
}

const server = setupServer();

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('actualCostIterator', () => {
  const client = new CostSourceClient({
    baseUrl: 'https://plugin-test.example.com',
  });

  it('iterates through multiple pages of actual cost results', async () => {
    server.use(paginatedActualCostHandler(200, 50));

    const request = create(GetActualCostRequestSchema, {
      resourceId: 'i-abc123',
      pageSize: 50,
    });

    const results: unknown[] = [];
    for await (const result of actualCostIterator(client, request)) {
      results.push(result);
    }

    expect(results).toHaveLength(200);
  });

  it('handles single-page results', async () => {
    server.use(paginatedActualCostHandler(10, 50));

    const request = create(GetActualCostRequestSchema, {
      resourceId: 'i-abc123',
      pageSize: 50,
    });

    const results: unknown[] = [];
    for await (const result of actualCostIterator(client, request)) {
      results.push(result);
    }

    expect(results).toHaveLength(10);
  });

  it('handles empty results', async () => {
    server.use(paginatedActualCostHandler(0, 50));

    const request = create(GetActualCostRequestSchema, {
      resourceId: 'i-abc123',
      pageSize: 50,
    });

    const results: unknown[] = [];
    for await (const result of actualCostIterator(client, request)) {
      results.push(result);
    }

    expect(results).toHaveLength(0);
  });

  it('propagates errors from the client', async () => {
    server.use(errorOnPageHandler(1));

    const request = create(GetActualCostRequestSchema, {
      resourceId: 'i-abc123',
      pageSize: 50,
    });

    await expect(async () => {
      for await (const _ of actualCostIterator(client, request)) {
        // consume
      }
    }).rejects.toThrow();
  });

  it('does not mutate the original request', async () => {
    server.use(paginatedActualCostHandler(100, 50));

    const request = create(GetActualCostRequestSchema, {
      resourceId: 'i-abc123',
      pageSize: 50,
      pageToken: '',
    });

    const originalToken = request.pageToken;

    for await (const _ of actualCostIterator(client, request)) {
      // consume all results
    }

    expect(request.pageToken).toBe(originalToken);
  });
});
