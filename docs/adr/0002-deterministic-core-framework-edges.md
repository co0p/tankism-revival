# ADR 0002: Keep gameplay deterministic and Ebitengine at the edges

- **Status:** Accepted
- **Date:** 2026-09-28

## Context

Collision, combat, AI, and lifecycle behavior need fast automated tests. Directly coupling those rules to Ebitengine input, images, windows, and audio makes tests fragile and hides state transitions.

## Decision

Put world data and gameplay systems under `game/` without Ebitengine imports. Translate live input into snapshots at the edge. Keep rendering, lighting, particles' Ebitengine draw code, resource handles, and audio playback in adapters. The executable composes these pieces.

## Consequences

Pure rules can be unit and integration tested in memory. Some adapter code is necessarily manually tested. Rendering components use stable resource identifiers rather than framework object handles.
