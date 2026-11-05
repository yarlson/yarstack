"""Build pipeline: check every diagram, write the HTML pages, and render the PNGs."""

from __future__ import annotations

import argparse
import re
import sys
from collections import Counter
from collections.abc import Iterator, Sequence
from dataclasses import dataclass
from html import escape
from pathlib import Path

from . import style
from .model import Diagram
from .icons import download_icon, missing_icons, unused_icon_files
from .rules import check

PAGE_FILE_PATTERN = re.compile(r"[a-z0-9]+(?:-[a-z0-9]+)*\.html")
ICONS_FOLDER = "icons"

PAGE_TEMPLATE = """<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{title}</title>
  <link href="{fonts}" rel="stylesheet">
  <style>
    *, *::before, *::after {{ box-sizing: border-box; margin: 0; padding: 0; }}
    body {{ background: {paper}; padding: 2rem; }}
    .frame {{ max-width: 1320px; margin: 0 auto; display: grid; gap: 2rem; }}
    .diagram-container {{ width: 100%; min-width: 0; overflow-x: auto; }}
    .diagram-container > svg {{ width: 100%; min-width: {width}px; display: block; }}
    @media print {{
      .diagram-container {{ overflow-x: visible; }}
      .diagram-container > svg {{ min-width: 0; }}
    }}
  </style>
</head>
<body>
  <div class="frame">
{diagrams}
  </div>
</body>
</html>
"""


@dataclass(frozen=True, slots=True)
class Page:
    """One HTML page that shows one or more diagrams, one PNG each."""

    file_name: str
    title: str
    diagrams: Sequence[Diagram]


def run(folder: Path, pages: Sequence[Page], argv: Sequence[str] | None = None) -> int:
    """Command line of a diagrams script. Return the exit status."""
    parser = argparse.ArgumentParser(description="Build the diagrams in this folder, or download their icons.")
    parser.add_argument(
        "command",
        nargs="?",
        default="build",
        choices=("build", "icons"),
        help="build (default): check, write the HTML pages, render the PNGs; icons: download missing icons",
    )
    args = parser.parse_args(argv)
    if args.command == "icons":
        return download_missing_icons(folder, pages)
    return build(folder, pages)


def build(folder: Path, pages: Sequence[Page]) -> int:
    icons_dir = folder / ICONS_FOLDER
    diagrams = [diagram for page in pages for diagram in page.diagrams]
    refs = [ref for diagram in diagrams for ref in diagram.icon_refs()]
    errors = list(_page_errors(pages))
    errors += [f"icon {ref}: {ref.path(icons_dir)} is missing; run the icons command" for ref in missing_icons(refs, icons_dir)]
    errors += [finding for diagram in diagrams for finding in check(diagram)]
    if errors:
        return _fail(errors)
    written = []
    for page in pages:
        path = folder / page.file_name
        path.write_text(_page_html(page, icons_dir), encoding="utf-8")
        written.append(path)
        print(f"wrote {page.file_name}")
    # Imported here so that the icons command and the checks work without Playwright.
    from .raster import render_pngs

    findings = render_pngs(written, folder)
    if findings:
        return _fail(findings)
    for warning in _stale_files(folder, pages, diagrams):
        print(f"warning: {warning}", file=sys.stderr)
    for path in unused_icon_files(refs, icons_dir):
        print(f"warning: no diagram uses {path.relative_to(folder)}; delete it", file=sys.stderr)
    return 0


def download_missing_icons(folder: Path, pages: Sequence[Page]) -> int:
    icons_dir = folder / ICONS_FOLDER
    refs = [ref for page in pages for diagram in page.diagrams for ref in diagram.icon_refs()]
    errors = []
    for ref in missing_icons(refs, icons_dir):
        try:
            path = download_icon(ref, icons_dir)
        except (LookupError, ValueError, OSError) as error:
            errors.append(str(error))
            continue
        print(f"downloaded {path.relative_to(folder)}")
    for path in unused_icon_files(refs, icons_dir):
        print(f"warning: no diagram uses {path.relative_to(folder)}; delete it", file=sys.stderr)
    return _fail(errors) if errors else 0


def _page_errors(pages: Sequence[Page]) -> Iterator[str]:
    for page in pages:
        if not PAGE_FILE_PATTERN.fullmatch(page.file_name):
            yield f"page {page.file_name!r}: the file name must be lowercase words joined by hyphens, ending in .html"
        if not page.diagrams:
            yield f"page {page.file_name!r} has no diagrams"
    for name, count in Counter(page.file_name for page in pages).items():
        if count > 1:
            yield f"page {name!r} is defined {count} times"
    for slug, count in Counter(diagram.slug for page in pages for diagram in page.diagrams).items():
        if count > 1:
            yield f"slug {slug!r} is used {count} times; it names the PNG, so it must be unique"


def _page_html(page: Page, icons_dir: Path) -> str:
    diagrams = "\n".join(
        f'    <div class="diagram-container">\n{diagram.svg(icons_dir)}\n    </div>' for diagram in page.diagrams
    )
    return PAGE_TEMPLATE.format(
        title=escape(page.title),
        fonts=escape(style.FONT_STYLESHEET),
        paper=style.PAPER,
        width=style.CANVAS_WIDTH,
        diagrams=diagrams,
    )


def _stale_files(folder: Path, pages: Sequence[Page], diagrams: Sequence[Diagram]) -> Iterator[str]:
    page_names = {page.file_name for page in pages}
    png_names = {f"{diagram.slug}.png" for diagram in diagrams}
    for path in sorted(folder.glob("*.html")):
        if path.name not in page_names:
            yield f"{path.name} is not built by any page; delete it if it is old"
    for path in sorted(folder.glob("*.png")):
        if path.name not in png_names:
            yield f"{path.name} matches no diagram slug; delete it if it is old"


def _fail(messages: Sequence[str]) -> int:
    for message in messages:
        print(f"error: {message}", file=sys.stderr)
    return 1
