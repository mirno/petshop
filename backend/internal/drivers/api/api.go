package api

import (
	"errors"
	"math"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/mirno/petshop/internal/entities"
	"github.com/mirno/petshop/internal/usecases"
	"github.com/mirno/petshop/pkg/petshop"
)

var _ petshop.ServerInterface = (*APIHandlers)(nil)

const (
	metadataRegionKey = "Region"

	errorCodeBadRequest = "bad_request"
	errorCodeNotFound   = "not_found"
	errorCodeInternal   = "internal_error"
)

// APIHandlers implements the generated Echo v5 server interface.
type APIHandlers struct {
	store usecases.Store[entities.Pet]
}

func NewAPIHandlers(store usecases.Store[entities.Pet]) *APIHandlers {
	return &APIHandlers{store: store}
}

func (handlers *APIHandlers) ListPets(ctx *echo.Context, params petshop.ListPetsParams) error {
	if params.Type != nil && !params.Type.Valid() {
		return badRequest(ctx, "type must be one of: dog, snake")
	}

	keys, err := handlers.store.Keys()
	if err != nil {
		return internalError(ctx)
	}
	slices.Sort(keys)

	pets := make([]petshop.Pet, 0, len(keys))
	for _, key := range keys {
		pet, err := handlers.store.Get(key)
		if err != nil {
			if errors.Is(err, usecases.ErrKeyNotFoundError) {
				continue
			}
			return internalError(ctx)
		}

		if params.Type != nil && petshop.AnimalType(pet.Type) != *params.Type {
			continue
		}
		if params.Region != nil && metadataRegion(pet.Metadata) != *params.Region {
			continue
		}

		pets = append(pets, pet)
	}

	return ctx.JSON(http.StatusOK, pets)
}

func (handlers *APIHandlers) CreatePet(ctx *echo.Context) error {
	var request petshop.CreatePetRequest
	if err := ctx.Bind(&request); err != nil {
		return badRequest(ctx, "request body must be valid JSON")
	}
	if message := validatePetRequest(request.Name, request.Type, request.Price, request.Status); message != "" {
		return badRequest(ctx, message)
	}

	pet := petFromRequest(petshop.PetId{}, request)
	pet.Id = uuid.New()
	if err := handlers.store.Save(pet.Id.String(), pet); err != nil {
		return internalError(ctx)
	}

	location := "/pets/" + pet.Id.String()
	ctx.Response().Header().Set(echo.HeaderLocation, location)
	return ctx.JSON(http.StatusCreated, pet)
}

func (handlers *APIHandlers) GetPet(ctx *echo.Context, petId petshop.PetId) error {
	pet, err := handlers.store.Get(petId.String())
	if err != nil {
		return handleStoreGetError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, pet)
}

func (handlers *APIHandlers) PatchPet(ctx *echo.Context, petId petshop.PetId) error {
	pet, err := handlers.store.Get(petId.String())
	if err != nil {
		return handleStoreGetError(ctx, err)
	}

	var request petshop.PatchPetRequest
	if err := ctx.Bind(&request); err != nil {
		return badRequest(ctx, "request body must be valid JSON")
	}
	if message := validatePatchPetRequest(request); message != "" {
		return badRequest(ctx, message)
	}

	if request.Name != nil {
		pet.Name = *request.Name
	}
	if request.Type != nil {
		pet.Type = *request.Type
	}
	if request.Price != nil {
		pet.Price = *request.Price
	}
	if request.Metadata != nil {
		pet.Metadata = *request.Metadata
	}

	if err := handlers.store.Save(pet.Id.String(), pet); err != nil {
		return internalError(ctx)
	}

	return ctx.JSON(http.StatusOK, pet)
}

func (handlers *APIHandlers) UpdatePet(ctx *echo.Context, petId petshop.PetId) error {
	if _, err := handlers.store.Get(petId.String()); err != nil {
		return handleStoreGetError(ctx, err)
	}

	var request petshop.UpdatePetRequest
	if err := ctx.Bind(&request); err != nil {
		return badRequest(ctx, "request body must be valid JSON")
	}
	if message := validatePetRequest(request.Name, request.Type, request.Price, request.Status); message != "" {
		return badRequest(ctx, message)
	}

	pet := petFromRequest(petId, request)
	if err := handlers.store.Save(pet.Id.String(), pet); err != nil {
		return internalError(ctx)
	}

	return ctx.JSON(http.StatusOK, pet)
}

func petFromRequest(id petshop.PetId, request petshop.CreatePetRequest) entities.Pet {
	pet := entities.Pet{
		Id:     id,
		Name:   request.Name,
		Type:   request.Type,
		Price:  request.Price,
		Status: request.Status,
	}
	if request.Metadata != nil {
		pet.Metadata = *request.Metadata
	}

	return pet
}

func validatePetRequest(name string, animalType petshop.AnimalType, price float64, status *petshop.PetStatus) string {
	if strings.TrimSpace(name) == "" {
		return "name is required"
	}
	if !animalType.Valid() {
		return "type must be one of: dog, snake"
	}
	if !validEntityPrice(price) {
		return "price must be zero or greater"
	}
	if status != nil && !status.Valid() {
		return "status must be one of: available, reserved, sold"
	}

	return ""
}

func validatePatchPetRequest(request petshop.PatchPetRequest) string {
	if request.Name != nil && strings.TrimSpace(*request.Name) == "" {
		return "name must not be empty"
	}
	if request.Type != nil && !request.Type.Valid() {
		return "type must be one of: dog, snake"
	}
	if request.Price != nil && !validEntityPrice(*request.Price) {
		return "price must be zero or greater"
	}
	if request.Status != nil && !request.Status.Valid() {
		return "status must be one of: available, reserved, sold"
	}

	return ""
}

func validEntityPrice(price float64) bool {
	return !math.IsNaN(price) &&
		!math.IsInf(price, 0) &&
		price >= 0
}

func metadataRegion(metadata petshop.Metadata) string {
	if metadata.Region != nil {
		return *metadata.Region
	}

	return metadata.AdditionalProperties[metadataRegionKey]
}

func handleStoreGetError(ctx *echo.Context, err error) error {
	if errors.Is(err, usecases.ErrKeyNotFoundError) {
		return notFound(ctx, "pet not found")
	}

	return internalError(ctx)
}

func badRequest(ctx *echo.Context, message string) error {
	return jsonError(ctx, http.StatusBadRequest, errorCodeBadRequest, message)
}

func notFound(ctx *echo.Context, message string) error {
	return jsonError(ctx, http.StatusNotFound, errorCodeNotFound, message)
}

func internalError(ctx *echo.Context) error {
	return jsonError(ctx, http.StatusInternalServerError, errorCodeInternal, "internal server error")
}

func jsonError(ctx *echo.Context, status int, code string, message string) error {
	return ctx.JSON(status, petshop.Error{
		Code:    code,
		Message: message,
	})
}
