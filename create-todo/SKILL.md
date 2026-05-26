---
name: create-todo
description: Create a committed Markdown todo in ~/.create-todo from a todo, task, follow-up, $create-todo request, or current conversation context.
---

# Create Todo

Create one durable Markdown todo from the user's request and current conversation context.

## Rules

- Do not read local project files unless the user explicitly references them.
- The todo must have a concrete action and clear acceptance criteria. If either is missing or vague, ask targeted clarification questions using the available question/user-input tool. If no such tool is available, ask directly. Ask as many times as needed before creating the todo.
- Synthesize a concise human title from the request and context. Do not use the user's raw wording if a clearer title is available.
- Write minimal Markdown: `# Title`, then a concise freeform body.
- If the bundled script returns a source repo URL, include it as the first body line: `Source: https://github.com/owner/repo`.
- Keep the body freeform, but make the action and acceptance criteria explicit.

## Workflow

1. Confirm the action and acceptance criteria are clear.
2. Generate the todo title and body.
3. Capture the user's current working directory, then run the bundled script from this skill directory with that source directory:

   ```bash
   source_dir="$PWD"
   skill_dir="<directory containing this SKILL.md>"
   "$skill_dir/create-todo.sh" --cwd "$source_dir" "<title>"
   ```

   The script prints two lines:
   - line 1: absolute path to the empty todo file
   - line 2: optional GitHub source repo URL; empty when unavailable

4. Write the Markdown content to the returned file path.
5. Format the todo:

   ```bash
   bunx prettier --write "$todo_file"
   ```

   If `bunx` is unavailable, use:

   ```bash
   npx prettier --write "$todo_file"
   ```

   If formatting fails, stop before committing.

6. Commit only the new todo file in `~/.create-todo`:

   ```bash
   todo_root="$HOME/.create-todo"
   todo_rel="${todo_file#"$todo_root"/}"
   git -C "$todo_root" add -- "$todo_rel"
   git -C "$todo_root" commit -m "todo: <short title>" -- "$todo_rel"
   ```

   Do not stage or commit any other files. If the commit fails, surface git's error.

7. If `~/.create-todo` has an `origin` remote, push `main`:

   ```bash
   git -C "$HOME/.create-todo" push origin main
   ```

   If no `origin` exists, report that the todo was committed locally and push was skipped. If push fails, report that the todo was created and committed, then surface the push error.

## Script Contract

`create-todo.sh` owns deterministic file creation and first-run todo-store setup:

- creates `~/.create-todo`
- initializes git on `main` when needed
- creates and commits minimal `README.md` and `.gitignore` as `initialize todo store` when the todo store has no commits
- detects the source project from `--cwd` git `origin`, git root, or folder name
- creates `~/.create-todo/<project_slug>/<YYYY-MM-DD>-<title_slug>-<HHMMSS>-<4hex>.md`
- returns the todo file path and optional GitHub source repo URL
