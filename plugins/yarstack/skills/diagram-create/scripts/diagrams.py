# /// script
# requires-python = ">=3.11"
# dependencies = ["playwright==1.61.0"]
# ///
"""The diagrams in this folder.

Build (check, write the HTML pages, render the PNGs):  uv run diagrams.py
Download the icons that the diagrams use:              uv run diagrams.py icons
"""

from pathlib import Path

from diagram import Diagram, Page, run

FOLDER = Path(__file__).resolve().parent


def web_shop_architecture() -> Diagram:
    d = Diagram(
        "example-architecture",
        title="Web shop architecture",
        description=(
            "Architecture of a web shop: customers reach the API through ingress-nginx in the Kubernetes cluster; "
            "the API stores orders in PostgreSQL, charges a payment provider, and queues jobs for a worker that "
            "writes product images to object storage."
        ),
    )
    d.zone(356, 160, 568, 304, label="CLUSTER · KUBERNETES")
    customers = d.node(60, 200, 200, kind="input", tag="USER", icon="tabler:user", name="Customers", sublabel="browser")
    ingress = d.node(
        380, 200, 200, kind="backend", tag="LB", icon="simple-icons:nginx", name="ingress-nginx", sublabel="TLS · routing"
    )
    api = d.node(700, 200, 200, kind="focal", tag="GO", icon="simple-icons:go", name="api", sublabel="REST · :8080")
    worker = d.node(
        700, 360, 200, kind="backend", tag="JOB", icon="tabler:clock", name="worker", sublabel="image resize · email"
    )
    payments = d.node(
        1020, 40, 200, kind="external", tag="EXT", icon="tabler:building-bank", name="Payment provider", sublabel="HTTPS API"
    )
    database = d.node(
        1020, 200, 200, kind="store", tag="PG", icon="simple-icons:postgresql", name="PostgreSQL",
        sublabel="orders · customers",
    )
    storage = d.node(
        1020, 360, 200, kind="store", tag="S3", icon="tabler:bucket", name="Object storage", sublabel="product images"
    )

    d.edge(customers.attach_point("right"), ingress.attach_point("left"), color="link", label="HTTPS")
    d.edge(ingress.attach_point("right"), api.attach_point("left"), label="ROUTE")
    d.edge(api.attach_point("right"), database.attach_point("left"), label="SQL")
    d.edge(api.attach_point("top"), (800, 80), payments.attach_point("left"), color="link", label="CHARGE")
    d.edge(api.attach_point("bottom"), worker.attach_point("top"), dashed=True, label="JOBS")
    d.edge(worker.attach_point("right"), storage.attach_point("left"), label="UPLOADS")
    return d


def web_shop_sign_in() -> Diagram:
    d = Diagram(
        "example-sign-in",
        title="Web shop sign-in",
        description=(
            "Sequence of one sign-in: the browser posts credentials to the API, the API verifies them with the "
            "identity provider, and then either creates a session in PostgreSQL and sets a cookie or returns a "
            "login error."
        ),
    )
    actors = [
        d.node(60, 40, 200, kind="input", tag="USER", icon="tabler:user", name="Browser", sublabel="customer"),
        d.node(380, 40, 200, kind="backend", tag="GO", icon="simple-icons:go", name="api", sublabel="REST · :8080"),
        d.node(
            700, 40, 200, kind="external", tag="IDP", icon="tabler:shield-lock", name="Identity provider", sublabel="OIDC"
        ),
        d.node(
            1020, 40, 200, kind="store", tag="PG", icon="simple-icons:postgresql", name="PostgreSQL", sublabel="sessions"
        ),
    ]
    for actor in actors:
        d.lifeline(actor, bottom=552)
    browser, api, identity, database = (actor.center_x for actor in actors)

    d.fragment(96, 236, 1088, 292, operator="ALT", dividers=(404,))
    d.guard("valid credentials", left=172, top=256)
    d.guard("else", left=172, top=420)

    d.edge((browser, 168), (api, 168), color="link", label="POST /LOGIN")
    d.edge((api, 208), (identity, 208), color="link", label="VERIFY")
    d.edge((identity, 296), (api, 296), color="accent", label="200 · USER ID")
    d.edge((api, 336), (database, 336), label="INSERT SESSION", label_at=(960, 336))
    d.edge((api, 376), (browser, 376), dashed=True, label="SET-COOKIE")
    d.edge((identity, 460), (api, 460), dashed=True, label="401")
    d.edge((api, 500), (browser, 500), dashed=True, label="LOGIN ERROR")
    return d


PAGES = [Page("example.html", "Example diagrams", [web_shop_architecture(), web_shop_sign_in()])]

if __name__ == "__main__":
    raise SystemExit(run(FOLDER, PAGES))
