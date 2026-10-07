"""Shared manim helpers for narrated explainer scenes.

Copy this file next to the scene file. Each scene subclasses Narrated and wraps
the animations for beat i in `with self.beat(i):`. The beat plays its narration
clip and waits until the clip ends, so keep each beat's animations shorter than
its clip.
"""

import json
from contextlib import contextmanager
from pathlib import Path

from manim import ORIGIN, RIGHT, WHITE, AnimationGroup, Arrow, FadeOut, RoundedRectangle, Scene, Text, VGroup, config

ROOT = Path(__file__).parent
AUDIO = ROOT / "audio"
BEAT_PAD = 0.45

BACKGROUND = "#0e1116"
PRIMARY = "#58C4DD"
ACCENT = "#F4D345"
GOOD = "#83C167"
BAD = "#FC6255"
NOTE = "#5CD0B3"
WARN = "#FF862F"
DIM = "#BBBBBB"
FONT = "Avenir Next"
MONO = "Menlo"

config.background_color = BACKGROUND


def txt(s, size=28, color=WHITE, **kw):
    return Text(s, font=FONT, font_size=size, color=color, **kw)


def mono(s, size=22, color=WHITE):
    return Text(s, font=MONO, font_size=size, color=color)


def card(label, color=NOTE, w=1.25, h=0.62, size=24):
    box = RoundedRectangle(corner_radius=0.12, width=w, height=h, color=color, stroke_width=3)
    box.set_fill(color, opacity=0.12)
    return VGroup(box, txt(label, size, color)).arrange(ORIGIN)


def flow(names, width=2.25, gap=0.22, color=DIM, size=20):
    boxes = VGroup(*[card(name, color, w=width, h=0.7, size=size) for name in names]).arrange(RIGHT, buff=gap)
    arrows = VGroup(*[
        Arrow(boxes[i].get_right(), boxes[i + 1].get_left(), buff=0.03, stroke_width=3, color=color,
              max_tip_length_to_length_ratio=0.5)
        for i in range(len(boxes) - 1)
    ])
    return VGroup(boxes, arrows)


def light_up(stage, color=GOOD):
    box, label = stage
    return AnimationGroup(box.animate.set_color(color).set_fill(color, opacity=0.15), label.animate.set_color(color))


class Narrated(Scene):
    @contextmanager
    def beat(self, i):
        name = type(self).__name__
        clip = json.loads((AUDIO / "durations.json").read_text())[name][i]
        self.add_sound(str(AUDIO / clip["file"]))
        start = self.renderer.time
        yield
        remaining = clip["seconds"] + BEAT_PAD - (self.renderer.time - start)
        if remaining > 0:
            self.wait(remaining)

    def clear_all(self, run_time=0.6):
        if self.mobjects:
            self.play(*[FadeOut(m) for m in self.mobjects], run_time=run_time)

