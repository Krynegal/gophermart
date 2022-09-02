package postgres

const (
	UsersTable = `
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE, 
	password TEXT NOT NULL
);
`

	UserTokens = `
CREATE TABLE IF NOT EXISTS tokens
(
	id serial PRIMARY KEY,
	user_id serial REFERENCES users(id) UNIQUE,
	token text,
	created_at timestamp without time zone default now()
);
`

	OrdersTable = `
CREATE TABLE IF NOT EXISTS orders
(
    id BIGSERIAL PRIMARY KEY,
	order_num BIGINT UNIQUE,
	user_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id)
);
`

	AccrualsTable = `
CREATE TABLE IF NOT EXISTS orders
(
	order_num BIGINT PRIMARY KEY,
	user_id INT NOT NULL,
	status TEXT NOT NULL DEFAULT 'NEW',
	amount REAL DEFAULT 0,
	uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (order_num) REFERENCES orders (order_num)
);
`

	WithdrawalsTable = `
CREATE TABLE IF NOT EXISTS withdrawals
(
    order_num BIGINT PRIMARY KEY,
	user_id INT NOT NULL,
	amount REAL DEFAULT 0,
	processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (order_num) REFERENCES orders (order_num)
);
`
)
