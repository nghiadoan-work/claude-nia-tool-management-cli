---
name: go-code-reviewer-archive
description: Move completed code review task checklists to the archive folder
---

# Go Code Reviewer - Task Archival

## Usage

Archive completed code review tasks to keep the active tasks directory clean and maintain a historical record of reviews.

## Command Behavior

When invoked, this command will:

1. **Verify completion** - Check that all checklist items are marked `[x]`
2. **Validate review summary** - Ensure review findings are documented
3. **Add completion date** - Update task file with completion timestamp
4. **Move to archive** - Relocate from `go-code-reviewer-progress/` to `go-code-reviewer-progress/archived/YYYY-MM/`
5. **Update index** - Add entry to `go-code-reviewer-progress/archived/INDEX.md`
6. **Confirm action** - Show success message with archive location

## Archive Directory Structure

```
go-code-reviewer-progress/
├── active-review-1.md
├── active-review-2.md
└── archived/
    ├── INDEX.md
    ├── 2025-11/
    │   ├── review-auth-package.md
    │   └── review-api-endpoints.md
    └── 2025-12/
        └── review-database-layer.md
```

## Archive Format

Tasks are organized by completion month:
- `go-code-reviewer-progress/archived/YYYY-MM/{task-name}.md`

The INDEX.md maintains a searchable list of all completed reviews:

```markdown
# Archived Code Reviews

## 2025-12

### review-database-layer.md
- **Completed**: 2025-12-15
- **Scope**: pkg/database package
- **Issues Found**: 8 (2 critical, 3 warnings, 3 suggestions)
- **Status**: All critical issues addressed
- **Summary**: Database layer review identified SQL injection vulnerabilities and missing connection pooling

## 2025-11

### review-api-endpoints.md
- **Completed**: 2025-11-20
- **Scope**: internal/api package
- **Issues Found**: 12 (0 critical, 5 warnings, 7 suggestions)
- **Status**: Clean - no critical issues
- **Summary**: API endpoints review found missing input validation and documentation gaps

### review-auth-package.md
- **Completed**: 2025-11-16
- **Scope**: internal/auth package
- **Issues Found**: 15 (2 critical, 5 warnings, 8 suggestions)
- **Status**: Critical issues fixed before merge
- **Summary**: Authentication review identified security vulnerabilities in token handling
```

## Example Usage

```bash
# User request
"Archive: review-auth-package task"

# Command:
# 1. Checks go-code-reviewer-progress/review-auth-package.md is 100% complete
# 2. Verifies review summary is documented
# 3. Adds "Completed: 2025-11-16" to file
# 4. Moves to go-code-reviewer-progress/archived/2025-11/review-auth-package.md
# 5. Updates go-code-reviewer-progress/archived/INDEX.md with review metadata
# 6. Removes from go-code-reviewer-progress/ directory
# 7. Confirms: "✓ Code review archived: archived/2025-11/review-auth-package.md"
```

## Completion Requirements

Before a review can be archived, verify:

### 1. All Checklist Items Completed

```markdown
✓ All items have [x]
✗ Cannot archive with [ ] items remaining
```

### 2. Review Summary Documented

```markdown
✓ Review Summary section is filled out:
  - Issues Found: {count}
  - Critical: {list}
  - Warnings: {list}
  - Suggestions: {list}

✗ Cannot archive without documented findings
```

### 3. Status is "Completed"

```markdown
✓ **Status**: Completed

✗ **Status**: In Progress  ← Cannot archive
✗ **Status**: Blocked      ← Cannot archive
```

## Archive Entry Format

Each archived review gets an entry in INDEX.md:

```markdown
### {task-name}.md
- **Completed**: YYYY-MM-DD
- **Scope**: {files/package/project}
- **Issues Found**: {total} ({critical} critical, {warnings} warnings, {suggestions} suggestions)
- **Status**: {resolution status}
- **Summary**: {one-line summary of review}
```

Example entries:

```markdown
### review-security-critical.md
- **Completed**: 2025-11-16
- **Scope**: Entire project
- **Issues Found**: 23 (5 critical, 8 warnings, 10 suggestions)
- **Status**: All critical issues fixed and verified
- **Summary**: Security audit before v2.0 release - identified auth vulnerabilities, missing input validation, and exposed secrets

### review-performance-optimization.md
- **Completed**: 2025-11-10
- **Scope**: pkg/cache, pkg/query packages
- **Issues Found**: 7 (0 critical, 2 warnings, 5 suggestions)
- **Status**: Clean - optimizations noted for future work
- **Summary**: Performance review identified memory allocation improvements and caching opportunities
```

## Implementation Guidelines

1. **Validate Completion**: Only archive 100% completed reviews with documented findings
2. **Preserve History**: Don't modify review content when archiving - keep original findings
3. **Maintain Index**: Always update INDEX.md with comprehensive metadata
4. **Date-based Organization**: Use YYYY-MM folder structure for chronological organization
5. **Confirm Action**: Show clear confirmation with archive path
6. **Searchable Metadata**: Include scope, issues count, and summary for easy searching

## Index Management

The INDEX.md file serves as:

1. **Quick Reference**: See all past reviews at a glance
2. **Metrics Tracking**: Track review frequency and issue trends
3. **Knowledge Base**: Learn from past reviews
4. **Audit Trail**: Maintain record of code quality over time

### Index Maintenance Rules

- **Newest First**: Most recent reviews at the top of each month section
- **Consistent Format**: Use exact same format for every entry
- **Complete Metadata**: All fields required (Completed, Scope, Issues Found, Status, Summary)
- **One Line Summary**: Keep summaries concise but informative
- **Issue Counts**: Include breakdown of critical/warnings/suggestions

## Handling Edge Cases

### Incomplete Reviews

If review is incomplete but needs to be archived (e.g., code deleted, project cancelled):

1. Add note in task file explaining why incomplete
2. Update status to "Incomplete - {reason}"
3. Archive with special notation in INDEX.md:

```markdown
### review-deprecated-module.md
- **Completed**: 2025-11-16 (Incomplete)
- **Scope**: pkg/oldmodule (deprecated)
- **Issues Found**: N/A
- **Status**: Review cancelled - module removed from codebase
- **Summary**: Module deprecated before review completion
```

### Blocked Reviews

If review is permanently blocked:

1. Document the blocking issue
2. Change status to "Blocked - {reason}"
3. Archive with explanation:

```markdown
### review-proprietary-lib.md
- **Completed**: 2025-11-16 (Blocked)
- **Scope**: vendor/proprietary package
- **Issues Found**: Unable to review
- **Status**: Blocked - no access to proprietary source code
- **Summary**: Third-party proprietary library - source not available for review
```

## Archive Statistics

Consider adding statistics section to INDEX.md:

```markdown
# Archive Statistics

## 2025-11
- **Total Reviews**: 12
- **Critical Issues Found**: 8
- **Warnings**: 24
- **Suggestions**: 47
- **Average Issues per Review**: 6.6
- **Most Common Issues**:
  - Missing error handling (8 occurrences)
  - Documentation gaps (7 occurrences)
  - Security concerns (5 occurrences)

## All Time
- **Total Reviews**: 47
- **Total Issues Found**: 312
- **Critical Issues**: 23
- **Average Resolution Time**: 2.3 days
```

## Archival Best Practices

1. **Archive Promptly**: Archive completed reviews within 24 hours
2. **Verify First**: Double-check completion before archiving
3. **Update Index Immediately**: Don't let INDEX.md get out of sync
4. **Preserve Context**: Include enough summary detail to understand the review later
5. **Track Outcomes**: Note in Status whether issues were fixed
6. **Learn from History**: Review archived tasks periodically to identify patterns

## Version Control Integration

If using git for the review tasks:

```bash
# Before archiving
git add go-code-reviewer-progress/archived/
git commit -m "Archive: review-auth-package - 2 critical issues fixed"

# This creates a permanent record in git history
```

## Search and Retrieval

Make archived reviews searchable:

```bash
# Find all reviews of a specific package
grep "pkg/auth" go-code-reviewer-progress/archived/INDEX.md

# Find all reviews with critical issues
grep "critical" go-code-reviewer-progress/archived/INDEX.md

# Find reviews from specific month
ls go-code-reviewer-progress/archived/2025-11/

# Search review content
grep -r "SQL injection" go-code-reviewer-progress/archived/
```

## Archive Retention

Consider establishing retention policies:

- **Recent Reviews** (< 6 months): Keep all details
- **Older Reviews** (6-24 months): Keep in archive
- **Ancient Reviews** (> 24 months): Consider compressing or summarizing
- **Always Keep**: INDEX.md with metadata for all reviews

## Confirmation Message Format

When archiving, display:

```
✓ Code Review Archived Successfully

Task: review-auth-package
Location: go-code-reviewer-progress/archived/2025-11/review-auth-package.md
Completed: 2025-11-16
Issues Found: 15 total (2 critical, 5 warnings, 8 suggestions)

The review has been added to the archive index.
To view: cat go-code-reviewer-progress/archived/INDEX.md
```

This provides clear confirmation and immediate access to the archived review.
