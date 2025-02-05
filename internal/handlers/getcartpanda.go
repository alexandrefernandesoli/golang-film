package handlers

import (
	"goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type CartpandaHandlerParams struct {
	ProductsStore store.ProductsStore
}

type CartpandaHandler struct {
	ProductsStore store.ProductsStore
}

func NewCartpandaHandler(params CartpandaHandlerParams) *CartpandaHandler {
	return &CartpandaHandler{
		ProductsStore: params.ProductsStore,
	}
}

func (h *CartpandaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	_, ok := r.Context().Value(middleware.UserKey).(*store.User)

	if !ok {
		// redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	if len(h.ProductsStore.GetProducts()) == 0 {
		err := h.ProductsStore.LoadProducts()
		if err != nil {
			http.Error(w, "Error loading products", http.StatusInternalServerError)
			return
		}
	}

	c := templates.Cartpanda(h.ProductsStore.GetProducts(), h.ProductsStore.GetLastTimeLoaded())
	err := templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
