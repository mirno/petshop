package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/mirno/petshop/internal/drivers/api"
	"github.com/mirno/petshop/internal/drivers/inmemorykvstore"
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/pkg/petshop"
)

func TestCreatePetStoresAndReturnsCreatedPet(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	router := newRouter(store)

	rec := performRequest(router, http.MethodPost, "/pets", `{
		"name": "Rex",
		"type": "dog",
		"price": 80,
		"metadata": {
			"Region": "Europe",
			"color": "brown"
		}
	}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var got petshop.Pet
	decodeJSON(t, rec.Body.String(), &got)
	if got.Id == uuid.Nil || got.Name != "Rex" || got.Type != petshop.Dog || got.Price != 80 {
		t.Fatalf("unexpected pet response: %#v", got)
	}
	if got.Metadata.Region == nil || *got.Metadata.Region != "Europe" {
		t.Fatalf("expected Region metadata to be returned, got %#v", got.Metadata)
	}
	if got.Metadata.AdditionalProperties["color"] != "brown" {
		t.Fatalf("expected additional metadata to be returned, got %#v", got.Metadata.AdditionalProperties)
	}
	if rec.Header().Get(echo.HeaderLocation) != "/pets/"+got.Id.String() {
		t.Fatalf("unexpected Location header: %q", rec.Header().Get(echo.HeaderLocation))
	}

	stored, err := store.Get(got.Id.String())
	if err != nil {
		t.Fatalf("expected pet to be stored: %v", err)
	}
	if stored.Name != got.Name || stored.Type != entities.Dog || stored.Price != got.Price {
		t.Fatalf("unexpected stored pet: %#v", stored)
	}
}

func TestListPetsFiltersByTypeAndRegion(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	router := newRouter(store)

	rexID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	slytherinID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	if err := store.Save(rexID.String(), entities.Pet{
		Id:       rexID,
		Name:     "Rex",
		Type:     entities.Dog,
		Price:    80,
		Metadata: metadataWithRegion("Europe"),
	}); err != nil {
		t.Fatalf("save Rex: %v", err)
	}
	if err := store.Save(slytherinID.String(), entities.Pet{
		Id:       slytherinID,
		Name:     "Slytherin",
		Type:     entities.Snake,
		Price:    400,
		Metadata: metadataWithRegion("North America"),
	}); err != nil {
		t.Fatalf("save Slytherin: %v", err)
	}

	rec := performRequest(router, http.MethodGet, "/pets?type=dog&region=Europe", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got []petshop.Pet
	decodeJSON(t, rec.Body.String(), &got)
	if len(got) != 1 || got[0].Id != rexID || got[0].Name != "Rex" {
		t.Fatalf("expected only Rex, got %#v", got)
	}
}

func TestGetPetReturnsNotFoundForMissingPet(t *testing.T) {
	router := newRouter(inmemorykvstore.NewInMemoryKVStore[entities.Pet]())
	missingID := uuid.MustParse("00000000-0000-0000-0000-000000000404")

	rec := performRequest(router, http.MethodGet, "/pets/"+missingID.String(), "")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}

	var got petshop.Error
	decodeJSON(t, rec.Body.String(), &got)
	if got.Code != "not_found" {
		t.Fatalf("expected not_found error, got %#v", got)
	}
}

func TestPatchPetUpdatesProvidedFields(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	router := newRouter(store)
	petID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if err := store.Save(petID.String(), entities.Pet{
		Id:       petID,
		Name:     "Rex",
		Type:     entities.Dog,
		Price:    80,
		Metadata: metadataWithRegion("Europe"),
	}); err != nil {
		t.Fatalf("save pet: %v", err)
	}

	rec := performRequest(router, http.MethodPatch, "/pets/"+petID.String(), `{
		"name": "Max",
		"price": 95,
		"metadata": {
			"Region": "Asia"
		}
	}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got petshop.Pet
	decodeJSON(t, rec.Body.String(), &got)
	if got.Name != "Max" || got.Type != petshop.Dog || got.Price != 95 {
		t.Fatalf("unexpected patched response: %#v", got)
	}
	if got.Metadata.Region == nil || *got.Metadata.Region != "Asia" {
		t.Fatalf("expected updated Region metadata, got %#v", got.Metadata)
	}

	stored, err := store.Get(petID.String())
	if err != nil {
		t.Fatalf("get stored pet: %v", err)
	}
	if stored.Name != "Max" || stored.Type != entities.Dog || stored.Price != 95 || stored.Metadata.Region == nil || *stored.Metadata.Region != "Asia" {
		t.Fatalf("unexpected stored pet after patch: %#v", stored)
	}
}

func TestUpdatePetReplacesExistingPet(t *testing.T) {
	store := inmemorykvstore.NewInMemoryKVStore[entities.Pet]()
	router := newRouter(store)
	petID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if err := store.Save(petID.String(), entities.Pet{
		Id:       petID,
		Name:     "Rex",
		Type:     entities.Dog,
		Price:    80,
		Metadata: metadataWithRegion("Europe", "color", "brown"),
	}); err != nil {
		t.Fatalf("save pet: %v", err)
	}

	rec := performRequest(router, http.MethodPut, "/pets/"+petID.String(), `{
		"name": "Slytherin",
		"type": "snake",
		"price": 400
	}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got petshop.Pet
	decodeJSON(t, rec.Body.String(), &got)
	if got.Id != petID || got.Name != "Slytherin" || got.Type != petshop.Snake || got.Price != 400 {
		t.Fatalf("unexpected update response: %#v", got)
	}

	stored, err := store.Get(petID.String())
	if err != nil {
		t.Fatalf("get stored pet: %v", err)
	}
	if stored.Name != "Slytherin" || stored.Type != entities.Snake || stored.Price != 400 || stored.Metadata.Region != nil || len(stored.Metadata.AdditionalProperties) != 0 {
		t.Fatalf("expected replaced pet without old metadata, got %#v", stored)
	}
}

func TestCreatePetRejectsInvalidRequest(t *testing.T) {
	router := newRouter(inmemorykvstore.NewInMemoryKVStore[entities.Pet]())

	rec := performRequest(router, http.MethodPost, "/pets", `{
		"name": "",
		"type": "cat",
		"price": -1
	}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var got petshop.Error
	decodeJSON(t, rec.Body.String(), &got)
	if got.Code != "bad_request" {
		t.Fatalf("expected bad_request error, got %#v", got)
	}
}

// Helper test functions...

func newRouter(store *inmemorykvstore.InMemoryKVStore[entities.Pet]) *echo.Echo {
	router := echo.New()
	petshop.RegisterHandlers(router, api.NewAPIHandlers(store))
	return router
}

func performRequest(router *echo.Echo, method string, path string, body string) *httptest.ResponseRecorder {
	var requestBody *strings.Reader
	if body == "" {
		requestBody = strings.NewReader("")
	} else {
		requestBody = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	return rec
}

func decodeJSON(t *testing.T, body string, dest any) {
	t.Helper()

	if err := json.Unmarshal([]byte(body), dest); err != nil {
		t.Fatalf("decode response %q: %v", body, err)
	}
}

func metadataWithRegion(region string, additional ...string) petshop.Metadata {
	metadata := petshop.Metadata{
		Region: &region,
	}
	for i := 0; i+1 < len(additional); i += 2 {
		metadata.Set(additional[i], additional[i+1])
	}

	return metadata
}
