# Pre-Session Workflow

**Purpose**: Generate a validated step guide for the next development increment.

## Prompt Sequence

### 1. Step Preparation
Execute: [**Step Guide Preparation**](.prompts/guide-prepare.md)
- Analyze current progress and identify decisions needed for the next step
- Review `checkpoint.md` and `review.md` for blockers or prerequisites

### 2. Step Guide Generation  
Execute: [**Step Guide Generation**](.prompts/guide-generate.md)
- Generate comprehensive `step-guide.md` for the **IMMEDIATE NEXT** step from `roadmap.md`
- Include implementation details, code examples, and completion criteria

### 3. Step Guide Validation
Execute: [**Step Guide Validation**](.prompts/guide-validate.md)
- Review step guide for errors, security issues, and pattern compliance
- Refine guide based on validation findings

**Output**: Production-ready `step-guide.md` for session execution