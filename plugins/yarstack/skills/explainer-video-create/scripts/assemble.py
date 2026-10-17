"""Inspect and join rendered scenes in narration order.

Usage:
    python assemble.py sheets narration.json VIDEO_DIR OUT_DIR
    python assemble.py concat narration.json VIDEO_DIR OUTPUT.mp4

`sheets` writes one 4x2 contact sheet per scene for visual review.
`concat` converts every scene to the first scene's size and frame rate, with
stereo 48 kHz audio padded to the scene's video length, and fades each scene in
and out through black so that the cut between an animation and a screen
recording is not abrupt. The fades keep every scene's length. Without the padding,
scenes that end in silence shift all later narration earlier. Without the
conversion, a scene from another source, such as a screen recording, can join
with silent audio.
"""

import argparse
import json
import subprocess
import sys
import tempfile
from pathlib import Path

MAX_DRIFT_SECONDS = 0.1
FADE_SECONDS = 0.4


def ffmpeg(*args):
    subprocess.run(["ffmpeg", "-loglevel", "error", "-y", *args], check=True)


def stream_seconds(path, stream):
    out = subprocess.run(
        ["ffprobe", "-v", "error", "-select_streams", stream, "-show_entries", "stream=duration", "-of", "csv=p=0", str(path)],
        capture_output=True, text=True, check=True,
    )
    first = out.stdout.split()[0] if out.stdout.split() else "0"
    return float(first.strip(",") or 0)


def video_format(path):
    out = subprocess.run(
        ["ffprobe", "-v", "error", "-select_streams", "v", "-show_entries", "stream=width,height,r_frame_rate", "-of", "csv=p=0", str(path)],
        capture_output=True, text=True, check=True,
    )
    width, height, rate = out.stdout.strip().split(",")
    return int(width), int(height), rate


def scene_videos(narration, video_dir):
    videos = [Path(video_dir) / f"{scene}.mp4" for scene in json.loads(Path(narration).read_text())]
    missing = [str(v) for v in videos if not v.exists()]
    if missing:
        sys.exit(f"missing rendered scenes: {', '.join(missing)}")
    return videos


def sheets(videos, out_dir):
    out = Path(out_dir)
    out.mkdir(parents=True, exist_ok=True)
    for video in videos:
        seconds = stream_seconds(video, "v")
        ffmpeg("-i", str(video), "-vf", f"fps=8/{seconds},scale=480:-1,tile=4x2", "-frames:v", "1", str(out / f"{video.stem}.png"))
        print(out / f"{video.stem}.png")


def concat(videos, output):
    width, height, rate = video_format(videos[0])
    fit = f"scale={width}:{height}:force_original_aspect_ratio=decrease,pad={width}:{height}:(ow-iw)/2:(oh-ih)/2,fps={rate},format=yuv420p"
    with tempfile.TemporaryDirectory() as tmp:
        listing = Path(tmp) / "list.txt"
        lines = []
        for video in videos:
            padded = Path(tmp) / video.name
            seconds = stream_seconds(video, "v")
            fades = f"fade=t=in:st=0:d={FADE_SECONDS},fade=t=out:st={seconds - FADE_SECONDS:.3f}:d={FADE_SECONDS}"
            ffmpeg("-i", str(video), "-vf", f"{fit},{fades}", "-af", f"apad,atrim=duration={seconds:.3f}", "-c:v", "libx264",
                   "-crf", "18", "-c:a", "aac", "-b:a", "192k", "-ar", "48000", "-ac", "2", str(padded))
            lines.append(f"file '{padded}'")
        listing.write_text("\n".join(lines) + "\n")
        ffmpeg("-f", "concat", "-safe", "0", "-i", str(listing), "-c", "copy", "-movflags", "+faststart", str(output))

    video_s, audio_s = stream_seconds(output, "v"), stream_seconds(output, "a")
    print(f"{output}: video {video_s:.2f}s, audio {audio_s:.2f}s")
    if abs(video_s - audio_s) > MAX_DRIFT_SECONDS:
        sys.exit(f"audio and video lengths differ by {abs(video_s - audio_s):.2f}s")


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("command", choices=["sheets", "concat"])
    parser.add_argument("narration")
    parser.add_argument("video_dir")
    parser.add_argument("target")
    args = parser.parse_args()

    videos = scene_videos(args.narration, args.video_dir)
    if args.command == "sheets":
        sheets(videos, args.target)
    else:
        concat(videos, args.target)


if __name__ == "__main__":
    main()
