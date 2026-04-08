# Agent Workflow

This project uses a Supervisor → Worker model for task distribution.

## Architecture

```
Supervisor Agent
     │
     ├── Triage Agent (creates issues)
     ├── Worker 1 → Issue #1
     ├── Worker 2 → Issue #3
     └── Worker N → Issue #N
```

## Agent Roles

### Supervisor Agent
- Assigns issues to workers
- Monitors progress
- Handles escalations
- Runs PR review loop

### Triage Agent
- Captures feature requests and bug reports
- Creates well-scoped issues
- Sets dependencies

### Worker Agent
- Picks up assigned issues
- Works in git worktrees
- Creates PRs and responds to reviews

## Label System

| Label | Meaning |
|-------|---------|
| `ready` | Issue ready to be picked up |
| `in-progress` | Worker has claimed this issue |
| `needs-review` | PR is ready for review |
| `needs-supervisor` | Blocked, needs escalation |

## GitHub Labels

```bash
gh label create "ready" --color 0E8A16 --description "Ready for pickup"
gh label create "in-progress" --color FBCA04 --description "Currently being worked on"
gh label create "needs-review" --color 1D76DB --description "PR ready for review"
gh label create "needs-supervisor" --color B60205 --description "Blocked, needs escalation"
```

## Starting a Task

```bash
# 1. Create worktree
git worktree add ../authz-key-validator-1 -b feature/issue-1

# 2. Claim issue
gh issue edit 1 --add-assignee @me
gh issue edit 1 --add-label "in-progress"

# 3. Work in isolation
cd ../authz-key-validator-1
git checkout -b feature/issue-1-key-cache

# 4. After commit + push
gh pr create --title "feat: Implement key cache (#1)" --body "Closes #1"
```

## PR Review Cycle
1. `gh copilot-review <pr>`
2. Address comments
3. Once approved, squash and merge
