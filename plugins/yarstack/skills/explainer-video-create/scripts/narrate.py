"""Generate one audio clip per narration beat and record clip durations.

Usage:
    python narrate.py narration.json --provider elevenlabs|say|espeak|silent
                      [--language CODE] [--env-file PATH] [--voice ID] [--out audio]

narration.json maps each scene class name to an ordered list of beat texts.
Clip names include a hash of the text, provider, model, and voice, so editing a
beat regenerates only that clip. The script writes <out>/durations.json, which
scenekit.py reads to time each beat. With elevenlabs, --language is required and
the script stops when the model does not support that language.
"""

import argparse
import hashlib
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path
from urllib import error, request

ELEVENLABS_URL = "https://api.elevenlabs.io/v1"
DEFAULT_VOICES = {"elevenlabs": "JBFqnCBsd6RMkjVDRZzb", "say": "Daniel", "espeak": "en-us", "silent": ""}
MODEL_ID = "eleven_v4"
SILENT_WORDS_PER_SECOND = 2.6


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


def run(*args):
    subprocess.run(args, check=True, capture_output=True)


def to_mp3(source, target):
    run("ffmpeg", "-y", "-i", str(source), "-ar", "44100", "-b:a", "128k", str(target))
    source.unlink()


def synthesize(provider, voice, key, beats, i, target):
    text = beats[i]
    if provider == "elevenlabs":
        body = {
            "text": text,
            "model_id": MODEL_ID,
            "previous_text": beats[i - 1] if i > 0 else None,
            "next_text": beats[i + 1] if i + 1 < len(beats) else None,
            "voice_settings": {"stability": 0.55, "similarity_boost": 0.75, "style": 0.15},
        }
        target.write_bytes(elevenlabs_call(key, f"/text-to-speech/{voice}?output_format=mp3_44100_128", body))
    elif provider == "say":
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
    parser.add_argument("narration")
    parser.add_argument("--provider", required=True, choices=sorted(DEFAULT_VOICES))
    parser.add_argument("--language", help="narration language as an ISO 639-1 code, such as en or lv")
    parser.add_argument("--env-file")
    parser.add_argument("--voice")
    parser.add_argument("--out", default="audio")
    parser.add_argument("--check-key", action="store_true", help="validate the ElevenLabs key and language, then exit")
    args = parser.parse_args()

    require_tools(args.provider)
    voice = args.voice or DEFAULT_VOICES[args.provider]
    key = read_env_key(args.env_file) if args.provider == "elevenlabs" else ""
    if args.provider == "elevenlabs":
        if not args.language:
            sys.exit("--language is required with --provider elevenlabs")
        if not key:
            sys.exit("ELEVENLABS_API_KEY not found in the environment or --env-file")
        check_elevenlabs_key(key)
        check_model_language(key, args.language)
        if args.check_key:
            return
    engine = f"elevenlabs/{MODEL_ID}" if args.provider == "elevenlabs" else args.provider

    narration = json.loads(Path(args.narration).read_text())
    out = Path(args.out)
    out.mkdir(exist_ok=True)
    durations = {}
    for scene, beats in narration.items():
        durations[scene] = []
        for i, text in enumerate(beats):
            digest = hashlib.sha256(f"{engine}|{voice}|{text}".encode()).hexdigest()[:10]
            path = out / f"{scene}_{i}_{digest}.mp3"
            if not path.exists():
                synthesize(args.provider, voice, key, beats, i, path)
                print("wrote", path.name, flush=True)
            durations[scene].append({"file": path.name, "seconds": clip_seconds(path)})

    (out / "durations.json").write_text(json.dumps(durations, indent=2))
    total = sum(beat["seconds"] for beats in durations.values() for beat in beats)
    print(f"{sum(len(b) for b in durations.values())} beats, {total:.1f} seconds of narration")


if __name__ == "__main__":
    main()
