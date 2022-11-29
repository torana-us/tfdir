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
