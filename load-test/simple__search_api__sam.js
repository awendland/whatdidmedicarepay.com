// Run w/ k6, e.g. `k6 run THIS_FILE_NAME.js`
import http from "k6/http";

const HOSTNAME = __ENV.HOSTNAME || "http://localhost:3000";

export default function () {
  const r = http.get(`${HOSTNAME}/api/search?query=name:sam`);
  check(r, {
    "status is 200": (r) => r.status === 200,
  });
}
