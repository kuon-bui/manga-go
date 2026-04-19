# API Implementation Status - Manga Go

Ngày tạo: 2026-04-19  
Trạng thái: Cần hoàn thiện 6 API chính

---

## 📊 Tóm tắt

| API | Trạng thái | Ghi chú |
|---|---|---|
| **1. Top Trending** | ❌ **Cần tạo** | API hoàn toàn mới |
| **2. Recent Updates** | ❌ **Cần tạo** | API hoàn toàn mới |
| **3. Comic Detail** | ⚠️ **Cần bổ sung** | Thêm uploaderId, translationGroup |
| **4. Comment Reactions** | ⚠️ **Cần sửa** | Response format không khớp yêu cầu |
| **5. Page-Level Reactions** | ❌ **Cần tạo** | Endpoint mới cho chapter pages |
| **6. Comment Report** | ❌ **Cần tạo** | Model + endpoint mới |

---

## 1️⃣ API Top Trending - `GET /api/v1/titles/trending`

### Trạng thái: ❌ **KHÔNG TỒN TẠI**

#### Cần thực hiện:
- [ ] Tạo endpoint `GET /comics/trending` (alias hoặc direct)
- [ ] Implement service method: `ComicService.GetTrendingComics(limit)`
- [ ] Sắp xếp theo: views/followers trong 7 ngày gần nhất, hoặc isHot flag

#### Yêu cầu response:
```typescript
{
  "data": [
    {
      "id": "string",
      "title": "string",           // từ comic.title
      "coverImage": "string",       // từ comic.thumbnail
      "author": "string",           // authors[0].name
      "synopsis": "string",         // từ comic.description
      "genres": [
        { "id": "string", "name": "string" }
      ],
      "latestChapter": {
        "number": "string",         // chapter.number
        "id": "string"              // chapter.id
      },
      "views": number              // NEW FIELD - cần track views
    }
  ]
}
```

#### Công việc thêm:
1. **Thêm field `views`** vào Comic model nếu chưa có
2. **Tạo migration** để track views
3. **Service method**: `ComicService.GetTrendingComics(ctx, limit)`
4. **Repository method**: `ComicRepo.FindTrendingComics(ctx, limit)`
5. **Route handler**: `ComicHandler.getTrendingComics()`

---

## 2️⃣ API Recent Updates - `GET /api/v1/chapters/recent-updates`

### Trạng thái: ❌ **KHÔNG TỒN TẠI**

#### Cần thực hiện:
- [ ] Tạo endpoint `GET /chapters/recent-updates` hoặc `GET /comics/recent-updates`
- [ ] Implement service để lấy chapters mới nhất kèm comic info
- [ ] Sắp xếp theo `chapter.created_at` DESC

#### Yêu cầu response:
```typescript
{
  "data": [
    {
      "title": {
        "id": "string",
        "name": "string",              // comic.title
        "coverImage": "string"         // comic.thumbnail
      },
      "chapter": {
        "id": "string",
        "number": "string",            // chapter.number
        "name": "string",              // chapter.title
        "createdAt": "2026-04-19T00:00:00Z"
      }
    }
  ],
  "pagination": {
    "total": number,
    "page": number,
    "limit": number
  }
}
```

#### Công việc:
1. **Service method**: `ChapterService.GetRecentUpdates(ctx, limit, offset)`
2. **Repository method**: `ChapterRepo.FindRecentUpdates(ctx, limit, offset)`
3. **Route handler**: `ChapterHandler.getRecentUpdates()` hoặc tạo mới endpoint

---

## 3️⃣ Comic Detail - Bổ sung fields

### Trạng thái: ⚠️ **TỒN TẠI nhưng THIẾU FIELD**

Endpoint hiện có: `GET /api/v1/comics/:comicSlug`

#### Cần bổ sung:
```typescript
{
  // Các field hiện tại...
  "uploaderId": "string",        // NEW - từ comic.uploaded_by_id
  "translationGroup": {          // NEW - từ comic.translation_group
    "id": "string",
    "name": "string",
    "slug": "string"
  } | null
}
```

#### Thực hiện:
1. **Sửa Comic model** - fields này đã tồn tại (UploadedByID, TranslationGroupID)
2. **Sửa ComicService.GetComic()** - include translationGroup khi fetch
3. **Kiểm tra response mapping** - đảm bảo uploaderId được include trong JSON response

#### File cần sửa:
- `internal/pkg/services/comic/get.go` - thêm `TranslationGroup` vào preload
- `internal/pkg/model/comic.go` - đã có fields, chỉ cần json tags

---

## 4️⃣ Comment Reactions - Sửa response format

### Trạng thái: ⚠️ **API TỒN TẠI nhưng RESPONSE SAI**

#### API hiện có:
- `POST /api/v1/comments/:id/reactions` - Add/remove reaction ✅
- `GET /api/v1/comments` - Get comments ⚠️ (thiếu reaction info)

#### Vấn đề:
Comment response hiện tại **KHÔNG CÓ**:
- `reactionCounts` - count by reaction type
- `userReaction` - reaction của user hiện tại

#### Yêu cầu response cho GET /comments:
```typescript
{
  "data": [
    {
      "id": "string",
      "content": "string",
      "author": {
        "id": "string",
        "name": "string",
        "avatar": "string"
      },
      "createdAt": "2026-04-19T00:00:00Z",
      "reactionCounts": {
        "LIKE": number,
        "LOVE": number,
        "HAHA": number,
        "SAD": number,
        "ANGRY": number
      },
      "userReaction": "LIKE" | null,    // Reaction của user đang login
      "replyCount": number
    }
  ]
}
```

#### Công việc:
1. **Tạo response DTO** `CommentResponse` với các field trên
2. **Repository method**: `ReactionRepo.CountByCommentId(ctx, commentId)` - count by type
3. **Repository method**: `ReactionRepo.GetUserReaction(ctx, commentId, userId)` - lấy reaction của user
4. **Service method**: Build comment response với reactions
5. **Sửa route handler**: Sử dụng DTO thay vì model trực tiếp

#### File cần tạo/sửa:
- `internal/pkg/request/comment/response.go` - NEW - DTO cho response
- `internal/pkg/services/comment/list.go` - Sửa để build response với reactions
- `internal/pkg/repo/reaction/reaction.go` - Thêm methods

---

## 5️⃣ Page-Level Comment Reactions - Endpoint mới

### Trạng thái: ❌ **KHÔNG TỒN TẠI**

#### Yêu cầu:
- `POST /api/v1/chapters/:chapterId/pages/:pageIndex/react`
- Body: `{ "type": "LIKE" | "LOVE" | "HAHA" | "WOW" | "SAD" | "ANGRY" }`

#### Hiện tại:
- Comment model đã hỗ trợ `pageIndex` ✅
- Có thể dùng comment reactions endpoint như hiện tại

#### Có 2 lựa chọn:
**A) Dùng comment reactions endpoint trực tiếp** (khuyến khích)
```
POST /api/v1/comments/:id/reactions
```
- Với comment có `pageIndex` được set

**B) Tạo endpoint riêng cho page**
```
POST /api/v1/chapters/:chapterId/pages/:pageIndex/react
```
- Wrapper xung quanh comment reactions
- Phức tạp hơn nhưng UI-friendly hơn

#### Đề nghị:
Nếu frontend hiện tại gọi option A, không cần làm gì thêm.  
Nếu cần option B, tạo route handler:
- `internal/app/api/route/chapter/react_page.go`
- Tìm hoặc tạo comment cho page đó
- Call comment reaction service

---

## 6️⃣ Comment Report - Chức năng mới

### Trạng thái: ❌ **KHÔNG CÓ MODEL/ENDPOINT**

#### Yêu cầu:
- `POST /api/v1/comments/:id/report`
- Body: `{ "reason": "SPAM" | "OFFENSIVE", "details": "string" }`

#### Cần tạo:

**1. Model** - `internal/pkg/model/comment_report.go`
```go
type CommentReport struct {
    common.SqlModel
    CommentId uuid.UUID
    UserId    uuid.UUID
    Reason    string  // SPAM, OFFENSIVE, HARASSMENT, etc.
    Details   *string // Optional description
    
    Comment *Comment `gorm:"foreignKey:CommentId"`
    User    *User    `gorm:"foreignKey:UserId"`
}
```

**2. Migration** - `migrations/YYYYMMDD_HHMMSS_create_comment_reports.sql`

**3. Repository** - `internal/pkg/repo/comment_report/repo.go`
- Method: `Create(ctx, report)`
- Method: `ExistsByCommentAndUser(ctx, commentId, userId)` - prevent duplicate reports

**4. Request DTO** - `internal/pkg/request/comment/report.go`
```go
type ReportCommentRequest struct {
    Reason  string  `json:"reason" binding:"required,oneof=SPAM OFFENSIVE HARASSMENT ADULT_CONTENT"`
    Details *string `json:"details" binding:"max=500"`
}
```

**5. Service method** - `internal/pkg/services/comment/report.go`
```go
func (s *CommentService) ReportComment(ctx, userId, commentId, reason, details) response.Result
```

**6. Route handler** - Sửa `internal/app/api/route/comment/route.go`
```go
idRg.POST("/report", cr.commentHandler.reportComment)
```

**7. Handler method** - `internal/app/api/route/comment/report_comment.go`

---

## 📋 Checklist Triển khai

### Phase 1: Bổ sung fields & Fix response (Ngắn nhất)
- [ ] Sửa Comic detail response (uploaderId, translationGroup)
- [ ] Sửa Comment list response (reactionCounts, userReaction)
- [ ] Thêm `views` field đến Comic (nếu chưa có)

### Phase 2: Tạo API mới (Trung bình)
- [ ] Implement Trending API
- [ ] Implement Recent Updates API
- [ ] Implement Comment Report

### Phase 3: Tùy chọn (Tuỳ UI yêu cầu)
- [ ] Tạo page-level reaction endpoint nếu không dùng comment endpoint

---

## 🔧 Dependencies & Integration

### Cần kiểm tra:
1. **Views tracking** - Có cần track views khi user mở comic?
2. **Trending logic** - Dựa trên views, followers, hoặc isHot flag?
3. **Recent updates** - Chỉ published chapters?
4. **Comment reactions** - Có cần real-time count update?

---

## 📝 Notes

- **Reaction types**: Hiện tại code chỉ support generic string type, cần định nghĩa constants (LIKE, LOVE, HAHA, etc.)
- **Comment pagination**: Query params hiện có `limit`, `page` - cần kiểm tra binding
- **Authorization**: POST comment/report endpoints có check user auth không? (Có trong middleware)
- **Soft delete**: Reactions/Reports sử dụng soft-delete, kiểm tra logic khi query
