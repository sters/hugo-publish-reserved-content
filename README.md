Moved to https://github.com/sters/hugo-utilities

---

# hugo-publish-reserved-content

```
go get github.com/sters/hugo-publish-reserved-content/cmd/hugo-publish-reserved-content
```
or download from [release page](https://github.com/sters/hugo-publish-reserved-content/releases).

```
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
reserved: true
draft: true
---
Hi everone! This content is reserved at `2021-01-01T00:00:00Z` !
```

When use this tools after `2021-01-01T00:00:00Z`, remove `reserved` and `draft` keys.

`-basePath XXX` is your hugo content directory.

```
hugo-publish-reserved-content -basePath XXX -reservedKey reserved -draftKey draft
```

```
---
date: 2021-01-01T00:00:00Z
title: Happy New Year!
---
Hi everone! This content is reserved at `2021-01-01T00:00:00Z`!
```

## In Github actions

You can use like this action for automated publish:

```
name: Check and publish reserved articles

on:
  schedule:
    - cron: '0 * * * *'

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      # See https://docs.github.com/en/free-pro-team@latest/actions/reference/events-that-trigger-workflows#triggering-new-workflows-using-a-personal-access-token
      - uses: actions/checkout@v2
        with:
          token: ${{ secrets.GH_PAT }}

      - run: go get github.com/sters/hugo-publish-reserved-content/cmd/hugo-publish-reserved-content

      - name: Check reserved content
        run: |
          export PATH=${PATH}:`go env GOPATH`/bin
          hugo-publish-reserved-content -basePath blog/content/

      - name: Commit and Push
        run: |
          set +e
          git config user.name github-actions
          git config user.email github-actions@github.com
          git add .
          git commit -m "publish reserved article"
          git push
```
