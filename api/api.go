package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "modernc.org/sqlite"

	"github.com/encero/reciper-api/ent"
	"github.com/encero/reciper-api/ent/recipe"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const workerQueue = "api-server"

func Run(ctx context.Context, url string) error {
	sqldb, err := sql.Open("sqlite", "file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	defer sqldb.Close()

	entc := ent.NewClient(ent.Driver(entsql.OpenDB("sqlite3", sqldb)))
	defer entc.Close()

	// Run the auto migration tool.
	if err := entc.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
	}

	conn, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}

	ec, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("nats encoded conn: %w", err)
	}

	h := handlers{
		entc: entc,
		ec:   ec,
	}

	_, err = ec.QueueSubscribe("recipes.upsert", workerQueue, h.Upsert)
	if err != nil {
		return fmt.Errorf("recipe.upsert subscription: %w", err)
	}

	_, err = ec.QueueSubscribe("recipes.detail.*", workerQueue, h.Detail)
	if err != nil {
		return fmt.Errorf("recipe.detail.* subscription: %w", err)
	}

	_, err = ec.QueueSubscribe("recipes.list", workerQueue, h.List)
	if err != nil {
		return fmt.Errorf("recipe.list subscription: %w", err)
	}

	_, err = ec.QueueSubscribe("recipes.delete.*", workerQueue, h.Delete)
	if err != nil {
		return fmt.Errorf("recipe.delete subscription: %w", err)
	}

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
}

const handlerTimeout = time.Second

func (h *handlers) List(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	recipes, err := h.entc.Recipe.Query().
		Order(ent.Asc(recipe.FieldTitle)).
		All(ctx)

	if err != nil {
		fmt.Println("recipe list", err)
		h.ec.Publish(msg.Reply, List{Status: StatusError})
		return
	}

	list := List{
		Status: StatusSuccess,
	}
	for _, r := range recipes {
		list.Recipes = append(list.Recipes, EntToRecipe(r))
	}

	h.ec.Publish(msg.Reply, list)
}

func (h *handlers) Delete(msg *nats.Msg) {
	id := uuid.MustParse(strings.Split(msg.Subject, ".")[2])

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	err := h.entc.Recipe.DeleteOneID(id).Exec(ctx)
	if err != nil {
		fmt.Println("recipe delete", err)
		h.ec.Publish(msg.Reply, Ack{Status: StatusError})
		return
	}

	h.ec.Publish(msg.Reply, Ack{Status: StatusSuccess})
}

func (h *handlers) Detail(msg *nats.Msg) {
	id := uuid.MustParse(strings.Split(msg.Subject, ".")[2])
	fmt.Println("get recipe with id", id)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
	defer cancel()

	eRecipe, err := h.entc.Recipe.Get(ctx, id)
	if err != nil {
		fmt.Println("get recipe", err)
		// todo: comunicate error
	}

	h.ec.Publish(msg.Reply, EntToRecipe(eRecipe))
}

func (h *handlers) Upsert(subject, reply string, r Recipe) {
	fmt.Printf("create recipe %+v\n", r)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	recipe, err := h.entc.Recipe.Get(ctx, r.ID)
	if err != nil && !ent.IsNotFound(err) {
		fmt.Println("recipe save error:", err)
		h.ec.Publish(reply, Ack{Status: StatusError})
		return
	}

	if ent.IsNotFound(err) {
		_, err = h.entc.Recipe.Create().
			SetTitle(r.Name).
			SetID(r.ID).
			Save(ctx)
	} else {
		_, err = recipe.Update().
			SetTitle(r.Name).
			Save(ctx)
	}

	if err != nil {
		fmt.Println("recipe save error:", err)
		h.ec.Publish(reply, Ack{Status: StatusError})
		return
	}

	h.ec.Publish(reply, Ack{Status: StatusSuccess})
}
