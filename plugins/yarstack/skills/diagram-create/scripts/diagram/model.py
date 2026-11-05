"""The diagram model: the elements of one SVG, their geometry, and their markup."""

from __future__ import annotations

import math
import re
from collections import defaultdict
from collections.abc import Sequence
from dataclasses import dataclass
from html import escape
from itertools import combinations
from pathlib import Path
from typing import Literal, NamedTuple, get_args

from . import style
from .geometry import Point, Rect, Segment
from .icons import IconRef, icon_markup

NodeKind = Literal["backend", "focal", "store", "external", "input", "optional"]
EdgeColor = Literal["muted", "link", "accent"]
Operator = Literal["ALT", "OPT", "LOOP"]
Side = Literal["left", "right", "top", "bottom"]
LabelSide = Literal["above", "below", "left", "right"]

SLUG_PATTERN = re.compile(r"[a-z0-9]+(?:-[a-z0-9]+)*")
MAX_TITLE_LENGTH = 60
TAB_HEIGHT = 24
OPERATOR_TAB = (40, 16)
MARKER_NAMES = {"muted": "arrow", "link": "arrow-link", "accent": "arrow-accent"}


class Lane(NamedTuple):
    top: int
    bottom: int
    name: str
    detail: str


class Chip(NamedTuple):
    icon: str
    name: str
    version: str


def _snap(value: float) -> int:
    return style.GRID * round(value / style.GRID)


def _grid_ceil(value: float) -> int:
    return style.GRID * math.ceil(value / style.GRID)


def _even_ceil(value: float) -> int:
    return 2 * math.ceil(value / 2)


def _choice(value: str, allowed: type, what: str) -> str:
    options = get_args(allowed)
    if value not in options:
        raise ValueError(f"{what} must be one of {', '.join(options)}, not {value!r}")
    return value


def label_width(text: str) -> int:
    return _even_ceil(style.LABEL.mono_width(text) + 10)


@dataclass(frozen=True, slots=True)
class Box(Rect):
    """A placed rectangle that edges can attach to."""

    def attach_point(self, side: Side, index: int = 1, count: int = 1) -> Point:
        """Point index of count points spread evenly along one side, snapped to the grid."""
        if not 1 <= index <= count:
            raise ValueError(f"attach point {index} of {count} does not exist")
        _choice(side, Side, "side")
        if side in ("left", "right"):
            return (self.x if side == "left" else self.right, self.y + _snap(self.h * index / (count + 1)))
        return (self.x + _snap(self.w * index / (count + 1)), self.y if side == "top" else self.bottom)


@dataclass(frozen=True, slots=True)
class Label:
    text: str
    rect: Rect
    color: str

    def markup(self) -> str:
        r = self.rect
        return (
            f'<g data-fit="2"><rect x="{r.x:g}" y="{r.y:g}" width="{r.w:g}" height="{r.h:g}" rx="2" fill="{style.PAPER}"/>'
            f'<text x="{r.center_x:g}" y="{r.y + 9:g}" {style.LABEL.attributes(self.color)} text-anchor="middle">'
            f"{escape(self.text)}</text></g>"
        )


def _node_frame(box: Box, kind: NodeKind, tag: str) -> list[str]:
    node_style = style.NODE_STYLES[kind]
    dash = ' stroke-dasharray="4,3"' if node_style.dashed else ""
    tag_width = _even_ceil(style.TAG.mono_width(tag) + 12)
    rect = f'x="{box.x:g}" y="{box.y:g}" width="{box.w:g}" height="{box.h:g}" rx="6"'
    return [
        f'<rect {rect} fill="{style.PAPER}"/>',
        f'<rect {rect} fill="{node_style.fill}" stroke="{node_style.stroke}" stroke-width="1"{dash}/>',
        f'<g data-fit="2"><rect x="{box.x + 8:g}" y="{box.y + 6:g}" width="{tag_width}" height="12" rx="2" fill="none" '
        f'stroke="{node_style.tag}" stroke-opacity="0.45" stroke-width="0.8"/>'
        f'<text x="{box.x + 8 + tag_width / 2:g}" y="{box.y + 15:g}" {style.TAG.attributes(node_style.tag)} '
        f'text-anchor="middle">{escape(tag)}</text></g>',
    ]


def _icon_color(kind: NodeKind) -> str:
    return style.ACCENT if kind == "focal" else style.INK


@dataclass(frozen=True, slots=True)
class Node:
    """An 80 px node: tag top left, 28 px icon centered above the name and the sublabel."""

    box: Box
    kind: NodeKind
    tag: str
    icon: IconRef
    name: str
    sublabel: str

    def describe(self) -> str:
        return f"node {self.name!r}"

    def icons(self) -> list[IconRef]:
        return [self.icon]

    def markup(self, icons_dir: Path) -> str:
        box, center = self.box, self.box.center_x
        parts = [
            *_node_frame(box, self.kind, self.tag),
            icon_markup(self.icon, icons_dir, center - 14, box.y + 10, 28, _icon_color(self.kind)),
            f'<text x="{center:g}" y="{box.y + 56:g}" {style.NAME.attributes(style.INK)} text-anchor="middle">'
            f"{escape(self.name)}</text>",
            f'<text x="{center:g}" y="{box.y + 70:g}" {style.SUBLABEL.attributes(style.MUTED)} text-anchor="middle">'
            f"{escape(self.sublabel)}</text>",
        ]
        return '<g data-fit="4">' + "".join(parts) + "</g>"


@dataclass(frozen=True, slots=True)
class ListNode:
    """A tall node: icon beside the name, then chips, then mono lines."""

    box: Box
    kind: NodeKind
    tag: str
    icon: IconRef
    name: str
    chips: tuple[Chip, ...]
    lines: tuple[str, ...]

    CHIPS_TOP = 64
    CHIP_STEP = 32
    LINES_TOP = 76
    LINES_AFTER_CHIPS = 20
    LINE_STEP = 18
    LINE_BOTTOM_PADDING = 14

    @classmethod
    def first_line_offset(cls, chip_count: int) -> int:
        if chip_count:
            return cls.CHIPS_TOP + chip_count * cls.CHIP_STEP + cls.LINES_AFTER_CHIPS
        return cls.LINES_TOP

    @classmethod
    def content_height(cls, chip_count: int, line_count: int) -> int:
        if line_count:
            last_baseline = cls.first_line_offset(chip_count) + (line_count - 1) * cls.LINE_STEP
            return _grid_ceil(last_baseline + cls.LINE_BOTTOM_PADDING)
        if chip_count:
            return cls.CHIPS_TOP + chip_count * cls.CHIP_STEP + 8
        return style.NODE_HEIGHT

    def describe(self) -> str:
        return f"node {self.name!r}"

    def icons(self) -> list[IconRef]:
        return [self.icon, *(IconRef.parse(chip.icon) for chip in self.chips)]

    def markup(self, icons_dir: Path) -> str:
        box = self.box
        parts = [
            *_node_frame(box, self.kind, self.tag),
            icon_markup(self.icon, icons_dir, box.x + 16, box.y + 26, 28, _icon_color(self.kind)),
            f'<text x="{box.x + 52:g}" y="{box.y + 45:g}" {style.NAME.attributes(style.INK)}>{escape(self.name)}</text>',
        ]
        for index, chip in enumerate(self.chips):
            x, y, w = box.x + 16, box.y + self.CHIPS_TOP + index * self.CHIP_STEP, box.w - 32
            parts.append(
                f'<g data-fit="4"><rect x="{x:g}" y="{y:g}" width="{w:g}" height="24" rx="4" fill="{style.CHIP_FILL}" '
                f'stroke="{style.MUTED}" stroke-width="0.8"/>'
                f"{icon_markup(IconRef.parse(chip.icon), icons_dir, x + 6, y + 4, 16, style.INK)}"
                f'<text x="{x + 28:g}" y="{y + 16:g}" {style.CHIP_NAME.attributes(style.INK)}>{escape(chip.name)}</text>'
                f'<text x="{x + w - 8:g}" y="{y + 16:g}" {style.SUBLABEL.attributes(style.MUTED)} text-anchor="end">'
                f"{escape(chip.version)}</text></g>"
            )
        first_line = box.y + self.first_line_offset(len(self.chips))
        for index, line in enumerate(self.lines):
            parts.append(
                f'<text x="{box.x + 16:g}" y="{first_line + index * self.LINE_STEP:g}" '
                f"{style.SUBLABEL.attributes(style.MUTED)}>{escape(line)}</text>"
            )
        return '<g data-fit="4">' + "".join(parts) + "</g>"


@dataclass(frozen=True, slots=True)
class Zone:
    box: Box
    label: str

    def describe(self) -> str:
        return f"zone {self.label!r}"

    @property
    def eyebrow(self) -> Rect:
        return Rect(self.box.x + 12, self.box.y + 6, _even_ceil(style.EYEBROW.mono_width(self.label) + 12), 12)

    def markup(self) -> str:
        box, eyebrow = self.box, self.eyebrow
        return (
            f'<rect x="{box.x:g}" y="{box.y:g}" width="{box.w:g}" height="{box.h:g}" rx="8" fill="{style.WASH}" '
            f'stroke="{style.FRAME_STROKE}" stroke-width="0.8" stroke-dasharray="4,4"/>'
            f'<g data-fit="2"><rect x="{eyebrow.x:g}" y="{eyebrow.y:g}" width="{eyebrow.w:g}" height="12" rx="2" '
            f'fill="{style.PAPER}"/><text x="{eyebrow.center_x:g}" y="{box.y + 15:g}" '
            f'{style.EYEBROW.attributes(style.MUTED)} text-anchor="middle">{escape(self.label)}</text></g>'
        )


@dataclass(frozen=True, slots=True)
class Tab:
    """A small labeled box on the edge of a zone. An edge that applies to the whole zone ends on it."""

    box: Box
    label: str

    def describe(self) -> str:
        return f"tab {self.label!r}"

    def markup(self) -> str:
        box = self.box
        return (
            f'<g data-fit="2"><rect x="{box.x:g}" y="{box.y:g}" width="{box.w:g}" height="{box.h:g}" rx="4" '
            f'fill="{style.PAPER}" stroke="{style.TAB_STROKE}" stroke-width="0.8"/>'
            f'<text x="{box.center_x:g}" y="{box.y + 15:g}" {style.EYEBROW.attributes(style.INK)} text-anchor="middle">'
            f"{escape(self.label)}</text></g>"
        )


@dataclass(frozen=True, slots=True)
class Lanes:
    x: int
    w: int
    label_width: int
    lanes: tuple[Lane, ...]

    @property
    def frame(self) -> Rect:
        return Rect(self.x, self.lanes[0].top, self.w, self.lanes[-1].bottom - self.lanes[0].top)

    def markup(self) -> str:
        frame, divider_x = self.frame, self.x + self.label_width
        parts = [
            f'<rect x="{frame.x:g}" y="{frame.y:g}" width="{frame.w:g}" height="{frame.h:g}" rx="8" fill="none" '
            f'stroke="{style.FRAME_STROKE}" stroke-width="0.8"/>',
            f'<line x1="{divider_x}" y1="{frame.top:g}" x2="{divider_x}" y2="{frame.bottom:g}" stroke="{style.RULE}" '
            f'stroke-width="0.8"/>',
        ]
        for lane in self.lanes:
            if lane.top > frame.top:
                parts.append(
                    f'<line x1="{self.x}" y1="{lane.top}" x2="{self.x + self.w}" y2="{lane.top}" stroke="{style.RULE}" '
                    f'stroke-width="0.8"/>'
                )
            middle = (lane.top + lane.bottom) / 2
            parts.append(
                f'<g data-fit="4"><rect x="{self.x}" y="{lane.top}" width="{self.label_width}" '
                f'height="{lane.bottom - lane.top}" fill="none"/>'
                f'<text x="{self.x + 16}" y="{middle - 2:g}" {style.EYEBROW.attributes(style.INK)}>{escape(lane.name)}</text>'
                f'<text x="{self.x + 16}" y="{middle + 12:g}" {style.LANE_DETAIL.attributes(style.SOFT)}>'
                f"{escape(lane.detail)}</text></g>"
            )
        return "".join(parts)


@dataclass(frozen=True, slots=True)
class Fragment:
    """A sequence combined fragment: frame, operator tab, and dashed region dividers."""

    box: Box
    operator: Operator
    dividers: tuple[int, ...]

    def describe(self) -> str:
        return f"{self.operator} fragment"

    @property
    def operator_tab(self) -> Rect:
        return Rect(self.box.x, self.box.y, *OPERATOR_TAB)

    def markup(self) -> str:
        box, tab = self.box, self.operator_tab
        parts = [
            f'<rect x="{box.x:g}" y="{box.y:g}" width="{box.w:g}" height="{box.h:g}" rx="4" fill="{style.WASH}" '
            f'stroke="{style.FRAGMENT_STROKE}" stroke-width="1"/>',
            f'<g data-fit="2"><rect x="{tab.x:g}" y="{tab.y:g}" width="{tab.w:g}" height="{tab.h:g}" rx="2" '
            f'fill="{style.PAPER}" stroke="{style.FRAGMENT_STROKE}" stroke-width="1"/>'
            f'<text x="{tab.center_x:g}" y="{tab.y + 12:g}" {style.OPERATOR.attributes(style.MUTED)} '
            f'text-anchor="middle">{self.operator}</text></g>',
        ]
        for y in self.dividers:
            parts.append(
                f'<line x1="{box.x + 8:g}" y1="{y}" x2="{box.right - 8:g}" y2="{y}" stroke="{style.DASHED_RULE}" '
                f'stroke-width="1" stroke-dasharray="4,3"/>'
            )
        return "".join(parts)


@dataclass(frozen=True, slots=True)
class Lifeline:
    x: float
    top: float
    bottom: float

    def describe(self) -> str:
        return f"lifeline at x={self.x:g}"

    @property
    def segment(self) -> Segment:
        return Segment((self.x, self.top), (self.x, self.bottom))

    def markup(self) -> str:
        return (
            f'<line x1="{self.x:g}" y1="{self.top:g}" x2="{self.x:g}" y2="{self.bottom:g}" stroke="{style.DASHED_RULE}" '
            f'stroke-width="1" stroke-dasharray="3,3"/>'
        )


@dataclass(frozen=True, slots=True)
class Hop:
    """Edge number edge jumps over edge number over, on its segment number segment, at coordinate at."""

    edge: int
    segment: int
    at: float
    over: int


@dataclass(frozen=True, slots=True)
class Edge:
    points: tuple[Point, ...]
    color: EdgeColor
    dashed: bool
    bidirectional: bool
    label: Label | None

    @property
    def segments(self) -> list[Segment]:
        return [Segment(a, b) for a, b in zip(self.points, self.points[1:])]

    @property
    def rank(self) -> tuple[int, int]:
        """Importance at a crossing: the lower rank hops."""
        return (0 if self.dashed else 1, list(MARKER_NAMES).index(self.color))

    def describe(self) -> str:
        start, end = self.points[0], self.points[-1]
        name = f" {self.label.text!r}" if self.label else ""
        return f"edge{name} ({start[0]:g},{start[1]:g})->({end[0]:g},{end[1]:g})"

    def markup(self, slug: str, hops: Sequence[Hop]) -> str:
        stroke = style.EDGE_COLORS[self.color]
        dash = ' stroke-dasharray="5,4"' if self.dashed else ""
        width = "1" if self.dashed else "1.2"
        marker = f"url(#{slug}-{MARKER_NAMES[self.color]})"
        start = f' marker-start="{marker}"' if self.bidirectional else ""
        return (
            f'<path d="{self._path_data(hops)}" fill="none" stroke="{stroke}" stroke-width="{width}"{dash}{start} '
            f'marker-end="{marker}"/>'
        )

    def _path_data(self, hops: Sequence[Hop]) -> str:
        r, hop_r = style.CORNER_RADIUS, style.HOP_RADIUS
        hops_by_segment = defaultdict(list)
        for hop in hops:
            hops_by_segment[hop.segment].append(hop.at)
        points, last = self.points, len(self.points) - 1
        commands = [f"M {points[0][0]:g},{points[0][1]:g}"]
        for index, segment in enumerate(self.segments):
            (ax, ay), (bx, by) = segment.a, segment.b
            horizontal = segment.horizontal
            step = 1 if (bx > ax if horizontal else by > ay) else -1
            command = "H" if horizontal else "V"
            for at in sorted(hops_by_segment[index], key=lambda value: value * step):
                if horizontal:
                    commands.append(f"H {at - hop_r * step:g} a {hop_r},{hop_r} 0 0,{1 if step > 0 else 0} {2 * hop_r * step},0")
                else:
                    commands.append(f"V {at - hop_r * step:g} a {hop_r},{hop_r} 0 0,{0 if step > 0 else 1} 0,{2 * hop_r * step}")
            corner = index + 1 < last
            end = (bx if horizontal else by) - (r * step if corner else 0)
            commands.append(f"{command} {end:g}")
            if corner:
                nx, ny = points[index + 2]
                if horizontal:
                    commands.append(f"Q {bx:g},{by:g} {bx:g},{by + (r if ny > by else -r):g}")
                else:
                    commands.append(f"Q {bx:g},{by:g} {bx + (r if nx > bx else -r):g},{by:g}")
        return " ".join(commands)


def _edge_points(points: Sequence[Point]) -> tuple[Point, ...]:
    """Validate the waypoints and drop the ones that do not turn."""
    if len(points) < 2:
        raise ValueError("an edge needs at least two points")
    segments = [Segment(a, b) for a, b in zip(points, points[1:])]
    kept = [points[0]]
    for before, after in zip(segments, segments[1:]):
        if before.horizontal != after.horizontal:
            kept.append(before.b)
            continue
        axis = 0 if before.horizontal else 1
        if (before.b[axis] - before.a[axis]) * (after.b[axis] - after.a[axis]) < 0:
            raise ValueError(f"edge turns back on itself at {before.b}")
    kept.append(points[-1])
    return tuple(kept)


def _edge_label(text: str, segments: Sequence[Segment], color: str, at: Point | None, side: LabelSide | None) -> Label:
    width, height, gap = label_width(text), style.LABEL_HEIGHT, style.LABEL_GAP
    if at is None:
        segment = max(segments, key=lambda candidate: candidate.length)
    else:
        segment = next((candidate for candidate in segments if candidate.holds(at)), None)
        if segment is None:
            raise ValueError(f"label_at {at} for {text!r} is not on its edge")
    if segment.horizontal:
        side = side or "above"
        if side not in ("above", "below"):
            raise ValueError(f"a label on a horizontal segment goes above or below, not {side!r}")
        center = at[0] if at else (segment.low + segment.high) / 2
        top = segment.position - gap - height if side == "above" else segment.position + gap
        return Label(text, Rect(center - width / 2, top, width, height), color)
    side = side or "right"
    if side not in ("left", "right"):
        raise ValueError(f"a label on a vertical segment goes left or right, not {side!r}")
    middle = at[1] if at else (segment.low + segment.high) / 2
    left = segment.position + gap if side == "right" else segment.position - gap - width
    return Label(text, Rect(left, middle - height / 2, width, height), color)


class Diagram:
    """One SVG, 1280 px wide. Its height fits the content plus the margin."""

    def __init__(self, slug: str, *, title: str, description: str) -> None:
        if not SLUG_PATTERN.fullmatch(slug):
            raise ValueError(f"slug {slug!r} must be lowercase words joined by hyphens")
        if not title.strip() or len(title) > MAX_TITLE_LENGTH:
            raise ValueError(f"{slug}: the title must have 1 to {MAX_TITLE_LENGTH} characters")
        if not description.strip():
            raise ValueError(f"{slug}: the description must not be empty")
        self.slug, self.title, self.description = slug, title, description
        self.zones: list[Zone] = []
        self.tabs: list[Tab] = []
        self.lane_frames: list[Lanes] = []
        self.fragments: list[Fragment] = []
        self.lifelines: list[Lifeline] = []
        self.nodes: list[Node | ListNode] = []
        self.edges: list[Edge] = []
        self.guards: list[Label] = []

    def zone(self, x: int, y: int, w: int, h: int, *, label: str) -> Box:
        zone = Zone(Box(x, y, w, h), label)
        if zone.eyebrow.right > zone.box.right - 12:
            raise ValueError(f"{self.slug}: label {label!r} is wider than its zone")
        self.zones.append(zone)
        return zone.box

    def tab(self, x: int, y: int, w: int, *, label: str) -> Box:
        tab = Tab(Box(x, y, w, TAB_HEIGHT), label)
        self.tabs.append(tab)
        return tab.box

    def lanes(self, x: int, w: int, *, label_width: int, lanes: Sequence[Lane]) -> None:
        lanes = tuple(Lane(*lane) for lane in lanes)
        if not lanes:
            raise ValueError(f"{self.slug}: a lanes frame needs at least one lane")
        for before, after in zip(lanes, lanes[1:]):
            if before.bottom != after.top:
                raise ValueError(f"{self.slug}: lane {after.name!r} must start where {before.name!r} ends")
        if any(lane.bottom <= lane.top for lane in lanes):
            raise ValueError(f"{self.slug}: every lane needs a bottom below its top")
        self.lane_frames.append(Lanes(x, w, label_width, lanes))

    def fragment(self, x: int, y: int, w: int, h: int, *, operator: Operator, dividers: Sequence[int] = ()) -> Box:
        box = Box(x, y, w, h)
        if any(not box.top < divider < box.bottom for divider in dividers):
            raise ValueError(f"{self.slug}: every divider of the {operator} fragment must be inside it")
        self.fragments.append(Fragment(box, _choice(operator, Operator, "operator"), tuple(dividers)))
        return box

    def lifeline(self, actor: Box, *, bottom: int) -> None:
        self.lifelines.append(Lifeline(actor.center_x, actor.bottom, bottom))

    def node(
        self, x: int, y: int, w: int, *, kind: NodeKind, tag: str, icon: str, name: str, sublabel: str
    ) -> Box:
        node = Node(Box(x, y, w, style.NODE_HEIGHT), _choice(kind, NodeKind, "kind"), tag, IconRef.parse(icon), name, sublabel)
        self.nodes.append(node)
        return node.box

    def list_node(
        self,
        x: int,
        y: int,
        w: int,
        *,
        kind: NodeKind,
        tag: str,
        icon: str,
        name: str,
        chips: Sequence[Chip] = (),
        lines: Sequence[str] = (),
        h: int | None = None,
    ) -> Box:
        chips = tuple(Chip(*chip) for chip in chips)
        for chip in chips:
            IconRef.parse(chip.icon)
        needed = ListNode.content_height(len(chips), len(lines))
        if h is not None and h < needed:
            raise ValueError(f"{self.slug}: node {name!r} needs h >= {needed} for its chips and lines")
        node = ListNode(
            Box(x, y, w, h or needed), _choice(kind, NodeKind, "kind"), tag, IconRef.parse(icon), name, chips, tuple(lines)
        )
        self.nodes.append(node)
        return node.box

    def edge(
        self,
        *points: Point,
        color: EdgeColor = "muted",
        dashed: bool = False,
        bidirectional: bool = False,
        label: str | None = None,
        label_at: Point | None = None,
        label_side: LabelSide | None = None,
    ) -> None:
        try:
            kept = _edge_points(points)
        except ValueError as error:
            raise ValueError(f"{self.slug}: edge {points}: {error}") from error
        _choice(color, EdgeColor, "color")
        segments = [Segment(a, b) for a, b in zip(kept, kept[1:])]
        placed = _edge_label(label, segments, style.EDGE_COLORS[color], label_at, label_side) if label else None
        self.edges.append(Edge(kept, color, dashed, bidirectional, placed))

    def guard(self, condition: str, *, left: int, top: int) -> None:
        text = f"[{condition}]"
        self.guards.append(Label(text, Rect(left, top, label_width(text), style.LABEL_HEIGHT), style.MUTED))

    @property
    def labels(self) -> list[Label]:
        return [edge.label for edge in self.edges if edge.label] + self.guards

    def icon_refs(self) -> list[IconRef]:
        return [ref for node in self.nodes for ref in node.icons()]

    def hops(self) -> list[Hop]:
        hops = []
        for (i, first), (j, second) in combinations(enumerate(self.edges), 2):
            loser = j if second.rank <= first.rank else i
            for si, s in enumerate(first.segments):
                for sj, t in enumerate(second.segments):
                    point = s.crossing(t)
                    if point is None:
                        continue
                    hopped = s if loser == i else t
                    hops.append(
                        Hop(
                            edge=loser,
                            segment=si if loser == i else sj,
                            at=point[0] if hopped.horizontal else point[1],
                            over=j if loser == i else i,
                        )
                    )
        return hops

    @property
    def height(self) -> int:
        bottoms = [node.box.bottom for node in self.nodes]
        bottoms += [zone.box.bottom for zone in self.zones] + [tab.box.bottom for tab in self.tabs]
        bottoms += [frame.frame.bottom for frame in self.lane_frames] + [fragment.box.bottom for fragment in self.fragments]
        bottoms += [lifeline.bottom for lifeline in self.lifelines] + [label.rect.bottom for label in self.labels]
        bottoms += [point[1] for edge in self.edges for point in edge.points]
        return _grid_ceil(max(bottoms, default=0) + style.MARGIN)

    def svg(self, icons_dir: Path) -> str:
        hops_by_edge = defaultdict(list)
        for hop in self.hops():
            hops_by_edge[hop.edge].append(hop)
        markers = "".join(
            f'<marker id="{self.slug}-{name}" markerWidth="8" markerHeight="6" refX="7" refY="3" '
            f'orient="auto-start-reverse"><polygon points="0 0, 8 3, 0 6" fill="{style.EDGE_COLORS[color]}"/></marker>'
            for color, name in MARKER_NAMES.items()
        )
        layers = [
            *(frame.markup() for frame in self.lane_frames),
            *(zone.markup() for zone in self.zones),
            *(fragment.markup() for fragment in self.fragments),
            *(lifeline.markup() for lifeline in self.lifelines),
            *(tab.markup() for tab in self.tabs),
            *(edge.markup(self.slug, hops_by_edge[index]) for index, edge in enumerate(self.edges)),
            *(label.markup() for label in self.labels),
            *(node.markup(icons_dir) for node in self.nodes),
        ]
        body = "\n".join(layers)
        return (
            f'<svg viewBox="0 0 {style.CANVAS_WIDTH} {self.height}" xmlns="http://www.w3.org/2000/svg" role="img" '
            f'aria-labelledby="{self.slug}-title {self.slug}-desc">\n'
            f'<title id="{self.slug}-title">{escape(self.title)}</title>\n'
            f'<desc id="{self.slug}-desc">{escape(self.description)}</desc>\n'
            f"<defs>{markers}</defs>\n"
            f'<rect width="100%" height="100%" fill="{style.PAPER}"/>\n{body}\n</svg>'
        )
