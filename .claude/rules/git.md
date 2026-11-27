# Git and GitHub

## When to commit
- Never run `git commit`, `git push`, `gh pr create` or `gh pr merge` unless the user explicitly asks for it in the current message
- Finishing a task is not a request to commit. Stop, summarize the changes, and wait
- Before committing, show the drafted message and wait for confirmation (the `/commit` command does this)
- Work on a feature branch (`feature/<name>`, `fix/<name>`), never commit directly to `main`

## No AI attribution
Applies to commit messages, PR titles and bodies, issues, and review comments:
- No `Co-Authored-By: Claude` or any other co-author trailer
- No "Generated with Claude Code", no robot emoji, no mention of AI or assistants

## Commit message format
- `<type>(<scope>): <description>`, max 72 chars, imperative mood, lowercase after the colon, no period
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `ci`, `chore`
- Scope is the package name when the change is confined to one package: `soap`, `handlers`, `config`, `models`
- Add a body only when the reason for the change is not obvious from the subject
- Reference the issue in the body or PR (`Closes #3`) when the work comes from the project board

## Writing style (commits, PRs, issues, docs)
- Plain, direct sentences. Write like a human, not a marketing page
- No emojis, no filler words ("seamlessly", "robust", "powerful", "leverages")
