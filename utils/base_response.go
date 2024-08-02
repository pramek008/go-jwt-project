package utils

import "github.com/gin-gonic/gin"

// BaseResponse adalah struktur generik yang digunakan untuk merespons permintaan API
type BaseResponse[T any] struct {
	StatusCode int    `json:"status_code"`     // Status HTTP dari respons
	IsSuccess  bool   `json:"is_success"`      // Menunjukkan apakah permintaan berhasil atau tidak
	Message    string `json:"message"`         // Pesan yang menyertai respons
	Limit      int64  `json:"limit,omitempty"` // Batas hasil untuk respons paginasi
	Page       int64  `json:"page,omitempty"`  // Nomor halaman untuk respons paginasi
	Total      int64  `json:"total,omitempty"` // Total item untuk respons paginasi
	Data       T      `json:"data,omitempty"`  // Data aktual yang dikirim dalam respons
}

// ResponseOption adalah tipe fungsi yang memodifikasi BaseResponse
type ResponseOption[T any] func(*BaseResponse[T])

// WithPagination adalah opsi untuk menambahkan informasi paginasi ke BaseResponse
func WithPagination[T any](limit, page, total int64) ResponseOption[T] {
	return func(r *BaseResponse[T]) {
		r.Limit = limit
		r.Page = page
		r.Total = total
	}
}

// SendResponse mengirim respons API ke klien
func SendResponse[T any](c *gin.Context, statusCode int, isSuccess bool, message string, data T, opts ...ResponseOption[T]) {
	// Membuat respons dasar dengan data yang diberikan
	response := BaseResponse[T]{
		StatusCode: statusCode,
		IsSuccess:  isSuccess,
		Message:    message,
		Data:       data,
	}

	// Menerapkan opsi tambahan ke respons jika ada
	for _, opt := range opts {
		opt(&response)
	}

	// Mengirim respons dalam format JSON
	c.JSON(statusCode, response)
}

// SendErrorResponse mengirim respons kesalahan API ke klien
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	// Menggunakan SendResponse dengan tipe data interface{} untuk mengirim respons kesalahan
	SendResponse[interface{}](c, statusCode, false, message, nil)
}

// SendPaginatedResponse mengirim respons paginasi API ke klien
func SendPaginatedResponse[T any](c *gin.Context, statusCode int, isSuccess bool, message string, data T, limit, page, total int64) {
	// Menggunakan SendResponse dengan opsi paginasi untuk mengirim respons paginasi
	SendResponse(c, statusCode, isSuccess, message, data, WithPagination[T](limit, page, total))
}
