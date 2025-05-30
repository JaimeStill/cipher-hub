# Roadmap Update

## Variables

- [STEP_NUMBER] - the roadmap step completed

## Prompt

Mark Step [STEP_NUMBER] as completed and update the `roadmap.md` artifact as follows:
- Identify any decisions needed or roadblocks and capture them as bullet points under the **Roadblocks** section.
- Add bullet points under the completed step describing key implementation decisions
- Add integration notes under future steps that will build upon this implementation

### Guidelines

- Use the Go project source code and comments as the source of truth
- Keep bullet point details concise (1-2 lines max)
- Avoid emojis
- Artifact only - no commentary

If the step doesn't exist, respond with "Step not found."