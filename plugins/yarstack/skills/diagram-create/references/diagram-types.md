# Diagram types

All four types use the style rules in SKILL.md §5 and the API in
[drawing-api.md](drawing-api.md). `scripts/diagrams.py` holds one complete
architecture diagram and one complete sequence diagram.

## Architecture

Use it for components and their connections in one snapshot: a system
overview, a runtime data flow, an integration map, an infrastructure
overview, or which hosts sit behind a login.

- Keep one main direction: left to right, or top to bottom.
- Group nodes by tier or trust boundary with zones: a cluster, a namespace, a
  network, or a provider. Use at most 3 zones. Each zone is a real boundary,
  not a layout aid.
- The focal node is the main integration point, the main store, or the key
  decision.
- Draw request and reply over one channel, for example a message bus, as one
  edge with `bidirectional=True` and one label, not as two edges.
- To show access across a trust boundary, label each allowed path. Draw a
  blocked path dashed, and stop it at the boundary.
- Do not draw a two-way arrow when one direction is clear from the context.

## Deployment

Use it for where software runs: environments, networks, clusters, hosts,
managed services, node pools, and versions. If the diagram shows no
placement decision (which host, which network, how many replicas, which
version), use Architecture.

- Zones are environment or network boundaries: a region, a private network,
  a cluster, or "outside the network". Use at most 3.
- Draw load balancers, gateways, databases, buckets, and key stores with
  `node()`. Draw a cluster or a VM with `list_node()`, and put what runs on
  it in chips.
- Each chip has a version: an image tag, a chart version, `@ SHA`, or a size.
  Put node pools and machine sizes in the `lines` of the list node.
- Label each network path with the protocol or the port: `HTTPS`, `SQL`,
  `S3 API`, `:5432`. Use `link` for a path to or from outside the zone, and
  `muted` inside it.
- Draw one environment in each diagram. For several environments, put one
  diagram for each on one page, with the same layout and the same list node
  `h`, so that readers can compare them.
- Put a replica count in the sublabel or the chip version. Do not draw one
  node for each replica.

## Swimlane

Use it for a process with handoffs between owners: CI/CD, release,
approval, or incident response.

- Draw one lane for each owner, a system or a team, with `lanes()`. The name
  and a detail, such as where it runs, sit in the left column. A
  `label_width` of 104 fits names of about 12 characters. Use at most 5
  lanes.
- Put each step in the lane of its owner. Never draw a step across two lanes.
- Keep the flow mostly left to right. Reorder the steps so that edges do not
  snake back and forth.
- Handoffs, the edges that cross lanes, are the important ones. Give the
  accent to the handoff that matters most.
- Give each lane an 80 px node plus 28 px or more above and below for labels.

```python
d.lanes(40, 1200, label_width=104, lanes=[
    Lane(40, 184, "SOURCE", "GIT HOSTING"),
    Lane(184, 392, "CI", "BUILD RUNNERS"),
    Lane(392, 528, "REGISTRY", "IMAGES"),
    Lane(528, 680, "CLUSTER", "EACH ENVIRONMENT"),
])
```

## Sequence

Use it for time-ordered requests and responses between 2 to 5 actors:
sign-in, token refresh, a webhook, or an API call trace.

- Put the actors in one row at the top with `node()`. Space their centers
  evenly, for example at x = 160, 480, 800, and 1120 for four actors. Give
  each a `lifeline()`.
- Draw each message as a horizontal edge between two lifelines. Time flows
  down. Keep messages 40 px apart (24 px minimum).
- A call is solid: `link` for HTTP, `muted` for internal calls. A return is
  dashed `muted`. The headline success is one solid `accent` message.
- A self-call is a short loop to the right of the lifeline:
  `d.edge((x, y), (x + 32, y), (x + 32, y + 32), (x, y + 32), label=..., label_at=(x + 32, y + 16))`.
- The label of a message centers between its two lifelines. If it would sit
  on a third lifeline, the build reports it. Move it with `label_at`.
- Draw a branch with one `fragment()`: `ALT` with at most 2 regions, or
  `OPT`, or `LOOP`. Keep the frame 12 px or more outside the outer lifelines
  that take part.
- Put each guard with `guard()` at the left of its region, 12 px right of the
  first lifeline and 20 px below the top of the region. Put the first message
  of a region 24 px or more below its guard. Keep a divider 16 px or more
  from the messages above and below it.
- Give the accent to one branch only.
- If the flow does not fit the budget, draw the happy path in one diagram and
  the failure path in a second.
