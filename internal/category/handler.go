package category

import (
	"errors"
	"go-shop/internal/models"
	"go-shop/pkg/webtool"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(r *gin.RouterGroup) {
	group := r.Group("/category")

	group.POST("", webtool.MakeHandler(h.Create))
	group.GET("/:id", webtool.MakeHandler(h.Get))
	group.PUT("/:id", webtool.MakeHandler(h.Update))
	group.DELETE("/:id", webtool.MakeHandler(h.Delete))
	group.PATCH("/:id/order", webtool.MakeHandler(h.UpdateOder))
	group.PATCH("/:id/status", webtool.MakeHandler(h.UpdateStatus))

	group.GET("/tree", webtool.MakeHandler(h.GetTree))
	group.GET("/flat", webtool.MakeHandler(h.GetFlat))
	group.GET("/slug/:slug", webtool.MakeHandler(h.GetBySlug))
	group.GET("/slug/:slug/products", webtool.MakeHandler(h.GetProductsBySlug))
}

// Группа 1: Управление категориями (Админка / Бэкенд)
// Эти эндпоинты нужны для заполнения и управления деревом категорий.
// GET /api/category/{id} — Получение данных конкретной категории (для формы редактирования).
// PUT /api/category/{id} — Обновление категории. Логика: Особо важная часть — обработка изменения parent_id. При перемещении категории в другую ветку нужно рекурсивно пересчитать level и path для всех её дочерних категорий. Это самая сложная операция.
// DELETE /api/category/{id} — Удаление категории. Логика: Каскадное удаление (как настроено в БД ON DELETE CASCADE) или перенос товаров в другую категорию (более безопасный вариант).
// PATCH /api/category/{id}/order — Изменение порядка отображения (display_order) в рамках одного уровня иерархии (среди "братьев и сестер").
// PATCH /api/category/{id}/status — Активация/деактивация категории (is_active). Деактивация обычно скрывает категорию и все её дочерние категории с фронта.

// POST /api/category — Создание новой категории.
// Логика: Проверка уникальности slug.
// Автоматический расчет level и path на основе parent_id. Если parent_id = 0 (или null) — это корневая категория (level=1, path=/id/).
func (h *Handler) Create(c *gin.Context) *webtool.APIError {
	var input CreateRequestDTO
	err := c.BindJSON(&input)
	if err != nil {
		return BadInput(err)
	}

	category, err := h.service.CreateCategory(CreateRequestDTOToCategory(&input))
	if err != nil {
		return CantCreate(err)
	}

	c.JSON(http.StatusOK, ToCreateResponseDTO(category))
	return nil
}
func (h *Handler) Get(c *gin.Context) *webtool.APIError {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return NoID(err)
	}

	logrus.Println("category", ID)
	category, err := h.service.GetCategory(&models.Category{ID: ID})
	if err != nil {
		return CantFind(err)
	}

	c.JSON(http.StatusOK, ToGetResponseDTO(category))
	return nil
}

func (h *Handler) Update(c *gin.Context) *webtool.APIError {
	var input UpdateRequestDTO
	err := c.BindJSON(&input)
	if err != nil {
		logrus.Errorf(err.Error())
		return BadInput(err)
	}
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return NoID(err)
	}

	category, err := h.service.UpdateCategory(UpdateRequestDTOToCategory(ID, &input))
	if err != nil {
		return CantUpdate(err)
	}

	c.JSON(http.StatusOK, ToUpdateResponseDTO(category))
	return nil
}
func (h *Handler) Delete(c *gin.Context) *webtool.APIError {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return NoID(err)
	}

	category, err := h.service.DeleteCategory(DeleteRequestDTOToCategory(ID))
	if err != nil {
		return CantDelete(err)
	}

	c.JSON(http.StatusOK, ToDeleteResponseDTO(category))
	return nil

}
func (h *Handler) UpdateOder(c *gin.Context) *webtool.APIError {
	var input UpdateOrderRequestDTO
	err := c.BindJSON(&input)
	if err != nil {
		return BadInput(err)
	}
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return NoID(err)
	}

	category, err := h.service.UpdateOrder(UpdateOrderRequestDTOToCategory(ID, &input))
	if err != nil {
		return CantUpdateOrder(err)
	}

	c.JSON(http.StatusOK, ToUpdateOrderResponseDTO(category))
	return nil
}
func (h *Handler) UpdateStatus(c *gin.Context) *webtool.APIError {
	var input UpdateStatusRequestDTO
	err := c.BindJSON(&input)
	if err != nil {
		return BadInput(err)
	}
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return NoID(err)
	}

	category, err := h.service.UpdateStatus(UpdateStatusRequestDTOToCategory(ID, &input))
	if err != nil {
		return CantUpdateStatus(err)
	}

	c.JSON(http.StatusOK, ToUpdateStatusResponseDTO(category))

	return nil
}

//Группа 2: Отображение и навигация (Публичный API для фронта)
//GET /api/category/tree — Самый важный эндпоинт. Получение всего дерева категорий в виде вложенной структуры (Nested Set). Ответ: Массив корневых категорий, каждая из которых содержит массив своих дочерних категорий (рекурсивно). Только активные (is_active=true). Для чего на фронте: Построение главного мега-меню, выпадающих списков в навигации, бокового меню фильтров.
//GET /api/category/flat — Получение списка категорий в виде плоского массива с указанием level и parent_id. Для чего: Для простых выпадающих списков в админке (при выборе родительской категории), для построения breadcrumbs (хлебных крошек) на клиенте.
//GET /api/category/{slug} — Получение данных категории по её SEO-ссылке (slug). Что возвращать: Поля категории + дополнительные данные для страницы каталога: Список дочерних категорий (для поднавигации на странице). Breadcrumbs (Хлебные крошки): Массив [{name: "Главная", slug: "/"}, {name: "Электроника", slug: "/elektronika"}, ...]. Легко строится из поля path. Прикрепленные фильтры/атрибуты (если у вас есть такая связь в БД). SEO-теги (meta_title, meta_description).
//GET /api/category/{slug}/products — Получение товаров этой категории. Критически важно: Товары должны включать товары из всех вложенных дочерних категорий. Пользователь, зайдя в "Электроника", ожидает увидеть и телефоны, и ноутбуки.Логика: Используйте поле path для эффективного поиска. Запрос типа: WHERE category.path LIKE '{current_category_path}%'.Дополнение: Должен поддерживать пагинацию, сортировку (по цене, новизне, рейтингу) и фильтрацию по брендам/атрибутам.

func (h *Handler) GetTree(c *gin.Context) *webtool.APIError {
	tree, err := h.service.GetCategoryTree()
	if err != nil {
		return CantGetTree(err)
	}

	c.JSON(http.StatusOK, tree)

	return nil
}

func (h *Handler) GetFlat(c *gin.Context) *webtool.APIError {
	categories, err := h.service.GetCategories()
	if err != nil {
		return CantGetList(err)
	}

	c.JSON(http.StatusOK, categories)

	return nil
}

func (h *Handler) GetBySlug(c *gin.Context) *webtool.APIError {
	slug := c.Param("slug")

	if slug == "" {
		return NoSlug(errors.New("no slug"))
	}
	category, err := h.service.GetCategoryBySlug(slug)
	if err != nil {
		return CantFind(err)
	}

	c.JSON(http.StatusOK, ToGetBySlugResponseDTO(category))
	return nil
}

func (h *Handler) GetProductsBySlug(c *gin.Context) *webtool.APIError {
	return nil
}
