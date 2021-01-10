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
