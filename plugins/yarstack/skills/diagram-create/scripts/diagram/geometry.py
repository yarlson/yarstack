"""Axis-aligned geometry: rectangles and horizontal or vertical segments."""

from __future__ import annotations

from dataclasses import dataclass

Point = tuple[float, float]

EPSILON = 0.5


@dataclass(frozen=True, slots=True)
class Rect:
    x: float
    y: float
    w: float
    h: float

    @property
    def left(self) -> float:
        return self.x

    @property
    def top(self) -> float:
        return self.y

    @property
    def right(self) -> float:
        return self.x + self.w

    @property
    def bottom(self) -> float:
        return self.y + self.h

    @property
    def center_x(self) -> float:
        return self.x + self.w / 2

    @property
    def center_y(self) -> float:
        return self.y + self.h / 2

    def __str__(self) -> str:
        return f"({self.x:g},{self.y:g} {self.w:g}x{self.h:g})"

    def inflated(self, margin: float) -> Rect:
        return Rect(self.x - margin, self.y - margin, self.w + 2 * margin, self.h + 2 * margin)

    def overlaps(self, other: Rect) -> bool:
        """True when the two interiors share an area."""
        return (
            min(self.right, other.right) - max(self.x, other.x) > EPSILON
            and min(self.bottom, other.bottom) - max(self.y, other.y) > EPSILON
        )

    def border_side(self, point: Point) -> str | None:
        """The side whose border holds the point, or None."""
        x, y = point
        on_x_span = self.x - EPSILON <= x <= self.right + EPSILON
        on_y_span = self.y - EPSILON <= y <= self.bottom + EPSILON
        if on_y_span and abs(x - self.x) <= EPSILON:
            return "left"
        if on_y_span and abs(x - self.right) <= EPSILON:
            return "right"
        if on_x_span and abs(y - self.y) <= EPSILON:
            return "top"
        if on_x_span and abs(y - self.bottom) <= EPSILON:
            return "bottom"
        return None

    def crossed_by(self, segment: Segment) -> bool:
        """True when the segment runs through the interior, not only along the border."""
        if segment.horizontal:
            inside = self.y + EPSILON < segment.position < self.bottom - EPSILON
            return inside and min(segment.high, self.right) - max(segment.low, self.x) > EPSILON
        inside = self.x + EPSILON < segment.position < self.right - EPSILON
        return inside and min(segment.high, self.bottom) - max(segment.low, self.y) > EPSILON


@dataclass(frozen=True, slots=True)
class Segment:
    """A horizontal or vertical segment from a to b."""

    a: Point
    b: Point

    def __post_init__(self) -> None:
        if self.a == self.b:
            raise ValueError(f"segment {self.a} has zero length")
        if self.a[0] != self.b[0] and self.a[1] != self.b[1]:
            raise ValueError(f"segment {self.a} -> {self.b} is diagonal")

    @property
    def horizontal(self) -> bool:
        return self.a[1] == self.b[1]

    @property
    def position(self) -> float:
        """The fixed coordinate: y of a horizontal segment, x of a vertical one."""
        return self.a[1] if self.horizontal else self.a[0]

    @property
    def low(self) -> float:
        return min(self.a[0], self.b[0]) if self.horizontal else min(self.a[1], self.b[1])

    @property
    def high(self) -> float:
        return max(self.a[0], self.b[0]) if self.horizontal else max(self.a[1], self.b[1])

    @property
    def length(self) -> float:
        return self.high - self.low

    def holds(self, point: Point) -> bool:
        along, across = (point[0], point[1]) if self.horizontal else (point[1], point[0])
        return abs(across - self.position) <= EPSILON and self.low - EPSILON <= along <= self.high + EPSILON

    def crossing(self, other: Segment) -> Point | None:
        """The point where two perpendicular segments cross inside both, or None."""
        if self.horizontal == other.horizontal:
            return None
        horizontal, vertical = (self, other) if self.horizontal else (other, self)
        x, y = vertical.position, horizontal.position
        if horizontal.low < x < horizontal.high and vertical.low < y < vertical.high:
            return (x, y)
        return None
