package fileprocess

type ImageProcessPayload struct {
	FilePath           string `json:"filePath"`
	TemporaryObjectKey string `json:"temporaryObjectKey"`
}

type ImageProcessCleanupPayload struct {
	TemporaryObjectKey string `json:"temporaryObjectKey"`
}
