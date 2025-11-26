---
description: Convert existing tasks into actionable, dependency-ordered GitHub issues for the feature based on available design artifacts.
tools: ["github/github-mcp-server/issue_write"]
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Outline

1. Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` from repo root
   and parse FEATURE_DIR and AVAILABLE_DOCS list. All paths must be absolute. For single quotes in args like
   "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").
1. From the executed script, extract the path to **tasks**.
1. Get the Git remote by running:

```bash
git config --get remote.origin.url
```

### Only Proceed to Next Steps if the Remote is a GitHub URL

1. For each task in the list, use the GitHub MCP server to create a new issue in the repository that is
   representative of the Git remote.

**Warning**: NEVER create issues in repositories that do not match the remote URL.
