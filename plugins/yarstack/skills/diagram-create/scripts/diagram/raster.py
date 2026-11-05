"""Render each diagram SVG to its own PNG in Chromium, after checking that its text fits.

Only this module needs Playwright. Run the diagrams script with `uv run` to get it.
"""

from __future__ import annotations

import struct
from collections.abc import Sequence
from importlib.metadata import version
from pathlib import Path

from playwright.sync_api import Browser, sync_playwright
from playwright.sync_api import Error as PlaywrightError

DIAGRAM_SELECTOR = ".diagram-container > svg"
FONT_CHECKS = ("600 12px Geist", "9px 'Geist Mono'")
SCALE = 2
VIEWPORT = {"width": 1400, "height": 900}

# Every [data-fit] group holds a frame (its first rect) and its texts. A text
# must stay inside the frame, inset by the data-fit value in px, and must not
# overlap another text of the same group. Fonts are real here, so sans-serif
# widths are exact.
TEXT_FIT_SCRIPT = """svg => {
  const slug = svg.getAttribute('aria-labelledby').split(' ')[0].replace(/-title$/, '');
  const findings = [];
  for (const group of svg.querySelectorAll('[data-fit]')) {
    const inset = Number(group.dataset.fit);
    const frame = group.querySelector(':scope > rect').getBBox();
    const texts = [...group.querySelectorAll(':scope > text')].map(t => ({text: t.textContent, box: t.getBBox()}));
    for (const {text, box} of texts) {
      const fits = box.x >= frame.x + inset - 0.5 && box.x + box.width <= frame.x + frame.width - inset + 0.5
        && box.y >= frame.y - 1 && box.y + box.height <= frame.y + frame.height + 1;
      if (!fits) {
        findings.push(`${slug}: text "${text}" is ${Math.ceil(box.width)} px wide and does not fit its ` +
          `${frame.width} px box with ${inset} px inset`);
      }
    }
    for (let i = 0; i < texts.length; i++) {
      for (let j = i + 1; j < texts.length; j++) {
        const a = texts[i].box, b = texts[j].box;
        if (a.x < b.x + b.width && b.x < a.x + a.width && a.y < b.y + b.height && b.y < a.y + a.height) {
          findings.push(`${slug}: text "${texts[i].text}" overlaps "${texts[j].text}"`);
        }
      }
    }
  }
  return findings;
}"""

SIZE_SCRIPT = """(svg, size) => {
  svg.style.width = `${size.width}px`;
  svg.style.minWidth = `${size.width}px`;
  svg.style.height = `${size.height}px`;
  for (let a = svg.parentElement; a; a = a.parentElement) a.style.setProperty('overflow', 'visible', 'important');
}"""


def render_pngs(pages: Sequence[Path], out_dir: Path) -> list[str]:
    """Write one PNG per diagram SVG into out_dir. Return findings; write no PNG for a page with findings."""
    findings: list[str] = []
    with sync_playwright() as playwright:
        try:
            browser = playwright.chromium.launch()
        except PlaywrightError as error:
            first_line = str(error).strip().splitlines()[0]
            return [
                f"Chromium did not start ({first_line}). Install it with: "
                f"uv run --with playwright=={version('playwright')} playwright install chromium"
            ]
        try:
            for page_path in pages:
                findings.extend(_render_page(browser, page_path, out_dir))
        finally:
            browser.close()
    return findings


def _render_page(browser: Browser, page_path: Path, out_dir: Path) -> list[str]:
    page = browser.new_page(device_scale_factor=SCALE, viewport=VIEWPORT)
    try:
        page.goto(page_path.resolve().as_uri())
        page.wait_for_load_state("networkidle")
        page.evaluate("document.fonts.ready")
        missing = [font for font in FONT_CHECKS if not page.evaluate("font => document.fonts.check(font)", font)]
        if missing:
            return [f"{page_path.name}: fonts {', '.join(missing)} did not load; the build needs Google Fonts access"]
        svgs = page.locator(DIAGRAM_SELECTOR).all()
        findings = [finding for svg in svgs for finding in svg.evaluate(TEXT_FIT_SCRIPT)]
        if findings:
            return findings
        for svg in svgs:
            slug = svg.get_attribute("aria-labelledby").split()[0].removesuffix("-title")
            width, height = map(int, svg.get_attribute("viewBox").split()[2:])
            svg.evaluate(SIZE_SCRIPT, {"width": width, "height": height})
            out = out_dir / f"{slug}.png"
            svg.screenshot(path=str(out))
            expected = (width * SCALE, height * SCALE)
            if _png_size(out) != expected:
                findings.append(f"{out.name}: size {_png_size(out)} is not the expected {expected}")
            else:
                print(f"wrote {out.name}  {expected[0]}x{expected[1]}")
        return findings
    finally:
        page.close()


def _png_size(path: Path) -> tuple[int, int]:
    with path.open("rb") as png:
        header = png.read(24)
    return struct.unpack(">II", header[16:24])
