package categories

import (
	"go-shop/pkg/webtool"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	errors  *Errors
}

func NewHandler(service *Service, errors *Errors) *Handler {
	return &Handler{service: service, errors: errors}
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	group := r.Group("/categories")

	group.POST("", webtool.MakeHandler(h.Create))
	group.GET("", webtool.MakeHandler(h.Get))
	group.PUT(":id", webtool.MakeHandler(h.Update))
	group.DELETE(":id", webtool.MakeHandler(h.Delete))
	group.PATCH(":id/order", webtool.MakeHandler(h.UpdateOder))
	group.PATCH(":id/status", webtool.MakeHandler(h.UpdateStatus))

	group.GET("tree", webtool.MakeHandler(h.GetTree))
	group.GET("flat", webtool.MakeHandler(h.GetFlat))
	group.GET(":slug", webtool.MakeHandler(h.GetBySlug))
	group.GET(":slug/products", webtool.MakeHandler(h.GetProductsBySlug))
}

// Группа 1: Управление категориями (Админка / Бэкенд)
// Эти эндпоинты нужны для заполнения и управления деревом категорий.
// POST /api/categories — Создание новой категории. Логика: Проверка уникальности slug. Автоматический расчет level и path на основе parent_id. Если parent_id = 0 (или null) — это корневая категория (level=1, path=/id/).
// GET /api/categories/{id} — Получение данных конкретной категории (для формы редактирования).
// PUT /api/categories/{id} — Обновление категории. Логика: Особо важная часть — обработка изменения parent_id. При перемещении категории в другую ветку нужно рекурсивно пересчитать level и path для всех её дочерних категорий. Это самая сложная операция.
// DELETE /api/categories/{id} — Удаление категории. Логика: Каскадное удаление (как настроено в БД ON DELETE CASCADE) или перенос товаров в другую категорию (более безопасный вариант).
// PATCH /api/categories/{id}/order — Изменение порядка отображения (display_order) в рамках одного уровня иерархии (среди "братьев и сестер").
// PATCH /api/categories/{id}/status — Активация/деактивация категории (is_active). Деактивация обычно скрывает категорию и все её дочерние категории с фронта.

func (h *Handler) Create(c *gin.Context) *webtool.APIError {
	var input CreateRequestDTO
	err := c.BindJSON(&input)
	if err != nil {
		return h.errors.BadInput
	}
	category, err := h.service.CreateCategory(CreateRequestDTOToCategory(&input))
	if err != nil {
		return h.errors.CantCreate
	}

	c.JSON(http.StatusOK, CategoryToCreateResponseDTO(category))
	return nil
}
func (h *Handler) Get(c *gin.Context) *webtool.APIError {
	return nil
}
func (h *Handler) Update(c *gin.Context) *webtool.APIError {
	return nil
}
func (h *Handler) Delete(c *gin.Context) *webtool.APIError {
	return nil
}
func (h *Handler) UpdateOder(c *gin.Context) *webtool.APIError {
	return nil
}
func (h *Handler) UpdateStatus(c *gin.Context) *webtool.APIError {
	return nil
}

//Группа 2: Отображение и навигация (Публичный API для фронта)
//GET /api/categories/tree — Самый важный эндпоинт. Получение всего дерева категорий в виде вложенной структуры (Nested Set). Ответ: Массив корневых категорий, каждая из которых содержит массив своих дочерних категорий (рекурсивно). Только активные (is_active=true). Для чего на фронте: Построение главного мега-меню, выпадающих списков в навигации, бокового меню фильтров.
//GET /api/categories/flat — Получение списка категорий в виде плоского массива с указанием level и parent_id. Для чего: Для простых выпадающих списков в админке (при выборе родительской категории), для построения breadcrumbs (хлебных крошек) на клиенте.
//GET /api/categories/{slug} — Получение данных категории по её SEO-ссылке (slug). Что возвращать: Поля категории + дополнительные данные для страницы каталога: Список дочерних категорий (для поднавигации на странице). Breadcrumbs (Хлебные крошки): Массив [{name: "Главная", slug: "/"}, {name: "Электроника", slug: "/elektronika"}, ...]. Легко строится из поля path. Прикрепленные фильтры/атрибуты (если у вас есть такая связь в БД). SEO-теги (meta_title, meta_description).
//GET /api/categories/{slug}/products — Получение товаров этой категории. Критически важно: Товары должны включать товары из всех вложенных дочерних категорий. Пользователь, зайдя в "Электроника", ожидает увидеть и телефоны, и ноутбуки.Логика: Используйте поле path для эффективного поиска. Запрос типа: WHERE categories.path LIKE '{current_category_path}%'.Дополнение: Должен поддерживать пагинацию, сортировку (по цене, новизне, рейтингу) и фильтрацию по брендам/атрибутам.

func (h *Handler) GetTree(c *gin.Context) *webtool.APIError {
	return nil
}

func (h *Handler) GetFlat(c *gin.Context) *webtool.APIError {
	return nil
}

func (h *Handler) GetBySlug(c *gin.Context) *webtool.APIError {
	return nil
}

func (h *Handler) GetProductsBySlug(c *gin.Context) *webtool.APIError {
	return nil
}
