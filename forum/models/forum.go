package models

// Информация о форуме.
type Forum struct {
	// Название форума.
	Title string `json:"title"`
	// Nickname пользователя, который отвечает за форум.
	User string `json:"user"`
	// Человекопонятный URL (https://ru.wikipedia.org/wiki/%D0%A1%D0%B5%D0%BC%D0%B0%D0%BD%D1%82%D0%B8%D1%87%D0%B5%D1%81%D0%BA%D0%B8%D0%B9_URL), уникальное поле.
	Slug string `json:"slug"`
	// Общее кол-во сообщений в данном форуме.
	Posts float32 `json:"posts,omitempty"`
	// Общее кол-во ветвей обсуждения в данном форуме.
	Threads float32 `json:"threads,omitempty"`
}
