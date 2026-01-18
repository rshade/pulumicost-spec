import { http, HttpResponse } from 'msw';

export const handlers = [
  // CostSourceService endpoints
  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/Name', () => {
    return HttpResponse.json({
      name: "AWS Cost Plugin"
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/Supports', () => {
    return HttpResponse.json({
      billingModes: ["FLAT", "HOURLY", "MONTHLY"],
      providers: ["AWS"]
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetActualCost', () => {
    return HttpResponse.json({
      total: {
        units: "100",
        nanos: 0,
        currencyCode: "USD"
      },
      lineItems: []
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetProjectedCost', () => {
    return HttpResponse.json({
      total: {
        units: "150",
        nanos: 0,
        currencyCode: "USD"
      },
      lineItems: []
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetPricingSpec', () => {
    return HttpResponse.json({
      pricingSpec: {
        provider: "AWS",
        resourceType: "ec2.instance",
        billingMode: "HOURLY",
        ratePerUnit: {
          units: "10",
          currencyCode: "USD"
        }
      }
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/EstimateCost', () => {
    return HttpResponse.json({
      estimatedCost: {
        units: "200",
        nanos: 0,
        currencyCode: "USD"
      }
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetRecommendations', () => {
    return HttpResponse.json({
      recommendations: [
        {
          id: "rec-1",
          title: "Downsize Instance",
          category: "RESOURCE_SIZING",
          actionType: "DOWNSIZE",
          estimatedSavings: {
            units: "50",
            nanos: 0,
            currencyCode: "USD"
          },
          priority: "HIGH"
        }
      ],
      nextPageToken: ""
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/DismissRecommendation', () => {
    return HttpResponse.json({
      success: true
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetBudgets', () => {
    return HttpResponse.json({
      budgets: [
        {
          id: "budget-1",
          name: "Monthly Budget",
          limit: {
            units: "1000",
            nanos: 0,
            currencyCode: "USD"
          },
          currentSpend: {
            units: "750",
            nanos: 0,
            currencyCode: "USD"
          }
        }
      ]
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/GetPluginInfo', () => {
    return HttpResponse.json({
      name: "AWS Cost Plugin",
      version: "1.0.0",
      specVersion: "v1"
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.CostSourceService/DryRun', () => {
    return HttpResponse.json({
      fieldMappings: [],
      resourceTypeSupported: true,
      configurationValid: true
    });
  }),

  // ObservabilityService endpoints
  http.post('https://plugin-aws.example.com/finfocus.v1.ObservabilityService/HealthCheck', () => {
    return HttpResponse.json({
      status: "HEALTHY"
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.ObservabilityService/GetMetrics', () => {
    return HttpResponse.json({
      metrics: []
    });
  }),

  http.post('https://plugin-aws.example.com/finfocus.v1.ObservabilityService/GetServiceLevelIndicators', () => {
    return HttpResponse.json({
      indicators: []
    });
  }),

  // PluginRegistryService endpoints
  http.post('https://plugin-registry.example.com/finfocus.v1.PluginRegistryService/DiscoverPlugins', () => {
    return HttpResponse.json({
      plugins: [
        {
          id: "plugin-aws",
          name: "AWS Cost Plugin",
          version: "1.0.0"
        }
      ]
    });
  }),

  http.post('https://plugin-registry.example.com/finfocus.v1.PluginRegistryService/ListInstalledPlugins', () => {
    return HttpResponse.json({
      plugins: [
        {
          id: "plugin-aws",
          name: "AWS Cost Plugin",
          version: "1.0.0"
        }
      ]
    });
  })
];