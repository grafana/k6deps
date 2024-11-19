"use k6 = 0.52";
"use k6 with k6/x/faker >= 0.3.0";
"use k6 with k6/x/sql >= 1.0.0";

import faker from "./faker.js";
import sql from "./sql.js";

export { setup, teardown } from "./sql.js";

export default () => {
  faker();
  sql();
};
