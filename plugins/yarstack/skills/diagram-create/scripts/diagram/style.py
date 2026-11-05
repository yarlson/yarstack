"""Visual tokens of the diagram style. Change the look here, and only here."""

from __future__ import annotations

from dataclasses import dataclass
from html import escape

PAPER = "#f8fafc"
INK = "#0f172a"
MUTED = "#475569"
SOFT = "#64748b"
ACCENT = "#ea580c"
ACCENT_TINT = "rgba(234,88,12,0.08)"
LINK = "#2563eb"

WASH = "rgba(15,23,42,0.02)"
CHIP_FILL = "rgba(15,23,42,0.05)"
RULE = "rgba(15,23,42,0.12)"
FRAME_STROKE = "rgba(15,23,42,0.16)"
DASHED_RULE = "rgba(15,23,42,0.20)"
FRAGMENT_STROKE = "rgba(15,23,42,0.22)"
TAB_STROKE = "rgba(15,23,42,0.30)"

FONT_STYLESHEET = (
    "https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600"
    "&family=Geist+Mono:wght@400;500;600&display=swap"
)
SANS = "'Geist', sans-serif"
MONO = "'Geist Mono', monospace"
MONO_ADVANCE_EM = 0.6

CANVAS_WIDTH = 1280
MARGIN = 40
GRID = 4
NODE_HEIGHT = 80
CORNER_RADIUS = 8
HOP_RADIUS = 8
ARROW_LENGTH = 8
LABEL_GAP = 8
LABEL_HEIGHT = 12


@dataclass(frozen=True, slots=True)
class TextStyle:
    family: str
    size: int
    weight: int = 400
    tracking_em: float = 0.0

    def attributes(self, color: str) -> str:
        weight = f' font-weight="{self.weight}"' if self.weight != 400 else ""
        tracking = f' letter-spacing="{self.tracking_em:g}em"' if self.tracking_em else ""
        return f'fill="{color}" font-size="{self.size}"{weight} font-family="{escape(self.family)}"{tracking}'

    def mono_width(self, text: str) -> float:
        """Exact advance width of text in Geist Mono, letter spacing included."""
        if self.family != MONO:
            raise ValueError("only monospace text has an exact width before rendering")
        return len(text) * self.size * (MONO_ADVANCE_EM + self.tracking_em)


NAME = TextStyle(SANS, 12, weight=600)
CHIP_NAME = TextStyle(SANS, 12)
SUBLABEL = TextStyle(MONO, 9)
LABEL = TextStyle(MONO, 8, tracking_em=0.06)
TAG = TextStyle(MONO, 7, tracking_em=0.08)
EYEBROW = TextStyle(MONO, 8, tracking_em=0.14)
LANE_DETAIL = TextStyle(MONO, 7, tracking_em=0.12)
OPERATOR = TextStyle(MONO, 8, tracking_em=0.12)


@dataclass(frozen=True, slots=True)
class NodeStyle:
    fill: str
    stroke: str
    tag: str
    dashed: bool = False


NODE_STYLES = {
    "backend": NodeStyle("#ffffff", INK, INK),
    "focal": NodeStyle(ACCENT_TINT, ACCENT, ACCENT),
    "store": NodeStyle("rgba(15,23,42,0.05)", MUTED, MUTED),
    "external": NodeStyle("rgba(15,23,42,0.03)", "rgba(15,23,42,0.30)", SOFT),
    "input": NodeStyle("rgba(71,85,105,0.10)", SOFT, SOFT),
    "optional": NodeStyle("rgba(15,23,42,0.02)", "rgba(15,23,42,0.20)", SOFT, dashed=True),
}

EDGE_COLORS = {"muted": MUTED, "link": LINK, "accent": ACCENT}
