# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is the **`.claude/agents/` directory** of the FinFocus Specification repository, containing specialized agent
configurations for the FinFocus ecosystem. This directory houses three specialized agent configurations that enable
context-aware assistance across the multi-repository FinFocus ecosystem.

## Agent Configuration Architecture

### Agent System Design

The FinFocus ecosystem uses **specialized agents** to handle different aspects of development across three repository types:

**Repository Detection Pattern:**

- **finfocus-spec**: Contains `proto/finfocus/costsource.proto` and `schemas/pricing_spec.schema.json`
- **finfocus-core**: Contains `cmd/finfocus/` and `internal/{pluginhost,engine,ingest,spec,cli}/`
- **finfocus-plugin-\***: Contains `cmd/finfocus-<name>/` and `plugin.manifest.json`

### Three Specialized Agents

### 1. FinFocus Senior Engineer (`finfocus-senior-engineer.md`)

**Primary Responsibilities:**

- **Architecture & Implementation**: Translate PM requirements into robust, testable Go code
- **Cross-Repository Consistency**: Maintain protocol alignment across spec, core, and plugin repos
- **Quality & Tooling Excellence**: CI/CD pipelines, buf configurations, Go module hygiene
- **Security & Stability**: Input validation, plugin isolation, error handling

**Key Capabilities:**

- gRPC interface and JSON schema consistency management
- Protocol versioning and backward compatibility
- Comprehensive testing with mocks and fixtures
- Developer ergonomics and tooling

**Workflow Pattern:**

1. Repository detection and context analysis
2. Engineering Status report (gaps, tech debt, architecture recommendations)
3. Implementation plan with pseudocode and dependency analysis
4. Complete implementation with tests and documentation
5. CI validation and PR documentation

### 2. FinFocus Technical Writer (`finfocus-technical-writer.md`)

**Primary Responsibilities:**

- **Documentation Creation**: README files, API docs, tutorials, examples
- **Developer Onboarding**: CONTRIBUTING guides, plugin author guides
- **Release Communication**: CHANGELOG entries, release notes, migration guides
- **Content Standards**: Copy-paste-ready examples, tested code snippets

**Key Capabilities:**

- Repository-aware documentation strategy
- Cross-repository terminology consistency
- Comprehensive example validation
- Visual documentation with mermaid diagrams

**Content Standards:**

- All code snippets must be complete and runnable
- Active voice with clear, specific instructions
- Expected outputs and error scenario documentation
- Mermaid diagrams for architecture illustrations

### 3. FinFocus Product Manager (`finfocus-product-manager.md`)

**Primary Responsibilities:**

- **Backlog Management**: Actionable tickets with user stories and acceptance criteria
- **Cross-Repo Coordination**: Dependencies and sequencing (spec → core → plugins)
- **MVP Focus**: 28-day timeline with tight scope management
- **Release Planning**: Version coordination and compatibility management

**Key Capabilities:**

- Repository detection and status assessment
- Cross-repository change protocol management
- Issue template standardization
- MVP roadmap tracking (Week 1-4 progression)

**Program Invariants:**

- No raw CUR parsing; vendor API integration only
- Plugin discovery at `~/.finfocus/plugins/<name>/<version>/<binary>`
- gRPC protocol as single source of truth
- Apache-2.0 licensing across all repos

## Agent Selection Guidelines

### Use Senior Engineer Agent When

- Implementing new cost calculation features
- Ensuring protocol consistency between repositories
- Architecture decisions and technical implementations
- Quality assurance and testing framework work
- Cross-repo technical debt resolution

### Use Technical Writer Agent When

- Creating or updating README files
- Writing API documentation or tutorials
- Documenting breaking changes or migration guides
- Creating plugin author guides
- Release documentation and changelogs

### Use Product Manager Agent When

- Planning sprint work and creating backlog items
- Coordinating releases across repositories
- Managing cross-repo dependencies
- Creating project roadmaps and timelines
- Issue tracking and milestone management

## Development Context

### Current Repository Status (finfocus-spec)

**Repository Type**: Specification repository
**Version**: v0.4.6 (production-ready)
**Key Components**:

- gRPC service definitions (`proto/finfocus/v1/costsource.proto`)
- JSON Schema validation (`schemas/pricing_spec.schema.json`)
- Production Go SDK (`sdk/go/`)
- Comprehensive testing framework (`sdk/go/testing/`)
- Cross-vendor examples (`examples/`)

**Architecture Highlights**:

- 8 RPC methods (Name, Supports, GetActualCost, GetProjectedCost, GetPricingSpec, EstimateCost, GetRecommendations, GetBudgets)
- 44+ billing modes across all major cloud providers
- Multi-level conformance testing (Basic/Standard/Advanced)
- Complete CI/CD pipeline with validation and benchmarks

### Agent Interaction Patterns

**Sequential Collaboration:**

1. **Product Manager** creates backlog items and defines requirements
2. **Senior Engineer** implements features with architectural considerations
3. **Technical Writer** documents implementation and creates user guides

**Cross-Repository Coordination:**

1. **Product Manager** identifies cross-repo dependencies
2. **Senior Engineer** ensures protocol consistency
3. **Technical Writer** maintains documentation alignment

**Quality Assurance Integration:**

- All agents emphasize comprehensive testing requirements
- Senior Engineer focuses on technical implementation
- Technical Writer ensures documentation completeness
- Product Manager tracks delivery against acceptance criteria

## Agent Invocation Patterns

### Practical Examples

**Senior Engineer Agent:**

```text
I need to implement support for calculating storage costs in the pricing engine
```

**Results in**: Architecture analysis, implementation plan, code + tests, CI validation

**Technical Writer Agent:**

```text
I just finished implementing the AWS cost plugin. Can you help me create proper documentation?
```

**Results in**: README, API docs, examples, plugin author guide

**Product Manager Agent:**

```text
I need to create issues for implementing the actual cost pipeline in finfocus-core
```

**Results in**: Structured issues, acceptance criteria, cross-repo dependencies

### Quick Agent Selection Decision Tree

**Question: What type of work?**

- Code/Architecture/Implementation → **Senior Engineer**
- Documentation/Examples/Guides → **Technical Writer**
- Planning/Issues/Coordination → **Product Manager**

**Question: What scope?**

- Single repository → Choose based on work type
- Cross-repository → Start with **Product Manager**
- Breaking changes → **All three agents** in sequence

## Repository-Specific Considerations

### In finfocus-spec (current repository)

- Changes here ripple to core and plugin repositories
- Breaking changes require coordinated releases across all repos
- buf.yaml and protobuf management is critical for ecosystem consistency
- Schema changes must maintain backward compatibility

### When Working with finfocus-core

- CLI changes need comprehensive documentation updates
- Engine modifications require plugin compatibility checks
- Integration tests span multiple repositories
- Plugin host changes affect all plugin implementations

### When Working with finfocus-plugin-\*

- Plugin manifest validation against spec requirements
- Conformance testing requirements (Basic/Standard/Advanced)
- Version compatibility with spec changes
- Vendor-specific integration patterns

## Build Commands Integration

All agents are aware of the standard build commands:

```bash
# Core development workflow
make generate    # Generate Go SDK from protobuf
make test        # Run comprehensive test suite
make validate    # Tests + linting + schema validation
make lint        # All linting (Go, buf, markdown, YAML)
make clean       # Clean generated files

# Testing variations
go test -bench=. -benchmem ./sdk/go/testing/  # Performance benchmarks
go test -v -run TestConformance               # Conformance validation
npm run validate                              # Schema validation

# Agent configuration management
# Validate agent YAML frontmatter structure
# Test agent configuration loading and parsing
```

## Quality Standards

### Code Quality (Senior Engineer)

- Idiomatic Go with comprehensive error handling
- Complete test coverage with mocks and fixtures
- Protocol versioning and backward compatibility
- Security best practices for data parsing

### Documentation Quality (Technical Writer)

- All examples must be tested and runnable
- Clear, actionable instructions with expected outputs
- Cross-repository terminology consistency
- Visual diagrams for complex workflows

### Product Quality (Product Manager)

- Clear acceptance criteria and definition of done
- Cross-repo dependency tracking
- MVP timeline adherence with scope management
- Version compatibility and upgrade path planning

## Agent Collaboration Framework

**Issue Creation Flow:**

1. Product Manager creates structured issue with acceptance criteria
2. Senior Engineer provides implementation plan and technical analysis
3. Technical Writer identifies documentation requirements
4. All agents coordinate on cross-repo impacts

**Release Coordination:**

1. Product Manager sequences changes across repositories
2. Senior Engineer implements with backward compatibility
3. Technical Writer creates migration guides and release notes
4. All agents validate against quality gates

**Quality Assurance:**

- Each agent validates their domain (code, docs, product requirements)
- Cross-agent review for completeness and consistency
- Integration testing across all agent outputs

## Agent System Troubleshooting

### If Agent Seems Unaware of Context

1. Ensure you're in the correct repository directory
2. Check that repository detection patterns match current structure
3. Verify agent has access to recent files and changes
4. Confirm the agent configuration YAML frontmatter is valid

### If Cross-Repository Coordination Fails

1. **Start with Product Manager agent** for dependency mapping and impact analysis
2. **Use Senior Engineer** for technical implementation impact across repos
3. **Use Technical Writer** for communication strategy and documentation alignment
4. **Validate sequence**: Always spec → core → plugins for breaking changes

### If Agent Selection Is Unclear

1. **Single focused task** → Choose by work type (code/docs/planning)
2. **Complex multi-faceted work** → Start with Product Manager for coordination
3. **Emergency fixes** → Senior Engineer for technical issues, Technical Writer for communication
4. **Release preparation** → All three agents in coordinated sequence

### Common Coordination Patterns

- **New feature**: Product Manager → Senior Engineer → Technical Writer
- **Bug fix**: Senior Engineer (+ Technical Writer if user-facing)
- **Breaking change**: Product Manager → Senior Engineer → Technical Writer (parallel execution)
- **Documentation update**: Technical Writer (+ Senior Engineer for technical review)

This agent configuration system ensures comprehensive coverage of the FinFocus ecosystem development needs while maintaining
high standards across technical implementation, documentation quality, and product delivery.
