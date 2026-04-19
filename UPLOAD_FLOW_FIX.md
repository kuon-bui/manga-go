# Upload Flow - Fixed for New Chapter Creation

## 🎯 Problem Solved

**Issue:** Khi tạo chapter mới, chưa có `chapterId` → không thể upload ảnh

**Solution:** Cho phép upload ảnh mà không cần `chapterId`
- Nếu có `chapterId`: `comics/{comicSlug}/chapters/{chapterSlug}/pages/{uuid}.ext`
- Nếu không có: `comics/{comicSlug}/temp-uploads/{uuid}.ext` (temp folder)

---

## 📌 Important: URL vs Path

Upload API response chứa 2 fields:
- **`url`**: `/files/content/comics/manga-slug/temp-uploads/uuid.jpg` 
  - Dùng để hiển thị ảnh trên FE (preview)
  - Đường dẫn đầy đủ để request ảnh từ server
  
- **`path`**: `comics/manga-slug/temp-uploads/uuid.jpg`
  - Chỉ là file path lưu trong storage
  - Dùng để gửi API tạo chapter
  - Server sẽ reconstruct thành `/files/content/{path}` khi cần

---

## ✅ Recommended Flow: Upload → Create Chapter

### **Step 1: Upload Images (WITHOUT chapterId)**

```http
POST /files/upload HTTP/1.1
Content-Type: multipart/form-data

file=<binary>
type=chapter
comicId=550e8400-e29b-41d4-a716-446655440000
(chapterId NOT needed - save to temp folder)
```

**Response:**
```json
{
  "data": {
    "url": "/files/content/comics/manga-slug/temp-uploads/uuid1.jpg",
    "path": "comics/manga-slug/temp-uploads/uuid1.jpg",
    "filename": "uuid1.jpg"
  }
}
```

**Frontend:**
- Dùng `url` để hiển thị ảnh preview: `<img src={data.url} />`
- Dùng `path` để gửi khi tạo chapter (collect all paths)

### **Step 2: Create Chapter with Images + Title**

```http
POST /comics/:comicId/chapters HTTP/1.1
Content-Type: application/json

{
  "number": "1",
  "title": "Chương 1: Bắt đầu",
  "pages": [
    {
      "imageUrl": "comics/manga-slug/temp-uploads/uuid1.jpg",
      "pageType": "image"
    },
    {
      "imageUrl": "comics/manga-slug/temp-uploads/uuid2.jpg",
      "pageType": "image"
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "id": "chapter-uuid",
    "number": "1",
    "title": "Chương 1: Bắt đầu",
    "slug": "chuong-1-bat-dau",  // ← Backend auto-generated!
    "uploadedById": "user-uuid",
    "uploadedBy": { "id": "...", "name": "..." },
    "pages": [
      {
        "imageUrl": "comics/manga-slug/temp-uploads/uuid1.jpg",
        "pageType": "image"
      },
      ...
    ]
  }
}
```

---

## 🚀 Frontend Implementation

```javascript
// Step 1: Upload images
async function uploadChapterImages(comicId, files) {
  const imagePaths = [];
  const imageUrls = [];  // For preview display
  
  for (const file of files) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', 'chapter');
    formData.append('comicId', comicId);
    // ← No chapterId needed!
    
    const response = await fetch('/api/v1/files/upload', {
      method: 'POST',
      body: formData
    });
    
    const data = await response.json();
    // Use 'url' for displaying image preview
    imageUrls.push(data.data.url);
    // Use 'path' for chapter creation API
    imagePaths.push({
      imageUrl: data.data.path,  // Path for DB storage
      pageType: 'image'
    });
  }
  
  return { imagePaths, imageUrls };  // Return both

// Step 2: Create chapter with images
async function createChapterWithImages(comicId, title, imagePages) {
  const response = await fetch(`/api/v1/comics/${comicId}/chapters`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      number: '1',
      title: title,
      // slug: NOT needed - auto-generated from title
      pages: imagePages  // From step 1: use 'path' field
    })
  });
  
  return await response.json();
}

// Usage
const comicId = '550e8400-e29b-41d4-a716-446655440000';
const files = [file1, file2, file3];
const { imagePaths, imageUrls } = await uploadChapterImages(comicId, files);

// Display previews
imageUrls.forEach(url => {
  console.log(`<img src="${url}" />`);  // Use 'url' for preview
});

// Create chapter with paths
const chapter = await createChapterWithImages(comicId, 'Chapter 1', imagePaths);
console.log(chapter.data.slug); // auto-generated: "chapter-1"
```

---

## 📁 Updated Folder Structure

```
S3/MinIO Bucket:
├── comics/
│   ├── manga-a-slug/
│   │   ├── cover/
│   │   │   └── uuid.jpg
│   │   ├── temp-uploads/          ← NEW: Temporary folder
│   │   │   ├── uuid1.jpg
│   │   │   ├── uuid2.jpg
│   │   │   └── uuid3.jpg
│   │   ├── chapters/
│   │   │   ├── ch-1-slug/
│   │   │   │   └── pages/
│   │   │   │       └── uuid.jpg
│   │   │   └── ch-2-slug/
│   │   │       └── pages/
│   │   │           └── uuid.jpg
│   │   └── ...
│   └── ...
```

---

## ⚙️ Backend Changes

### Endpoint `/files/upload`

**Parameter changes:**
- `chapterId`: Now **OPTIONAL** (not required)

**Logic:**
```
if type=chapter:
  if chapterId provided:
    path = "comics/{comicSlug}/chapters/{chapterSlug}/pages/{uuid}.ext"
  else:
    path = "comics/{comicSlug}/temp-uploads/{uuid}.ext"  ← NEW
else if type=cover:
  path = "comics/{comicSlug}/cover/{uuid}.ext"
```

### Create Chapter API

**Accepts:**
- `number` (required)
- `title` (required)
- `slug` (OPTIONAL - backend auto-generates if not provided)
- `pages` (required, array)

**Slug auto-generation:**
```
If slug not provided:
  slug = slugify(title)
  // "Chương 1: Bắt đầu" → "chuong-1-bat-dau"
```

---

## 📋 Summary

| Item | Before | After |
|------|--------|-------|
| Upload without chapterId | ❌ Error | ✅ Saved to temp-uploads |
| Slug input | ⚠️ Manual | ✅ Auto-generated from title |
| chapterId required | ✅ Always | ⚠️ Optional |
| Temp folder | ❌ No | ✅ temp-uploads/ |

---

## 🔄 Update Chapter Images (After Creation)

Once chapter is created, you can **add/remove/reorder** pages:

```http
PUT /comics/:comicId/chapters/:chapterId/pages HTTP/1.1

{
  "pages": [
    { "imageUrl": "comics/.../pages/uuid1.jpg", "pageType": "image" },
    { "imageUrl": "comics/.../pages/uuid2.jpg", "pageType": "image" },
    { "imageUrl": "comics/.../temp-uploads/uuid3.jpg", "pageType": "image" }
  ]
}
```

Images can be from `temp-uploads` or `pages` folder - doesn't matter!
