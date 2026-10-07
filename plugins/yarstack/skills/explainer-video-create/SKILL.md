---
name: explainer-video-create
description: Create a narrated 3Blue1Brown-style animated explainer video about how an existing project feature works, using manim with ElevenLabs or local text-to-speech. Use when the user asks for an explainer, walkthrough, or teaching video about code, architecture, or system behavior; not for promotional clips or slide decks.
---

# Explainer Video Creation

Produce one narrated MP4 that teaches how a feature works. Every narrated claim must come from investigated evidence.

## Workflow

### 1. Set the scope

Extract the project, feature, focus, audience, and output path from the request. Default to engineers who know the domain, a 3 to 5 minute length, and `~/Desktop/<feature>-explainer.mp4`. Ask only when the feature or project cannot be identified.

### 2. Investigate the feature

Run `system-investigate` on the feature in read-only mode. Record a short fact sheet in the scratch directory:

- each mechanism, state change, rule, and failure path the video will show, with a `file:line` reference;
- one concrete example from repository history, such as a real commit, PR, or log line;
- facts that the evidence does not settle.

Narrate only facts from the sheet. Leave out inconclusive facts or present them as open questions.

### 3. Choose the voice

Select the first available provider:

1. `elevenlabs`: `ELEVENLABS_API_KEY` is in the environment, or in a `.env` file that the user named or that exists in the working directory or project root. Validate it with `narrate.py --provider elevenlabs --check-key --env-file <path>` before you write narration that depends on it.
2. `say` on macOS, or `espeak` when `espeak-ng` is installed.
3. `silent`, which only sets the timing.

If no key is found, or the key fails validation, ask the user once. Offer three choices: give a key, use the best local voice, or make a silent video. Continue with the local voice only if the user agrees. Never print, log, or copy a key. ElevenLabs receives the narration text, so keep secrets, customer data, and internal hostnames out of the narration.

### 4. Prepare tools

Require `ffmpeg` and `ffprobe`. Create an isolated environment in the scratch directory, with Python 3.12 for manim wheel support. Prefer `uv venv --python 3.12` followed by `uv pip install manim`. Ask before you make any system-wide installation. Use `Text` objects so that LaTeX is not needed.

### 5. Write the storyboard and narration

Plan 5 to 8 scenes. Give each scene one idea and split it into beats of 1 to 3 spoken sentences. A beat is the unit of audio-to-picture sync. A typical arc:

1. Problem: show the failure that the feature prevents.
2. Core rule: state the key invariant in one sentence.
3. One scene for each mechanism.
4. Failure and recovery paths.
5. Recap: the invariant again, now earned.

Follow 3Blue1Brown style. Use a dark background and one visual metaphor for each idea. Build the picture step by step and use motion to show change. Give each color one meaning for the whole video. Keep text to short labels. Write narration as plain spoken English. Spell out identifiers that are hard to say, and leave file paths and long hashes on screen only.

Save the narration as `narration.json` in scene order: `{"S1Problem": ["beat", ...], ...}`.

### 6. Generate the audio

Copy `scripts/scenekit.py` from this skill's directory into the scratch directory. Run `scripts/narrate.py narration.json --provider <provider>` there. It writes one clip per beat and `audio/durations.json`. Clip names hash the text, so an edited beat regenerates only its own clip.

### 7. Write the scenes

Write `scenes.py` with `from manim import *` and `from scenekit import *`. Give each scene class the same name as its `narration.json` key, subclass `Narrated`, and wrap the animations for beat `i` in `with self.beat(i):`. Keep the animations in each beat shorter than the clip, because the beat waits until the clip ends. Call `self.clear_all()` between visual sections.

### 8. Render a draft and inspect it

Render with `manim -ql --disable_caching scenes.py <scene names>`. Then run `scripts/assemble.py sheets narration.json media/videos/scenes/480p15 sheets` and view every sheet. Fix text cut off at the frame edge, overlapping labels, unreadable sizes, and objects that remain from an earlier beat. Render again only the scenes you changed.

### 9. Render the final video

Render with `manim -qh --fps 30 --disable_caching`. Then run `scripts/assemble.py concat narration.json media/videos/scenes/1080p30 <output>`. The script pads each scene's audio before it joins the scenes, and it fails when the audio and video lengths differ.

## Result

Report the output path, the length, the voice provider, the scenes and the claims each one makes with its `file:line` evidence, the facts left out, and the scratch path of `narration.json` and `scenes.py` for later edits. Do not publish or upload the video unless the user asks.
