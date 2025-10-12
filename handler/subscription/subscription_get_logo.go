package subscription

import (
	"net/http"
	"os"
	"path/filepath"
	"subscriptionplus/server/pkg/httpx/httperr"

	"github.com/gorilla/mux"
)

func (h *Handler) GetSubscriptionLogoHandler(w http.ResponseWriter, r *http.Request) error {
	logo := mux.Vars(r)["logo"]

	if logo == "" || filepath.Clean(logo) != logo {
		return httperr.BadRequest("filename not specified")
	}

	filePath := "media/subscriptions/" + logo

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return httperr.NotFound("logo file not found")
	}

	http.ServeFile(w, r, filePath)
	return nil
}
