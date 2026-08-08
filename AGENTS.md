# AGENTS.md for asciichgolangpublic

## Constitution

The [constitution.md](constitution.md) has to be applied to the whole repo and its codebase.

## Editing files

- Always use `sudo` to make you owner of the files or directories to adjust.
- Use the `file_editor` to modify files. There is no need to use shell commands for this.
- After file edit is done check all `README.md` besides the files modified. Keep all of them up to date.


## Specifications

- Specifications are directly written into the package directories. They are always named `<packagename>.spec.md`.
- All specifications have to be linked in the `README.md` of the same package as:
    ```
    ## Specifications

    For specifications see [<packagename>.spec.md](<packagename>.spec.md)
    ```
- Specifications are hierarchically organized. Every <packagename>.spec.md counts for the subdirectories/ subpackages as well.

## Consistency check

When asking for checking the repo consistency:
- Correct spelling errors in comments and all MarkDown files.
- Validate all specifications are met.
