"""Join a recorded walkthrough with its narration and inspect the result.

Usage:
    python walkthrough.py mux RECORDING AUDIO_DIR BEAT_STARTS OUTPUT.mp4
    python walkthrough.py sheet VIDEO OUTPUT.png

AUDIO_DIR holds the clips and durations.json that narrate.py writes for a
narration file with one scene named Walkthrough. BEAT_STARTS is a JSON list with
the start of each beat in seconds from the start of the recording.

`mux` places each clip at its beat start, trims the recording to the narrated
span, and fails when clips overlap, the recording ends before the narration, or
the output audio and video lengths differ. `sheet` writes a 4x2 contact sheet.
"""

import argparse
import json
import subprocess
import sys
import tempfile
from pathlib import Path

SCENE = "Walkthrough"
LEAD_SECONDS = 1.0
TAIL_SECONDS = 1.5
MAX_DRIFT_SECONDS = 0.1


def ffmpeg(*args):
    subprocess.run(["ffmpeg", "-loglevel", "error", "-y", *args], check=True)


def stream_seconds(path, stream):
    out = subprocess.run(
        ["ffprobe", "-v", "error", "-select_streams", stream, "-show_entries", "stream=duration", "-of", "csv=p=0", str(path)],
        capture_output=True, text=True, check=True,
    )
    return float(out.stdout.strip() or 0)


def narrated_beats(audio_dir, beat_starts):
    clips = json.loads((Path(audio_dir) / "durations.json").read_text())[SCENE]
    starts = json.loads(Path(beat_starts).read_text())
    if len(starts) != len(clips):
        sys.exit(f"{len(starts)} beat starts for {len(clips)} narration clips")
    for i in range(1, len(starts)):
        previous_end = starts[i - 1] + clips[i - 1]["seconds"]
        if starts[i] < previous_end:
            sys.exit(f"beat {i} starts at {starts[i]:.2f}s, before beat {i - 1} narration ends at {previous_end:.2f}s")
    return [(start, Path(audio_dir) / clip["file"], clip["seconds"]) for start, clip in zip(starts, clips)]


def mux(recording, audio_dir, beat_starts, output):
    beats = narrated_beats(audio_dir, beat_starts)
    with tempfile.TemporaryDirectory() as tmp:
        video = Path(tmp) / "recording.mp4"
        ffmpeg("-i", str(recording), "-an", "-c:v", "libx264", "-pix_fmt", "yuv420p", str(video))
        narration_end = beats[-1][0] + beats[-1][2]
        recorded = stream_seconds(video, "v")
        if recorded < narration_end:
            sys.exit(f"recording ends at {recorded:.2f}s, before the narration ends at {narration_end:.2f}s")

        start = max(0.0, beats[0][0] - LEAD_SECONDS)
        length = min(recorded, narration_end + TAIL_SECONDS) - start
        inputs, delays = [], []
        for k, (beat_start, clip, _) in enumerate(beats, start=1):
            inputs += ["-i", str(clip)]
            delays.append(f"[{k}:a]adelay={round((beat_start - start) * 1000)}:all=1[a{k}]")
        mix = "".join(f"[a{k}]" for k in range(1, len(beats) + 1))
        audio = ";".join(delays) + f";{mix}amix=inputs={len(beats)}:normalize=0,apad,atrim=0:{length:.3f}[aout]"
        ffmpeg(
            "-ss", f"{start:.3f}", "-t", f"{length:.3f}", "-i", str(video), *inputs,
            "-filter_complex", audio, "-map", "0:v", "-map", "[aout]",
            "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", "-ar", "44100",
            "-movflags", "+faststart", str(output),
        )

    video_s, audio_s = stream_seconds(output, "v"), stream_seconds(output, "a")
    print(f"{output}: video {video_s:.2f}s, audio {audio_s:.2f}s")
    if abs(video_s - audio_s) > MAX_DRIFT_SECONDS:
        sys.exit(f"audio and video lengths differ by {abs(video_s - audio_s):.2f}s")


def sheet(video, output):
    seconds = stream_seconds(video, "v")
    ffmpeg("-i", str(video), "-vf", f"fps=8/{seconds},scale=480:-1,tile=4x2", "-frames:v", "1", str(output))
    print(output)


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    commands = parser.add_subparsers(dest="command", required=True)
    mux_args = commands.add_parser("mux")
    for name in ("recording", "audio_dir", "beat_starts", "output"):
        mux_args.add_argument(name)
    sheet_args = commands.add_parser("sheet")
    for name in ("video", "output"):
        sheet_args.add_argument(name)
    args = parser.parse_args()

    if args.command == "mux":
        mux(args.recording, args.audio_dir, args.beat_starts, args.output)
    else:
        sheet(args.video, args.output)


if __name__ == "__main__":
    main()
