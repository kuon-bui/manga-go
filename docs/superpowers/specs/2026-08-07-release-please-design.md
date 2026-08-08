# Release automation với release-please (BE + FE)

Ngày: 2026-08-07
Phạm vi: `kuon-bui/manga-go` (BE, Go) và `kuon-bui/manga-go-fe` (FE, Next.js)

## Mục tiêu

Tự động hoá versioning, changelog và GitHub Release cho cả hai repo bằng
[googleapis/release-please](https://github.com/googleapis/release-please). Riêng BE, mỗi
release phải build và push Docker image cho hai service `api` và `queue`.

## Bối cảnh

| | BE `manga-go` | FE `manga-go-fe` |
|---|---|---|
| Stack | Go 1.25.4 | Next.js 15 / React 19, yarn |
| Conventional commits | Có (65/118 commit) | Có (46/63 commit) |
| Tag hiện tại | Không có | Không có |
| Breaking change trong history | Không | Không |
| Workflow hiện có | `test.yml`, `docker-publish.yml` | Không có |
| CHANGELOG.md | Không có | Có, nhưng là changelog của template "Vibe-Coding Prompt Template" |

## Vấn đề cốt lõi

release-please tạo tag bằng `GITHUB_TOKEN`. GitHub cố tình không cho event sinh ra từ token
mặc định trigger workflow khác, nên `docker-publish.yml` hiện tại (trigger `push: tags: v*`)
sẽ **không tự chạy** khi release-please cắt release.

Ba cách giải quyết đã cân nhắc:

1. **Gọi trong cùng workflow (đã chọn).** Chuyển `docker-publish.yml` thành reusable workflow
   (`workflow_call`); `release.yml` gọi nó khi `release_created == 'true'`. Không cần secret
   mới, luồng deterministic, đọc log ở một chỗ.
2. **PAT cho release-please.** Đơn giản về code nhưng thêm một secret phải xoay vòng; PAT hết
   hạn thì release vẫn chạy mà image không được build — lỗi im lặng.
3. **`workflow_run`.** Không cần PAT nhưng phải tự query API xem có release nào vừa tạo không;
   phức tạp và khó debug.

## Quyết định

| Quyết định | Chọn |
|---|---|
| Nối release → docker | Reusable workflow (`workflow_call`) |
| Version khởi điểm | `0.1.0` cho cả hai repo |
| Phạm vi changelog đầu tiên | Quét toàn bộ history (không set `bootstrap-sha`) |
| CHANGELOG.md cũ của FE | Archive sang `docs/CHANGELOG.legacy.md` |
| CI cho FE | Thêm `ci.yml` (lint + test + build) |
| Stamp version vào Go binary | Không — ngoài scope |

Release đầu tiên của cả hai repo sẽ là **`v0.2.0`**: baseline `0.1.0` cộng với các commit
`feat` có trong history, không có breaking change. Cả hai repo đều dưới ngưỡng
`commit-search-depth` mặc định (500) nên quét full history an toàn.

## Thiết kế

### Cấu hình chung cho cả hai repo

Dùng manifest-driven config thay vì khai `release-type` inline trong workflow — cấu hình được
versioned và review như code.

Hai file mỗi repo:

- `release-please-config.json`
- `.release-please-manifest.json` — nội dung `{".": "0.1.0"}`

```json
{
  "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
  "release-type": "go",
  "include-component-in-tag": false,
  "bump-minor-pre-major": true,
  "changelog-sections": [
    { "type": "feat", "section": "Features" },
    { "type": "fix", "section": "Bug Fixes" },
    { "type": "perf", "section": "Performance" },
    { "type": "refactor", "section": "Refactoring" },
    { "type": "docs", "section": "Documentation" },
    { "type": "test", "hidden": true },
    { "type": "chore", "hidden": true },
    { "type": "ci", "hidden": true },
    { "type": "build", "hidden": true }
  ],
  "packages": { ".": { "package-name": "manga-go" } }
}
```

Ba tuỳ chọn cần giải thích:

- **`include-component-in-tag: false`** — bắt buộc. Trong manifest mode giá trị này default
  `true`, sinh tag `manga-go-v0.2.0`. Docker image tag và mọi thứ downstream đều giả định
  tag dạng `v0.2.0`.
- **`bump-minor-pre-major: true`** — giữ semantics 0.x: breaking change bump lên `0.3.0` chứ
  không nhảy `1.0.0` ngoài ý muốn.
- **`changelog-sections`** — default của release-please chỉ hiện `feat`/`fix`/`perf`. Cả hai
  repo có nhiều commit `refactor(...)` và `docs(...)` đáng vào changelog.

FE khác BE hai chỗ: `"release-type": "node"` (tự bump `version` trong `package.json`) và
`"package-name": "manga-go-fe"`.

### BE — `.github/workflows/release.yml` (mới)

```yaml
name: Release
on:
  push: { branches: [main] }
permissions:
  contents: write
  pull-requests: write
  issues: write
concurrency:
  group: release-${{ github.ref }}
  cancel-in-progress: false
jobs:
  release-please:
    runs-on: ubuntu-latest
    outputs:
      release_created: ${{ steps.release.outputs.release_created }}
      tag_name: ${{ steps.release.outputs.tag_name }}
      version: ${{ steps.release.outputs.version }}
    steps:
      - uses: googleapis/release-please-action@v4
        id: release

  docker:
    needs: release-please
    if: needs.release-please.outputs.release_created == 'true'
    uses: ./.github/workflows/docker-publish.yml
    with:
      ref: ${{ needs.release-please.outputs.tag_name }}
      version: ${{ needs.release-please.outputs.version }}
    secrets: inherit
```

`cancel-in-progress: false` vì release có side effect (tạo tag, push image) — không được cắt
giữa đường. Action không cần `with:` vì tự đọc `release-please-config.json`.

### BE — `.github/workflows/docker-publish.yml` (sửa)

Giữ nguyên matrix, target, cache và bước build. Sửa ba chỗ:

| Chỗ | Trước | Sau |
|---|---|---|
| Trigger | `push: tags: v*` + `workflow_dispatch` | `workflow_call` + `workflow_dispatch`, cả hai nhận input `ref` và `version` |
| Checkout | HEAD | `ref: ${{ inputs.ref }}` — build đúng commit của tag |
| Image tags | `type=semver,pattern=...` | `type=raw,value=<prefix><version>` |

`type=semver` cần git ref là tag mới parse được. Khi gọi qua `workflow_call`, `github.ref` là
`refs/heads/main` nên semver không sinh ra tag nào. Truyền version tường minh vừa đúng vừa
đọc được.

Bỏ `push: tags` để chỉ còn một đường vào — muốn build tay thì dùng `workflow_dispatch` và
nhập `ref` / `version`.

**Bug được fix cùng lúc:** dòng
`type=raw,value=${{ matrix.latest_tag }},enable=${{ github.ref == 'refs/heads/main' }}`.
Workflow chỉ chạy trên tag push, khi đó `github.ref` luôn là `refs/tags/v*`, nên điều kiện
`enable=` chưa bao giờ đúng và `latest` / `queue-latest` chưa bao giờ được push. Vì workflow
giờ chỉ chạy khi có release thật, bỏ `enable=` đi.

Image tag sau mỗi release:

```
<user>/manga-go:0.2.0       <user>/manga-go:queue-0.2.0
<user>/manga-go:latest      <user>/manga-go:queue-latest
<user>/manga-go:sha-<sha>   <user>/manga-go:queue-sha-<sha>
```

### FE — `.github/workflows/release.yml` (mới)

Giống BE, bỏ job `docker`. `release-type: node` bump `package.json` từ `0.1.0` lên `0.2.0`.
`package.json` có `"private": true` nên không có bước publish npm.

### FE — `.github/workflows/ci.yml` (mới)

Release PR của release-please cần có gác cổng, tương đương `test.yml` của BE.

```yaml
name: CI
on:
  push: { branches: [main] }
  pull_request: { branches: [main] }
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: 22, cache: yarn }
      - run: yarn install --frozen-lockfile
      - run: yarn lint
      - run: yarn test
      - run: cp .env.example .env
      - run: yarn build
```

Không có step `tsc --noEmit` riêng vì `next build` đã typecheck. `cp .env.example .env` để
`next build` có sẵn các biến `NEXT_PUBLIC_*`.

### FE — archive changelog cũ

`git mv CHANGELOG.md docs/CHANGELOG.legacy.md`, thêm một đoạn note ở đầu file nói rõ đây là
changelog của template gốc và changelog hiện hành do release-please quản lý ở `CHANGELOG.md`.
release-please sẽ tạo `CHANGELOG.md` mới sạch.

## Việc phải làm thủ công trên GitHub

**Bắt buộc, cả hai repo:** Settings → Actions → General → Workflow permissions → bật
*"Allow GitHub Actions to create and approve pull requests"*. Thiếu setting này thì
release-please chạy xong nhưng không mở được Release PR và fail khá im lặng.

Còn lại:

- Secrets `DOCKERHUB_USERNAME` / `DOCKERHUB_TOKEN` trong environment `docker-publish` của BE
  đã có sẵn, không cần thêm.
- Nếu environment `docker-publish` đang bật *required reviewers*, job `docker` sẽ treo chờ
  approve sau mỗi release. Cần kiểm tra và quyết định giữ hay bỏ.
- FE không cần secret nào.

## Luồng vận hành sau khi xong

1. Merge commit conventional vào `main`.
2. release-please mở hoặc cập nhật Release PR — `chore(main): release 0.2.0`.
3. Merge Release PR đó.
4. release-please tạo tag `v0.2.0`, GitHub Release, và cập nhật `CHANGELOG.md`.
5. Chỉ BE: job `docker` build và push 6 image tag.

## Ngoài scope

- Không đụng vào `test.yml` của BE.
- Không stamp version vào Go binary; Docker image tag đã đủ để truy vết.
- Không Docker image cho FE.
- Không publish npm (FE là `private`).
