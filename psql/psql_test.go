package psql

import (
	"os"
	"strings"
	"testing"

	"github.com/dosco/super-graph/qcode"
)

const (
	errNotExpected = "Generated SQL did not match what was expected"
)

var (
	qcompile *qcode.Compiler
	pcompile *Compiler
)

func TestMain(m *testing.M) {
	fm := qcode.NewFilterMap(map[string]string{
		"users": "{ id: { _eq: $user_id } }",
		"posts": "{ account_id: { _eq: $account_id } }",
	})

	bl := qcode.NewBlacklist([]string{
		"secret",
		"password",
		"token",
	})

	qcompile = qcode.NewCompiler(fm, bl)

	tables := []*DBTable{
		&DBTable{Name: "customers", Type: "table"},
		&DBTable{Name: "users", Type: "table"},
		&DBTable{Name: "products", Type: "table"},
		&DBTable{Name: "purchases", Type: "table"},
	}

	columns := [][]*DBColumn{
		[]*DBColumn{
			&DBColumn{ID: 1, Name: "id", Type: "bigint", NotNull: true, PrimaryKey: true, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 2, Name: "full_name", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 3, Name: "phone", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 4, Name: "email", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 5, Name: "encrypted_password", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 6, Name: "reset_password_token", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 7, Name: "reset_password_sent_at", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 8, Name: "remember_created_at", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 9, Name: "created_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 10, Name: "updated_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)}},
		[]*DBColumn{
			&DBColumn{ID: 1, Name: "id", Type: "bigint", NotNull: true, PrimaryKey: true, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 2, Name: "full_name", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 3, Name: "phone", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 4, Name: "avatar", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 5, Name: "email", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 6, Name: "encrypted_password", Type: "character varying", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 7, Name: "reset_password_token", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 8, Name: "reset_password_sent_at", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 9, Name: "remember_created_at", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 10, Name: "created_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 11, Name: "updated_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)}},
		[]*DBColumn{
			&DBColumn{ID: 1, Name: "id", Type: "bigint", NotNull: true, PrimaryKey: true, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 2, Name: "name", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 3, Name: "description", Type: "text", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 4, Name: "price", Type: "numeric(7,2)", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 5, Name: "user_id", Type: "bigint", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "users", FKeyColID: []int{1}},
			&DBColumn{ID: 6, Name: "created_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 7, Name: "updated_at", Type: "timestamp without time zone", NotNull: true, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)}},
		[]*DBColumn{
			&DBColumn{ID: 1, Name: "id", Type: "bigint", NotNull: true, PrimaryKey: true, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 2, Name: "customer_id", Type: "bigint", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "customers", FKeyColID: []int{1}},
			&DBColumn{ID: 3, Name: "product_id", Type: "bigint", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "products", FKeyColID: []int{1}},
			&DBColumn{ID: 4, Name: "sale_type", Type: "character varying", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 5, Name: "quantity", Type: "integer", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 6, Name: "due_date", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)},
			&DBColumn{ID: 7, Name: "returned", Type: "timestamp without time zone", NotNull: false, PrimaryKey: false, Uniquekey: false, FKeyTable: "", FKeyColID: []int(nil)}},
	}

	schema := initSchema()

	for i, t := range tables {
		updateSchema(schema, t, columns[i])
	}

	vars := NewVariables(map[string]string{
		"account_id": "select account_id from users where id = $user_id",
	})

	pcompile = NewCompiler(schema, vars)

	os.Exit(m.Run())
}

func compileGQLToPSQL(t *testing.T, gql string) string {
	qc, err := qcompile.CompileQuery(gql)
	if err != nil {
		t.Fatal(err)
	}

	var sqlStmt strings.Builder

	if err := pcompile.Compile(&sqlStmt, qc); err != nil {
		t.Fatal(err)
	}

	return sqlStmt.String()
}

func TestCompileGQLWithArgs(t *testing.T) {
	gql := `query {
		products(
			# returns only 30 items
			limit: 30,
	
			# starts from item 10, commented out for now
			# offset: 10,
	
			# orders the response items by highest price
			order_by: { price: desc },
	
			# no duplicate prices returned
			distinct: [ price ]
			
			# only items with an id >= 30 and < 30 are returned
			where: { id: { and: { greater_or_equals: 20, lt: 28 } } }) {
			id
			name
			price
		}
	}`

	sql := `SELECT json_object_agg('products', products) FROM (SELECT coalesce(json_agg("products" ORDER BY "products_0.ob.price" DESC), '[]') AS "products" FROM (SELECT  DISTINCT ON ("price") row_to_json((SELECT "sel_0" FROM (SELECT "products_0"."id" AS "id", "products_0"."name" AS "name", "products_0"."price" AS "price") AS "sel_0")) AS "products", "products_0"."price" AS "products_0.ob.price" FROM (SELECT "products"."id", "products"."name", "products"."price" FROM "products" WHERE ((("products"."id") < (28)) AND (("products"."id") >= (20)) )  LIMIT ('30') :: integer) AS "products_0" ORDER BY "products_0.ob.price" DESC LIMIT ('30') :: integer) AS "products_0") AS "done_1337";`

	resSQL := compileGQLToPSQL(t, gql)

	if resSQL != sql {
		t.Fatal(errNotExpected)
	}
}

func TestCompileGQLOneToMany(t *testing.T) {
	gql := `query {
		users {
			email
			products {
				name
				price
			}
		}
	}`

	sql := `SELECT json_object_agg('users', users) FROM (SELECT coalesce(json_agg("users"), '[]') AS "users" FROM (SELECT row_to_json((SELECT "sel_0" FROM (SELECT "users_0"."email" AS "email", "products_1.join"."products" AS "products") AS "sel_0")) AS "users" FROM (SELECT "users"."email", "users"."id" FROM "users" WHERE ((("users"."id") = ('{{user_id}}')) )  LIMIT ('20') :: integer) AS "users_0" LEFT OUTER JOIN LATERAL (SELECT coalesce(json_agg("products"), '[]') AS "products" FROM (SELECT row_to_json((SELECT "sel_1" FROM (SELECT "products_1"."name" AS "name", "products_1"."price" AS "price") AS "sel_1")) AS "products" FROM (SELECT "products"."name", "products"."price" FROM "products" WHERE ( (("products"."user_id") = ("users_0"."id")) )  LIMIT ('20') :: integer) AS "products_1" LIMIT ('20') :: integer) AS "products_1") AS "products_1.join" ON ('true')  LIMIT ('20') :: integer) AS "users_0") AS "done_1337";`

	resSQL := compileGQLToPSQL(t, gql)

	if resSQL != sql {
		t.Fatal(errNotExpected)
	}
}

func TestCompileGQLBelongTo(t *testing.T) {
	gql := `query {
		products {
			name
			price
			users {
				email
			}
		}
	}`

	sql := `SELECT json_object_agg('products', products) FROM (SELECT coalesce(json_agg("products"), '[]') AS "products" FROM (SELECT row_to_json((SELECT "sel_0" FROM (SELECT "products_0"."name" AS "name", "products_0"."price" AS "price", "users_1.join"."users" AS "users") AS "sel_0")) AS "products" FROM (SELECT "products"."name", "products"."price", "products"."user_id" FROM "products" LIMIT ('20') :: integer) AS "products_0" LEFT OUTER JOIN LATERAL (SELECT coalesce(json_agg("users"), '[]') AS "users" FROM (SELECT row_to_json((SELECT "sel_1" FROM (SELECT "users_1"."email" AS "email") AS "sel_1")) AS "users" FROM (SELECT "users"."email" FROM "users" WHERE ( (("users"."id") = ("products_0"."user_id")) )  LIMIT ('20') :: integer) AS "users_1" LIMIT ('20') :: integer) AS "users_1") AS "users_1.join" ON ('true')  LIMIT ('20') :: integer) AS "products_0") AS "done_1337";`

	resSQL := compileGQLToPSQL(t, gql)

	if resSQL != sql {
		t.Fatal(errNotExpected)
	}
}

func TestCompileGQLManyToMany(t *testing.T) {
	gql := `query { 
		products {
			name
			customers {
				email
				full_name
			}
		}
	}`

	sql := `SELECT json_object_agg('products', products) FROM (SELECT coalesce(json_agg("products"), '[]') AS "products" FROM (SELECT row_to_json((SELECT "sel_0" FROM (SELECT "products_0"."name" AS "name", "customers_1.join"."customers" AS "customers") AS "sel_0")) AS "products" FROM (SELECT "products"."name", "products"."id" FROM "products" LIMIT ('20') :: integer) AS "products_0" LEFT OUTER JOIN LATERAL (SELECT coalesce(json_agg("customers"), '[]') AS "customers" FROM (SELECT row_to_json((SELECT "sel_1" FROM (SELECT "customers_1"."email" AS "email", "customers_1"."full_name" AS "full_name") AS "sel_1")) AS "customers" FROM (SELECT "customers"."email", "customers"."full_name" FROM "customers" LEFT OUTER JOIN "purchases" ON (("purchases"."product_id") = ("products_0"."id")) WHERE ( (("customers"."id") = ("purchases"."customer_id")) )  LIMIT ('20') :: integer) AS "customers_1" LIMIT ('20') :: integer) AS "customers_1") AS "customers_1.join" ON ('true')  LIMIT ('20') :: integer) AS "products_0") AS "done_1337";`

	resSQL := compileGQLToPSQL(t, gql)

	if resSQL != sql {
		t.Fatal(errNotExpected)
	}
}

func TestCompileGQLManyToManyReverse(t *testing.T) {
	gql := `query {
		customers {
			email
			full_name
			products {
				name
			}
		}
	}`

	sql := `SELECT json_object_agg('customers', customers) FROM (SELECT coalesce(json_agg("customers"), '[]') AS "customers" FROM (SELECT row_to_json((SELECT "sel_0" FROM (SELECT "customers_0"."email" AS "email", "customers_0"."full_name" AS "full_name", "products_1.join"."products" AS "products") AS "sel_0")) AS "customers" FROM (SELECT "customers"."email", "customers"."full_name", "customers"."id" FROM "customers" LIMIT ('20') :: integer) AS "customers_0" LEFT OUTER JOIN LATERAL (SELECT coalesce(json_agg("products"), '[]') AS "products" FROM (SELECT row_to_json((SELECT "sel_1" FROM (SELECT "products_1"."name" AS "name") AS "sel_1")) AS "products" FROM (SELECT "products"."name" FROM "products" LEFT OUTER JOIN "purchases" ON (("purchases"."customer_id") = ("customers_0"."id")) WHERE ( (("products"."id") = ("purchases"."product_id")) )  LIMIT ('20') :: integer) AS "products_1" LIMIT ('20') :: integer) AS "products_1") AS "products_1.join" ON ('true')  LIMIT ('20') :: integer) AS "customers_0") AS "done_1337";`

	resSQL := compileGQLToPSQL(t, gql)

	if resSQL != sql {
		t.Fatal(errNotExpected)
	}
}
