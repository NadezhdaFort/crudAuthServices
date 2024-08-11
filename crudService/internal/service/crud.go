package service

import (
	"crud/internal/domain"
	"crud/internal/repository/recipedb"
	"github.com/google/uuid"
)

var recipes recipedb.DB

func Init(DB recipedb.DB) {
	recipes = DB
}

func Get(id string) (*domain.Recipe, error) {
	recipe, err := recipes.Get(id)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

func Delete(id string) error {
	return recipes.Delete(id)
}

func AddOrUpd(r *domain.Recipe) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	newRecipe := domain.Recipe{
		ID:          r.ID,
		AuthorID:    r.AuthorID,
		Name:        r.Name,
		Ingredients: r.Ingredients,
		Temperature: r.Temperature,
	}
	if err := recipes.Set(newRecipe.ID, &newRecipe); err != nil {
		return err
	}
	return nil
}
