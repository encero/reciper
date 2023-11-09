package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/encero/reciper/ent"
	"github.com/encero/reciper/ent/recipe"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const workerQueue = "api-server"

func Run(ctx context.Context, entc *ent.Client, lg *zap.Logger, natsURL string) error {
	// Run the auto migration tool.
	if err := entc.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %w", err)
	}

	conn, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("nats connect url: %s: %w", natsURL, err)
	}

	ec, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats encoded conn: %w", err)
	}

	h := handlers{
		entc: entc,
		ec:   ec,
		lg:   lg,
	}

	_, err = ec.QueueSubscribe(HandlersRecipesUpsert, workerQueue, h.Upsert)
	if err != nil {
		return fmt.Errorf("recipe.upsert subscription: %w", err)
	}

	_, err = ec.QueueSubscribe(HandlersRecipeDetail, workerQueue, h.Detail)
	if err != nil {
		return fmt.Errorf("recipe.detail.* subscription: %w", err)
	}

	_, err = ec.QueueSubscribe(HandlersRecipeList, workerQueue, h.List)
	if err != nil {
		return fmt.Errorf("recipe.list subscription: %w", err)
	}

	_, err = ec.QueueSubscribe(HandlersRecipeDelete, workerQueue, h.Delete)
	if err != nil {
		return fmt.Errorf("recipe.delete subscription: %w", err)
	}

	_, err = ec.QueueSubscribe(HandlersMarkAsPlanned, workerQueue, h.MarksAsPlanned)
	if err != nil {
		return fmt.Errorf("recipe.planned.* subscription: %w", err)
	}

	_, err = ec.QueueSubscribe(HandlersMarkAsCooked, workerQueue, h.MarkAsCooked)
	if err != nil {
		return fmt.Errorf("recipe.cooked.* subscription: %w", err)
	}

	lg.Info("Api server started")
	<-ctx.Done()

	err = conn.Drain()
	if err != nil {
		return fmt.Errorf("nats drain: %w", err)
	}

	return nil
}

type handlers struct {
	entc *ent.Client
	ec   *nats.EncodedConn
	lg   *zap.Logger
}

const handlerTimeout = time.Second

const HandlersRecipeList = "recipes.list"

func (h *handlers) List(msg *nats.Msg) {
	lg := h.lg.With(ZapRequestID(), ZapHandler(HandlersRecipeList))
	lg.Debug("Incomming recipe list request")

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	recipes, err := h.entc.Recipe.Query().
		Order(ent.Asc(recipe.FieldTitle)).
		All(ctx)

	if err != nil {
		lg.Error("Retrieving list of recipes", zap.Error(err))

		err = h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	list := Envelope[List]{
		Status: StatusSuccess,
		Data:   List{},
	}

	for _, r := range recipes {
		recipe, err := EntToRecipe(ctx, r)
		if err != nil {
			lg.Error("loading recipe data", zap.Error(err))

			err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
			logNatsPublishError(lg, err)

			return
		}

		list.Data = append(list.Data, recipe)
	}

	err = h.ec.Publish(msg.Reply, list)
	logNatsPublishError(lg, err)
}

const HandlersRecipeDelete = "recipes.delete.*"

func (h *handlers) Delete(msg *nats.Msg) {
	lg := h.lg.With(ZapRequestID(), ZapHandler(HandlersRecipeDelete))
	id := uuid.MustParse(strings.Split(msg.Subject, ".")[2])

	lg.Debug("Incomming recipe delete request")

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	err := h.entc.Recipe.DeleteOneID(id).Exec(ctx)
	if err != nil {
		lg.Error("Deleting recipe", zap.Error(err))

		err = h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = h.ec.Publish(msg.Reply, Ack{Status: StatusSuccess})
	logNatsPublishError(lg, err)
}

const HandlersRecipeDetail = "recipes.detail.*"

func (h *handlers) Detail(msg *nats.Msg) {
	id := uuid.MustParse(strings.Split(msg.Subject, ".")[2])

	lg := h.lg.With(ZapHandler(HandlersRecipesUpsert), ZapRecipeID(id), ZapRequestID())
	lg.Debug("Incomming recipe list request")

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	eRecipe, err := h.entc.Recipe.Get(ctx, id)
	if err != nil && ent.IsNotFound(err) {
		lg.Info("recipe not found")

		err = h.ec.Publish(msg.Reply, Ack{Status: StatusNotFound})
		logNatsPublishError(lg, err)

		return
	}

	if err != nil {
		lg.Error("Retrieve recipe", zap.Error(err))

		err = h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	recipe, err := EntToRecipe(ctx, eRecipe)
	if err != nil {
		lg.Error("loading recipe data", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = h.ec.Publish(msg.Reply, Envelope[Recipe]{
		Status: StatusSuccess,
		Data:   recipe,
	})
	logNatsPublishError(lg, err)
}

const HandlersRecipesUpsert = "recipes.upsert"

func (h *handlers) Upsert(subject, reply string, r Recipe) {
	lg := h.lg.With(ZapHandler(HandlersRecipesUpsert), ZapRecipeID(r.ID), ZapRequestID())

	lg.Debug("Incomming recipe upsert", ZapRecipe(r))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	recipe, err := h.entc.Recipe.Get(ctx, r.ID)
	if err != nil && !ent.IsNotFound(err) {
		lg.Error("Retrieving recipe for upsert", zap.Error(err))

		err = h.ec.Publish(reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	if ent.IsNotFound(err) {
		lg.Info("About to create new recipe")

		recipe, err = h.entc.Recipe.Create().
			SetTitle(r.Name).
			SetID(r.ID).
			Save(ctx)
	} else {
		lg.Info("About to update recipe")

		recipe, err = recipe.Update().
			SetTitle(r.Name).
			Save(ctx)
	}

	if err != nil {
		lg.Error("Recipe upsert", zap.Error(err))

		err = h.ec.Publish(reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	out, err := EntToRecipe(ctx, recipe)
	if err != nil {
		lg.Error("loading recipe data", zap.Error(err))

		err := h.ec.Publish(reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = h.ec.Publish(reply, Envelope[Recipe]{
		Status: StatusSuccess,
		Data:   out,
	})
	logNatsPublishError(lg, err)
}

const HandlersMarkAsPlanned = "recipes.planned.*"

func (h handlers) MarksAsPlanned(subject, reply string, req RequestPlanned) {
	id := uuid.MustParse(strings.Split(subject, ".")[2])
	lg := h.lg.With(ZapRequestID(), ZapHandler(HandlersMarkAsPlanned), ZapRecipeID(id))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	recipe, err := h.entc.Recipe.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			lg.Info("recipe not found")

			err := h.ec.Publish(reply, Ack{Status: StatusNotFound})
			logNatsPublishError(lg, err)

			return
		}

		lg.Error("load recipe", zap.Error(err))

		err := h.ec.Publish(reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	lg.Sugar().Infof("about to set recipe as planned:%v", req.Planned)

	_, err = recipe.Update().
		SetPlanned(req.Planned).
		Save(ctx)
	if err != nil {
		lg.Error("save updated recipe", zap.Error(err))

		err := h.ec.Publish(reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = h.ec.Publish(reply, Ack{Status: StatusSuccess})
	logNatsPublishError(lg, err)
}

const HandlersMarkAsCooked = "recipes.cooked.*"

func (h *handlers) MarkAsCooked(msg *nats.Msg) {
	id := uuid.MustParse(strings.Split(msg.Subject, ".")[2])
	lg := h.lg.With(ZapRequestID(), ZapHandler(HandlersMarkAsCooked), ZapRecipeID(id))

	ctx, canncel := context.WithTimeout(context.Background(), time.Second)
	defer canncel()

	tx, err := h.entc.BeginTx(ctx, nil)
	if err != nil {
		lg.Error("open transaction", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	defer tx.Rollback() // nolint: errcheck

	// todo: refactor to helper
	recipe, err := tx.Recipe.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			lg.Info("recipe not found")

			err := h.ec.Publish(msg.Reply, Ack{Status: StatusNotFound})
			logNatsPublishError(lg, err)

			return
		}

		lg.Error("load recipe", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	if !recipe.Planned {
		lg.Info("recipe is not planned for cooking, skipping")

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusSuccess})
		logNatsPublishError(lg, err)

		_ = tx.Commit()

		return
	}

	_, err = tx.CookingHistory.Create().
		SetID(uuid.New()).
		SetRecipe(recipe).
		SetCookedAt(time.Now()).
		Save(ctx)
	if err != nil {
		lg.Error("create recipe history", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	_, err = recipe.Update().
		SetPlanned(false).
		Save(ctx)
	if err != nil {
		lg.Error("setting recipe as not planned", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = tx.Commit()
	if err != nil {
		lg.Error("transaction commit", zap.Error(err))

		err := h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		logNatsPublishError(lg, err)

		return
	}

	err = h.ec.Publish(msg.Reply, Ack{Status: StatusSuccess})
	logNatsPublishError(lg, err)
}

func logNatsPublishError(lg *zap.Logger, err error) {
	if err != nil {
		lg.Error("Publishing to nats", zap.Error(err))
	}
}
