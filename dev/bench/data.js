window.BENCHMARK_DATA = {
  "lastUpdate": 1764634627817,
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
      },
      {
        "commit": {
          "author": {
            "email": "131337+rshade@users.noreply.github.com",
            "name": "Richard Shade",
            "username": "rshade"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0ffcdc48e132b1eece31a1c51280cd250c608c23",
          "message": "feat: run concurrent benchmark for EstimateCost (#113)\n\n* feat: run concurrent benchmark for EstimateCost\n\nThis change implements benchmarks in ci.\n\nfixes #87\n\n* fix(markdown): Remove extra newline from tasks.md\n\n* fix(ci): Use benchmark-data branch for performance history",
          "timestamp": "2025-11-30T08:55:55-06:00",
          "tree_id": "8b16b1c214b3e5fda7900192b0b9b64905d06d21",
          "url": "https://github.com/rshade/pulumicost-spec/commit/0ffcdc48e132b1eece31a1c51280cd250c608c23"
        },
        "date": 1764514607047,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkName",
            "value": 41839,
            "unit": "ns/op\t    8540 B/op\t     143 allocs/op",
            "extra": "28780 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 41839,
            "unit": "ns/op",
            "extra": "28780 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8540,
            "unit": "B/op",
            "extra": "28780 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "28780 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 53749,
            "unit": "ns/op\t    9496 B/op\t     172 allocs/op",
            "extra": "21103 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 53749,
            "unit": "ns/op",
            "extra": "21103 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9496,
            "unit": "B/op",
            "extra": "21103 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "21103 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 80184,
            "unit": "ns/op\t   18445 B/op\t     294 allocs/op",
            "extra": "14683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 80184,
            "unit": "ns/op",
            "extra": "14683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18445,
            "unit": "B/op",
            "extra": "14683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "14683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 54896,
            "unit": "ns/op\t    9701 B/op\t     176 allocs/op",
            "extra": "20384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 54896,
            "unit": "ns/op",
            "extra": "20384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9701,
            "unit": "B/op",
            "extra": "20384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 66162,
            "unit": "ns/op\t   12918 B/op\t     242 allocs/op",
            "extra": "17218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 66162,
            "unit": "ns/op",
            "extra": "17218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12918,
            "unit": "B/op",
            "extra": "17218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "17218 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 46368,
            "unit": "ns/op\t    8769 B/op\t     149 allocs/op",
            "extra": "25444 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 46368,
            "unit": "ns/op",
            "extra": "25444 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8769,
            "unit": "B/op",
            "extra": "25444 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "25444 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 359177,
            "unit": "ns/op\t   67910 B/op\t    1176 allocs/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 359177,
            "unit": "ns/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67910,
            "unit": "B/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3272 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 21277,
            "unit": "ns/op\t    8367 B/op\t     134 allocs/op",
            "extra": "53905 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 21277,
            "unit": "ns/op",
            "extra": "53905 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8367,
            "unit": "B/op",
            "extra": "53905 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "53905 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 22310,
            "unit": "ns/op\t    8589 B/op\t     140 allocs/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 22310,
            "unit": "ns/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8589,
            "unit": "B/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 677510,
            "unit": "ns/op\t  431012 B/op\t    6930 allocs/op",
            "extra": "1744 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 677510,
            "unit": "ns/op",
            "extra": "1744 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431012,
            "unit": "B/op",
            "extra": "1744 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6930,
            "unit": "allocs/op",
            "extra": "1744 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 672175,
            "unit": "ns/op\t  431574 B/op\t    6931 allocs/op",
            "extra": "1681 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 672175,
            "unit": "ns/op",
            "extra": "1681 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431574,
            "unit": "B/op",
            "extra": "1681 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1681 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 50373,
            "unit": "ns/op\t    9353 B/op\t     156 allocs/op",
            "extra": "23286 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 50373,
            "unit": "ns/op",
            "extra": "23286 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9353,
            "unit": "B/op",
            "extra": "23286 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "23286 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 78632,
            "unit": "ns/op\t   18445 B/op\t     294 allocs/op",
            "extra": "14806 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 78632,
            "unit": "ns/op",
            "extra": "14806 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18445,
            "unit": "B/op",
            "extra": "14806 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "14806 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 172198,
            "unit": "ns/op\t   77380 B/op\t    1161 allocs/op",
            "extra": "7694 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 172198,
            "unit": "ns/op",
            "extra": "7694 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77380,
            "unit": "B/op",
            "extra": "7694 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "7694 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 508634,
            "unit": "ns/op\t  316253 B/op\t    4489 allocs/op",
            "extra": "2310 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 508634,
            "unit": "ns/op",
            "extra": "2310 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 316253,
            "unit": "B/op",
            "extra": "2310 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2310 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 54080,
            "unit": "ns/op\t    9694 B/op\t     176 allocs/op",
            "extra": "21829 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 54080,
            "unit": "ns/op",
            "extra": "21829 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9694,
            "unit": "B/op",
            "extra": "21829 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21829 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 54887,
            "unit": "ns/op\t    9720 B/op\t     176 allocs/op",
            "extra": "21638 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 54887,
            "unit": "ns/op",
            "extra": "21638 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9720,
            "unit": "B/op",
            "extra": "21638 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21638 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 54253,
            "unit": "ns/op\t    9743 B/op\t     176 allocs/op",
            "extra": "21582 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 54253,
            "unit": "ns/op",
            "extra": "21582 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9743,
            "unit": "B/op",
            "extra": "21582 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21582 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 54960,
            "unit": "ns/op\t    9763 B/op\t     175 allocs/op",
            "extra": "21567 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 54960,
            "unit": "ns/op",
            "extra": "21567 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9763,
            "unit": "B/op",
            "extra": "21567 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "21567 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "rshade98@hotmail.com",
            "name": "Richard Shade",
            "username": "rshade"
          },
          "committer": {
            "email": "rshade98@hotmail.com",
            "name": "Richard Shade",
            "username": "rshade"
          },
          "distinct": true,
          "id": "8944316a5337a12652efaa700999b7fd400517de",
          "message": "feat(ci): add performance regression testing workflow\n\n- Add .github/workflows/benchmarks.yml for automated regression testing\n- Configure 10% regression threshold\n- Run benchmarks on all SDK packages\n- Remove redundant benchmark job from ci.yml\n- Update README with regression testing details",
          "timestamp": "2025-11-30T10:05:36-06:00",
          "tree_id": "fe1e3abb3ed21e183ea980ecfe28f2787859f016",
          "url": "https://github.com/rshade/pulumicost-spec/commit/8944316a5337a12652efaa700999b7fd400517de"
        },
        "date": 1764518835373,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFocusRecordBuilder",
            "value": 241.7,
            "unit": "ns/op\t     528 B/op\t       6 allocs/op",
            "extra": "5079478 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - ns/op",
            "value": 241.7,
            "unit": "ns/op",
            "extra": "5079478 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "5079478 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5079478 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder",
            "value": 32.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36424184 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - ns/op",
            "value": 32.96,
            "unit": "ns/op",
            "extra": "36424184 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "36424184 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "36424184 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build",
            "value": 494.5,
            "unit": "ns/op\t    1128 B/op\t       8 allocs/op",
            "extra": "2436032 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - ns/op",
            "value": 494.5,
            "unit": "ns/op",
            "extra": "2436032 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - B/op",
            "value": 1128,
            "unit": "B/op",
            "extra": "2436032 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2436032 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild",
            "value": 238.2,
            "unit": "ns/op\t     352 B/op\t       6 allocs/op",
            "extra": "5020234 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - ns/op",
            "value": 238.2,
            "unit": "ns/op",
            "extra": "5020234 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "5020234 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5020234 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord",
            "value": 544.3,
            "unit": "ns/op\t     928 B/op\t       8 allocs/op",
            "extra": "2198461 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - ns/op",
            "value": 544.3,
            "unit": "ns/op",
            "extra": "2198461 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - B/op",
            "value": 928,
            "unit": "B/op",
            "extra": "2198461 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2198461 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger",
            "value": 245.4,
            "unit": "ns/op\t     544 B/op\t       3 allocs/op",
            "extra": "4875565 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - ns/op",
            "value": 245.4,
            "unit": "ns/op",
            "extra": "4875565 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "4875565 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4875565 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall",
            "value": 235.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5104062 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - ns/op",
            "value": 235.5,
            "unit": "ns/op",
            "extra": "5104062 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5104062 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5104062 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor",
            "value": 773.6,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1551891 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - ns/op",
            "value": 773.6,
            "unit": "ns/op",
            "extra": "1551891 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1551891 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1551891 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation",
            "value": 770.3,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1560910 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - ns/op",
            "value": 770.3,
            "unit": "ns/op",
            "extra": "1560910 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1560910 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1560910 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration",
            "value": 197.8,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6059686 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - ns/op",
            "value": 197.8,
            "unit": "ns/op",
            "extra": "6059686 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6059686 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6059686 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider",
            "value": 4.681,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "253405921 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - ns/op",
            "value": 4.681,
            "unit": "ns/op",
            "extra": "253405921 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "253405921 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "253405921 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource",
            "value": 3.276,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "365198514 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - ns/op",
            "value": 3.276,
            "unit": "ns/op",
            "extra": "365198514 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "365198514 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "365198514 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus",
            "value": 5.006,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "239404257 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - ns/op",
            "value": 5.006,
            "unit": "ns/op",
            "extra": "239404257 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "239404257 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "239404257 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel",
            "value": 4.691,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256999395 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - ns/op",
            "value": 4.691,
            "unit": "ns/op",
            "extra": "256999395 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256999395 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256999395 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod",
            "value": 4.208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "285027435 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - ns/op",
            "value": 4.208,
            "unit": "ns/op",
            "extra": "285027435 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "285027435 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "285027435 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability",
            "value": 8.484,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141257542 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - ns/op",
            "value": 8.484,
            "unit": "ns/op",
            "extra": "141257542 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141257542 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141257542 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission",
            "value": 5.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "235928636 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - ns/op",
            "value": 5.09,
            "unit": "ns/op",
            "extra": "235928636 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "235928636 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "235928636 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod",
            "value": 5.161,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "233579269 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - ns/op",
            "value": 5.161,
            "unit": "ns/op",
            "extra": "233579269 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "233579269 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "233579269 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased",
            "value": 11.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - ns/op",
            "value": 11.44,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values",
            "value": 3.281,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "357787929 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - ns/op",
            "value": 3.281,
            "unit": "ns/op",
            "extra": "357787929 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "357787929 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "357787929 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values",
            "value": 4.675,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256864508 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - ns/op",
            "value": 4.675,
            "unit": "ns/op",
            "extra": "256864508 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256864508 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256864508 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values",
            "value": 5.002,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "240274602 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - ns/op",
            "value": 5.002,
            "unit": "ns/op",
            "extra": "240274602 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "240274602 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "240274602 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values",
            "value": 4.831,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248030904 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - ns/op",
            "value": 4.831,
            "unit": "ns/op",
            "extra": "248030904 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248030904 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248030904 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values",
            "value": 8.494,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141348872 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - ns/op",
            "value": 8.494,
            "unit": "ns/op",
            "extra": "141348872 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141348872 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141348872 times\n4 procs"
          },
          {
            "name": "BenchmarkName",
            "value": 40004,
            "unit": "ns/op\t    8538 B/op\t     143 allocs/op",
            "extra": "29757 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 40004,
            "unit": "ns/op",
            "extra": "29757 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8538,
            "unit": "B/op",
            "extra": "29757 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "29757 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 51842,
            "unit": "ns/op\t    9492 B/op\t     172 allocs/op",
            "extra": "21862 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 51842,
            "unit": "ns/op",
            "extra": "21862 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9492,
            "unit": "B/op",
            "extra": "21862 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "21862 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 76816,
            "unit": "ns/op\t   18439 B/op\t     294 allocs/op",
            "extra": "15192 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 76816,
            "unit": "ns/op",
            "extra": "15192 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18439,
            "unit": "B/op",
            "extra": "15192 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15192 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 54817,
            "unit": "ns/op\t    9697 B/op\t     176 allocs/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 54817,
            "unit": "ns/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9697,
            "unit": "B/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21070 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 64804,
            "unit": "ns/op\t   12914 B/op\t     242 allocs/op",
            "extra": "17768 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 64804,
            "unit": "ns/op",
            "extra": "17768 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12914,
            "unit": "B/op",
            "extra": "17768 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "17768 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 45862,
            "unit": "ns/op\t    8767 B/op\t     149 allocs/op",
            "extra": "25888 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 45862,
            "unit": "ns/op",
            "extra": "25888 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8767,
            "unit": "B/op",
            "extra": "25888 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "25888 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 355433,
            "unit": "ns/op\t   67935 B/op\t    1176 allocs/op",
            "extra": "3164 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 355433,
            "unit": "ns/op",
            "extra": "3164 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67935,
            "unit": "B/op",
            "extra": "3164 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3164 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 21264,
            "unit": "ns/op\t    8366 B/op\t     134 allocs/op",
            "extra": "55622 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 21264,
            "unit": "ns/op",
            "extra": "55622 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8366,
            "unit": "B/op",
            "extra": "55622 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "55622 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 22249,
            "unit": "ns/op\t    8590 B/op\t     140 allocs/op",
            "extra": "50996 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 22249,
            "unit": "ns/op",
            "extra": "50996 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8590,
            "unit": "B/op",
            "extra": "50996 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "50996 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 678141,
            "unit": "ns/op\t  431017 B/op\t    6930 allocs/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 678141,
            "unit": "ns/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431017,
            "unit": "B/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6930,
            "unit": "allocs/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 694673,
            "unit": "ns/op\t  431521 B/op\t    6931 allocs/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 694673,
            "unit": "ns/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431521,
            "unit": "B/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 49119,
            "unit": "ns/op\t    9352 B/op\t     156 allocs/op",
            "extra": "23499 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 49119,
            "unit": "ns/op",
            "extra": "23499 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9352,
            "unit": "B/op",
            "extra": "23499 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "23499 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 77175,
            "unit": "ns/op\t   18442 B/op\t     294 allocs/op",
            "extra": "15230 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 77175,
            "unit": "ns/op",
            "extra": "15230 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18442,
            "unit": "B/op",
            "extra": "15230 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15230 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 174758,
            "unit": "ns/op\t   77424 B/op\t    1161 allocs/op",
            "extra": "6607 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 174758,
            "unit": "ns/op",
            "extra": "6607 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77424,
            "unit": "B/op",
            "extra": "6607 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "6607 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 505174,
            "unit": "ns/op\t  315055 B/op\t    4489 allocs/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 505174,
            "unit": "ns/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 315055,
            "unit": "B/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 53565,
            "unit": "ns/op\t    9693 B/op\t     176 allocs/op",
            "extra": "22045 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 53565,
            "unit": "ns/op",
            "extra": "22045 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9693,
            "unit": "B/op",
            "extra": "22045 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22045 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 54419,
            "unit": "ns/op\t    9721 B/op\t     176 allocs/op",
            "extra": "21666 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 54419,
            "unit": "ns/op",
            "extra": "21666 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9721,
            "unit": "B/op",
            "extra": "21666 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21666 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 53846,
            "unit": "ns/op\t    9743 B/op\t     176 allocs/op",
            "extra": "21650 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 53846,
            "unit": "ns/op",
            "extra": "21650 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9743,
            "unit": "B/op",
            "extra": "21650 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21650 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 53869,
            "unit": "ns/op\t    9762 B/op\t     175 allocs/op",
            "extra": "21819 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 53869,
            "unit": "ns/op",
            "extra": "21819 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9762,
            "unit": "B/op",
            "extra": "21819 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "21819 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "131337+rshade@users.noreply.github.com",
            "name": "Richard Shade",
            "username": "rshade"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "45f4b2e7c2a37c9414aada68343731a2f0e7913c",
          "message": "docs(006-estimate-cost): update data-model.md with actual decimal type (T054) (#114)\n\n## Summary\n\n- Replace placeholder `decimal` type with actual `double` type in EstimateCostResponse documentation\n- Add rationale referencing consistency with `GetProjectedCostResponse.cost_per_month`\n- Fix examples to use numeric format instead of string format\n\n## Changes\n\n### data-model.md Updates\n\n1. **Table constraints (line 57)**: Updated to explicitly reference `GetProjectedCostResponse.cost_per_month` consistency\n2. **Validation rules (line 63)**: Clarified `double` type choice with rationale\n3. **Examples (lines 71, 80)**: Changed from string format (`\"7.30\"`) to numeric format (`7.30`)\n4. **Diagram (line 137)**: Changed `cost_monthly: decimal` to `cost_monthly: double`\n\n## Verification\n\nVerified against actual implementation:\n\n- Proto definition (`costsource.proto:487`): `double cost_monthly = 2;`\n- Generated Go code (`costsource.pb.go:2416`): `CostMonthly float64`\n- Reference field (`costsource.pb.go:693`): `GetProjectedCostResponse.CostPerMonth float64`\n\n## Test Plan\n\n- [x] `make validate` passes\n- [x] Markdown linting passes\n- [x] Documentation matches proto definition\n- [x] Documentation matches generated Go code\n\nCloses #89",
          "timestamp": "2025-11-30T10:28:11-06:00",
          "tree_id": "2c80c00366b39bfa8990a76a7ac9f70ead9055d5",
          "url": "https://github.com/rshade/pulumicost-spec/commit/45f4b2e7c2a37c9414aada68343731a2f0e7913c"
        },
        "date": 1764520187595,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFocusRecordBuilder",
            "value": 240.7,
            "unit": "ns/op\t     528 B/op\t       6 allocs/op",
            "extra": "5160020 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - ns/op",
            "value": 240.7,
            "unit": "ns/op",
            "extra": "5160020 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "5160020 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5160020 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder",
            "value": 32.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36200720 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - ns/op",
            "value": 32.92,
            "unit": "ns/op",
            "extra": "36200720 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "36200720 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "36200720 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build",
            "value": 492.4,
            "unit": "ns/op\t    1128 B/op\t       8 allocs/op",
            "extra": "2450437 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - ns/op",
            "value": 492.4,
            "unit": "ns/op",
            "extra": "2450437 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - B/op",
            "value": 1128,
            "unit": "B/op",
            "extra": "2450437 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2450437 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild",
            "value": 239.7,
            "unit": "ns/op\t     352 B/op\t       6 allocs/op",
            "extra": "5021064 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - ns/op",
            "value": 239.7,
            "unit": "ns/op",
            "extra": "5021064 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "5021064 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5021064 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord",
            "value": 542.6,
            "unit": "ns/op\t     928 B/op\t       8 allocs/op",
            "extra": "2126230 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - ns/op",
            "value": 542.6,
            "unit": "ns/op",
            "extra": "2126230 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - B/op",
            "value": 928,
            "unit": "B/op",
            "extra": "2126230 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2126230 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger",
            "value": 270.9,
            "unit": "ns/op\t     544 B/op\t       3 allocs/op",
            "extra": "4518409 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - ns/op",
            "value": 270.9,
            "unit": "ns/op",
            "extra": "4518409 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "4518409 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4518409 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall",
            "value": 233.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5113381 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - ns/op",
            "value": 233.3,
            "unit": "ns/op",
            "extra": "5113381 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5113381 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5113381 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor",
            "value": 742.5,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1609939 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - ns/op",
            "value": 742.5,
            "unit": "ns/op",
            "extra": "1609939 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1609939 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1609939 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation",
            "value": 737.7,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1623421 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - ns/op",
            "value": 737.7,
            "unit": "ns/op",
            "extra": "1623421 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1623421 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1623421 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration",
            "value": 191.5,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6219361 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - ns/op",
            "value": 191.5,
            "unit": "ns/op",
            "extra": "6219361 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6219361 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6219361 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider",
            "value": 4.679,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256979935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - ns/op",
            "value": 4.679,
            "unit": "ns/op",
            "extra": "256979935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256979935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256979935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource",
            "value": 3.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "359585971 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - ns/op",
            "value": 3.39,
            "unit": "ns/op",
            "extra": "359585971 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "359585971 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "359585971 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus",
            "value": 4.989,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "240527412 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - ns/op",
            "value": 4.989,
            "unit": "ns/op",
            "extra": "240527412 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "240527412 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "240527412 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel",
            "value": 4.686,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256458440 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - ns/op",
            "value": 4.686,
            "unit": "ns/op",
            "extra": "256458440 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256458440 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256458440 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod",
            "value": 4.218,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "265324699 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - ns/op",
            "value": 4.218,
            "unit": "ns/op",
            "extra": "265324699 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "265324699 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "265324699 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability",
            "value": 8.483,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "140861103 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - ns/op",
            "value": 8.483,
            "unit": "ns/op",
            "extra": "140861103 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "140861103 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "140861103 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission",
            "value": 4.828,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248202469 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - ns/op",
            "value": 4.828,
            "unit": "ns/op",
            "extra": "248202469 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248202469 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248202469 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod",
            "value": 5.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "229162268 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - ns/op",
            "value": 5.14,
            "unit": "ns/op",
            "extra": "229162268 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "229162268 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "229162268 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased",
            "value": 11.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - ns/op",
            "value": 11.46,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased",
            "value": 9.829,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "121928300 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - ns/op",
            "value": 9.829,
            "unit": "ns/op",
            "extra": "121928300 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "121928300 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "121928300 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values",
            "value": 3.283,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "364862913 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - ns/op",
            "value": 3.283,
            "unit": "ns/op",
            "extra": "364862913 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "364862913 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "364862913 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values",
            "value": 4.814,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248731324 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - ns/op",
            "value": 4.814,
            "unit": "ns/op",
            "extra": "248731324 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248731324 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248731324 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values",
            "value": 4.994,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "233744911 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - ns/op",
            "value": 4.994,
            "unit": "ns/op",
            "extra": "233744911 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "233744911 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "233744911 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values",
            "value": 4.888,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "247960357 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - ns/op",
            "value": 4.888,
            "unit": "ns/op",
            "extra": "247960357 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "247960357 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "247960357 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values",
            "value": 8.582,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141197358 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - ns/op",
            "value": 8.582,
            "unit": "ns/op",
            "extra": "141197358 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141197358 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141197358 times\n4 procs"
          },
          {
            "name": "BenchmarkName",
            "value": 41777,
            "unit": "ns/op\t    8541 B/op\t     143 allocs/op",
            "extra": "28326 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 41777,
            "unit": "ns/op",
            "extra": "28326 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8541,
            "unit": "B/op",
            "extra": "28326 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "28326 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 54960,
            "unit": "ns/op\t    9496 B/op\t     172 allocs/op",
            "extra": "21052 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 54960,
            "unit": "ns/op",
            "extra": "21052 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9496,
            "unit": "B/op",
            "extra": "21052 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "21052 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 85753,
            "unit": "ns/op\t   18454 B/op\t     294 allocs/op",
            "extra": "13874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 85753,
            "unit": "ns/op",
            "extra": "13874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18454,
            "unit": "B/op",
            "extra": "13874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "13874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 57471,
            "unit": "ns/op\t    9704 B/op\t     176 allocs/op",
            "extra": "19837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 57471,
            "unit": "ns/op",
            "extra": "19837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9704,
            "unit": "B/op",
            "extra": "19837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "19837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 67438,
            "unit": "ns/op\t   12922 B/op\t     242 allocs/op",
            "extra": "16698 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 67438,
            "unit": "ns/op",
            "extra": "16698 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12922,
            "unit": "B/op",
            "extra": "16698 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "16698 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 48310,
            "unit": "ns/op\t    8771 B/op\t     149 allocs/op",
            "extra": "24823 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 48310,
            "unit": "ns/op",
            "extra": "24823 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8771,
            "unit": "B/op",
            "extra": "24823 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "24823 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 367631,
            "unit": "ns/op\t   67969 B/op\t    1176 allocs/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 367631,
            "unit": "ns/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67969,
            "unit": "B/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3037 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 21479,
            "unit": "ns/op\t    8367 B/op\t     134 allocs/op",
            "extra": "54536 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 21479,
            "unit": "ns/op",
            "extra": "54536 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8367,
            "unit": "B/op",
            "extra": "54536 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "54536 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 22289,
            "unit": "ns/op\t    8589 B/op\t     140 allocs/op",
            "extra": "52640 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 22289,
            "unit": "ns/op",
            "extra": "52640 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8589,
            "unit": "B/op",
            "extra": "52640 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "52640 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 654806,
            "unit": "ns/op\t  431007 B/op\t    6929 allocs/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 654806,
            "unit": "ns/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431007,
            "unit": "B/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6929,
            "unit": "allocs/op",
            "extra": "1730 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 677212,
            "unit": "ns/op\t  431569 B/op\t    6931 allocs/op",
            "extra": "1747 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 677212,
            "unit": "ns/op",
            "extra": "1747 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431569,
            "unit": "B/op",
            "extra": "1747 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1747 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 52216,
            "unit": "ns/op\t    9356 B/op\t     156 allocs/op",
            "extra": "22556 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 52216,
            "unit": "ns/op",
            "extra": "22556 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9356,
            "unit": "B/op",
            "extra": "22556 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "22556 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 84132,
            "unit": "ns/op\t   18455 B/op\t     294 allocs/op",
            "extra": "14058 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 84132,
            "unit": "ns/op",
            "extra": "14058 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18455,
            "unit": "B/op",
            "extra": "14058 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "14058 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 176212,
            "unit": "ns/op\t   77418 B/op\t    1161 allocs/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 176212,
            "unit": "ns/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77418,
            "unit": "B/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "6627 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 514854,
            "unit": "ns/op\t  314600 B/op\t    4489 allocs/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 514854,
            "unit": "ns/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 314600,
            "unit": "B/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2137 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 56953,
            "unit": "ns/op\t    9699 B/op\t     176 allocs/op",
            "extra": "20758 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 56953,
            "unit": "ns/op",
            "extra": "20758 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9699,
            "unit": "B/op",
            "extra": "20758 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20758 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 57325,
            "unit": "ns/op\t    9726 B/op\t     176 allocs/op",
            "extra": "20604 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 57325,
            "unit": "ns/op",
            "extra": "20604 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9726,
            "unit": "B/op",
            "extra": "20604 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20604 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 56936,
            "unit": "ns/op\t    9747 B/op\t     176 allocs/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 56936,
            "unit": "ns/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9747,
            "unit": "B/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 57188,
            "unit": "ns/op\t    9768 B/op\t     175 allocs/op",
            "extra": "20637 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 57188,
            "unit": "ns/op",
            "extra": "20637 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9768,
            "unit": "B/op",
            "extra": "20637 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "20637 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "41898282+github-actions[bot]@users.noreply.github.com",
            "name": "github-actions[bot]",
            "username": "github-actions[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a04a156f04899db0c4e70b4dd51c730f0519ca47",
          "message": "chore(main): release 0.4.2 (#110)\n\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-30T20:41:42-06:00",
          "tree_id": "61545219e8df0a5dcb1f08dda32652f3cd547b51",
          "url": "https://github.com/rshade/pulumicost-spec/commit/a04a156f04899db0c4e70b4dd51c730f0519ca47"
        },
        "date": 1764556996950,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFocusRecordBuilder",
            "value": 240.1,
            "unit": "ns/op\t     528 B/op\t       6 allocs/op",
            "extra": "5121912 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - ns/op",
            "value": 240.1,
            "unit": "ns/op",
            "extra": "5121912 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "5121912 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5121912 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder",
            "value": 32.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "35863947 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - ns/op",
            "value": 32.96,
            "unit": "ns/op",
            "extra": "35863947 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "35863947 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "35863947 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build",
            "value": 491.6,
            "unit": "ns/op\t    1128 B/op\t       8 allocs/op",
            "extra": "2425456 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - ns/op",
            "value": 491.6,
            "unit": "ns/op",
            "extra": "2425456 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - B/op",
            "value": 1128,
            "unit": "B/op",
            "extra": "2425456 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2425456 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild",
            "value": 237.8,
            "unit": "ns/op\t     352 B/op\t       6 allocs/op",
            "extra": "4949978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - ns/op",
            "value": 237.8,
            "unit": "ns/op",
            "extra": "4949978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "4949978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "4949978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord",
            "value": 545.8,
            "unit": "ns/op\t     928 B/op\t       8 allocs/op",
            "extra": "2188791 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - ns/op",
            "value": 545.8,
            "unit": "ns/op",
            "extra": "2188791 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - B/op",
            "value": 928,
            "unit": "B/op",
            "extra": "2188791 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2188791 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger",
            "value": 242.3,
            "unit": "ns/op\t     544 B/op\t       3 allocs/op",
            "extra": "4984584 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - ns/op",
            "value": 242.3,
            "unit": "ns/op",
            "extra": "4984584 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "4984584 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4984584 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall",
            "value": 235.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5072913 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - ns/op",
            "value": 235.1,
            "unit": "ns/op",
            "extra": "5072913 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5072913 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5072913 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor",
            "value": 740.2,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1629219 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - ns/op",
            "value": 740.2,
            "unit": "ns/op",
            "extra": "1629219 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1629219 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1629219 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation",
            "value": 740.2,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1622318 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - ns/op",
            "value": 740.2,
            "unit": "ns/op",
            "extra": "1622318 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1622318 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1622318 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration",
            "value": 189.4,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6341764 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - ns/op",
            "value": 189.4,
            "unit": "ns/op",
            "extra": "6341764 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6341764 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6341764 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider",
            "value": 4.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256702320 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - ns/op",
            "value": 4.67,
            "unit": "ns/op",
            "extra": "256702320 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256702320 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256702320 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource",
            "value": 3.286,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366661838 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - ns/op",
            "value": 3.286,
            "unit": "ns/op",
            "extra": "366661838 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366661838 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366661838 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus",
            "value": 5.125,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "239942092 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - ns/op",
            "value": 5.125,
            "unit": "ns/op",
            "extra": "239942092 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "239942092 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "239942092 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel",
            "value": 4.672,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256749421 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - ns/op",
            "value": 4.672,
            "unit": "ns/op",
            "extra": "256749421 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256749421 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256749421 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod",
            "value": 4.799,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "285352034 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - ns/op",
            "value": 4.799,
            "unit": "ns/op",
            "extra": "285352034 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "285352034 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "285352034 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability",
            "value": 8.483,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "138608710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - ns/op",
            "value": 8.483,
            "unit": "ns/op",
            "extra": "138608710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "138608710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "138608710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission",
            "value": 4.826,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248698461 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - ns/op",
            "value": 4.826,
            "unit": "ns/op",
            "extra": "248698461 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248698461 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248698461 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod",
            "value": 5.154,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "231772548 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - ns/op",
            "value": 5.154,
            "unit": "ns/op",
            "extra": "231772548 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "231772548 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "231772548 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased",
            "value": 11.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - ns/op",
            "value": 11.14,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased",
            "value": 10.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - ns/op",
            "value": 10.3,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values",
            "value": 3.265,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366891780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - ns/op",
            "value": 3.265,
            "unit": "ns/op",
            "extra": "366891780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366891780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366891780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values",
            "value": 4.668,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "254188671 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - ns/op",
            "value": 4.668,
            "unit": "ns/op",
            "extra": "254188671 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "254188671 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "254188671 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values",
            "value": 5.071,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "240366664 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - ns/op",
            "value": 5.071,
            "unit": "ns/op",
            "extra": "240366664 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "240366664 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "240366664 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values",
            "value": 4.825,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248502331 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - ns/op",
            "value": 4.825,
            "unit": "ns/op",
            "extra": "248502331 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248502331 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248502331 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values",
            "value": 8.503,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141281018 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - ns/op",
            "value": 8.503,
            "unit": "ns/op",
            "extra": "141281018 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141281018 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141281018 times\n4 procs"
          },
          {
            "name": "BenchmarkName",
            "value": 42114,
            "unit": "ns/op\t    8544 B/op\t     143 allocs/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 42114,
            "unit": "ns/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8544,
            "unit": "B/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "27334 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 54511,
            "unit": "ns/op\t    9497 B/op\t     172 allocs/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 54511,
            "unit": "ns/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9497,
            "unit": "B/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "20786 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 83197,
            "unit": "ns/op\t   18452 B/op\t     294 allocs/op",
            "extra": "13970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 83197,
            "unit": "ns/op",
            "extra": "13970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18452,
            "unit": "B/op",
            "extra": "13970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "13970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 57667,
            "unit": "ns/op\t    9705 B/op\t     176 allocs/op",
            "extra": "19536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 57667,
            "unit": "ns/op",
            "extra": "19536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9705,
            "unit": "B/op",
            "extra": "19536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "19536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 68430,
            "unit": "ns/op\t   12924 B/op\t     242 allocs/op",
            "extra": "16476 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 68430,
            "unit": "ns/op",
            "extra": "16476 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12924,
            "unit": "B/op",
            "extra": "16476 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "16476 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 48012,
            "unit": "ns/op\t    8771 B/op\t     149 allocs/op",
            "extra": "24866 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 48012,
            "unit": "ns/op",
            "extra": "24866 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8771,
            "unit": "B/op",
            "extra": "24866 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "24866 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 373323,
            "unit": "ns/op\t   67973 B/op\t    1176 allocs/op",
            "extra": "3025 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 373323,
            "unit": "ns/op",
            "extra": "3025 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67973,
            "unit": "B/op",
            "extra": "3025 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3025 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 21439,
            "unit": "ns/op\t    8368 B/op\t     134 allocs/op",
            "extra": "52815 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 21439,
            "unit": "ns/op",
            "extra": "52815 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "52815 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "52815 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 22707,
            "unit": "ns/op\t    8590 B/op\t     140 allocs/op",
            "extra": "50857 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 22707,
            "unit": "ns/op",
            "extra": "50857 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8590,
            "unit": "B/op",
            "extra": "50857 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "50857 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 697960,
            "unit": "ns/op\t  431102 B/op\t    6931 allocs/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 697960,
            "unit": "ns/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431102,
            "unit": "B/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 696427,
            "unit": "ns/op\t  431684 B/op\t    6931 allocs/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 696427,
            "unit": "ns/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431684,
            "unit": "B/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 51735,
            "unit": "ns/op\t    9357 B/op\t     156 allocs/op",
            "extra": "22356 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 51735,
            "unit": "ns/op",
            "extra": "22356 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9357,
            "unit": "B/op",
            "extra": "22356 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "22356 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 81966,
            "unit": "ns/op\t   18452 B/op\t     294 allocs/op",
            "extra": "14436 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 81966,
            "unit": "ns/op",
            "extra": "14436 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18452,
            "unit": "B/op",
            "extra": "14436 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "14436 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 174662,
            "unit": "ns/op\t   77428 B/op\t    1161 allocs/op",
            "extra": "6771 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 174662,
            "unit": "ns/op",
            "extra": "6771 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77428,
            "unit": "B/op",
            "extra": "6771 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "6771 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 509055,
            "unit": "ns/op\t  314994 B/op\t    4489 allocs/op",
            "extra": "2217 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 509055,
            "unit": "ns/op",
            "extra": "2217 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 314994,
            "unit": "B/op",
            "extra": "2217 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2217 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 56421,
            "unit": "ns/op\t    9698 B/op\t     176 allocs/op",
            "extra": "20994 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 56421,
            "unit": "ns/op",
            "extra": "20994 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9698,
            "unit": "B/op",
            "extra": "20994 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20994 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 57307,
            "unit": "ns/op\t    9727 B/op\t     176 allocs/op",
            "extra": "20467 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 57307,
            "unit": "ns/op",
            "extra": "20467 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9727,
            "unit": "B/op",
            "extra": "20467 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20467 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 56554,
            "unit": "ns/op\t    9748 B/op\t     176 allocs/op",
            "extra": "20772 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 56554,
            "unit": "ns/op",
            "extra": "20772 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9748,
            "unit": "B/op",
            "extra": "20772 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20772 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 57059,
            "unit": "ns/op\t    9768 B/op\t     175 allocs/op",
            "extra": "20678 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 57059,
            "unit": "ns/op",
            "extra": "20678 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9768,
            "unit": "B/op",
            "extra": "20678 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "20678 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a2c22fa359fed8b4e932d6c637a3841b1f62bb3d",
          "message": "chore(deps): Update prettier to 3.7.3 (#115)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-12-01T06:50:49-06:00",
          "tree_id": "ef15a1aba11617d6be988dc73d83760799e8e5c2",
          "url": "https://github.com/rshade/pulumicost-spec/commit/a2c22fa359fed8b4e932d6c637a3841b1f62bb3d"
        },
        "date": 1764593539511,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkFocusRecordBuilder",
            "value": 242.1,
            "unit": "ns/op\t     528 B/op\t       6 allocs/op",
            "extra": "5127231 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - ns/op",
            "value": 242.1,
            "unit": "ns/op",
            "extra": "5127231 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "5127231 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5127231 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder",
            "value": 32.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36194487 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - ns/op",
            "value": 32.92,
            "unit": "ns/op",
            "extra": "36194487 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "36194487 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "36194487 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build",
            "value": 553,
            "unit": "ns/op\t    1128 B/op\t       8 allocs/op",
            "extra": "2402106 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - ns/op",
            "value": 553,
            "unit": "ns/op",
            "extra": "2402106 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - B/op",
            "value": 1128,
            "unit": "B/op",
            "extra": "2402106 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2402106 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild",
            "value": 238.9,
            "unit": "ns/op\t     352 B/op\t       6 allocs/op",
            "extra": "4885694 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - ns/op",
            "value": 238.9,
            "unit": "ns/op",
            "extra": "4885694 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "4885694 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "4885694 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord",
            "value": 546.9,
            "unit": "ns/op\t     928 B/op\t       8 allocs/op",
            "extra": "2172391 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - ns/op",
            "value": 546.9,
            "unit": "ns/op",
            "extra": "2172391 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - B/op",
            "value": 928,
            "unit": "B/op",
            "extra": "2172391 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2172391 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger",
            "value": 244.7,
            "unit": "ns/op\t     544 B/op\t       3 allocs/op",
            "extra": "4956097 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - ns/op",
            "value": 244.7,
            "unit": "ns/op",
            "extra": "4956097 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "4956097 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4956097 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall",
            "value": 232.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4991290 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - ns/op",
            "value": 232.5,
            "unit": "ns/op",
            "extra": "4991290 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4991290 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4991290 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor",
            "value": 742.7,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1613864 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - ns/op",
            "value": 742.7,
            "unit": "ns/op",
            "extra": "1613864 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1613864 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1613864 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation",
            "value": 737,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1620068 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - ns/op",
            "value": 737,
            "unit": "ns/op",
            "extra": "1620068 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1620068 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1620068 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration",
            "value": 191.4,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6258622 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - ns/op",
            "value": 191.4,
            "unit": "ns/op",
            "extra": "6258622 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6258622 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6258622 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider",
            "value": 4.737,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "253985077 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - ns/op",
            "value": 4.737,
            "unit": "ns/op",
            "extra": "253985077 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "253985077 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "253985077 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource",
            "value": 3.267,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "367372519 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - ns/op",
            "value": 3.267,
            "unit": "ns/op",
            "extra": "367372519 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "367372519 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "367372519 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus",
            "value": 5.002,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "239878935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - ns/op",
            "value": 5.002,
            "unit": "ns/op",
            "extra": "239878935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "239878935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "239878935 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel",
            "value": 4.721,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "254666653 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - ns/op",
            "value": 4.721,
            "unit": "ns/op",
            "extra": "254666653 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "254666653 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "254666653 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod",
            "value": 4.201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "280865877 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - ns/op",
            "value": 4.201,
            "unit": "ns/op",
            "extra": "280865877 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "280865877 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "280865877 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability",
            "value": 8.499,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "139062710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - ns/op",
            "value": 8.499,
            "unit": "ns/op",
            "extra": "139062710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "139062710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "139062710 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission",
            "value": 4.829,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "236206362 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - ns/op",
            "value": 4.829,
            "unit": "ns/op",
            "extra": "236206362 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "236206362 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "236206362 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod",
            "value": 5.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "233178884 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - ns/op",
            "value": 5.13,
            "unit": "ns/op",
            "extra": "233178884 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "233178884 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "233178884 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased",
            "value": 11.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - ns/op",
            "value": 11.17,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased",
            "value": 10.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - ns/op",
            "value": 10.27,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values",
            "value": 3.281,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366624427 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - ns/op",
            "value": 3.281,
            "unit": "ns/op",
            "extra": "366624427 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366624427 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366624427 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values",
            "value": 4.674,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256157469 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - ns/op",
            "value": 4.674,
            "unit": "ns/op",
            "extra": "256157469 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256157469 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256157469 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values",
            "value": 4.992,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "240017562 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - ns/op",
            "value": 4.992,
            "unit": "ns/op",
            "extra": "240017562 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "240017562 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "240017562 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values",
            "value": 4.838,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248775394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - ns/op",
            "value": 4.838,
            "unit": "ns/op",
            "extra": "248775394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248775394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248775394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values",
            "value": 8.497,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141577208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - ns/op",
            "value": 8.497,
            "unit": "ns/op",
            "extra": "141577208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141577208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141577208 times\n4 procs"
          },
          {
            "name": "BenchmarkName",
            "value": 38807,
            "unit": "ns/op\t    8536 B/op\t     143 allocs/op",
            "extra": "30535 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 38807,
            "unit": "ns/op",
            "extra": "30535 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8536,
            "unit": "B/op",
            "extra": "30535 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "30535 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 50832,
            "unit": "ns/op\t    9489 B/op\t     172 allocs/op",
            "extra": "22404 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 50832,
            "unit": "ns/op",
            "extra": "22404 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9489,
            "unit": "B/op",
            "extra": "22404 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "22404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 75617,
            "unit": "ns/op\t   18435 B/op\t     294 allocs/op",
            "extra": "15672 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 75617,
            "unit": "ns/op",
            "extra": "15672 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18435,
            "unit": "B/op",
            "extra": "15672 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15672 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 53375,
            "unit": "ns/op\t    9698 B/op\t     176 allocs/op",
            "extra": "20845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 53375,
            "unit": "ns/op",
            "extra": "20845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9698,
            "unit": "B/op",
            "extra": "20845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "20845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 63599,
            "unit": "ns/op\t   12913 B/op\t     242 allocs/op",
            "extra": "17944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 63599,
            "unit": "ns/op",
            "extra": "17944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12913,
            "unit": "B/op",
            "extra": "17944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "17944 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 44720,
            "unit": "ns/op\t    8765 B/op\t     149 allocs/op",
            "extra": "26415 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 44720,
            "unit": "ns/op",
            "extra": "26415 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8765,
            "unit": "B/op",
            "extra": "26415 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "26415 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 349531,
            "unit": "ns/op\t   67898 B/op\t    1176 allocs/op",
            "extra": "3361 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 349531,
            "unit": "ns/op",
            "extra": "3361 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67898,
            "unit": "B/op",
            "extra": "3361 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3361 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 20658,
            "unit": "ns/op\t    8366 B/op\t     134 allocs/op",
            "extra": "56401 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 20658,
            "unit": "ns/op",
            "extra": "56401 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8366,
            "unit": "B/op",
            "extra": "56401 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "56401 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 21507,
            "unit": "ns/op\t    8588 B/op\t     140 allocs/op",
            "extra": "53557 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 21507,
            "unit": "ns/op",
            "extra": "53557 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8588,
            "unit": "B/op",
            "extra": "53557 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "53557 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 683525,
            "unit": "ns/op\t  431072 B/op\t    6931 allocs/op",
            "extra": "1748 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 683525,
            "unit": "ns/op",
            "extra": "1748 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431072,
            "unit": "B/op",
            "extra": "1748 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1748 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 690638,
            "unit": "ns/op\t  431592 B/op\t    6931 allocs/op",
            "extra": "1722 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 690638,
            "unit": "ns/op",
            "extra": "1722 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431592,
            "unit": "B/op",
            "extra": "1722 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1722 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 46980,
            "unit": "ns/op\t    9349 B/op\t     156 allocs/op",
            "extra": "24411 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 46980,
            "unit": "ns/op",
            "extra": "24411 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9349,
            "unit": "B/op",
            "extra": "24411 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "24411 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 73505,
            "unit": "ns/op\t   18436 B/op\t     294 allocs/op",
            "extra": "15760 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 73505,
            "unit": "ns/op",
            "extra": "15760 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18436,
            "unit": "B/op",
            "extra": "15760 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15760 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 166890,
            "unit": "ns/op\t   77424 B/op\t    1161 allocs/op",
            "extra": "6784 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 166890,
            "unit": "ns/op",
            "extra": "6784 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77424,
            "unit": "B/op",
            "extra": "6784 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "6784 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 484943,
            "unit": "ns/op\t  313766 B/op\t    4489 allocs/op",
            "extra": "2184 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 484943,
            "unit": "ns/op",
            "extra": "2184 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 313766,
            "unit": "B/op",
            "extra": "2184 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2184 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 53222,
            "unit": "ns/op\t    9692 B/op\t     176 allocs/op",
            "extra": "22129 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 53222,
            "unit": "ns/op",
            "extra": "22129 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9692,
            "unit": "B/op",
            "extra": "22129 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22129 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 52798,
            "unit": "ns/op\t    9717 B/op\t     176 allocs/op",
            "extra": "22441 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 52798,
            "unit": "ns/op",
            "extra": "22441 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9717,
            "unit": "B/op",
            "extra": "22441 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22441 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 53236,
            "unit": "ns/op\t    9740 B/op\t     176 allocs/op",
            "extra": "22188 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 53236,
            "unit": "ns/op",
            "extra": "22188 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9740,
            "unit": "B/op",
            "extra": "22188 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22188 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 52344,
            "unit": "ns/op\t    9759 B/op\t     175 allocs/op",
            "extra": "22519 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 52344,
            "unit": "ns/op",
            "extra": "22519 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9759,
            "unit": "B/op",
            "extra": "22519 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "22519 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "131337+rshade@users.noreply.github.com",
            "name": "Richard Shade",
            "username": "rshade"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "97e34f5fba0f1a635ab9651ed2bc80510898b962",
          "message": "feat(sdk/go/currency): extract ISO 4217 validation as reusable package (T101) (#116)\n\n- Add shared `sdk/go/currency` package with complete ISO 4217 metadata and zero-allocation validation helpers\n- Update pluginsdk Focus validation to consume the new package instead of inline maps\n- Document usage, benchmarks, and defensive copy behavior for external consumers\n\n1. **New package**: `Currency` struct with code, name, numeric code, and minor units plus `IsValid`, `GetCurrency`, and `AllCurrencies`\n2. **Performance**: Zero-allocation linear scan validation and defensive copies for metadata access\n3. **Benchmarks**: Added benchmarks targeting `<15 ns/op` for validation and coverage for metadata integrity\n4. **Docs**: README with usage examples and performance table; package docs in `doc.go`\n\n1. **FOCUS currency validation**: Replaced inline ISO 4217 map with calls to `currency.IsValid` in `focus_conformance.go`\n2. **Docs**: Referenced shared currency utilities in pluginsdk README\n\n- [ ] `go test ./...`\n- [ ] `make validate`\n- [ ] Markdown linting (`npm run lint:md` or equivalent)\n\nCloses #101",
          "timestamp": "2025-12-01T18:15:20-06:00",
          "tree_id": "615a88d545fccfe4687526bbc4276b1fcaa56b80",
          "url": "https://github.com/rshade/pulumicost-spec/commit/97e34f5fba0f1a635ab9651ed2bc80510898b962"
        },
        "date": 1764634626927,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkIsValid/USD",
            "value": 1385,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/USD - ns/op",
            "value": 1385,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/USD - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/USD - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/MXN",
            "value": 429,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "2747091 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/MXN - ns/op",
            "value": 429,
            "unit": "ns/op",
            "extra": "2747091 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/MXN - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "2747091 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/MXN - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "2747091 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/ZWL",
            "value": 719.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1665750 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/ZWL - ns/op",
            "value": 719.3,
            "unit": "ns/op",
            "extra": "1665750 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/ZWL - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1665750 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/ZWL - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1665750 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/Invalid",
            "value": 832.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1424821 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/Invalid - ns/op",
            "value": 832.5,
            "unit": "ns/op",
            "extra": "1424821 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/Invalid - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1424821 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValid/Invalid - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1424821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/USD",
            "value": 10.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/USD - ns/op",
            "value": 10.56,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/USD - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/USD - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/JPY",
            "value": 10.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/JPY - ns/op",
            "value": 10.52,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/JPY - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/JPY - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/Invalid",
            "value": 8.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "149637079 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/Invalid - ns/op",
            "value": 8.02,
            "unit": "ns/op",
            "extra": "149637079 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/Invalid - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "149637079 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCurrency/Invalid - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "149637079 times\n4 procs"
          },
          {
            "name": "BenchmarkAllCurrencies",
            "value": 0.3117,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkAllCurrencies - ns/op",
            "value": 0.3117,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkAllCurrencies - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkAllCurrencies - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCurrency_String",
            "value": 0.3117,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCurrency_String - ns/op",
            "value": 0.3117,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCurrency_String - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCurrency_String - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder",
            "value": 237.2,
            "unit": "ns/op\t     528 B/op\t       6 allocs/op",
            "extra": "5207194 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - ns/op",
            "value": 237.2,
            "unit": "ns/op",
            "extra": "5207194 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "5207194 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5207194 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder",
            "value": 32.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "36349264 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - ns/op",
            "value": 32.93,
            "unit": "ns/op",
            "extra": "36349264 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "36349264 times\n4 procs"
          },
          {
            "name": "BenchmarkNewFocusRecordBuilder - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "36349264 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build",
            "value": 485.3,
            "unit": "ns/op\t    1128 B/op\t       8 allocs/op",
            "extra": "2369088 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - ns/op",
            "value": 485.3,
            "unit": "ns/op",
            "extra": "2369088 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - B/op",
            "value": 1128,
            "unit": "B/op",
            "extra": "2369088 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_Build - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2369088 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild",
            "value": 236.6,
            "unit": "ns/op\t     352 B/op\t       6 allocs/op",
            "extra": "5055978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - ns/op",
            "value": 236.6,
            "unit": "ns/op",
            "extra": "5055978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - B/op",
            "value": 352,
            "unit": "B/op",
            "extra": "5055978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_ChainedBuild - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "5055978 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord",
            "value": 543.7,
            "unit": "ns/op\t     928 B/op\t       8 allocs/op",
            "extra": "2204103 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - ns/op",
            "value": 543.7,
            "unit": "ns/op",
            "extra": "2204103 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - B/op",
            "value": 928,
            "unit": "B/op",
            "extra": "2204103 times\n4 procs"
          },
          {
            "name": "BenchmarkFocusRecordBuilder_FullRecord - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "2204103 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger",
            "value": 244.8,
            "unit": "ns/op\t     544 B/op\t       3 allocs/op",
            "extra": "4860967 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - ns/op",
            "value": 244.8,
            "unit": "ns/op",
            "extra": "4860967 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - B/op",
            "value": 544,
            "unit": "B/op",
            "extra": "4860967 times\n4 procs"
          },
          {
            "name": "BenchmarkNewPluginLogger - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4860967 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall",
            "value": 220.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5243196 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - ns/op",
            "value": 220.7,
            "unit": "ns/op",
            "extra": "5243196 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5243196 times\n4 procs"
          },
          {
            "name": "BenchmarkLogCall - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5243196 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor",
            "value": 731.2,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1637185 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - ns/op",
            "value": 731.2,
            "unit": "ns/op",
            "extra": "1637185 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1637185 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptor - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1637185 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation",
            "value": 733.6,
            "unit": "ns/op\t     608 B/op\t       9 allocs/op",
            "extra": "1635493 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - ns/op",
            "value": 733.6,
            "unit": "ns/op",
            "extra": "1635493 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - B/op",
            "value": 608,
            "unit": "B/op",
            "extra": "1635493 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorValidation - allocs/op",
            "value": 9,
            "unit": "allocs/op",
            "extra": "1635493 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration",
            "value": 191.4,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6263745 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - ns/op",
            "value": 191.4,
            "unit": "ns/op",
            "extra": "6263745 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6263745 times\n4 procs"
          },
          {
            "name": "BenchmarkInterceptorGeneration - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6263745 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider",
            "value": 4.666,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256514622 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - ns/op",
            "value": 4.666,
            "unit": "ns/op",
            "extra": "256514622 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256514622 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256514622 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource",
            "value": 3.274,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366967654 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - ns/op",
            "value": 3.274,
            "unit": "ns/op",
            "extra": "366967654 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366967654 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidDiscoverySource - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366967654 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus",
            "value": 5.055,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "239572363 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - ns/op",
            "value": 5.055,
            "unit": "ns/op",
            "extra": "239572363 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "239572363 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginStatus - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "239572363 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel",
            "value": 4.679,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "250413096 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - ns/op",
            "value": 4.679,
            "unit": "ns/op",
            "extra": "250413096 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "250413096 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSecurityLevel - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "250413096 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod",
            "value": 4.748,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "255107254 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - ns/op",
            "value": 4.748,
            "unit": "ns/op",
            "extra": "255107254 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "255107254 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidInstallationMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "255107254 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability",
            "value": 8.488,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141208396 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - ns/op",
            "value": 8.488,
            "unit": "ns/op",
            "extra": "141208396 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141208396 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141208396 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission",
            "value": 5.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "235679281 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - ns/op",
            "value": 5.1,
            "unit": "ns/op",
            "extra": "235679281 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "235679281 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidSystemPermission - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "235679281 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod",
            "value": 5.137,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "233060354 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - ns/op",
            "value": 5.137,
            "unit": "ns/op",
            "extra": "233060354 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "233060354 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidAuthMethod - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "233060354 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased",
            "value": 11.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - ns/op",
            "value": 11.23,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidProvider_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased",
            "value": 10.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - ns/op",
            "value": 10.84,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkIsValidPluginCapability_MapBased - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values",
            "value": 3.272,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366886780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - ns/op",
            "value": 3.272,
            "unit": "ns/op",
            "extra": "366886780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366886780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_4Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366886780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values",
            "value": 4.675,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "256722234 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - ns/op",
            "value": 4.675,
            "unit": "ns/op",
            "extra": "256722234 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "256722234 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_5Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "256722234 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values",
            "value": 4.996,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "240215178 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - ns/op",
            "value": 4.996,
            "unit": "ns/op",
            "extra": "240215178 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "240215178 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_6Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "240215178 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values",
            "value": 4.829,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "248886984 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - ns/op",
            "value": 4.829,
            "unit": "ns/op",
            "extra": "248886984 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "248886984 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_9Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "248886984 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values",
            "value": 8.509,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "141430992 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - ns/op",
            "value": 8.509,
            "unit": "ns/op",
            "extra": "141430992 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "141430992 times\n4 procs"
          },
          {
            "name": "BenchmarkValidation_14Values - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "141430992 times\n4 procs"
          },
          {
            "name": "BenchmarkName",
            "value": 38877,
            "unit": "ns/op\t    8535 B/op\t     143 allocs/op",
            "extra": "30712 times\n4 procs"
          },
          {
            "name": "BenchmarkName - ns/op",
            "value": 38877,
            "unit": "ns/op",
            "extra": "30712 times\n4 procs"
          },
          {
            "name": "BenchmarkName - B/op",
            "value": 8535,
            "unit": "B/op",
            "extra": "30712 times\n4 procs"
          },
          {
            "name": "BenchmarkName - allocs/op",
            "value": 143,
            "unit": "allocs/op",
            "extra": "30712 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports",
            "value": 50178,
            "unit": "ns/op\t    9489 B/op\t     172 allocs/op",
            "extra": "22368 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - ns/op",
            "value": 50178,
            "unit": "ns/op",
            "extra": "22368 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - B/op",
            "value": 9489,
            "unit": "B/op",
            "extra": "22368 times\n4 procs"
          },
          {
            "name": "BenchmarkSupports - allocs/op",
            "value": 172,
            "unit": "allocs/op",
            "extra": "22368 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost",
            "value": 75406,
            "unit": "ns/op\t   18436 B/op\t     294 allocs/op",
            "extra": "15548 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - ns/op",
            "value": 75406,
            "unit": "ns/op",
            "extra": "15548 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - B/op",
            "value": 18436,
            "unit": "B/op",
            "extra": "15548 times\n4 procs"
          },
          {
            "name": "BenchmarkGetActualCost - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15548 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost",
            "value": 53474,
            "unit": "ns/op\t    9696 B/op\t     176 allocs/op",
            "extra": "21184 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - ns/op",
            "value": 53474,
            "unit": "ns/op",
            "extra": "21184 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - B/op",
            "value": 9696,
            "unit": "B/op",
            "extra": "21184 times\n4 procs"
          },
          {
            "name": "BenchmarkGetProjectedCost - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "21184 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec",
            "value": 64238,
            "unit": "ns/op\t   12796 B/op\t     242 allocs/op",
            "extra": "17845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - ns/op",
            "value": 64238,
            "unit": "ns/op",
            "extra": "17845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - B/op",
            "value": 12796,
            "unit": "B/op",
            "extra": "17845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPricingSpec - allocs/op",
            "value": 242,
            "unit": "allocs/op",
            "extra": "17845 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost",
            "value": 44690,
            "unit": "ns/op\t    8765 B/op\t     149 allocs/op",
            "extra": "26570 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - ns/op",
            "value": 44690,
            "unit": "ns/op",
            "extra": "26570 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - B/op",
            "value": 8765,
            "unit": "B/op",
            "extra": "26570 times\n4 procs"
          },
          {
            "name": "BenchmarkEstimateCost - allocs/op",
            "value": 149,
            "unit": "allocs/op",
            "extra": "26570 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods",
            "value": 347285,
            "unit": "ns/op\t   67934 B/op\t    1176 allocs/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - ns/op",
            "value": 347285,
            "unit": "ns/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - B/op",
            "value": 67934,
            "unit": "B/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkAllMethods - allocs/op",
            "value": 1176,
            "unit": "allocs/op",
            "extra": "3183 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests",
            "value": 20710,
            "unit": "ns/op\t    8365 B/op\t     134 allocs/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - ns/op",
            "value": 20710,
            "unit": "ns/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - B/op",
            "value": 8365,
            "unit": "B/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentRequests - allocs/op",
            "value": 134,
            "unit": "allocs/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost",
            "value": 21551,
            "unit": "ns/op\t    8588 B/op\t     140 allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - ns/op",
            "value": 21551,
            "unit": "ns/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - B/op",
            "value": 8588,
            "unit": "B/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost - allocs/op",
            "value": 140,
            "unit": "allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50",
            "value": 675115,
            "unit": "ns/op\t  431044 B/op\t    6931 allocs/op",
            "extra": "1706 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - ns/op",
            "value": 675115,
            "unit": "ns/op",
            "extra": "1706 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - B/op",
            "value": 431044,
            "unit": "B/op",
            "extra": "1706 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCost50 - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1706 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency",
            "value": 673988,
            "unit": "ns/op\t  431586 B/op\t    6931 allocs/op",
            "extra": "1693 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - ns/op",
            "value": 673988,
            "unit": "ns/op",
            "extra": "1693 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - B/op",
            "value": 431586,
            "unit": "B/op",
            "extra": "1693 times\n4 procs"
          },
          {
            "name": "BenchmarkConcurrentEstimateCostLatency - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "1693 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour",
            "value": 47783,
            "unit": "ns/op\t    9348 B/op\t     156 allocs/op",
            "extra": "24457 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - ns/op",
            "value": 47783,
            "unit": "ns/op",
            "extra": "24457 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - B/op",
            "value": 9348,
            "unit": "B/op",
            "extra": "24457 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/1Hour - allocs/op",
            "value": 156,
            "unit": "allocs/op",
            "extra": "24457 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours",
            "value": 75441,
            "unit": "ns/op\t   18438 B/op\t     294 allocs/op",
            "extra": "15586 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - ns/op",
            "value": 75441,
            "unit": "ns/op",
            "extra": "15586 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - B/op",
            "value": 18438,
            "unit": "B/op",
            "extra": "15586 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/24Hours - allocs/op",
            "value": 294,
            "unit": "allocs/op",
            "extra": "15586 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days",
            "value": 173568,
            "unit": "ns/op\t   77405 B/op\t    1161 allocs/op",
            "extra": "6930 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - ns/op",
            "value": 173568,
            "unit": "ns/op",
            "extra": "6930 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - B/op",
            "value": 77405,
            "unit": "B/op",
            "extra": "6930 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/7Days - allocs/op",
            "value": 1161,
            "unit": "allocs/op",
            "extra": "6930 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days",
            "value": 496444,
            "unit": "ns/op\t  317681 B/op\t    4489 allocs/op",
            "extra": "2176 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - ns/op",
            "value": 496444,
            "unit": "ns/op",
            "extra": "2176 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - B/op",
            "value": 317681,
            "unit": "B/op",
            "extra": "2176 times\n4 procs"
          },
          {
            "name": "BenchmarkActualCostDataSizes/30Days - allocs/op",
            "value": 4489,
            "unit": "allocs/op",
            "extra": "2176 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS",
            "value": 53512,
            "unit": "ns/op\t    9693 B/op\t     176 allocs/op",
            "extra": "22024 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - ns/op",
            "value": 53512,
            "unit": "ns/op",
            "extra": "22024 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - B/op",
            "value": 9693,
            "unit": "B/op",
            "extra": "22024 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/AWS - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22024 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure",
            "value": 53014,
            "unit": "ns/op\t    9718 B/op\t     176 allocs/op",
            "extra": "22150 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - ns/op",
            "value": 53014,
            "unit": "ns/op",
            "extra": "22150 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - B/op",
            "value": 9718,
            "unit": "B/op",
            "extra": "22150 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Azure - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22150 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP",
            "value": 52892,
            "unit": "ns/op\t    9741 B/op\t     176 allocs/op",
            "extra": "22070 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - ns/op",
            "value": 52892,
            "unit": "ns/op",
            "extra": "22070 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - B/op",
            "value": 9741,
            "unit": "B/op",
            "extra": "22070 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/GCP - allocs/op",
            "value": 176,
            "unit": "allocs/op",
            "extra": "22070 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes",
            "value": 53268,
            "unit": "ns/op\t    9761 B/op\t     175 allocs/op",
            "extra": "22291 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - ns/op",
            "value": 53268,
            "unit": "ns/op",
            "extra": "22291 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - B/op",
            "value": 9761,
            "unit": "B/op",
            "extra": "22291 times\n4 procs"
          },
          {
            "name": "BenchmarkDifferentProviders/Kubernetes - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "22291 times\n4 procs"
          }
        ]
      }
    ]
  }
}