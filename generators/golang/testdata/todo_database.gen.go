const itemSelectQuery = `SELECT "id", "text", "created_at" FROM "item" `

func ItemSelectFirst(ctx context.Context, q QueryerContext, extraQuery string, args ...any) (todo.Item, error) {
	row := q.QueryRowContext(ctx, itemSelectQuery+extraQuery, args...)
	var o todo.Item
	err := row.Scan(
		&o.ID,
		&o.Text,
		&o.CreatedAt,
	)
	return o, err
}

func ItemSelect(ctx context.Context, q QueryerContext, extraQuery string, args ...any) ([]todo.Item, error) {
	rows, err := q.QueryContext(ctx, itemSelectQuery+extraQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var objs []todo.Item
	for rows.Next() {
		var obj todo.Item
		if err := rows.Scan(
			&obj.ID,
			&obj.Text,
			&obj.CreatedAt,
		); err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return objs, nil
}

func ItemSelectWithID(ctx context.Context, q QueryerContext, id string) (todo.Item, error) {
	return ItemSelectFirst(ctx, q, `WHERE "id" = $1`, id)
}

func ItemSelectAll(ctx context.Context, q QueryerContext) ([]todo.Item, error) {
	return ItemSelect(ctx, q, "")
}

const itemInsertQuery = `INSERT INTO "item" ("id", "text", "created_at") VALUES ($1, $2, $3);`

func ItemInsert(ctx context.Context, e ExecerContext, obj todo.Item) error {
	_, err := e.ExecContext(ctx, itemInsertQuery,
		obj.ID,
		obj.Text,
		obj.CreatedAt,
	)
	return err
}

const itemUpdateQuery = `UPDATE "item"
SET "text" = $2, "created_at" = $3
WHERE id = $1;`

func ItemUpdate(ctx context.Context, e ExecerContext, obj todo.Item) error {
	_, err := e.ExecContext(ctx, itemUpdateQuery,
		obj.ID,
		obj.Text,
		obj.CreatedAt,
	)
	return err
}

const itemUpsertQuery = `INSERT INTO "item" ("id", "text", "created_at") VALUES ($1, $2, $3)
ON CONFLICT ("id")
DO UPDATE SET "id" = $1, "text" = $2, "created_at" = $3;`

func ItemUpsert(ctx context.Context, e ExecerContext, obj todo.Item) error {
	_, err := e.ExecContext(ctx, itemUpsertQuery,
		obj.ID,
		obj.Text,
		obj.CreatedAt,
	)
	return err
}

const itemDeleteQuery = `DELETE FROM "item" WHERE "id" = $1;`

func ItemDelete(ctx context.Context, e ExecerContext, id string) error {
	_, err := e.ExecContext(ctx, itemDeleteQuery, id)
	return err
}

const itemCountQuery = "SELECT COUNT(*) FROM item "

func ItemCount(ctx context.Context, q QueryerContext, extraQuery string, args ...any) (int, error) {
	row := q.QueryRowContext(ctx, itemCountQuery+extraQuery, args...)
	var c int
	err := row.Scan(&c)
	return c, err
}
