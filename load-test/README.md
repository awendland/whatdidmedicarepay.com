# Load Test

The [k6](https://k6.io/docs/) framework is being used for load testing.

Each `*.js` file represents a script that Virtual Users will execute during the load test. You can run a 10 second load test with 5 users using the command `k6 run -u 5 -d 10s SOME_SCRIPT.js`.

The scripts should default to using `http://localhost:3000` unless the `HOSTNAME` environment variable is set.
