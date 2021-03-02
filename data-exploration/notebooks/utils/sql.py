import re
import statistics


def sql_regexp(x, y):
    """
    Register with sqlite using `con.create_function("regexp", 2, sql_regexp)`.

    You can then sanitize data using SQL such as
    ```sql
    SELECT zipcode FROM locations WHERE zipcode REGEXP "[0-9]{5}"
    ```

    Note that using this will force a table scan, even if an index exists. The
    LIKE operation can leverage indices.
    """
    if x is None or y is None:
        return 0
    else:
        return 1 if re.match(x, y) else 0


class DistributionSummary:
    """
    Register with sqlite using `con.create_aggregate("dist_summary", 3, DistributionSummary)`.

    You can then run SQL such as:
    ```sql
    SELECT state, dist_summary(population, 4, 'all') FROM cities GROUP BY state
    ```

    to get a quartile breakdown of the city populations for each state.
    """

    def __init__(self):
        self.values = []
        self.n = 4
        self.i = "all"

    # i can be a specific quantile to report, or can be "all" to indicate that all should be returned
    def step(self, value, n, i):
        try:
            self.values.append(float(value))
        except Exception:
            pass
        self.n = n
        self.i = i

    def finalize(self):
        self.values.sort()
        if len(self.values) == 0:
            return None if self.i == "all" else None
        if len(self.values) == 1:
            return self.values[0] if self.i == "all" else None
        if len(self.values) == 2:
            return ", ".join([str(v) for v in self.values]) if self.i == "all" else None
        quantiles = (
            [str(round(self.values[0], 1))]
            + [str(round(q, 1)) for q in statistics.quantiles(self.values, n=self.n)]
            + [str(round(self.values[-1], 1))]
        )
        return ", ".join(quantiles) if self.i == "all" else float(quantiles[self.i])


class Last:
    """
    Register with sqlite using `con.create_aggregate("last", 1, Last)`.

    You can then run SQL such as:
    ```sql
    SELECT zipcode, last(state) FROM locations GROUP BY zipcode
    ```
    """

    def __init__(self):
        self.value = None

    def step(self, value):
        self.value = value

    def finalize(self):
        return self.value
