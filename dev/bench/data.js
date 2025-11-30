window.BENCHMARK_DATA = {
  "lastUpdate": 1764514607439,
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
      }
    ]
  }
}