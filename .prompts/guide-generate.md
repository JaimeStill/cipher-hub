# Step Guide Generation Prompt

## Request

Based on the provided Cipher Hub project context, generate a complete `step-guide.md` file for the next uncompleted step in the roadmap.

## Instructions

1. **Identify Next Step**: Review the roadmap to determine the next unmarked step that should be implemented
2. **Follow Format**: Use the existing `step-guide.md` as the exact structural template
3. **Leverage Context**: Build upon the current Go source code, established patterns, and completed infrastructure
4. **Maintain Standards**: Ensure 1-hour implementation scope, Go best practices, security-first design, and >95% test coverage

## Output Requirements

Generate a complete step guide with:
- Accurate step identification from roadmap progression
- Technical specifications appropriate for the identified step
- Implementation details that build on existing codebase patterns
- Comprehensive testing and verification procedures
- Security considerations relevant to the step's functionality
- Modern Go idioms, error handling, and code organization
- Precise but concise writing that respects context and token usage
- Avoid emojis
- Artifact only - no commentary

The guide should be immediately executable as a standalone development session following your established development philosophy.