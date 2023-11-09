package api

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapHandler(handler string) zap.Field {
	return zap.String("handler", handler)
}

func ZapRecipeID(recipeID uuid.UUID) zap.Field {
	return zap.String("recipeID", recipeID.String())
}

func ZapRequestID() zap.Field {
	return zap.String("requestId", uuid.New().String())
}

func ZapRecipe(recipe Recipe) zap.Field {
	return zap.Object("recipe", zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		oe.AddString("ID", recipe.ID.String())
		oe.AddString("title", recipe.Name)

		return nil
	}))
}
