window.BENCHMARK_DATA = {
  "lastUpdate": 1764514339108,
  "repoUrl": "https://github.com/rshade/pulumicost-spec",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "name": "rshade",
            "username": "rshade"
          },
          "committer": {
            "name": "rshade",
            "username": "rshade"
          },
          "id": "52650472520295d4c8d449dae631f63563850cf5",
          "message": "feat: run concurrent benchmark for EstimateCost",
          "timestamp": "2025-11-30T13:24:27Z",
          "url": "https://github.com/rshade/pulumicost-spec/pull/113/commits/52650472520295d4c8d449dae631f63563850cf5"
        },
        "date": 1764514338318,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkName",
            "value": 31192,
            "unit": "ns/op\t    8524 B/op\t     143 allocs/op",
            "extra": "36588 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 31192,
            "unit": "ns/op",
            "extra": "36588 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8524,
            "unit": "B/op",
            "extra": "36588 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "36588 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 39321,
            "unit": "ns/op\t    9465 B/op\t     172 allocs/op",
            "extra": "29338 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 39321,
            "unit": "ns/op",
            "extra": "29338 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9465,
            "unit": "B/op",
            "extra": "29338 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "29338 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 63187,
            "unit": "ns/op\t   18409 B/op\t     294 allocs/op",
            "extra": "18848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 63187,
            "unit": "ns/op",
            "extra": "18848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18409,
            "unit": "B/op",
            "extra": "18848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "18848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 40430,
            "unit": "ns/op\t    9667 B/op\t     176 allocs/op",
            "extra": "29246 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 40430,
            "unit": "ns/op",
            "extra": "29246 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9667,
            "unit": "B/op",
            "extra": "29246 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "29246 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 52495,
            "unit": "ns/op\t   12895 B/op\t     242 allocs/op",
            "extra": "20886 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 52495,
            "unit": "ns/op",
            "extra": "20886 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12895,
            "unit": "B/op",
            "extra": "20886 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "20886 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 35130,
            "unit": "ns/op\t    8752 B/op\t     149 allocs/op",
            "extra": "31494 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 35130,
            "unit": "ns/op",
            "extra": "31494 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8752,
            "unit": "B/op",
            "extra": "31494 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "31494 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 278639,
            "unit": "ns/op\t   67737 B/op\t    1176 allocs/op",
            "extra": "4335 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 278639,
            "unit": "ns/op",
            "extra": "4335 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67737,
            "unit": "B/op",
            "extra": "4335 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "4335 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 19391,
            "unit": "ns/op\t    8362 B/op\t     134 allocs/op",
            "extra": "62028 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 19391,
            "unit": "ns/op",
            "extra": "62028 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8362,
            "unit": "B/op",
            "extra": "62028 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "62028 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 20412,
            "unit": "ns/op\t    8584 B/op\t     140 allocs/op",
            "extra": "57829 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 20412,
            "unit": "ns/op",
            "extra": "57829 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8584,
            "unit": "B/op",
            "extra": "57829 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "57829 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 696745,
            "unit": "ns/op\t  431140 B/op\t    6930 allocs/op",
            "extra": "1612 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 696745,
            "unit": "ns/op",
            "extra": "1612 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431140,
            "unit": "B/op",
            "extra": "1612 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6930,
            "unit": "allocs/op",
            "extra": "1612 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 711936,
            "unit": "ns/op\t  431639 B/op\t    6932 allocs/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 711936,
            "unit": "ns/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431639,
            "unit": "B/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6932,
            "unit": "allocs/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 38079,
            "unit": "ns/op\t    9326 B/op\t     156 allocs/op",
            "extra": "32378 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 38079,
            "unit": "ns/op",
            "extra": "32378 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9326,
            "unit": "B/op",
            "extra": "32378 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "32378 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 63352,
            "unit": "ns/op\t   18413 B/op\t     294 allocs/op",
            "extra": "18840 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 63352,
            "unit": "ns/op",
            "extra": "18840 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18413,
            "unit": "B/op",
            "extra": "18840 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "18840 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 154140,
            "unit": "ns/op\t   77390 B/op\t    1161 allocs/op",
            "extra": "7419 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 154140,
            "unit": "ns/op",
            "extra": "7419 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77390,
            "unit": "B/op",
            "extra": "7419 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "7419 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 560002,
            "unit": "ns/op\t  314816 B/op\t    4489 allocs/op",
            "extra": "1867 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 560002,
            "unit": "ns/op",
            "extra": "1867 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 314816,
            "unit": "B/op",
            "extra": "1867 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "1867 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 40720,
            "unit": "ns/op\t    9671 B/op\t     176 allocs/op",
            "extra": "28065 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 40720,
            "unit": "ns/op",
            "extra": "28065 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9671,
            "unit": "B/op",
            "extra": "28065 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "28065 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 40650,
            "unit": "ns/op\t    9694 B/op\t     176 allocs/op",
            "extra": "29046 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 40650,
            "unit": "ns/op",
            "extra": "29046 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9694,
            "unit": "B/op",
            "extra": "29046 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "29046 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 40809,
            "unit": "ns/op\t    9716 B/op\t     176 allocs/op",
            "extra": "28962 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 40809,
            "unit": "ns/op",
            "extra": "28962 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9716,
            "unit": "B/op",
            "extra": "28962 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "28962 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 40859,
            "unit": "ns/op\t    9738 B/op\t     175 allocs/op",
            "extra": "28420 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 40859,
            "unit": "ns/op",
            "extra": "28420 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9738,
            "unit": "B/op",
            "extra": "28420 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "28420 times\n4 procs"
          }
        ]
      }
    ]
  }
}