## tfdir

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

### Github Actions

```yml
- uses: actions/checkout@v3
  with:
    fetch-depth: 0
- name: install tfdir
  env:
    GITHUB_TOKEN: ${{ secrets.PAT }}
  run: |
    curl "https://$GITHUB_TOKEN@raw.githubusercontent.com/torana-us/tfdir/master/installer.sh" | bash
- uses: technote-space/get-diff-action@v6
- name: get target dir
  run: echo ${{ env.GIT_DIFF }} | tr ' ' '\n' | ./tfdir get
```

## Release

`v[0-9]+.[0-9]+.[0-9]+`形式のタグをpushするとrelease workflowが動きます
