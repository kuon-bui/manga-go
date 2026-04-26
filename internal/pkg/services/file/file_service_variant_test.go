package fileservice

import (
	"bytes"
	"context"
	"testing"
)

func TestGetFileByVariant_ReturnVariantFile(t *testing.T) {
	fake := newFakeStorage()
	service := &FileService{objectStorage: fake}

	canonical := "comics/demo/cover/abc.webp"
	variantKey := BuildVariantObjectKey(canonical, ImageVariantSmall)
	fake.files[variantKey] = []byte("small-data")

	content, resolvedKey, err := service.GetFileByVariant(context.Background(), canonical, "small")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(content) != "small-data" {
		t.Fatalf("unexpected content: %q", string(content))
	}

	if resolvedKey != variantKey {
		t.Fatalf("want resolved key %q, got %q", variantKey, resolvedKey)
	}

	if len(fake.getCalls) != 1 || fake.getCalls[0] != variantKey {
		t.Fatalf("expected single lookup to %q, got %#v", variantKey, fake.getCalls)
	}
}

func TestGetFileByVariant_FallbackToCanonicalWhenVariantMissing(t *testing.T) {
	fake := newFakeStorage()
	service := &FileService{objectStorage: fake}

	canonical := "comics/demo/chapters/ch-1/pages/abc.webp"
	variantKey := BuildVariantObjectKey(canonical, ImageVariantMedium)
	fake.files[canonical] = []byte("normal-data")

	content, resolvedKey, err := service.GetFileByVariant(context.Background(), canonical, "medium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(content) != "normal-data" {
		t.Fatalf("unexpected content: %q", string(content))
	}

	if resolvedKey != canonical {
		t.Fatalf("want fallback key %q, got %q", canonical, resolvedKey)
	}

	if len(fake.getCalls) != 2 || fake.getCalls[0] != variantKey || fake.getCalls[1] != canonical {
		t.Fatalf("unexpected lookup sequence: %#v", fake.getCalls)
	}
}

func TestGetFileByVariant_InvalidInputs(t *testing.T) {
	service := &FileService{objectStorage: newFakeStorage()}

	if _, _, err := service.GetFileByVariant(context.Background(), "../secret.webp", "normal"); err == nil {
		t.Fatalf("expected invalid filename error")
	}

	if _, _, err := service.GetFileByVariant(context.Background(), "comics/demo/cover/abc.webp", "bad"); err == nil {
		t.Fatalf("expected invalid variant error")
	}
}

func TestUploadImageVariants_Success(t *testing.T) {
	fake := newFakeStorage()
	service := &FileService{objectStorage: fake}

	canonical := "comics/demo/cover/abc.webp"
	pngData := makePNG(1200, 600)

	result, err := service.UploadImageVariants(context.Background(), canonical, bytes.NewReader(pngData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Path != canonical {
		t.Fatalf("want canonical path %q, got %q", canonical, result.Path)
	}

	if result.URL != "/files/content/"+canonical {
		t.Fatalf("unexpected canonical url: %q", result.URL)
	}

	if result.ContentType != webpContentType {
		t.Fatalf("unexpected content type: %q", result.ContentType)
	}

	if len(result.Variants) != len(imageVariantOrder) {
		t.Fatalf("want %d variants, got %d", len(imageVariantOrder), len(result.Variants))
	}

	if len(fake.uploaded) != len(imageVariantOrder) {
		t.Fatalf("want %d uploaded objects, got %d", len(imageVariantOrder), len(fake.uploaded))
	}

	widthByVariant := map[string]int{}
	for _, item := range result.Variants {
		widthByVariant[item.Variant] = item.Width
		if item.ContentType != webpContentType {
			t.Fatalf("variant %s has unexpected content type %q", item.Variant, item.ContentType)
		}
		if item.URL != "/files/content/"+canonical {
			t.Fatalf("variant %s has unexpected url %q", item.Variant, item.URL)
		}
	}

	if widthByVariant[string(ImageVariantSmall)] != 480 {
		t.Fatalf("small width should be 480, got %d", widthByVariant[string(ImageVariantSmall)])
	}
	if widthByVariant[string(ImageVariantMedium)] != 720 {
		t.Fatalf("medium width should be 720, got %d", widthByVariant[string(ImageVariantMedium)])
	}
	if widthByVariant[string(ImageVariantLarge)] != 1080 {
		t.Fatalf("large width should be 1080, got %d", widthByVariant[string(ImageVariantLarge)])
	}
	if widthByVariant[string(ImageVariantNormal)] != 1200 {
		t.Fatalf("normal width should be 1200, got %d", widthByVariant[string(ImageVariantNormal)])
	}

	for _, variant := range imageVariantOrder {
		key := BuildVariantObjectKey(canonical, variant)
		if _, ok := fake.uploaded[key]; !ok {
			t.Fatalf("missing uploaded key %q", key)
		}
		if fake.uploadedContentType[key] != webpContentType {
			t.Fatalf("uploaded key %q has unexpected content type %q", key, fake.uploadedContentType[key])
		}
	}
}

func TestUploadImageVariants_InvalidInput(t *testing.T) {
	service := &FileService{objectStorage: newFakeStorage()}

	if _, err := service.UploadImageVariants(context.Background(), "comics/demo/cover/abc.webp", nil); err == nil {
		t.Fatalf("expected error for nil body")
	}

	if _, err := service.UploadImageVariants(context.Background(), "../invalid.webp", bytes.NewReader(makePNG(10, 10))); err == nil {
		t.Fatalf("expected error for invalid object key")
	}

	if _, err := service.UploadImageVariants(context.Background(), "comics/demo/cover/abc.webp", bytes.NewReader([]byte("not-an-image"))); err == nil {
		t.Fatalf("expected decode error for non-image payload")
	}
}
