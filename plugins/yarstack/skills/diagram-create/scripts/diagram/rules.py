"""Layout rules that a diagram must pass before it is written."""

from __future__ import annotations

from collections import defaultdict
from collections.abc import Iterator
from itertools import combinations

from . import style
from .model import Diagram, Edge
from .geometry import EPSILON, Point, Rect, Segment

MIN_PARALLEL_GAP = 12
MIN_ATTACH_GAP = 12
MIN_STROKE_CLEARANCE = 4
MIN_HOP_ROOM = 16
MIN_HOP_SPACING = 24
MAX_ACCENTS = 2


def check(diagram: Diagram) -> list[str]:
    """All rule findings for one diagram. An empty list means the diagram passes."""
    rules = (
        _off_grid,
        _outside_margin,
        _short_segments,
        _loose_ends,
        _edges_through_boxes,
        _crowded_edges,
        _crowded_attach_points,
        _cramped_hops,
        _covered_labels,
        _overlapping_boxes,
        _too_many_accents,
    )
    return [f"{diagram.slug}: {finding}" for rule in rules for finding in rule(diagram)]


def _point(point: Point) -> str:
    return f"({point[0]:g},{point[1]:g})"


def _boxes(diagram: Diagram) -> list[tuple[str, Rect]]:
    """Opaque boxes: nodes and tabs. Edges must not run through them, and labels must not touch them."""
    return [(node.describe(), node.box) for node in diagram.nodes] + [(tab.describe(), tab.box) for tab in diagram.tabs]


def _frames(diagram: Diagram) -> list[Rect]:
    """Outlines that an edge may end on: zones, lane frames, and fragments."""
    return (
        [zone.box for zone in diagram.zones]
        + [frame.frame for frame in diagram.lane_frames]
        + [fragment.box for fragment in diagram.fragments]
    )


def _tags(diagram: Diagram) -> list[tuple[str, Rect]]:
    """Masked text that zones and fragments paint first: edges must not cross it."""
    return [(f"the label of {zone.describe()}", zone.eyebrow) for zone in diagram.zones] + [
        (f"the operator tab of the {fragment.describe()}", fragment.operator_tab) for fragment in diagram.fragments
    ]


def _off_grid(diagram: Diagram) -> Iterator[str]:
    grid = style.GRID
    placed: list[tuple[str, dict[str, float]]] = []
    for node in diagram.nodes:
        placed.append((node.describe(), {"x": node.box.x, "y": node.box.y, "w": node.box.w, "h": node.box.h}))
    for zone in diagram.zones:
        placed.append((zone.describe(), {"x": zone.box.x, "y": zone.box.y, "w": zone.box.w, "h": zone.box.h}))
    for tab in diagram.tabs:
        placed.append((tab.describe(), {"x": tab.box.x, "y": tab.box.y, "w": tab.box.w}))
    for frame in diagram.lane_frames:
        values = {"x": frame.x, "w": frame.w, "label_width": frame.label_width}
        values.update({f"lane {lane.name!r} top": lane.top for lane in frame.lanes})
        values["bottom"] = frame.lanes[-1].bottom
        placed.append(("lanes", values))
    for fragment in diagram.fragments:
        box = fragment.box
        values = {"x": box.x, "y": box.y, "w": box.w, "h": box.h}
        values.update({f"divider {index + 1}": y for index, y in enumerate(fragment.dividers)})
        placed.append((fragment.describe(), values))
    for lifeline in diagram.lifelines:
        placed.append((lifeline.describe(), {"x": lifeline.x, "top": lifeline.top, "bottom": lifeline.bottom}))
    for edge in diagram.edges:
        values = {}
        for x, y in edge.points:
            values[f"x of point {_point((x, y))}"] = x
            values[f"y of point {_point((x, y))}"] = y
        placed.append((edge.describe(), values))
    for guard in diagram.guards:
        placed.append((f"guard {guard.text!r}", {"left": guard.rect.x, "top": guard.rect.y}))
    for what, values in placed:
        for name, value in values.items():
            if value % grid:
                yield f"{what}: {name} = {value:g} is off the {grid} px grid"


def _outside_margin(diagram: Diagram) -> Iterator[str]:
    low, high = style.MARGIN, style.CANVAS_WIDTH - style.MARGIN
    drawn = _boxes(diagram) + [(zone.describe(), zone.box) for zone in diagram.zones]
    drawn += [("lanes", frame.frame) for frame in diagram.lane_frames]
    drawn += [(fragment.describe(), fragment.box) for fragment in diagram.fragments]
    drawn += [(f"label {label.text!r}", label.rect) for label in diagram.labels]
    drawn += [(edge.describe(), Rect(x, y, 0, 0)) for edge in diagram.edges for x, y in edge.points]
    for what, rect in drawn:
        if rect.left < low or rect.top < low or rect.right > high:
            yield f"{what} {rect} is outside the {style.MARGIN} px margin of the {style.CANVAS_WIDTH} px canvas"


def _short_segments(diagram: Diagram) -> Iterator[str]:
    for edge in diagram.edges:
        segments = edge.segments
        last = len(segments) - 1
        for index, segment in enumerate(segments):
            needed = style.CORNER_RADIUS * ((index > 0) + (index < last))
            needed += style.ARROW_LENGTH * ((index == last) + (index == 0 and edge.bidirectional))
            if segment.length < needed:
                yield (
                    f"{edge.describe()}: segment {_point(segment.a)}->{_point(segment.b)} is {segment.length:g} px; "
                    f"its corners and arrowheads need {needed} px"
                )


def _loose_ends(diagram: Diagram) -> Iterator[str]:
    outlines = [box for _, box in _boxes(diagram)] + _frames(diagram)
    lifelines = [lifeline.segment for lifeline in diagram.lifelines]
    for edge in diagram.edges:
        for end in (edge.points[0], edge.points[-1]):
            on_outline = any(rect.border_side(end) for rect in outlines)
            if not on_outline and not any(lifeline.holds(end) for lifeline in lifelines):
                yield f"{edge.describe()}: end {_point(end)} touches no node, tab, lifeline, or frame"


def _edges_through_boxes(diagram: Diagram) -> Iterator[str]:
    obstacles = _boxes(diagram) + _tags(diagram)
    for edge in diagram.edges:
        for what, rect in obstacles:
            if any(rect.crossed_by(segment) for segment in edge.segments):
                yield f"{edge.describe()} runs through {what}; route it around"


def _crowded_edges(diagram: Diagram) -> Iterator[str]:
    for first, second in combinations(diagram.edges, 2):
        for s in first.segments:
            for t in second.segments:
                if s.horizontal != t.horizontal:
                    continue
                gap = abs(s.position - t.position)
                shared = min(s.high, t.high) - max(s.low, t.low)
                if gap < MIN_PARALLEL_GAP and shared > EPSILON:
                    yield (
                        f"{first.describe()} and {second.describe()} run {gap:g} px apart for {shared:g} px; "
                        f"keep parallel runs {MIN_PARALLEL_GAP} px or more apart"
                    )


def _crowded_attach_points(diagram: Diagram) -> Iterator[str]:
    for what, box in _boxes(diagram):
        ends: dict[str, list[tuple[float, Edge]]] = defaultdict(list)
        for edge in diagram.edges:
            for end in (edge.points[0], edge.points[-1]):
                side = box.border_side(end)
                if side:
                    ends[side].append((end[1] if side in ("left", "right") else end[0], edge))
        for side, attached in ends.items():
            for (a, first), (b, second) in combinations(attached, 2):
                if abs(a - b) < MIN_ATTACH_GAP:
                    yield (
                        f"{first.describe()} and {second.describe()} attach {abs(a - b):g} px apart on the {side} "
                        f"side of {what}; spread them {MIN_ATTACH_GAP} px or more apart"
                    )


def _cramped_hops(diagram: Diagram) -> Iterator[str]:
    edges, hops = diagram.edges, diagram.hops()
    by_segment: dict[tuple[int, int], list[float]] = defaultdict(list)
    for hop in hops:
        segment = edges[hop.edge].segments[hop.segment]
        room = min(hop.at - segment.low, segment.high - hop.at)
        if room < MIN_HOP_ROOM:
            yield (
                f"{edges[hop.edge].describe()} crosses {edges[hop.over].describe()} {room:g} px from a corner or an end; "
                f"the hop needs {MIN_HOP_ROOM} px"
            )
        by_segment[(hop.edge, hop.segment)].append(hop.at)
    for (edge, _), positions in by_segment.items():
        positions.sort()
        for a, b in zip(positions, positions[1:]):
            if b - a < MIN_HOP_SPACING:
                yield f"{edges[edge].describe()} has two hops {b - a:g} px apart; they need {MIN_HOP_SPACING} px"


def _covered_labels(diagram: Diagram) -> Iterator[str]:
    labels = [(f"label {label.text!r}", label.rect) for label in diagram.labels]
    strokes: list[tuple[str, Segment]] = [
        (edge.describe(), segment) for edge in diagram.edges for segment in edge.segments
    ] + [(lifeline.describe(), lifeline.segment) for lifeline in diagram.lifelines]
    for what, rect in labels:
        for box_what, box in _boxes(diagram):
            if rect.overlaps(box):
                yield f"{what} {rect} is covered by {box_what}"
        clearance = rect.inflated(MIN_STROKE_CLEARANCE)
        for stroke_what, segment in strokes:
            if clearance.crossed_by(segment):
                yield f"{what} {rect} sits within {MIN_STROKE_CLEARANCE} px of {stroke_what}"
    masked = labels + _tags(diagram)
    for (first_what, first), (second_what, second) in combinations(masked, 2):
        if first.overlaps(second):
            yield f"{first_what} {first} overlaps {second_what} {second}"


def _overlapping_boxes(diagram: Diagram) -> Iterator[str]:
    for (first_what, first), (second_what, second) in combinations(_boxes(diagram), 2):
        if first.overlaps(second):
            yield f"{first_what} {first} overlaps {second_what} {second}"


def _too_many_accents(diagram: Diagram) -> Iterator[str]:
    accents = [node.describe() for node in diagram.nodes if node.kind == "focal"]
    accents += [edge.describe() for edge in diagram.edges if edge.color == "accent"]
    if len(accents) > MAX_ACCENTS:
        yield f"{len(accents)} accent elements ({', '.join(accents)}); use {MAX_ACCENTS} or fewer"
