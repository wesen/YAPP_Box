---
Title: Reference — Debate Questions: Builder Contract Migration
Ticket: YAPP-DSL-GAPS-001
Status: active
Topics:
    - yapp
    - architecture
    - codegen
    - debate
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/debate/01-debate-builder-contract-migration-format-and-candidates.md
      Note: Candidate profiles and debate format
    - Path: /home/manuel/code/others/YAPP_Box/ttmp/2025/11/15/YAPP-DSL-GAPS-001-dsl-feature-gaps-analysis-missing-yapp-arrays/analysis/04-2025-11-17-builder-contract-and-codegen-options.md
      Note: Source analysis with Options A-E
ExternalSources: []
Summary: 8-question set mapped to candidates and decision points for builder contract migration debate.
LastUpdated: 2025-11-17
---

# Reference — Debate Questions: Builder Contract Migration

## Question Set

### Question 1: Go/No-Go — Is refactoring the builder contract urgent?

**Context:** Current system works but has pain points (YAML re-marshal, untyped maps, [][]any).

**Decision point:** Should we invest in refactoring now or defer?

**Primary candidates:**
- Alex Chen (The Pragmatist)
- Jordan Rivera (The Architect)
- Casey Thompson (The New Hire)

**Key tensions:**
- Working code vs technical debt
- Immediate ROI vs long-term maintainability
- Onboarding friction vs migration risk

---

### Question 2: Which path first — Path 1 (A→B) or Path 2 (A→C)?

**Context:** Two staged approaches with different end goals.

**Path 1 (A→B):**
- Option A: Generate `Decode()` per module, keep `[][]any` output
- Option B: Add `ArrayDecl` IR to handle multi-array modules

**Path 2 (A→C [+D]):**
- Option A: Generate `Decode()` per module
- Option C: Store typed slices in `Model`, typed collectors
- Option D (optional): Capability interfaces for special outputs

**Decision point:** Which path delivers more value sooner?

**Primary candidates:**
- Jordan Rivera (The Architect)
- Alex Chen (The Pragmatist)
- Sam Park (The Codegen Maintainer)

**Key tensions:**
- Incremental vs comprehensive change
- Interface stability vs type safety
- Codegen complexity vs runtime clarity

---

### Question 3: Developer ergonomics — Which path reduces friction fastest?

**Context:** Current pain: YAML marshal/unmarshal in every builder, untyped map juggling.

**Decision point:** Which approach most improves the module authoring experience?

**Primary candidates:**
- Sam Park (The Codegen Maintainer)
- Casey Thompson (The New Hire)
- The Orchestrator (`features.go`)

**Key tensions:**
- Generated code readability vs hand-written clarity
- Fewer files to edit vs more magic
- Learning curve for new patterns

---

### Question 5: Special-cases — How do we retire `cutoutFeatureModule`?

**Context:** Cutouts return `map[string][][]any` not `[][]any`, handled by custom `FeatureModule` impl.

**Decision point:** How do we unify cutouts with other modules?

**Primary candidates:**
- Six-Faced Friend (`cutouts/`)
- The Orchestrator (`features.go`)
- The Contract (`registry/schema.go`)

**Key tensions:**
- Single interface vs capability-based
- IR abstraction vs direct emission
- Migration risk for working module

---

### Question 6: Codegen scope — What additions are required?

**Context:** Current schemagen generates structs, tests, validators. New needs: `Decode()`, wrappers, possibly IR.

**Decision point:** How much new generation is feasible without template spaghetti?

**Primary candidates:**
- The Generator (`schemagen/`)
- Sam Park (The Codegen Maintainer)
- Jordan Rivera (The Architect)

**Key tensions:**
- Template complexity vs output quality
- Generation time vs developer time saved
- Debugging generated code vs hand-written code

---

### Question 7: Type-safety end-to-end — Move typed slices into `Model` now or later?

**Context:** Option C proposes `Model` stores `[]PcbStandsItem` not `[]map[string]any`.

**Decision point:** Do we need full type safety now or is `Decode()` enough?

**Primary candidates:**
- Jordan Rivera (The Architect)
- The Orchestrator (`features.go`)
- The Contract (`registry/schema.go`)

**Key tensions:**
- Type safety vs interface churn
- Compile-time checks vs runtime flexibility
- Migration scope (7 modules + Model + features.go)

---

### Question 9: Enforceability — How do we keep modules honest post-change?

**Context:** After refactor, how do we prevent regression to old patterns?

**Decision point:** What tooling/docs/tests ensure new patterns stick?

**Primary candidates:**
- The Contract (`registry/schema.go`)
- Sam Park (The Codegen Maintainer)
- Casey Thompson (The New Hire)

**Key tensions:**
- Linter rules vs trust
- Documentation vs examples
- Compile-time enforcement vs runtime checks

---

### Question 10: Future extensibility — Accommodating new output shapes?

**Context:** What if future modules need outputs beyond single/multi arrays?

**Decision point:** Does our chosen path handle unknown future needs?

**Primary candidates:**
- Jordan Rivera (The Architect)
- Six-Faced Friend (`cutouts/`)
- The Generator (`schemagen/`)

**Key tensions:**
- Flexibility vs simplicity
- Anticipating needs vs YAGNI
- Interface evolution vs stability

---

## Question Flow

```
Q1 (Go/No-Go) → Q2 (Path choice)
                    ↓
        ┌───────────┴───────────┐
        ↓                       ↓
    Q3 (Ergonomics)         Q5 (Special-cases)
        ↓                       ↓
    Q6 (Codegen scope)      Q7 (Type-safety)
        ↓                       ↓
        └───────────┬───────────┘
                    ↓
            Q9 (Enforceability)
                    ↓
            Q10 (Extensibility)
```

## Candidate Participation Matrix

| Question | Primary Candidates | Supporting Cast |
|----------|-------------------|-----------------|
| Q1 | Pragmatist, Architect, New Hire | - |
| Q2 | Architect, Pragmatist, Codegen Maintainer | Contract, Orchestrator |
| Q3 | Codegen Maintainer, New Hire, Orchestrator | Generator |
| Q5 | Six-Faced Friend, Orchestrator, Contract | Pragmatist |
| Q6 | Generator, Codegen Maintainer, Architect | - |
| Q7 | Architect, Orchestrator, Contract | Pragmatist |
| Q9 | Contract, Codegen Maintainer, New Hire | - |
| Q10 | Architect, Six-Faced Friend, Generator | Contract |

## References

- [Candidate Profiles](../debate/01-debate-builder-contract-migration-format-and-candidates.md)
- [Source Analysis](../analysis/04-2025-11-17-builder-contract-and-codegen-options.md)
