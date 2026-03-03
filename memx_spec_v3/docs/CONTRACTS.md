---
owner: memx-core
status: active
last_reviewed_at: 2026-03-03
next_review_due: 2026-06-03
---

# CONTRACTS

API/CLI の機械可読契約（フィールド単位）の正本。

## API Contracts

### POST /v1/notes:ingest
#### Request
| field | type | required | constraints |
| --- | --- | --- | --- |
| store | string | yes | `short\|chronicle\|memopedia\|archive` |
| title | string | yes | non-empty |
| body | string | yes | non-empty |

#### Response 200
| field | type | required | constraints |
| --- | --- | --- | --- |
| id | string | yes | note identifier |
| store | string | yes | request.store と同値 |
| created_at | string | yes | RFC3339 |

### POST /v1/notes:search
#### Request
| field | type | required | constraints |
| --- | --- | --- | --- |
| store | string | yes | `short\|chronicle\|memopedia\|archive` |
| query | string | yes | non-empty |
| limit | integer | no | `1..100` |

#### Response 200
| field | type | required | constraints |
| --- | --- | --- | --- |
| items | array<object> | yes | 検索ヒット一覧 |
| items[].id | string | yes | note identifier |
| items[].title | string | yes | - |
| items[].snippet | string | yes | - |
| total | integer | yes | `>=0` |

### GET /v1/notes/{id}
#### Response 200
| field | type | required | constraints |
| --- | --- | --- | --- |
| id | string | yes | path id と同値 |
| store | string | yes | `short\|chronicle\|memopedia\|archive` |
| title | string | yes | - |
| body | string | yes | - |
| created_at | string | yes | RFC3339 |

## CLI Contracts（`--json`）
- `mem in short --json` は `POST /v1/notes:ingest` の Response 200 と同型。
- `mem out search --json` は `POST /v1/notes:search` の Response 200 と同型。
- `mem out show --json` は `GET /v1/notes/{id}` の Response 200 と同型。

## Error Contracts
| http_status | code | retryable | description |
| --- | --- | --- | --- |
| 400 | INVALID_ARGUMENT | false | 入力検証エラー |
| 403 | POLICY_DENIED | false | ポリシー拒否（fail-closed） |
| 404 | NOT_FOUND | false | ノート未検出 |
| 500 | INTERNAL | false | 内部障害 |
