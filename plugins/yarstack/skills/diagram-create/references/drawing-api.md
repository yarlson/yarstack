# Drawing API

`diagram` exports `Diagram`, `Box`, `Lane`, `Chip`, `Page`, and
`run`. All coordinates are absolute SVG pixels on the 4 px grid. The canvas
is 1280 px wide, with a 40 px margin.

```python
from pathlib import Path

from diagram import Chip, Diagram, Lane, Page, run

def overview() -> Diagram:
    d = Diagram("shop-overview", title="Web shop overview", description="One sentence about the content.")
    api = d.node(380, 200, 200, kind="focal", tag="GO", icon="simple-icons:go", name="api", sublabel="REST · :8080")
    db = d.node(700, 200, 200, kind="store", tag="PG", icon="simple-icons:postgresql", name="PostgreSQL", sublabel="orders")
    d.edge(api.attach_point("right"), db.attach_point("left"), label="SQL")
    return d

PAGES = [Page("overview.html", "Web shop overview", [overview()])]

if __name__ == "__main__":
    raise SystemExit(run(Path(__file__).resolve().parent, PAGES))
```

The paint order is fixed, whatever the call order: lane frames, zones,
fragments, lifelines, tabs, then edges, then labels, then nodes.

## Diagram

`Diagram(slug, *, title, description)`. The slug is lowercase words joined by
hyphens. It names the PNG and prefixes every ID, so it must be unique in the
folder. The title has 60 characters or fewer.

| Method                                                                                                        | Draws                                                                | Notes                                                                                                                                      |
| ------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `node(x, y, w, *, kind, tag, icon, name, sublabel)` → `Box`                                                   | Standard node, 80 px tall                                            | Tag top left, 28 px icon centered above the name and the mono sublabel.                                                                    |
| `list_node(x, y, w, *, kind, tag, icon, name, chips=(), lines=(), h=None)` → `Box`                            | Tall node with the icon beside the name, then chips, then mono lines | Without `h`, the height fits the content. With `h`, it must be at least that height. Use the same `h` to align list nodes across diagrams. |
| `zone(x, y, w, h, *, label)` → `Box`                                                                          | Dashed boundary with an uppercase eyebrow at the top left            | Leave 36 px or more from the top to the first node, and 20 px or more on the other sides.                                                  |
| `tab(x, y, w, *, label)` → `Box`                                                                              | Small solid box, 24 px tall, on the edge of a zone                   | An edge that applies to every node in the zone ends on the tab.                                                                            |
| `lanes(x, w, *, label_width, lanes)`                                                                          | Swimlane frame                                                       | `lanes` is a list of `Lane(top, bottom, name, detail)`. Each lane starts where the one before it ends.                                     |
| `fragment(x, y, w, h, *, operator, dividers=())` → `Box`                                                      | Sequence frame: `ALT`, `OPT`, or `LOOP`                              | A dashed line at each divider y.                                                                                                           |
| `lifeline(actor, *, bottom)`                                                                                  | Dashed lifeline from the bottom center of an actor box               |                                                                                                                                            |
| `guard(condition, *, left, top)`                                                                              | `[condition]` on a mask, for a fragment region                       |                                                                                                                                            |
| `edge(*points, color="muted", dashed=False, bidirectional=False, label=None, label_at=None, label_side=None)` | Connector with an arrowhead                                          | See "Edges".                                                                                                                               |

`Chip(icon, name, version)` is an artifact inside a list node: the name on
the left, the version on the right. Chips start 64 px below the top of the
node, every 32 px. The lines follow 20 px after the last chip, every 18 px.

`Box` has `left`, `right`, `top`, `bottom`, `center_x`, `center_y`, and
`attach_point(side, index=1, count=1)`. That method returns point `index` of
`count` points spread evenly along `side`, snapped to the grid. Use it to fan
several edges out of one side: `api.attach_point("bottom", 2, 3)`.

## Node kinds

| Kind       | Use for                                                                |
| ---------- | ---------------------------------------------------------------------- |
| `backend`  | What you run: services, jobs, gateways, load balancers, clusters       |
| `focal`    | The one or two elements the diagram is about. The icon turns orange.   |
| `store`    | Databases, buckets, queues, registries, secret stores                  |
| `external` | Systems you do not own: partners, SaaS, providers, the public internet |
| `input`    | People and clients: users, engineers, browsers                         |
| `optional` | Optional, legacy, or minor workloads. The border is dashed.            |

## Edges

- **Waypoints.** Pass the start, each turn, and the end. Each pair shares x
  or y. The package draws 8 px rounded corners and drops waypoints that do
  not turn. Start and end on a box side (use `attach_point`), on a
  lifeline, or on a frame.
- **Color.** `muted` by default, `link` for HTTP and external calls, `accent`
  for the headline path. `dashed=True` for optional, async, return, or
  passive flows. `bidirectional=True` puts an arrowhead at both ends, for
  request and reply over one channel.
- **Hops.** Where two edges cross, the package draws an 8 px hop on the less
  important one. You do not draw hops.
- **Label.** The label goes on the longest segment: 8 px above a horizontal
  segment, or 8 px to the right of a vertical one, at its middle. It takes
  the color of the edge. To move it, pass `label_at`, a point on the edge
  where the label centers. To flip it, pass `label_side`: `below` for a
  horizontal segment, `left` for a vertical one.

## Pages and the command line

`Page(file_name, title, diagrams)` is one HTML page with one or more
diagrams. Each diagram gets its own PNG, named after its slug. Put the
diagrams of one environment set or one topic on one page.

`run(folder, pages)` is the command line. It takes `build` (the default) or
`icons`.

## Icons

Every node and chip has an icon. The icon is ink, or orange on a focal node.

| Reference             | File in `D/icons/`        | Source                                                                                                                                             |
| --------------------- | ------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `tabler:<name>`       | `tabler/<name>.svg`       | [Tabler Icons](https://tabler.io/icons) outline set, MIT. Generic concepts: user, world, key, server, database, clock, bucket.                     |
| `simple-icons:<slug>` | `simple-icons/<slug>.svg` | [Simple Icons](https://simpleicons.org), CC0. Products and tools: github, kubernetes, postgresql, nginx, grafana.                                  |
| `<folder>:<name>`     | `<folder>/<name>.png`     | A vendor icon pack, for example cloud product icons. Add the file by hand, only when the vendor terms allow it, and name the source in the README. |

- `uv run D/diagrams.py icons` downloads the missing Tabler and Simple Icons
  files at pinned versions, and the Tabler `LICENSE`. The pinned versions
  are in `diagram/icons.py`.
- A "no icon named" error means that the name is wrong. Look it up on the
  icon site. Do not guess. The Simple Icons slug is the file name on the
  site, for example `nodedotjs` for Node.js. If a brand has no mark, use a
  generic Tabler icon.
- The package keeps only plain shapes from each SVG file. It rejects a file
  with any other element, such as a script.
- Delete icon files that the build reports as unused.
