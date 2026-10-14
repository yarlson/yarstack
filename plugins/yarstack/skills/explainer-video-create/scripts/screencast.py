"""Turn one screen recording into a narrated scene.

Usage:
    python screencast.py SCENE RECORDING STEP_ENDS AUDIO_DIR OUTPUT.mp4

SCENE is the narration.json key, and its beats describe the recorded steps in
order. STEP_ENDS is a JSON list with the end of each step in seconds from the
start of the recording, one entry for each beat. Step i plays from the end of
step i-1 to its own end, then holds its last frame until beat i's clip ends.
The recording therefore does not depend on the voice: a new narration changes
only how long each step holds.
"""

import argparse
import json
import subprocess
import sys
import tempfile
from pathlib import Path

BEAT_PAD = 0.45


def ffmpeg(*args):
    subprocess.run(["ffmpeg", "-loglevel", "error", "-y", *args], check=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("scene")
    parser.add_argument("recording")
    parser.add_argument("step_ends")
    parser.add_argument("audio_dir")
    parser.add_argument("output")
    args = parser.parse_args()

    audio = Path(args.audio_dir)
    clips = json.loads((audio / "durations.json").read_text())[args.scene]
    ends = json.loads(Path(args.step_ends).read_text())
    if len(ends) != len(clips):
        sys.exit(f"{len(ends)} step ends for {len(clips)} beats of {args.scene}")
    if any(b <= a for a, b in zip([0.0, *ends], ends)):
        sys.exit("step ends must increase")

    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory() as tmp:
        parts = []
        for i, (start, end, clip) in enumerate(zip([0.0, *ends], ends, clips)):
            seconds = max(end - start, clip["seconds"] + BEAT_PAD)
            part = Path(tmp) / f"step{i}.mp4"
            ffmpeg(
                "-ss", f"{start:.3f}", "-to", f"{end:.3f}", "-i", args.recording, "-i", str(audio / clip["file"]),
                "-filter_complex",
                f"[0:v]tpad=stop_mode=clone:stop_duration={seconds},trim=duration={seconds},setpts=PTS-STARTPTS,format=yuv420p[v];"
                f"[1:a]apad,atrim=duration={seconds},aresample=48000,aformat=channel_layouts=stereo[a]",
                "-map", "[v]", "-map", "[a]", "-c:v", "libx264", "-crf", "16", "-c:a", "aac", "-b:a", "192k", str(part),
            )
            parts.append(part)
        listing = Path(tmp) / "list.txt"
        listing.write_text("".join(f"file '{part}'\n" for part in parts))
        ffmpeg("-f", "concat", "-safe", "0", "-i", str(listing), "-c", "copy", str(output))
    print(output)


if __name__ == "__main__":
    main()
