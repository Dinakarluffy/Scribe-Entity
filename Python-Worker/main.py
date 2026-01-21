#!/usr/bin/env python3
# python_worker_full_corrected.py

import os
import re
import json
import uuid
import logging
import mimetypes
import tempfile
import argparse
import sys
from datetime import datetime
from typing import Dict, Any, List

from dotenv import load_dotenv

# Optional imports
try:
    import whisper
except Exception:
    whisper = None

try:
    from moviepy.video.io.VideoFileClip import VideoFileClip
except Exception:
    VideoFileClip = None

try:
    import google.genai as genai
except Exception:
    genai = None

try:
    import psycopg2
except Exception:
    psycopg2 = None

try:
    from jsonschema import validate as json_validate, ValidationError
except Exception:
    json_validate = None
    ValidationError = Exception


# -------------------------------------------------
# ENV
# -------------------------------------------------
load_dotenv()

GEMINI_API_KEY = os.getenv("GEMINI_API_KEY", "").strip()
DB_HOST = os.getenv("DB_HOST", "").strip()
DB_PORT = os.getenv("DB_PORT", "").strip()
DB_USER = os.getenv("DB_USER", "").strip()
DB_PASSWORD = os.getenv("DB_PASSWORD", "").strip()
DB_NAME = os.getenv("DB_NAME", "").strip()
DB_TABLE = os.getenv("DB_TABLE", "scribe_entity_classification_dev").strip()
RESULTS_DIR = os.getenv("RESULTS_DIR", "results").strip()
MODEL_NAME = os.getenv("GEMINI_MODEL", "gemini-2.5-flash").strip()
MAX_CHUNK_CHARS = int(os.getenv("MAX_CHUNK_CHARS", "16000"))


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(message)s",
    stream=sys.stderr,
)
logger = logging.getLogger(__name__)


# -------------------------------------------------
# UTILITIES
# -------------------------------------------------
EXPECTED_SCHEMA = {
    "type": "object",
    "required": ["entities", "tone", "style", "safety_flags"]
}


def extract_json_substring(text: str) -> str:
    match = re.search(r"(\{[\s\S]*\})", text)
    return match.group(1) if match else text


def validate_schema(doc: Dict[str, Any]) -> bool:
    if json_validate:
        try:
            json_validate(instance=doc, schema=EXPECTED_SCHEMA)
            return True
        except ValidationError:
            return False
    return all(k in doc for k in EXPECTED_SCHEMA["required"])


def chunk_transcript(text: str) -> List[str]:
    if len(text) <= MAX_CHUNK_CHARS:
        return [text]
    chunks = []
    start = 0
    while start < len(text):
        end = min(start + MAX_CHUNK_CHARS, len(text))
        split = text.rfind(" ", start, end)
        split = split if split > start else end
        chunks.append(text[start:split])
        start = split
    return chunks


# -------------------------------------------------
# TRANSCRIPTION
# -------------------------------------------------
def extract_transcript(path: str) -> str:
    if not os.path.exists(path):
        raise FileNotFoundError(f"File not found: {path}")

    mime, _ = mimetypes.guess_type(path)
    ext = os.path.splitext(path)[1].lower()

    if mime and mime.startswith("text") or ext == ".txt":
        with open(path, "r", encoding="utf-8") as f:
            return f.read()

    if whisper is None:
        raise RuntimeError("Whisper is not installed")

    if mime and mime.startswith("audio"):
        return whisper.load_model("base").transcribe(path)["text"]

    if mime and mime.startswith("video"):
        if VideoFileClip is None:
            raise RuntimeError("moviepy not installed")

        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as tmp:
            tmp_path = tmp.name

        clip = VideoFileClip(path)
        clip.audio.write_audiofile(tmp_path, logger=None)
        clip.close()

        try:
            return whisper.load_model("base").transcribe(tmp_path)["text"]
        finally:
            os.remove(tmp_path)

    raise ValueError("Unsupported file type")


# -------------------------------------------------
# GEMINI
# -------------------------------------------------
def build_prompt(transcript: str) -> str:
    return f"""
Return STRICT JSON ONLY.

{{
  "entities": {{
    "people": [],
    "tools": [],
    "brands": [],
    "products": [],
    "organizations": []
  }},
  "tone": {{
    "primary": "",
    "secondary": [],
    "confidence": 0.0
  }},
  "style": {{
    "primary": "",
    "confidence": 0.0
  }},
  "safety_flags": {{
    "sensitive_domains": [],
    "severity": "",
    "requires_review": false
  }}
}}

Transcript:
{transcript}
"""


def analyze_with_gemini(transcript: str) -> Dict[str, Any]:
    if genai is None or not GEMINI_API_KEY:
        raise RuntimeError("Gemini is not configured")

    client = genai.Client(api_key=GEMINI_API_KEY)

    merged = {
        "entities": {"people": [], "tools": [], "brands": [], "products": [], "organizations": []},
        "tone": {"primary": "", "secondary": [], "confidence": 0.0},
        "style": {"primary": "", "confidence": 0.0},
        "safety_flags": {"sensitive_domains": [], "severity": "", "requires_review": False},
    }

    for chunk in chunk_transcript(transcript):
        raw = client.models.generate_content(
            model=MODEL_NAME,
            contents=build_prompt(chunk),
        ).text

        doc = json.loads(extract_json_substring(raw))
        if not validate_schema(doc):
            raise ValueError("Invalid Gemini schema")

        for k in merged["entities"]:
            merged["entities"][k].extend(x for x in doc["entities"][k] if x not in merged["entities"][k])

        if doc["tone"]["confidence"] > merged["tone"]["confidence"]:
            merged["tone"] = doc["tone"]

        if doc["style"]["confidence"] > merged["style"]["confidence"]:
            merged["style"] = doc["style"]

        merged["safety_flags"]["requires_review"] |= doc["safety_flags"]["requires_review"]

    return merged


# -------------------------------------------------
# DB
# -------------------------------------------------
def insert_db(result: Dict[str, Any], analysis_id: str, transcript_id: str, creator_id: str):
    if not psycopg2 or not DB_HOST:
        return

    conn = psycopg2.connect(
        host=DB_HOST,
        port=DB_PORT,
        user=DB_USER,
        password=DB_PASSWORD,
        dbname=DB_NAME,
    )
    cur = conn.cursor()
    cur.execute(
        f"""
        INSERT INTO {DB_TABLE}
        (analysis_id, transcript_id, creator_id, entities, tone, style, safety_flags, created_at, updated_at)
        VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)
        """,
        (
            analysis_id,
            transcript_id,
            creator_id,
            json.dumps(result["entities"]),
            json.dumps(result["tone"]),
            json.dumps(result["style"]),
            json.dumps(result["safety_flags"]),
            datetime.utcnow(),
            datetime.utcnow(),
        ),
    )
    conn.commit()
    cur.close()
    conn.close()


# -------------------------------------------------
# MAIN
# -------------------------------------------------
def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    parser.add_argument("--analysis-id", required=True)
    parser.add_argument("--transcript-id")
    parser.add_argument("--creator-id")
    parser.add_argument("--no-db", action="store_true")

    args = parser.parse_args()

    analysis_id = args.analysis_id
    transcript_id = args.transcript_id or str(uuid.uuid4())
    creator_id = args.creator_id or str(uuid.uuid4())

    transcript = extract_transcript(args.file)
    analysis = analyze_with_gemini(transcript)

    payload = {
        "analysis_id": analysis_id,
        "transcript_id": transcript_id,
        "creator_id": creator_id,
        **analysis,
        "created_at": datetime.utcnow().isoformat() + "Z",
        "status": "success",
    }

    os.makedirs(RESULTS_DIR, exist_ok=True)
    with open(os.path.join(RESULTS_DIR, f"{analysis_id}.json"), "w", encoding="utf-8") as f:
        json.dump(payload, f, indent=2)

    if not args.no_db:
        insert_db(analysis, analysis_id, transcript_id, creator_id)

    print(json.dumps(payload, ensure_ascii=False))


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(json.dumps({"status": "failed", "error": str(e)}))
        sys.exit(1)
