# rsync CLI

A simple CLI tool to synchronize files and directories between given source and destination folder.

## Usage

```sh
go run ./cmd/cli/main.go --src <source-folder> --dst <destination-folder> [--delete-missing]
```

### Mandatory Arguments

- `--src <source-folder>`  
  Path to the source directory to sync from.

- `--dst <destination-folder>`  
  Path to the destination directory to sync to.

### Optional Arguments

- `--delete-missing`  
  If provided (set to `true`), files and directories present in the destination but missing in the source will be deleted from the destination.

## Examples

Sync from `./data/source` to `./data/dest`:

```sh
go run ./cmd/cli/main.go --src ./data/source --dst ./data/dest
```

Sync and delete files in destination that are not present in source:

```sh
go run ./cmd/cli/main.go --src ./data/source --dst ./data/dest --delete-missing
```

## Notes

- Both `--src` and `--dst` are required.
- The tool preserves file modification times.
- Directory structure is recreated in the destination.
- If `--delete-missing` is not specified, extra files in the destination are left untouched.
