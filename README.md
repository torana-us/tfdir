# tfdir

terraformを実行するべきディレクトリを取得します

moduleを変更したときにそのmoduleに依存しているterraformを検知できます

[テスト](https://github.com/torana-us/tfdir/blob/master/cmd/get_test.go)を見ると分かりやすい

## Usage

### CLI

```
git diff --name-only | tfdir get
```

### config

参照 [config.yaml](./config.yaml)

### GitHub Actions

```yml
jobs:
  get_target_dirs:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
      pull-requests: write
    steps:
      - uses: actions/checkout@v4
      - name: install tfdir
        run: |
          curl "https://raw.githubusercontent.com/torana-us/tfdir/master/installer.sh" | bash
      - name: get target dir
        id: target_dirs
        env:
          HEAD_REF: ${{ github.head_ref }}
        run: |
          git diff --diff-filter=AMRCD \
            --name-only \
            ${{ env.HEAD_REF }}..origin/main \
            "terraform/**.tf" \
            | xargs \
            | ./tfdir get
```

## Release

`v[0-9]+.[0-9]+.[0-9]+`形式のタグをpushするとrelease workflowが動きます

## License

MIT