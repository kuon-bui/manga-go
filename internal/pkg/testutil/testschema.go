package testutil

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const assignUUIDCallbackName = "testutil:assign-uuid-primary-key"

var uuidType = reflect.TypeFor[uuid.UUID]()

func MustSyncSchemas(t testing.TB, db *gorm.DB, schemas ...any) {
	t.Helper()

	if err := db.Callback().Create().Before("gorm:create").Register(assignUUIDCallbackName, assignUUIDPrimaryKey); err != nil {
		if !strings.Contains(err.Error(), "duplicated") {
			t.Fatalf("failed to register test uuid callback: %v", err)
		}
	}

	if err := db.AutoMigrate(schemas...); err != nil {
		t.Fatalf("failed to sync test schema: %v", err)
	}
}

func NewSQLiteDB(t testing.TB) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(NewSQLiteDSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	return db
}

func MustReadUUID(t testing.TB, db *gorm.DB, query string, args ...any) uuid.UUID {
	t.Helper()

	var raw string
	if err := db.Raw(query, args...).Scan(&raw).Error; err != nil {
		t.Fatalf("failed to query uuid: %v", err)
	}

	id, err := uuid.Parse(raw)
	if err != nil {
		t.Fatalf("failed to parse uuid %q: %v", raw, err)
	}

	return id
}

func NewSQLiteDSN() string {
	return ":memory:"
}

func assignUUIDPrimaryKey(tx *gorm.DB) {
	assignUUIDValue(tx.Statement.ReflectValue)
}

func assignUUIDValue(value reflect.Value) {
	for value.IsValid() && value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	if !value.IsValid() {
		return
	}

	switch value.Kind() {
	case reflect.Struct:
		assignUUIDToStruct(value)
	case reflect.Slice, reflect.Array:
		for index := 0; index < value.Len(); index++ {
			assignUUIDValue(value.Index(index))
		}
	}
}

func assignUUIDToStruct(value reflect.Value) {
	if field := value.FieldByName("ID"); field.IsValid() && field.CanSet() && field.Type() == uuidType {
		if current, ok := field.Interface().(uuid.UUID); ok && current == uuid.Nil {
			field.Set(reflect.ValueOf(uuid.New()))
		}
	}

	for index := 0; index < value.NumField(); index++ {
		fieldType := value.Type().Field(index)
		if !fieldType.Anonymous {
			continue
		}

		field := value.Field(index)
		if !field.IsValid() {
			continue
		}

		for field.Kind() == reflect.Ptr {
			if field.IsNil() {
				break
			}
			field = field.Elem()
		}

		if field.IsValid() && field.Kind() == reflect.Struct {
			assignUUIDToStruct(field)
		}
	}
}

type SQLModel struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;primaryKey"`
	CreatedAt *time.Time     `gorm:"column:created_at"`
	UpdatedAt *time.Time     `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

type Author struct {
	SQLModel
	Name string `gorm:"column:name"`
}

func (Author) TableName() string {
	return "authors"
}

type Genre struct {
	SQLModel
	Name        string  `gorm:"column:name"`
	Slug        string  `gorm:"column:slug"`
	Description *string `gorm:"column:description"`
	Thumbnail   *string `gorm:"column:thumbnail"`
}

func (Genre) TableName() string {
	return "genres"
}

type Tag struct {
	SQLModel
	Name string `gorm:"column:name"`
	Slug string `gorm:"column:slug"`
}

func (Tag) TableName() string {
	return "tags"
}

type Comic struct {
	SQLModel
	Title              string     `gorm:"column:title"`
	Slug               string     `gorm:"column:slug"`
	AlternativeTitles  []byte     `gorm:"column:alternative_titles"`
	Description        *string    `gorm:"column:description"`
	Thumbnail          *string    `gorm:"column:thumbnail"`
	Banner             *string    `gorm:"column:banner"`
	Type               string     `gorm:"column:type"`
	Status             string     `gorm:"column:status"`
	AgeRating          string     `gorm:"column:age_rating"`
	IsPublished        bool       `gorm:"column:is_published"`
	IsHot              bool       `gorm:"column:is_hot"`
	IsFeatured         bool       `gorm:"column:is_featured"`
	PublishedYear      *int       `gorm:"column:published_year"`
	LastChapterAt      *time.Time `gorm:"column:last_chapter_at"`
	TranslationGroupID *uuid.UUID `gorm:"column:translation_group_id;type:uuid"`
	UploadedByID       *uuid.UUID `gorm:"column:uploaded_by_id;type:uuid"`
}

func (Comic) TableName() string {
	return "comics"
}

type Chapter struct {
	SQLModel
	ComicID      uuid.UUID  `gorm:"column:comic_id;type:uuid"`
	Number       string     `gorm:"column:number"`
	ChapterIdx   uint       `gorm:"column:chapter_idx"`
	Title        string     `gorm:"column:title"`
	Slug         string     `gorm:"column:slug"`
	IsPublished  bool       `gorm:"column:is_published"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	UploadedByID *uuid.UUID `gorm:"column:uploaded_by_id;type:uuid"`
}

func (Chapter) TableName() string {
	return "chapters"
}

type Page struct {
	SQLModel
	ChapterID  uuid.UUID `gorm:"column:chapter_id;type:uuid"`
	PageNumber int       `gorm:"column:page_number"`
	PageType   string    `gorm:"column:page_type"`
	ImageURL   string    `gorm:"column:image_url"`
	Content    string    `gorm:"column:content"`
}

func (Page) TableName() string {
	return "pages"
}

type ReadingProgress struct {
	SQLModel
	UserID        uuid.UUID `gorm:"column:user_id;type:uuid"`
	ComicID       uuid.UUID `gorm:"column:comic_id;type:uuid"`
	ChapterID     uuid.UUID `gorm:"column:chapter_id;type:uuid"`
	ScrollPercent int       `gorm:"column:scroll_percent"`
}

func (ReadingProgress) TableName() string {
	return "reading_progresses"
}

type ReadingHistory struct {
	SQLModel
	UserID     uuid.UUID  `gorm:"column:user_id;type:uuid"`
	ChapterID  uuid.UUID  `gorm:"column:chapter_id;type:uuid"`
	ComicID    uuid.UUID  `gorm:"column:comic_id;type:uuid"`
	LastReadAt *time.Time `gorm:"column:last_read_at"`
}

func (ReadingHistory) TableName() string {
	return "reading_histories"
}

type Rating struct {
	SQLModel
	UserID  uuid.UUID `gorm:"column:user_id;type:uuid"`
	ComicID uuid.UUID `gorm:"column:comic_id;type:uuid"`
	Score   int       `gorm:"column:score"`
	Comment *string   `gorm:"column:comment"`
}

func (Rating) TableName() string {
	return "ratings"
}

type Permission struct {
	SQLModel
	Name string `gorm:"column:name"`
}

func (Permission) TableName() string {
	return "permissions"
}

type Role struct {
	SQLModel
	Name string `gorm:"column:name"`
}

func (Role) TableName() string {
	return "roles"
}

type User struct {
	SQLModel
	Name                  string     `gorm:"column:name"`
	Avatar                *string    `gorm:"column:avatar"`
	Email                 string     `gorm:"column:email"`
	Password              string     `gorm:"column:password"`
	ResetPasswordToken    string     `gorm:"column:reset_password_token"`
	ResetPasswordExpiryAt *time.Time `gorm:"column:reset_password_expiry_at"`
	TranslationGroupID    *uuid.UUID `gorm:"column:translation_group_id;type:uuid"`
	UserConfig            []byte     `gorm:"column:user_config"`
}

func (User) TableName() string {
	return "users"
}

type UserRole struct {
	UserID uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey"`
	RoleID uuid.UUID `gorm:"column:role_id;type:uuid;primaryKey"`
}

func (UserRole) TableName() string {
	return "users_roles"
}

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"column:role_id;type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"column:permission_id;type:uuid;primaryKey"`
}

func (RolePermission) TableName() string {
	return "roles_permissions"
}

type TranslationGroup struct {
	SQLModel
	Name    string    `gorm:"column:name"`
	Slug    string    `gorm:"column:slug"`
	OwnerID uuid.UUID `gorm:"column:owner_id;type:uuid"`
	LogoURL *string   `gorm:"column:logo_url"`
}

func (TranslationGroup) TableName() string {
	return "translation_groups"
}

type ComicFollow struct {
	SQLModel
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid"`
	ComicID      uuid.UUID `gorm:"column:comic_id;type:uuid"`
	FollowStatus string    `gorm:"column:follow_status"`
}

func (ComicFollow) TableName() string {
	return "comic_follows"
}

type UserComicRead struct {
	SQLModel
	UserID   uuid.UUID `gorm:"column:user_id;type:uuid"`
	ComicID  uuid.UUID `gorm:"column:comic_id;type:uuid"`
	ReadData []byte    `gorm:"column:read_data"`
}

func (UserComicRead) TableName() string {
	return "user_comic_reads"
}

type Comment struct {
	SQLModel
	UserID    uuid.UUID  `gorm:"column:user_id;type:uuid"`
	ChapterID *uuid.UUID `gorm:"column:chapter_id;type:uuid"`
	ComicID   uuid.UUID  `gorm:"column:comic_id;type:uuid"`
	ParentID  *uuid.UUID `gorm:"column:parent_id;type:uuid"`
	PageIndex *int       `gorm:"column:page_index"`
	Content   string     `gorm:"column:content"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentReaction struct {
	SQLModel
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid"`
	CommentID uuid.UUID `gorm:"column:comment_id;type:uuid"`
	Type      string    `gorm:"column:type"`
}

func (CommentReaction) TableName() string {
	return "comment_reactions"
}

type PageReaction struct {
	SQLModel
	UserID uuid.UUID `gorm:"column:user_id;type:uuid"`
	PageID uuid.UUID `gorm:"column:page_id;type:uuid"`
	Type   string    `gorm:"column:type"`
}

func (PageReaction) TableName() string {
	return "page_reactions"
}
