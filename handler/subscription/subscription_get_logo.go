package subscription

import (
	"net/http"
	"os"
	"subscriptionplus/server/pkg/httpx/httperr"
	"subscriptionplus/server/pkg/machinelearning"
)

func (h *Handler) GetSubscriptionLogoHandler(w http.ResponseWriter, r *http.Request) error {
	nlp := machinelearning.NewNLPBuilder()

	name := r.URL.Query().Get("name")
	if name == "" {
		return httperr.NotFound("image file not found")
	}

	img := nlp.GetSubscriptionImage(name)

	if img == "" {
		return httperr.BadRequest("filename not specified")
	}

	filePath := "media/subscriptions/" + img

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return httperr.NotFound("logo file not found")
	}

	http.ServeFile(w, r, filePath)
	return nil
}
