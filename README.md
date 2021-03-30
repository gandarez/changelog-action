# Changelog Action

This actions prints a formatted changelog.

## Example usage

### Basic

```yaml
- id: changelog
  uses: gandarez/changelog-action@v{latest}
- name: "Print changelog"
  run: echo "${{ steps.changelog.outputs.changelog }}"
```

### Advanced

```yaml
- id: changelog
  uses: gandarez/changelog-action@v{latest}
  with:
    current_tag: "v0.2.8"
    previous_tag: "v0.2.2"
    main_branch_name: "trunk"
    exclude: |
        "^Merge pull request .*"
        "Fix .*"
    debug: "true"
- name: "Print changelog"
  run: echo "${{ steps.changelog.outputs.changelog }}"
```

## Inputs

| parameter           | required | description                                                                      | default     |
| ---                 | ---      | ---                                                                              | ---         |
| current_tag         |          | The current tag to be used instead of auto detecting.                            |             |
| previous_tag        |          | The previous tag to be used instead of auto detecting.                           |             |
| exclude             |          | Commit messages matching the regexp listed here will be removed from the output. |             |
| debug               |          | Enables debug mode.                                                              | false       |

## Outpus

| parameter           | description              |
| ---                 | ---                      |
| changelog           | The formatted changelog. |
