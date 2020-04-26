# File Action

This action writes a file from an input value optionally decoding the content.

## Inputs

### `path`

**optional** 

The path to the file to write. Will be generated if omitted.

### `content` 

**required**

Content of the file.

### `encoding`

**optional** 

- `base64`
- `base58`
- `base32`

If present, the action will decode     

## Outputs

### `path`

The path of resulting file relative to the workspace.

## Example usage

```yaml
- id:   file
  uses: ./
  with:
    encoding: "base64"
    content: "aGVsbG8gd29ybGQ="
- run: |
     [[ -f ${{ steps.file.outputs.path }} ]]
     [[ "$(cat ${{ steps.file.outputs.path }})" == "hello world" ]]
```
