# Command Template Guide

This guide explains the command workflow structure for managing tasks through specification, execution, and archival.

## Overview

The command workflow consists of three phases:

1. **Append** - Add tasks to the checklist
2. **Apply** - Execute the task checklist
3. **Archive** - Move completed tasks to archive

## File Structure

### Command Directory
```
.claude/commands/{command-name}/
├── append.md            # Add tasks to checklist
├── apply.md             # Task execution command
├── archive.md           # Task archival command
└── TEMPLATE_GUIDE.md
```

### Working Directory (Created by append.md)
```
{command-name}-progress/
├── active-task-1.md
├── active-task-2.md
└── archived/            # Completed tasks
    ├── INDEX.md
    ├── 2025-01/
    │   └── completed-task-1.md
    └── 2025-02/
        └── completed-task-2.md
```

---

## 1. append.md - Add Tasks to Checklist

**Purpose**: Add new tasks to the checklist with detailed breakdown.

### Template Structure

```markdown
---
name: {command-name}-append
description: Add {task-type} tasks to the checklist with detailed breakdown
---

# {Command Name} - Add Tasks

## Usage

Add a new task with detailed checklist items for {task-type}.

## Command Behavior

When invoked, this command will:

1. **Create working directory** - Initialize `{command-name}-progress/` folder (if not exists)
2. **Analyze the request** - Understand the task requirements
3. **Break down into steps** - Create a logical sequence of subtasks
4. **Create checklist** - Generate a markdown checklist file
5. **Save to progress folder** - Store in `{command-name}-progress/{task-name}.md`

## Checklist Format

```markdown
# Task: {Task Name}

**Created**: {date}
**Status**: Not Started
**Priority**: {High/Medium/Low}

## Overview
{Brief description of the task}

## Checklist

- [ ] Step 1: {Description}
- [ ] Step 2: {Description}
- [ ] Step 3: {Description}

## Notes
{Any additional context or considerations}
\```

## Example Usage

\```bash
# User request
"Append: Implement user authentication feature"

# Command generates
go-code-reviewer-progress/
├── implement-user-authentication.md
└── archived/  (created if doesn't exist)
\```

## Implementation Guidelines

1. **Be Specific** - Each checklist item should be actionable
2. **Logical Order** - Steps should follow a natural progression
3. **Reasonable Scope** - Break large tasks into manageable chunks
4. **Include Context** - Add notes for important considerations
```

---

## 2. apply.md - Task Execution Command

**Purpose**: Execute tasks from the checklist and mark items as completed.

### Template Structure

```markdown
---
name: {command-name}-apply
description: Execute tasks from the checklist and check off completed items
---

# {Command Name} - Task Execution

## Usage

Work through the task checklist and mark items as completed.

## Command Behavior

When invoked, this command will:

1. **Read checklist** - Load the task file from `{command-name}-progress/{task-name}.md`
2. **Execute tasks** - Work through unchecked items sequentially
3. **Update progress** - Mark completed items with `[x]`
4. **Update status** - Change status field (Not Started → In Progress → Completed)
5. **Save changes** - Write updates back to the task file

## Workflow

### Initial State
\```markdown
- [ ] Step 1: Create database schema
- [ ] Step 2: Implement API endpoints
- [ ] Step 3: Write tests
\```

### During Execution
\```markdown
- [x] Step 1: Create database schema
- [x] Step 2: Implement API endpoints ← Currently working on
- [ ] Step 3: Write tests
\```

### Completion
\```markdown
- [x] Step 1: Create database schema
- [x] Step 2: Implement API endpoints
- [x] Step 3: Write tests
\```

## Status Updates

Update the status field as work progresses:

- **Not Started** - No items completed
- **In Progress** - Some items completed
- **Completed** - All items checked off
- **Blocked** - Cannot proceed (add reason in notes)

## Example Usage

\```bash
# User request
"Apply: Work on implement-user-authentication task"

# Command:
# 1. Opens go-code-reviewer-progress/implement-user-authentication.md
# 2. Shows current progress
# 3. Works on next unchecked item
# 4. Marks as [x] when done
# 5. Continues until all complete or user stops
\```

## Implementation Guidelines

1. **Show Progress** - Display current completion percentage
2. **One at a Time** - Focus on one checklist item at a time
3. **Verify Completion** - Ensure step is fully done before checking off
4. **Add Notes** - Document any issues or decisions made
5. **Update Timestamp** - Add "Last Updated" field when saving
```

---

## 3. archive.md - Task Archival Command

**Purpose**: Move completed task checklists to the archive for record-keeping.

### Template Structure

```markdown
---
name: {command-name}-archive
description: Move completed task checklists to the archive folder
---

# {Command Name} - Task Archival

## Usage

Archive completed tasks to keep the active tasks directory clean.

## Command Behavior

When invoked, this command will:

1. **Verify completion** - Check that all items are marked `[x]`
2. **Add completion date** - Update task file with completion timestamp
3. **Move to archive** - Relocate from `{command-name}-progress/` to `{command-name}-progress/archived/`
4. **Update index** - Add entry to `{command-name}-progress/archived/INDEX.md`
5. **Confirm action** - Show success message with archive location

## Directory Structure

\```
tasks/
├── active-task-1.md
└── active-task-2.md

archive/
├── INDEX.md
├── 2025-01/
│   └── completed-task-1.md
└── 2025-02/
    └── completed-task-2.md
\```

## Archive Format

Tasks are organized by completion month:
- `{command-name}-progress/archived/YYYY-MM/{task-name}.md`

The INDEX.md maintains a searchable list:

\```markdown
# Archived Tasks

## 2025-02

- [Implement User Authentication](2025-02/implement-user-authentication.md) - Completed: 2025-02-15
- [Add Dark Mode](2025-02/add-dark-mode.md) - Completed: 2025-02-20

## 2025-01

- [Initial Setup](2025-01/initial-setup.md) - Completed: 2025-01-30
\```

## Example Usage

\```bash
# User request
"Archive: implement-user-authentication task"

# Command:
# 1. Checks {command-name}-progress/implement-user-authentication.md is 100% complete
# 2. Adds "Completed: 2025-02-15" to file
# 3. Moves to {command-name}-progress/archived/2025-02/implement-user-authentication.md
# 4. Updates {command-name}-progress/archived/INDEX.md
# 5. Removes from {command-name}-progress/ directory
\```

## Implementation Guidelines

1. **Validate Completion** - Only archive 100% completed tasks
2. **Preserve History** - Don't modify task content when archiving
3. **Maintain Index** - Always update INDEX.md for searchability
4. **Date-based Organization** - Use YYYY-MM folder structure
5. **Confirm Action** - Show clear confirmation of archival
```

---

## Complete Workflow Example

### Phase 1: Append Task

\```bash
User: "I need to implement a code review workflow"

Command: append.md
Output: go-code-reviewer-progress/code-review-workflow.md

# Task: Code Review Workflow

**Created**: 2025-02-15
**Status**: Not Started
**Priority**: High

## Checklist
- [ ] Design review process
- [ ] Implement automated checks
- [ ] Create review templates
- [ ] Write documentation
- [ ] Test the workflow
\```

### Phase 2: Execution

\```bash
User: "Apply the code review workflow task"

Command: apply.md
Actions:
1. Opens go-code-reviewer-progress/code-review-workflow.md
2. Works on "Design review process" ✓
3. Updates file with [x]
4. Status: In Progress
5. Continues with next items...

Updated file shows:
- [x] Design review process
- [x] Implement automated checks
- [x] Create review templates
- [x] Write documentation
- [x] Test the workflow

**Status**: Completed
\```

### Phase 3: Archival

\```bash
User: "Archive the code review workflow task"

Command: archive.md
Actions:
1. Verifies all items checked ✓
2. Adds completion timestamp
3. Moves to go-code-reviewer-progress/archived/2025-02/code-review-workflow.md
4. Updates go-code-reviewer-progress/archived/INDEX.md
5. Confirms: "Task archived successfully"
\```

---

## Best Practices

### For append.md
- ✅ Break tasks into 3-10 actionable items
- ✅ Use clear, specific language
- ✅ Include context and dependencies
- ❌ Don't create overly granular checklists
- ❌ Don't mix unrelated tasks

### For apply.md
- ✅ Work sequentially through checklist
- ✅ Verify each step before checking off
- ✅ Update status regularly
- ❌ Don't skip ahead without completing prior steps
- ❌ Don't check off incomplete items

### For archive.md
- ✅ Only archive 100% completed tasks
- ✅ Maintain organized archive structure
- ✅ Keep INDEX.md updated
- ❌ Don't archive incomplete tasks
- ❌ Don't lose task history

---

## File Naming Conventions

### Task Files
- Use kebab-case: `implement-feature.md`
- Be descriptive: `add-user-authentication.md` not `auth.md`
- Avoid dates in names (use creation date field)

### Archive Organization
- Year-Month folders: `2025-01/`, `2025-02/`
- Preserve original task name in archive
- Single INDEX.md at archive root

---

## Integration with Other Tools

This command workflow integrates with:

- **Git** - Tasks can be tracked in version control
- **CI/CD** - Checklist completion can trigger pipelines
- **Project Management** - Export checklists to external tools
- **Documentation** - Archive serves as project history

---

## Customization

Adapt this template for specific use cases:

- **Code Reviews**: Add "Reviewer" and "PR Link" fields
- **Bug Fixes**: Include "Bug Report", "Root Cause", "Fix Applied"
- **Features**: Add "Design Doc", "Testing Plan", "Deployment Steps"
- **Releases**: Include "Version", "Release Notes", "Rollback Plan"

---

## Summary

| Command | Purpose | Input | Output |
|---------|---------|-------|--------|
| append.md | Add task to checklist | Task description | {command-name}-progress/{name}.md |
| apply.md | Execute & track progress | Task name | Updated task file |
| archive.md | Move completed task | Task name | {command-name}-progress/archived/YYYY-MM/{name}.md |

This workflow provides a systematic approach to task management through clear phases of appending tasks, execution, and archival.
