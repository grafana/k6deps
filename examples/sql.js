// source: https://github.com/grafana/xk6-sql/blob/v0.4.0/examples/sqlite3_test.js
import sql from "k6/x/sql";
import driver from "k6/x/sql/driver/ramsql";

const db = sql.open(driver);

export function setup() {
  db.exec(`CREATE TABLE IF NOT EXISTS namevalues (
           id integer PRIMARY KEY AUTOINCREMENT,
           name varchar NOT NULL,
           value varchar);`);
}

export function teardown() {
  db.close();
}

export default function () {
  db.exec("INSERT INTO namevalues (name, value) VALUES('plugin-name', 'k6-plugin-sql');");

  let results = db.query("SELECT * FROM namevalues WHERE name = $1;", "plugin-name");
  for (const row of results) {
    console.log(`name: ${row.name}, value: ${row.value}`);
  }
}
