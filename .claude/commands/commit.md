Analyze staged changes and create a properly formatted commit.

**Message format:** `<type>(<scope>): <description>`
- Types: `feat`, `fix`, `refactor`, `test`, `docs`, `ci`, `chore`
- Scope is the package name when the change is in one package: `soap`, `handlers`, `config`, `models`
- Description: imperative mood, lowercase, no period, max 72 chars total
- Add a blank line + body only when the *why* is non-obvious from the subject
- No "Co-Authored-By", no AI mentions, no attribution lines

**Steps:**
1. Run `git diff --cached --stat` then `git diff --cached` to read all staged changes
2. Run `git status` to check for relevant untracked files worth staging — never stage `PLAN.md` or any file listed in `.gitignore`
3. Check if `README.md` or `CLAUDE.md` need updating based on the staged changes:
   - New env vars, endpoints, commands, or config → update README.md
   - Changed interfaces, file structure, patterns, setup steps, or slash commands → update CLAUDE.md
   - If updates are needed, make them and stage them before committing
4. Draft the commit message, show it, and ask for confirmation before committing
5. Run `git commit -m "<subject>"` (or with `-m "<subject>" -m "<body>"` if a body is needed)

$ARGUMENTS
