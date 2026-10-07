"""Inspect and join rendered scenes in narration order.

Usage:
    python assemble.py sheets narration.json VIDEO_DIR OUT_DIR
    python assemble.py concat narration.json VIDEO_DIR OUTPUT.mp4

`sheets` writes one 4x2 contact sheet per scene for visual review.
`concat` pads each scene's audio to its video length before joining. Without the
padding, scenes that end in silence shift all later narration earlier.
"""

import argparse
import json
import subprocess
import sys
import tempfile
from pathlib import Path

MAX_DRIFT_SECONDS = 0.1


def ffmpeg(*args):
    subprocess.run(["ffmpeg", "-loglevel", "error", "-y", *args], check=True)


def stream_seconds(path, stream):
    out = subprocess.run(
        ["ffprobe", "-v", "error", "-select_streams", stream, "-show_entries", "stream=duration", "-of", "csv=p=0", str(path)],
        capture_output=True, text=True, check=True,
    )
    return float(out.stdout.strip() or 0)


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
    with tempfile.TemporaryDirectory() as tmp:
        listing = Path(tmp) / "list.txt"
        lines = []
        for video in videos:
            padded = Path(tmp) / video.name
            ffmpeg("-i", str(video), "-af", "apad", "-shortest", "-c:v", "copy", "-c:a", "aac", "-b:a", "192k", "-ar", "44100", str(padded))
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
