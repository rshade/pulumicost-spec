# TypeScript Client SDK for FinFocus Plugin Ecosystem

**Feature Branch**: `028-typescript-sdk`
**Status**: Specification Complete - Ready for Planning Phase
**Created**: 2026-01-16

## Overview

This feature specification defines a comprehensive browser-first TypeScript Client SDK for
the FinFocus plugin ecosystem. The SDK enables frontend applications to directly integrate
with FinFocus cost source plugins using the Connect protocol with JSON/HTTP encoding.

## Key Deliverables

### Core SDK Components

1. **Core Client Library** (`finfocus-client`)
   - Three client classes covering all 22 RPC methods
   - CostSourceClient (11 methods)
   - ObservabilityClient (3 methods)
   - RegistryClient (8 methods)

2. **Builder Pattern** for complex message construction
   - ResourceDescriptorBuilder (10 fluent methods)
   - RecommendationFilterBuilder (16 filter fields)
   - FocusRecordBuilder (FOCUS 1.2/1.3 support)

3. **Comprehensive Validation Layer**
   - Request validators for all RPC methods
   - Response validators with FOCUS compliance checking
   - ISO 4217 currency validation (180+ codes)
   - Billing mode validation (44+ modes)

4. **Helper Utilities**
   - Cost calculations (hourly ↔ monthly conversion)
   - Growth projections (linear and exponential)
   - Async pagination iterator with configurable limits
   - Recommendation filtering (16 field support)
   - Configurable sorting

5. **REST API Wrapper** (`finfocus-rest`)
   - 22 HTTP endpoint mappings
   - JSON ↔ protobuf automatic transformation
   - Standard HTTP error code mapping

6. **Framework Plugins**
   - Express middleware (`finfocus-express`)
   - Fastify plugin (`finfocus-fastify`)
   - NestJS module (`finfocus-nestjs`)

### Documentation & Examples

- Quick start guide (5-line working example)
- Comprehensive API reference with JSDoc
- Browser compatibility guide (Chrome 60+, Firefox 55+, Safari 12+)
- REST integration guide
- Framework-specific integration examples

## Architecture Decisions

### Connect Protocol (vs gRPC/gRPC-Web)

- **Rationale**: Browser-native protocol with JSON/HTTP encoding
- **Benefits**: CORS-friendly, debuggable, smaller bundle size
- **Trade-off**: ~30% larger payloads than binary gRPC

### ES2018 Compilation Target

- **Rationale**: Balance modern TypeScript with broad browser support
- **Coverage**: Chrome 60+, Firefox 55+, Safari 12+
- **Trade-off**: ~10% larger bundle vs no polyfills needed for most features

### npm Workspaces Architecture

- **Structure**: 5 packages (client, rest, express, fastify, nestjs) + examples
- **Benefits**: Single source of truth, consistent versioning, simplified CI/CD
- **Alternative Rejected**: Monorepo (lerna), separate repos

### Version-Controlled Generated Code

- **Location**: `sdk/typescript/packages/client/src/generated/`
- **Benefits**: Fast npm installs (no post-install build), better debugging, clear PR diffs
- **Trade-off**: Slightly larger repository size (~200 KB gzipped)

## Technical Specifications

### Bundle Size Targets

- Core client only: ≤ 40 KB (minified + gzipped)
- Includes builders, validation, and helpers

### Performance SLAs

- Name() RPC: < 100ms
- Supports() RPC: < 50ms
- GetActualCost() RPC: < 2s (24h), < 10s (30d)
- GetProjectedCost() RPC: < 200ms
- DryRun() RPC: < 100ms p99

### Test Coverage

- Minimum 80% code coverage via Vitest
- Unit tests for builders, validators, helpers
- Integration tests for all 22 RPC methods with MSW mock server

### Browser Compatibility

- Chrome 60+
- Firefox 55+
- Safari 12+
- IE not supported (ES2018 features assumed)

## File Structure

```tree
specs/028-typescript-sdk/
├── spec.md                              # Main specification document
├── checklists/
│   └── requirements.md                  # Comprehensive validation checklist
└── README.md                            # This file
```

## Specification Validation

✅ All mandatory sections complete and validated
✅ 3 prioritized user stories with independent tests
✅ 14 functional requirements with clear success criteria
✅ 14 measurable outcomes with concrete targets
✅ 6 edge cases documented
✅ 7 key entities fully defined
✅ No unclear requirements (NEEDS CLARIFICATION markers)
✅ Markdown validation passing (0 errors)
✅ All requirements implementable and technology-specific

See `checklists/requirements.md` for detailed validation checklist.

## Next Steps

### Planning Phase

The specification is complete and ready for the planning phase:

```bash
/speckit.plan
```

This will:

1. Create a detailed implementation plan with task breakdown
2. Identify critical path and dependencies
3. Generate actionable tasks with timeline
4. Establish quality gates and success metrics

### Implementation Phases (11 Weeks)

**Phase 1**: Foundation (Week 1)

- Initialize monorepo with npm workspaces ✅ (already started)
- Configure buf code generation ✅ (updated)
- Set up TypeScript configurations ✅ (created)

**Phase 2-11**: See detailed plan in `spec.md`

## Key Dates & Milestones

- **Specification**: Complete (2026-01-16)
- **Planning**: Pending (`/speckit.plan`)
- **Implementation**: To be scheduled

## Related Documentation

- Main spec: `spec.md`
- Requirements checklist: `checklists/requirements.md`
- Architecture decisions: See "Architecture Decisions" section above
- Implementation timeline: See Phase details in `spec.md`

## Contacts & References

**Feature Champion**: FinFocus Product Team
**Related Specs**:

- 026-focus-1-3-migration (FOCUS 1.3 support)
- 032-plugin-dry-run (DryRun capability)
- 029-plugin-info-rpc (GetPluginInfo support)

## Quality Gates

All quality gates passed:

- [x] Specification completeness (all sections filled)
- [x] Requirement clarity (no [NEEDS CLARIFICATION] markers)
- [x] Success criteria measurability (concrete targets)
- [x] Architecture decisions documented (rationale provided)
- [x] Implementation feasibility (all requirements implementable)
- [x] Markdown validation (0 linting errors)
- [x] Checklist validation (comprehensive validation passing)

**Status**: ✅ READY FOR PLANNING PHASE
