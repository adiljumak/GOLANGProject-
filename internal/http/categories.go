package http

import (
	"Lecture6/internal/models"
	"Lecture6/internal/store"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
)

type CategoryResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewCategoryResource(store store.Store, cache *lru.TwoQueueCache) *CategoryResource {
	return &CategoryResource{
		store: store,
		cache: cache,
	}
}

func (cr *CategoryResource) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", cr.CreateCategory)
	r.Get("/", cr.AllCategories)
	return r
}

func (cr *CategoryResource) CreateCategory(writer http.ResponseWriter, request *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(request.Body).Decode(category); err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(writer, "Error", err)
		return
	}

	if err := cr.store.Categories().Create(request.Context(), category); err != nil {
		fmt.Fprintf(writer, "Error", err)
		return
	}

	cr.cache.Purge()
	writer.WriteHeader(http.StatusCreated)
}

func (cr *CategoryResource) AllCategories(writer http.ResponseWriter, request *http.Request) {
	queryValues := request.URL.Query()
	filter := &models.CategoriesFilter{}
	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		categoriesFromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(writer, request, categoriesFromCache)
			return
		}
		filter.Query = &searchQuery
	}

	categories, err := cr.store.Categories().All(request.Context(), filter)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "Error", err)
		return
	}
	if searchQuery != "" {
		cr.cache.Add(searchQuery, categories)
	}

	render.JSON(writer, request, categories)
}

//r.Post("/categories", )
//r.Get("/categories", )
//r.Get("/categories/{id}", func(writer http.ResponseWriter, request *http.Request) {
//	idStr := chi.URLParam(request, "id")
//	id, err := strconv.Atoi(idStr)
//
//	if err != nil {
//
//		writer.WriteHeader(http.StatusBadRequest)
//		fmt.Fprintf(writer, "Error", err)
//		return
//	}
//
//	category, err := s.store.Categories().ByID(request.Context(), id)
//	if err != nil {
//		writer.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(writer, "DB err: #{err}")
//		return
//	}
//	render.JSON(writer, request, category)
//})
//r.Put("/categories", func(writer http.ResponseWriter, request *http.Request) {
//	category := new(models.Category)
//	if err := json.NewDecoder(request.Body).Decode(category); err != nil {
//		writer.WriteHeader(http.StatusUnprocessableEntity)
//		fmt.Fprintf(writer, "Error", err)
//		return
//	}
//	err := validation.ValidateStruct(category,
//		validation.Field(&category.ID, validation.Required),
//		validation.Field(&category.Name, validation.Required),
//	)
//	validation.Field(&category.ID, validation.Required)
//	if err != nil {
//		writer.WriteHeader(http.StatusUnprocessableEntity)
//		fmt.Fprintf(writer, "Unknown")
//		return
//	}
//	if err := s.store.Categories().Update(request.Context(), category); err != nil {
//		writer.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(writer, "Error", err)
//		return
//	}
//})
//r.Delete("/categories/{id}", func(writer http.ResponseWriter, request *http.Request) {
//	idStr := chi.URLParam(request, "id")
//	id, err := strconv.Atoi(idStr)
//	if err != nil {
//		writer.WriteHeader(http.StatusBadRequest)
//		fmt.Fprintf(writer, "Error", err)
//		return
//	}
//	if err := s.store.Categories().Delete(request.Context(), id); err != nil {
//		writer.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(writer, "Error", err)
//		return
//	}
//})
