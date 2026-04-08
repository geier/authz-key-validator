# Agent Workflow

## Branding

- Product name is **prokube.ai** (lowercase, with ".ai" suffix)
- Never write "ProKube", "Prokube", or "prokube.ai" (missing dot)

## Worker Agent

Workers process a single issue, then exit.

```
┌──────────────────────────────────────────────────────┐
│  1. RECEIVE TASK                                    │
│     - Supervisor assigns a specific issue number      │
└───────────────────────────────────────────────────┬───┘
                                                  │
┌─────────────────────────────────────────────────────▼───┐
│  2. CREATE WORKTREE                                   │
│     git fetch origin main                             │
│     git worktree add worktrees/issue-N -b feature/issue-N origin/main │
│     cd worktrees/issue-N                              │
└───────────────────────────────────────────────┬───────┘
                                            │
┌───────────────────────────────────────────────▼───────┐
│  3. CLAIM THE ISSUE                                  │
│     gh issue edit <N> --add-assignee @me             │
│     gh issue edit <N> --add-label in-progress        │
└───────────────────────────────────────────────┬───────┘
                                            │
┌───────────────────────────────────────────────▼───────┐
│  4. IMPLEMENT                                         │
│     - Follow AGENTS.md for coding standards           │
│     - Run tests: make test                           │
│     - Commit specific files (never `git add .`)       │
└───────────────────────────────────────────────┬───────┘
                                            │
┌───────────────────────────────────────────────▼───────┐
│  5. CREATE PR                                         │
│     git push -u origin feature/issue-N               │
│     gh pr create --title "..." --body "Closes #N"   │
│     gh copilot-review <pr-number>                     │
└───────────────────────────────────────────────────────┘
```

### Worker Rules

1. Work only in `worktrees/issue-N/`
2. Stage specific files (`git add file1 file2`)
3. Never use `git add .` or `git add -A`

## Supervisor Loop

```
┌─────────────────────────────────────────────────────────┐
│  1. ASSIGN WORK                                     │
│     gh issue list --search "no:assignee label:ready no:blocked-by" │
│     Spawn one worker per issue                       │
├─────────────────────────────────────────────────────────┤
│  2. MONITOR                                         │
│     - Check worker outputs                           │
│     - Handle escalations (needs-supervisor label)    │
│     - Verify PR quality                              │
├─────────────────────────────────────────────────────────┤
│  3. REVIEW PRs                                      │
│     Phase 1: Hygiene (base branch, files changed)   │
│     Phase 2: Deep review (spawn review workers)     │
│     Phase 3: Copilot review cycle                   │
├─────────────────────────────────────────────────────────┤
│  4. CLEANUP                                         │
│     - Merge PRs, update dependencies               │
│     - Remove resolved blockers                       │
└─────────────────────────────────────────────────────────┘
```

### PR Review Cycle

**Phase 1: Hygiene Check**
- Verify base branch (should be main)
- Check files changed (should only touch relevant files)
- Detect cross-PR conflicts

**Phase 2: Deep Review**
- Spawn review workers to check:
  - Issue requirements met
  - Code follows AGENTS.md
  - Edge cases handled

**Phase 3: Copilot Review**
```bash
gh copilot-review <pr-number>
```
┌─────────────────────────────────────────────────────────┐
│  1. ASSIGN WORK                                     │
│     gh issue list --search "no:assignee label:ready no:blocked-by" │
│     Spawn one worker per issue                       │
├────────────────────────────────────────────────────────┤
│  2. MONITOR                                         │
│     - Check worker outputs                           │
│     - Handle escalations (needs-supervisor label)    │
│     - Verify PR quality                              │
├───────────────────────────────────────────────────────┤
│  3. REVIEW PRs                                      │
│     Phase 1: Hygiene (base branch, files changed)   │
│     Phase 2: Deep review (spawn review workers)     │
│     Phase 3: Copilot review cycle                   │
├───────────────────────────────────────────────────────┤
│  4. CLEANUP                                         │
│     - Merge PRs, update dependencies                │
│     - Remove resolved blockers                       │
└───────────────────────────────────────────────────────┘
```

## Triage Agent

Converts user requests into GitHub issues with:
1. Title and context
2. Acceptance criteria
3. Files affected (if known)
4. Priority label

## Labels

| Label | Meaning |
|-------|---------|
| `ready` | Well-scoped, ready to work |
| `in-progress` | Currently being worked on |
| `needs-supervisor` | Escalation — supervisor must review |
| `task` | Implementation task |
| `priority:high` | Time-sensitive |

## Escalation

Workers escalate by:
1. Adding `needs-supervisor` label
2. Commenting with specific questions
3. Unassigning themselves

## Required GitHub CLI Extensions

Install the Copilot review extension:

```bash
gh extension install ChrisCarini/gh-copilot-review
```

Usage:
```bash
gh copilot-review <pr-number>
```

## Quick Reference

```bash
# Find work
gh issue list --search "no:assignee label:ready no:blocked-by"

# Claim issue
gh issue edit <N> --add-assignee @me --add-label in-progress

# Check escalations
gh issue list --label needs-supervisor

# Request Copilot review (requires gh-copilot-review extension)
gh copilot-review <pr-number>
```

## Project-Specific Details

| Item | Path/Command |
|------|-------------|
| Tests | `make test` |
| Build | `make build` |
| Run | `make run` |
| Lint | `golangci-lint run` |

File structure:
```
.
├── cmd/server/main.go   # Entry point
├── internal/
│   ├── auth/           # Authorization logic
│   ├── cache/          # Key cache
│   ├── config/         # Configuration
│   └── crdwatcher/     # CRD watcher
```
