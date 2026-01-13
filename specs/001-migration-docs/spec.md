# Feature Specification: Add Migration Documentation

**Feature Branch**: `001-migration-docs`
**Created**: 2026-01-12
**Status**: Draft
**Input**: User description: "Add Migration Documentation (High Priority) - The PR renames
environment variables from PULUMICOST*\* to FINFOCUS*\*, but migration guidance is missing.
This could break existing deployments."

## Clarifications

### Session 2026-01-12

- Q: Should the SDK provide temporary backwards compatibility for old environment variable
  names? → A: Full backwards compatibility - SDK reads both old and new names during
  transition period

## User Scenarios & Testing _(mandatory)_

Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
you should still have a viable MVP (Minimum Viable Product) that delivers value.

Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
Think of each story as a standalone slice of functionality that can be:

- Developed independently
- Tested independently
- Deployed independently
- Demonstrated to users independently

### User Story 1 - Environment Variable Migration Guide (Priority: P1)

As a DevOps engineer with existing PulumiCost deployments, I need clear documentation on
how to update environment variables so that my plugins continue functioning after the rename
to FinFocus.

**Why this priority**: Environment variables are critical for plugin operation. Without
clear guidance, deployments will fail, causing immediate business impact.

**Independent Test**: Can be fully tested by verifying that all 6 environment variable
renames are documented with before/after examples and that users can successfully update
their configurations.

**Acceptance Scenarios**:

1. **Given** a user has `PULUMICOST_PLUGIN_PORT=50051` set in their environment, **When**
   they consult the migration guide, **Then** they can identify that this should be changed
   to `FINFOCUS_PLUGIN_PORT`.

2. **Given** a user has multiple environment variables set (PORT, LEVEL, FILE), **When**
   they review the migration guide, **Then** they understand which variables must be renamed
   and which have fallback behavior.

3. **Given** a user is updating their Kubernetes manifests, **When** they follow the
   migration steps, **Then** they can successfully deploy their plugins with the new
   environment variables.

---

### User Story 2 - Plugin Directory Migration (Priority: P1)

As a plugin administrator, I need clear instructions on migrating plugin directories so that
existing plugins are discovered correctly after the rename.

**Why this priority**: Plugin discovery paths are fundamental to the plugin architecture.
Misplaced plugins will cause startup failures.

**Independent Test**: Can be fully tested by verifying that the migration command correctly
moves plugins from the old directory to the new one and that plugins are subsequently
discovered.

**Acceptance Scenarios**:

1. **Given** a user has plugins installed in `~/.pulumicost/plugins/`, **When** they follow
   the migration guide, **Then** their plugins are moved to `~/.finfocus/plugins/`.

2. **Given** a user runs the migration command on a system where the old directory does not
   exist, **When** they follow the guide, **Then** they receive a clear message that no
   migration is needed.

3. **Given** a user has custom plugin configuration files referencing the old path, **When**
   they consult the migration guide, **Then** they understand which files need updating.

---

### User Story 3 - LLM-Friendly Migration Manifest (Priority: P2)

As an AI coding assistant helping with repository migrations, I need a machine-readable
manifest so that I can automatically identify and apply migration changes in downstream
repositories.

**Why this priority**: Enables automated migration assistance across multiple repositories,
reducing manual effort for users with many deployments.

**Independent Test**: Can be fully tested by verifying that an AI assistant can parse the
manifest and correctly identify all required changes without human intervention.

**Acceptance Scenarios**:

1. **Given** an AI assistant has access to the `llm-migration.json` manifest, **When** it
   analyzes a repository, **Then** it can identify all environment variables that need
   renaming.

2. **Given** an AI assistant processes the migration manifest, **When** it generates changes,
   **Then** the changes match the documented migration steps exactly.

3. **Given** a user wants to verify the migration manifest, **When** they inspect the JSON
   file, **Then** it is valid JSON and contains all expected migration fields.

---

### User Story 4 - Changelog and README Updates (Priority: P2)

As a project maintainer, I need the migration information integrated into existing
documentation so that users naturally discover it when reading release notes.

**Why this priority**: Users expect migration information in standard locations (CHANGELOG,
README). Missing documentation leads to support burden and deployment failures.

**Independent Test**: Can be fully tested by verifying that users can find migration
guidance from the CHANGELOG entry for the rename release and from the README support
section.

**Acceptance Scenarios**:

1. **Given** a user reads the CHANGELOG for version 0.5.0, **When** they look for migration
   information, **Then** they find a clear migration section with links to detailed guidance.

2. **Given** a user reads the README, **When** they look for help with upgrading, **Then**
   they find a link to the migration guide in the Support section.

3. **Given** a user follows the link from CHANGELOG, **When** they navigate to the migration
   guide, **Then** they can complete their migration successfully.

---

### Edge Cases

- **Mixed variable usage**: Custom scripts referencing both PULUMICOST*\* and FINFOCUS*\*
  variables simultaneously - SDK supports both during transition period.
- **Plugin installation methods**: Plugins installed via manual copy, package manager, or CLI
  - migration guide provides commands for all methods.
- **Multiple plugin directories**: Users with plugins in custom directories - guide includes consolidation steps.
- **Environment-specific configs**: Dev/staging/production configurations - guide recommends updating all environments systematically.

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide migration documentation for all renamed environment
  variables (PORT, LOG_LEVEL, LOG_FILE, LOG_FORMAT, TRACE_ID, TEST_MODE).
- **FR-001a**: SDK MUST support backwards compatibility by reading both old
  (`PULUMICOST_*`) and new (`FINFOCUS_*`) environment variable names during migration
  transition period.
- **FR-002**: System MUST document the plugin directory path change from
  `~/.pulumicost/plugins/` to `~/.finfocus/plugins/`.
- **FR-003**: Migration guide MUST include step-by-step instructions that users can follow
  without additional research.
- **FR-004**: System MUST provide a machine-readable JSON manifest for AI-assisted migration.
- **FR-005**: Migration information MUST be accessible from the CHANGELOG for the release
  that introduces the rename.
- **FR-006**: Migration information MUST be linked from the README Support section.
- **FR-007**: Migration documentation MUST be validated against markdown linting standards.

### Key Entities

- **MigrationGuide**: Human-readable documentation explaining the rename and migration steps.
- **LLMManifest**: Machine-readable JSON file containing migration metadata for AI
  assistants.
- **EnvironmentVariableMapping**: Documented mapping from old `PULUMICOST_*` variables to
  new `FINFOCUS_*` variables.
- **PluginPathMapping**: Documented mapping from old plugin discovery path to new path.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Users can complete environment variable migration in under 5 minutes following
  the documented steps (measured by timed manual execution).
- **SC-002**: 100% of required environment variable renames are documented with clear
  before/after examples (verified by checklist of all 6 variables).
- **SC-003**: Plugin directory migration can be completed with a single documented command.
- **SC-004**: AI assistants can parse the migration manifest and generate correct migration
  changes without human intervention.
- **SC-005**: Migration documentation passes markdown linting validation with zero errors.
- **SC-006**: Users can find migration guidance from either CHANGELOG or README within 3
  navigation steps.

### Assumptions

- Users have shell access to run migration commands.
- Users are familiar with their deployment configuration locations.
- SDK provides backwards compatibility for old variable names during transition.
- Users will follow documented steps rather than attempting custom migration approaches.
