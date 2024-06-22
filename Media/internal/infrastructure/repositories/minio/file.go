package minioRepo

import (
	"Media/internal/domain"
	"Media/internal/usecases"
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log/slog"
	"strings"
	"unicode"
)

const (
	NameKey     = "Name"
	AuthorIDKey = "Authorid"
)

var _ usecases.FileRepositoryInterface = &FileRepository{}

type FileRepository struct {
	bucketName string
	client     *minio.Client
	logger     *slog.Logger
}

func (f *FileRepository) CreateFile(ctx context.Context, file domain.File) error {
	_, err := f.client.PutObject(ctx, f.bucketName, file.ID.String(), strings.NewReader(string(file.Content)), int64(len(file.Content)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
		UserMetadata: map[string]string{
			NameKey:     cleanName(file.Name),
			AuthorIDKey: file.AuthorID.String(),
		},
	})
	if err != nil {
		f.logger.Error("FileRepository.CreateFile", slog.Any("error", err.Error()))
		return err
	}
	return nil
}

func (f *FileRepository) GetFile(ctx context.Context, id uuid.UUID) (domain.File, error) {
	const op = "FileRepository.GetFile"
	logger := f.logger.With("op", op)

	object, err := f.client.GetObject(ctx, f.bucketName, id.String(), minio.GetObjectOptions{})
	if err != nil {
		logger.Error("error while getting object", slog.Any("error", err.Error()))
		return domain.File{}, err
	}
	defer func(object *minio.Object) {
		err := object.Close()
		if err != nil {
			logger.Error("error while closing object", slog.Any("error", err.Error()))
		}
	}(object)

	info, err := object.Stat()
	if err != nil {
		logger.Error("error while getting object info", slog.Any("error", err.Error()))
		return domain.File{}, err
	}

	content := make([]byte, info.Size)
	_, err = object.Read(content)
	if err != nil && err.Error() != "EOF" {
		logger.Error("error while reading object", slog.Any("error", err.Error()))
		return domain.File{}, err
	}

	rawAuthorID := info.UserMetadata[AuthorIDKey]
	if rawAuthorID == "" {
		logger.Error("AuthorID not found")
		return domain.File{}, err
	}

	authorID, err := uuid.Parse(rawAuthorID)
	if err != nil {
		logger.Error("error while parsing AuthorID", slog.Any("error", err.Error()))
		return domain.File{}, err
	}

	return domain.File{
		ID:       id,
		AuthorID: authorID,
		Name:     info.UserMetadata[NameKey],
		Size:     info.Size,
		Content:  content,
	}, nil
}

func (f *FileRepository) DeleteFile(ctx context.Context, id uuid.UUID) error {
	err := f.client.RemoveObject(ctx, f.bucketName, id.String(), minio.RemoveObjectOptions{})
	if err != nil {
		f.logger.Error("FileRepository.DeleteFile", slog.Any("error", err.Error()))
		return err
	}
	return nil
}

func NewFileRepository(client *minio.Client, bucketName string, logger *slog.Logger) *FileRepository {
	return &FileRepository{
		bucketName: bucketName,
		client:     client,
		logger:     logger,
	}
}

func cleanName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, name)
	return result
}
