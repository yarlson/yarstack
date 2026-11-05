"""Icon files: the trust boundary for downloaded SVG and PNG icons.

Icons come from third-party packages, so every SVG is parsed and reduced to
plain shapes before a diagram inlines it.
"""

from __future__ import annotations

import base64
import re
import urllib.error
import urllib.request
import xml.etree.ElementTree as ElementTree
from collections.abc import Iterable
from dataclasses import dataclass
from html import escape
from pathlib import Path

MAX_ICON_BYTES = 64 * 1024
DOWNLOAD_TIMEOUT_SECONDS = 30
NAME_PATTERN = re.compile(r"[a-z0-9]+(?:-[a-z0-9]+)*")
SVG_NAMESPACE = "{http://www.w3.org/2000/svg}"
SHAPES = {"path", "circle", "ellipse", "line", "polyline", "polygon", "rect"}
SHAPE_ATTRIBUTES = {
    "d", "cx", "cy", "r", "rx", "ry", "x", "y", "width", "height",
    "x1", "y1", "x2", "y2", "points", "fill-rule", "clip-rule",
}
PAINT_ATTRIBUTES = {"fill", "stroke"}
PAINT_VALUES = {"none", "currentColor"}
PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"


@dataclass(frozen=True, slots=True)
class IconPackage:
    url_template: str
    license_url: str | None
    stroked: bool
    catalog_url: str


PACKAGES = {
    "tabler": IconPackage(
        "https://cdn.jsdelivr.net/npm/@tabler/icons@3.49.0/icons/outline/{name}.svg",
        "https://cdn.jsdelivr.net/npm/@tabler/icons@3.49.0/LICENSE",
        stroked=True,
        catalog_url="https://tabler.io/icons",
    ),
    "simple-icons": IconPackage(
        "https://cdn.jsdelivr.net/npm/simple-icons@16.35.0/icons/{name}.svg",
        None,
        stroked=False,
        catalog_url="https://simpleicons.org",
    ),
}


@dataclass(frozen=True, slots=True)
class IconRef:
    """An icon reference: tabler:<name>, simple-icons:<slug>, or <folder>:<name> for a PNG."""

    source: str
    name: str

    @classmethod
    def parse(cls, ref: str) -> IconRef:
        source, _, name = ref.partition(":")
        if not NAME_PATTERN.fullmatch(source) or not NAME_PATTERN.fullmatch(name):
            raise ValueError(f"icon {ref!r} must be <source>:<name> in lowercase letters, digits, and hyphens")
        return cls(source, name)

    def __str__(self) -> str:
        return f"{self.source}:{self.name}"

    def path(self, icons_dir: Path) -> Path:
        suffix = ".svg" if self.source in PACKAGES else ".png"
        return icons_dir / self.source / f"{self.name}{suffix}"


def icon_markup(ref: IconRef, icons_dir: Path, x: float, y: float, size: int, color: str) -> str:
    path = ref.path(icons_dir)
    data = _read_icon(path)
    position = f'x="{x:g}" y="{y:g}" width="{size}" height="{size}"'
    if ref.source not in PACKAGES:
        if not data.startswith(PNG_SIGNATURE):
            raise ValueError(f"{path}: not a PNG file")
        encoded = base64.b64encode(data).decode("ascii")
        return f'<image {position} href="data:image/png;base64,{encoded}" aria-hidden="true"/>'
    shapes = shapes_markup(data, str(path))
    if PACKAGES[ref.source].stroked:
        return (
            f'<svg aria-hidden="true" {position} viewBox="0 0 24 24" fill="none" stroke="currentColor" '
            f'stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" color="{color}">{shapes}</svg>'
        )
    return f'<svg aria-hidden="true" {position} viewBox="0 0 24 24" fill="{color}">{shapes}</svg>'


def shapes_markup(data: bytes, origin: str) -> str:
    """Reduce a 24 x 24 icon SVG to its plain shapes. Reject anything that is not a shape."""
    try:
        root = ElementTree.fromstring(data)
    except ElementTree.ParseError as error:
        raise ValueError(f"{origin}: not valid SVG ({error})") from error
    if root.tag.removeprefix(SVG_NAMESPACE) != "svg" or root.get("viewBox") != "0 0 24 24":
        raise ValueError(f"{origin}: expected an <svg> with viewBox 0 0 24 24")
    shapes = []
    for element in root:
        tag = element.tag.removeprefix(SVG_NAMESPACE)
        if tag in {"title", "desc"}:
            continue
        if tag not in SHAPES or len(element):
            raise ValueError(f"{origin}: unsupported element <{tag}>")
        paint = {key: value for key, value in element.attrib.items() if key in PAINT_ATTRIBUTES}
        if any(value not in PAINT_VALUES for value in paint.values()):
            raise ValueError(f"{origin}: <{tag}> sets a paint other than none or currentColor")
        if paint.get("stroke") == "none" and paint.get("fill", "none") == "none":
            continue
        attributes = "".join(
            f' {key}="{escape(value)}"'
            for key, value in element.attrib.items()
            if key in SHAPE_ATTRIBUTES or key in paint
        )
        shapes.append(f"<{tag}{attributes}/>")
    if not shapes:
        raise ValueError(f"{origin}: the icon has no shapes")
    return "".join(shapes)


def missing_icons(refs: Iterable[IconRef], icons_dir: Path) -> list[IconRef]:
    return sorted({ref for ref in refs if not ref.path(icons_dir).is_file()}, key=str)


def unused_icon_files(refs: Iterable[IconRef], icons_dir: Path) -> list[Path]:
    used = {ref.path(icons_dir) for ref in refs}
    return sorted(
        path for path in icons_dir.glob("*/*") if path.is_file() and path.name != "LICENSE" and path not in used
    )


def download_icon(ref: IconRef, icons_dir: Path) -> Path:
    """Download one icon from its pinned package. Write it only if it passes shapes_markup."""
    package = PACKAGES.get(ref.source)
    if package is None:
        raise LookupError(
            f"{ref}: add {ref.path(icons_dir)} by hand; only {' and '.join(PACKAGES)} icons can be downloaded"
        )
    data = _download(
        package.url_template.format(name=ref.name),
        f"{ref}: {package.catalog_url} has no icon named {ref.name!r}",
    )
    shapes_markup(data, str(ref))
    target = ref.path(icons_dir)
    target.parent.mkdir(parents=True, exist_ok=True)
    target.write_bytes(data)
    license_file = target.parent / "LICENSE"
    if package.license_url and not license_file.exists():
        license_file.write_bytes(_download(package.license_url, f"{package.license_url}: license not found"))
    return target


def _read_icon(path: Path) -> bytes:
    if path.stat().st_size > MAX_ICON_BYTES:
        raise ValueError(f"{path}: larger than {MAX_ICON_BYTES} bytes")
    return path.read_bytes()


def _download(url: str, not_found_message: str) -> bytes:
    try:
        with urllib.request.urlopen(url, timeout=DOWNLOAD_TIMEOUT_SECONDS) as response:
            data = response.read(MAX_ICON_BYTES + 1)
    except urllib.error.HTTPError as error:
        if error.code == 404:
            raise LookupError(not_found_message) from error
        raise
    if len(data) > MAX_ICON_BYTES:
        raise ValueError(f"{url}: larger than {MAX_ICON_BYTES} bytes")
    return data
