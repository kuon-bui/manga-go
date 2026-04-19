# Upload Image API - Redesigned Flow

## Overview
Endpoint `/files/upload` telah didesain ulang untuk:
- ✅ Nhận `comicId` và `chapterId` (UUIDs) thay vì slug
- ✅ Backend tự resolve ID → Slug
- ✅ Phân biệt loại upload: chapter images vs comic cover
- ✅ Unique filename + organized folder structure

---

## Endpoint

### `POST /files/upload`

**Form Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file` | file | ✅ | Image file (max 10MB, image/* only) |
| `type` | string | ✅ | Upload type: `"chapter"` hoặc `"cover"` |
| `comicId` | string | ✅ | Comic ID (UUID format) |
| `chapterId` | string | ⚠️ | Chapter ID (UUID) - required nếu `type=chapter` |

---

## Use Cases

### 1️⃣ Upload Chapter Images

**When:** Sau khi tạo chapter (hoặc cùng lúc)

**Request:**
```http
POST /files/upload HTTP/1.1
Content-Type: multipart/form-data

file=<binary>
type=chapter
comicId=550e8400-e29b-41d4-a716-446655440000
chapterId=660e8400-e29b-41d4-a716-446655440001
```

**Response:**
```json
{
  "message": "Upload image successfully",
  "data": {
    "url": "/files/content/comics/manga-slug/chapters/ch-1-slug/pages/123e4567-e89b-12d3-a456-426614174000.jpg",
    "filename": "123e4567-e89b-12d3-a456-426614174000.jpg",
    "path": "comics/manga-slug/chapters/ch-1-slug/pages/123e4567-e89b-12d3-a456-426614174000.jpg",
    "content_type": "image/jpeg",
    "size": 245632
  }
}
```

**Store `path` trong database (để dùng lúc cập nhật chapter)**

---

### 2️⃣ Upload Comic Cover

**When:** Tạo/update comic cover

**Request:**
```http
POST /files/upload HTTP/1.1
Content-Type: multipart/form-data

file=<binary>
type=cover
comicId=550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "data": {
    "url": "/files/content/comics/manga-slug/cover/789f1234-e89b-12d3-a456-426614174111.jpg",
    "path": "comics/manga-slug/cover/789f1234-e89b-12d3-a456-426614174111.jpg"
  }
}
```

**Update comic:** `PATCH /comics/:comicId { thumbnail: path }`

---

## Folder Structure

```
S3/MinIO Bucket:
├── comics/
│   ├── manga-a-slug/
│   │   ├── cover/
│   │   │   ├── 123e4567-e89b-12d3.jpg
│   │   │   └── 456f8901-a23b-45d6.jpg
│   │   ├── chapters/
│   │   │   ├── ch-1-slug/
│   │   │   │   └── pages/
│   │   │   │       ├── 789f1234-b89c-12d3.jpg
│   │   │   │       ├── 890a2345-c89d-23e4.jpg
│   │   │   │       └── ...
│   │   │   └── ch-2-slug/
│   │   │       └── pages/
│   │   │           └── ...
│   │   └── ...
│   └── manga-b-slug/
│       └── ...
├── translation-groups/
│   └── ...
└── ...
```

---

## Frontend Integration Example

### Scenario: Create Chapter with Images

```javascript
// Step 1: Create chapter (images sẽ được add sau)
const createResponse = await fetch('/api/v1/comics/:comicId/chapters', {
  method: 'POST',
  body: JSON.stringify({
    number: '1',
    title: 'Chapter 1',
    slug: 'ch-1',
    pages: [] // Hoặc có thể create rỗng rồi update sau
  })
});
const chapter = await createResponse.json();
const chapterId = chapter.data.id;

// Step 2: Upload images
const imagePaths = [];
for (const imageFile of files) {
  const formData = new FormData();
  formData.append('file', imageFile);
  formData.append('type', 'chapter');
  formData.append('comicId', comicId);
  formData.append('chapterId', chapterId);
  
  const uploadResponse = await fetch('/api/v1/files/upload', {
    method: 'POST',
    body: formData
  });
  const uploadData = await uploadResponse.json();
  imagePaths.push({
    imageUrl: uploadData.data.path,
    pageType: 'image'
  });
}

// Step 3: Update chapter with pages
await fetch(`/api/v1/comics/:comicId/chapters/:chapterId/pages`, {
  method: 'PUT',
  body: JSON.stringify({
    pages: imagePaths
  })
});
```

---

## Manage Chapter Images

### Add/Remove/Reorder Pages

**Endpoint:** `PUT /comics/:comicId/chapters/:chapterId/pages`

**Request:**
```json
{
  "pages": [
    {
      "pageType": "image",
      "imageUrl": "comics/manga-a/chapters/ch-1/pages/uuid1.jpg"
    },
    {
      "pageType": "image",
      "imageUrl": "comics/manga-a/chapters/ch-1/pages/uuid2.jpg"
    },
    {
      "pageType": "image",
      "imageUrl": "comics/manga-a/chapters/ch-1/pages/uuid3.jpg"
    }
  ]
}
```

- **Add:** Append new page object với imageUrl
- **Remove:** Exclude từ list
- **Reorder:** Sắp xếp lại thứ tự trong array

---

## Error Cases

```
❌ Invalid comicId format
   → Response: "invalid comicId format"

❌ Comic not found
   → Response: "comic not found"

❌ Chapter not found
   → Response: "chapter not found or doesn't belong to this comic"

❌ Missing required parameter
   → Response: "'chapterId' is required when type=chapter"

❌ Invalid type
   → Response: "'type' must be 'chapter' or 'cover'"

❌ File too large
   → Response: "File size exceeds 10MB"

❌ Not an image
   → Response: "Only image files are allowed"
```

---

## Notes

- ✅ **Unique filenames:** UUID4 ensures no overwrites
- ✅ **Organized structure:** Dễ track, dễ cleanup
- ✅ **ID-based:** Safe hơn slug (slug can change)
- ✅ **Flexible:** Support both chapter images và cover
- ✅ **Validation:** Backend verify comic-chapter relationship
- ⚠️ **Old slugs:** Nếu có ảnh dùng slug path cũ, cần migration script

---

## Database Considerations

**Image URL Format:**
- Stored in `chapters.pages[].image_url` (string)
- Format: `comics/{comicSlug}/chapters/{chapterSlug}/pages/{uuid}.{ext}`
- Có thể reconstruct URL khi cần: `/files/content/{path}`

**Cleanup:**
- Khi delete chapter → Xóa ảnh trong S3
- Khi delete comic → Xóa toàn bộ folder
