"""Generate one audio clip per narration beat and record clip durations.

Usage:
    python narrate.py narration.json --provider elevenlabs|say|espeak|silent
                      [--language CODE] [--env-file PATH] [--voice ID] [--speed 1.0]
                      [--out audio]
    python narrate.py --provider elevenlabs --list-voices [--env-file PATH]

narration.json maps each scene class name to an ordered list of beat texts.
ElevenLabs reads each scene in one take with a long pause between beats, and
the script cuts the take into one clip per beat inside those pauses. Separate
requests per beat give each beat its own pace, which sounds like random speed
changes. Clip names include a hash of the scene's beats, provider, model,
voice, and speed, so editing a beat regenerates only its scene. --speed changes
the tempo of each clip after synthesis without changing its pitch, because
Eleven v4 ignores the speed voice setting. The script writes
<out>/durations.json, which scenekit.py reads to time each beat. With
elevenlabs, --language is required and the script stops when the model does not
support that language.
"""

import argparse
import base64
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path
from urllib import error, request

ELEVENLABS_URL = "https://api.elevenlabs.io/v1"
DEFAULT_VOICES = {"elevenlabs": "7EzWGsX10sAS4c9m9cPf", "say": "Daniel", "espeak": "en-us", "silent": ""}
MODEL_ID = "eleven_v4"
SILENT_WORDS_PER_SECOND = 2.6
BEAT_SEPARATOR = " [long pause] "
MIN_PAUSE_SECONDS = 0.1
PAUSE_SEARCH_SECONDS = 0.5
EDGE_SILENCE_SECONDS = 0.15
LEADING_TAGS = re.compile(r"(\s*\[[^\]]*\])*\s*")


def read_env_key(env_file):
    key = os.environ.get("ELEVENLABS_API_KEY", "").strip()
    if key or not env_file:
        return key
    for line in Path(env_file).read_text().splitlines():
        name, _, value = line.removeprefix("export ").partition("=")
        if name.strip() == "ELEVENLABS_API_KEY":
            return value.strip().strip("\"'")
    return ""


def elevenlabs_call(key, path, body=None, allow_missing_permission=False):
    data = json.dumps(body).encode() if body is not None else None
    req = request.Request(f"{ELEVENLABS_URL}{path}", data=data, headers={"xi-api-key": key, "Content-Type": "application/json"})
    try:
        with request.urlopen(req, timeout=120) as response:
            return response.read()
    except error.HTTPError as exc:
        detail = exc.read().decode(errors="replace")[:300]
        if allow_missing_permission and "missing_permissions" in detail:
            return None
        sys.exit(f"ElevenLabs {exc.code} on {path.split('?')[0]}: {detail}")


def check_elevenlabs_key(key):
    quota = elevenlabs_call(key, "/user/subscription", allow_missing_permission=True)
    if quota is None:
        print("ElevenLabs key accepted; it cannot read the quota, so the first clip confirms text-to-speech access")
        return
    quota = json.loads(quota)
    remaining = quota.get("character_limit", 0) - quota.get("character_count", 0)
    print(f"ElevenLabs key valid; {remaining} characters left")


def check_model_language(key, language):
    models = elevenlabs_call(key, "/models", allow_missing_permission=True)
    if models is None:
        print(f"ElevenLabs key cannot read the model list, so support for language {language} in {MODEL_ID} is unverified")
        return
    languages = {model["model_id"]: {lang["language_id"] for lang in model.get("languages", [])} for model in json.loads(models)}
    if language not in languages.get(MODEL_ID, set()):
        others = sorted(model for model, codes in languages.items() if language in codes)
        sys.exit(f"{MODEL_ID} does not support language {language}; models that do: {', '.join(others) or 'none'}")
    print(f"{MODEL_ID} supports language {language}")


def list_elevenlabs_voices(key):
    for voice in json.loads(elevenlabs_call(key, "/voices"))["voices"]:
        print(f"{voice['voice_id']}  {voice['name']}  ({voice.get('category', '')})")


def run(*args):
    subprocess.run(args, check=True, capture_output=True)


def to_mp3(source, target):
    run("ffmpeg", "-y", "-i", str(source), "-ar", "44100", "-b:a", "128k", str(target))
    source.unlink()


def silences(path):
    out = subprocess.run(
        ["ffmpeg", "-i", str(path), "-af", f"silencedetect=n=-40dB:d={MIN_PAUSE_SECONDS}", "-f", "null", "-"],
        capture_output=True, text=True, check=True,
    ).stderr
    starts = [float(s) for s in re.findall(r"silence_start: ([0-9.]+)", out)]
    ends = [float(s) for s in re.findall(r"silence_end: ([0-9.]+)", out)]
    return list(zip(starts, ends))


def boundary_pause(window_start, window_end, pauses):
    near = [(start, end) for start, end in pauses if start <= window_end and end >= window_start]
    if not near:
        sys.exit(f"no pause in the narration between {window_start:.2f} and {window_end:.2f} seconds")
    return max(near, key=lambda pause: pause[1] - pause[0])


def elevenlabs_scene(voice, key, beats, paths):
    text = BEAT_SEPARATOR.join(beats)
    body = {
        "text": text,
        "model_id": MODEL_ID,
        "voice_settings": {"stability": 1.0, "similarity_boost": 0.75, "style": 0.0},
    }
    take = json.loads(elevenlabs_call(key, f"/text-to-speech/{voice}/with-timestamps?output_format=mp3_44100_128", body))
    starts = take["alignment"]["character_start_times_seconds"]
    ends = take["alignment"]["character_end_times_seconds"]
    if len(starts) != len(text):
        sys.exit(f"ElevenLabs aligned {len(starts)} characters for {len(text)} characters of text")
    whole = paths[0].with_suffix(".take.mp3")
    whole.write_bytes(base64.b64decode(take["audio_base64"]))
    pauses = silences(whole)
    spans, clip_start, offset = [], 0.0, 0
    for i in range(1, len(beats)):
        offset += len(beats[i - 1])
        next_beat = offset + len(BEAT_SEPARATOR)
        first = next_beat + LEADING_TAGS.match(beats[i]).end()
        pause_start, pause_end = boundary_pause(ends[offset - 1] - PAUSE_SEARCH_SECONDS, starts[first] + PAUSE_SEARCH_SECONDS, pauses)
        spans.append(["-ss", f"{clip_start:.3f}", "-to", f"{pause_start + EDGE_SILENCE_SECONDS:.3f}"])
        clip_start = pause_end - EDGE_SILENCE_SECONDS
        offset = next_beat
    spans.append(["-ss", f"{clip_start:.3f}"])
    for span, path in zip(spans, paths):
        run("ffmpeg", "-y", "-i", str(whole), *span, "-b:a", "128k", str(path))
    whole.unlink()


def local_beat(provider, voice, text, target):
    if provider == "say":
        aiff = target.with_suffix(".aiff")
        run("say", "-v", voice, "-r", "175", "-o", str(aiff), text)
        to_mp3(aiff, target)
    elif provider == "espeak":
        wav = target.with_suffix(".wav")
        run("espeak-ng", "-v", voice, "-s", "160", "-w", str(wav), text)
        to_mp3(wav, target)
    else:
        seconds = max(1.5, len(text.split()) / SILENT_WORDS_PER_SECOND)
        run("ffmpeg", "-y", "-f", "lavfi", "-i", "anullsrc=r=44100:cl=mono", "-t", f"{seconds:.2f}", "-b:a", "128k", str(target))


def change_tempo(path, speed):
    faster = path.with_suffix(".tempo.mp3")
    run("ffmpeg", "-y", "-i", str(path), "-filter:a", f"atempo={speed}", "-b:a", "128k", str(faster))
    faster.replace(path)


def clip_seconds(path):
    out = subprocess.run(
        ["ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", str(path)],
        capture_output=True, text=True, check=True,
    )
    return round(float(out.stdout.strip()), 3)


def require_tools(provider):
    needed = {"ffmpeg", "ffprobe"} | {"say": {"say"}, "espeak": {"espeak-ng"}}.get(provider, set())
    missing = sorted(tool for tool in needed if not shutil.which(tool))
    if missing:
        sys.exit(f"missing tools for provider {provider}: {', '.join(missing)}")


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("narration", nargs="?")
    parser.add_argument("--provider", required=True, choices=sorted(DEFAULT_VOICES))
    parser.add_argument("--language", help="narration language as an ISO 639-1 code, such as en or lv")
    parser.add_argument("--env-file")
    parser.add_argument("--voice")
    parser.add_argument("--speed", type=float, default=1.0, help="tempo factor from 0.5 to 2.0, applied after synthesis")
    parser.add_argument("--out", default="audio")
    parser.add_argument("--check-key", action="store_true", help="validate the ElevenLabs key and language, then exit")
    parser.add_argument("--list-voices", action="store_true", help="print the ElevenLabs voices this key can use and exit")
    args = parser.parse_args()

    if not 0.5 <= args.speed <= 2.0:
        parser.error("--speed must be from 0.5 to 2.0")
    require_tools(args.provider)
    voice = args.voice or DEFAULT_VOICES[args.provider]
    key = read_env_key(args.env_file) if args.provider == "elevenlabs" else ""
    if args.provider == "elevenlabs":
        if not key:
            sys.exit("ELEVENLABS_API_KEY not found in the environment or --env-file")
        if args.list_voices:
            list_elevenlabs_voices(key)
            return
        if not args.language:
            sys.exit("--language is required with --provider elevenlabs")
        check_elevenlabs_key(key)
        check_model_language(key, args.language)
        if args.check_key:
            return
    if not args.narration:
        parser.error("narration is required")
    engine = f"elevenlabs/{MODEL_ID}" if args.provider == "elevenlabs" else args.provider

    narration = json.loads(Path(args.narration).read_text())
    out = Path(args.out)
    out.mkdir(exist_ok=True)
    durations = {}
    for scene, beats in narration.items():
        speed = "" if args.speed == 1.0 else f"|{args.speed}"
        digest = hashlib.sha256(f"{engine}|{voice}{speed}|{json.dumps(beats)}".encode()).hexdigest()[:10]
        paths = [out / f"{scene}_{i}_{digest}.mp3" for i in range(len(beats))]
        if not all(path.exists() for path in paths):
            if args.provider == "elevenlabs":
                elevenlabs_scene(voice, key, beats, paths)
            else:
                for text, path in zip(beats, paths):
                    local_beat(args.provider, voice, text, path)
            for path in paths:
                if args.speed != 1.0:
                    change_tempo(path, args.speed)
            print("wrote", scene, flush=True)
        durations[scene] = [{"file": path.name, "seconds": clip_seconds(path)} for path in paths]
        words = sum(len(re.sub(r"\[[^\]]*\]", " ", text).split()) for text in beats)
        print(f"{scene}: {words / sum(beat['seconds'] for beat in durations[scene]):.2f} words per second")

    (out / "durations.json").write_text(json.dumps(durations, indent=2))
    total = sum(beat["seconds"] for beats in durations.values() for beat in beats)
    print(f"{sum(len(b) for b in durations.values())} beats, {total:.1f} seconds of narration")


if __name__ == "__main__":
    main()
