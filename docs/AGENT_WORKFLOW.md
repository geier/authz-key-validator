# Agent Workflow

## Architecture

```
Supervisor Agent
 ├── Triage Agent (creates issues)
 ├── Worker 1 → Issue #1
 └── Worker 2 → Issue #3
```

## Agent Roles

| Role | Responsibility |
|------|---------------|
| **Supervisor** | Assigns work, monitors progress, handles escalations |
| **Worker** | Implements assigned issues in isolated worktrees |
| **Triage** | Converts requests into well-scoped GitHub issues |

## Worker Agent Workflow

```bash
# 1. Receive issue number from supervisor
N=1

# 2. Create isolated worktree
git fetch origin main
git worktree add worktrees/issue-$N -b feature/issue-$N origin/main

# 3. Claim the issue
gh issue edit $N --add-assignee @me
gh issue edit $N --add-label in-progress

# 4. Work in isolation
cd worktrees/issue-$N
# ... implement ...

# 5. Commit specific files (never use . or -A)
git add internal/cache/cache.go internal/cache/cache_test.go
git commit -m "feat: implement key cache"
git push -u origin feature/issue-$N

# 6. Create PR
gh pr create --title "feat: implement key cache" --body "Closes #$N"

# 7. Request Copilot review
gh copilot-review <pr-number>
```

### Worker Rules
- Only work in `worktrees/issue-N/` directory
- Stage specific files only (never `git add .`)
- Let the "Closes #N" in PR body close the issue on merge

## Supervisor Agent Loop

```
┌──────────────────────────────────────────────────────────┐
│ 1. ASSIGN WORK                                         │
│    - Find ready issues:                                 │
│      gh issue list --search "no:assignee no:blocked-by label:ready" │
│    - Spawn one worker per issue                         │
├─────────────────────────────────────────────────────────┤
│ 2. MONITOR                                             │
│    - Check worker outputs                               │
│    - Watch for `needs-supervisor` label                │
│    - Spawn new workers for completed workers            │
├─────────────────────────────────────────────────────────┤
│ 3. REVIEW PRs                                          │
│    - Phase 1: Hygiene (base branch, files changed)    │
│    - Phase 2: Deep review (spawn review workers)       │
│    - Phase 3: Copilot review cycle                     │
├─────────────────────────────────────────────────────────┤
│ 4. CLEANUP                                             │
│    - Remove resolved dependencies after PRs merge        │
│    - Update labels (remove in-progress, etc.)          │
└──────────────────────────────────────────────────────────┘
```

## Labels

| Label | Meaning |
|-------|---------|
| `ready` | Well-scoped, ready to pick up |
| `in-progress` | Currently being worked on |
| `needs-supervisor` | Escalated, needs attention |
| `task` | Type: implementation task |
| `priority:high` | Important or time-sensitive |

## Escalation

When blocked, the worker:
1. Comments on the issue explaining the blocker
2. Adds `needs-supervisor` label
3. Removes self as assignee

## Triage Agent

Converts user requests into GitHub issues with:
- Clear title
- Context section (why this matters)
- Acceptance criteria
- Files affected (if known)
- Priority label

## Quick Reference

```bash
# Find available work
gh issue list --search "no:assignee label:ready no:blocked-by"

# Claim an issue
gh issue edit $N --add-assignee @me --add-label in-progress

# Check escalations
gh issue list --label needs-supervisor

# Start a PR review cycle
gh copilot-review <pr-number>
```
