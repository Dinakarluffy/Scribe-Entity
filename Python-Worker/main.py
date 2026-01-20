import os
import json
import psycopg2
from dotenv import load_dotenv
import uuid

load_dotenv()

conn = psycopg2.connect(
    host=os.getenv("DB_HOST"),
    port=os.getenv("DB_PORT"),
    user=os.getenv("DB_USER"),
    password=os.getenv("DB_PASSWORD"),
    dbname=os.getenv("DB_NAME")
)

cursor = conn.cursor()

# Mandatory per EM
cursor.execute("SET ROLE dinakaran_dev;")

TABLE_NAME = os.getenv("DB_TABLE", "scribe_entity_classification_dev")

query = f"""
INSERT INTO {TABLE_NAME}
(
  transcript_id,
  creator_id,
  entities,
  tone,
  style,
  safety_flags
)
VALUES (%s, %s, %s, %s, %s, %s)
"""

data = (
    str(uuid.uuid4()),   # transcript_id
    str(uuid.uuid4()),   # creator_id
    json.dumps({
        "people": ["Asta"],
        "tools": ["Notion"],
        "brands": ["Notion"],
        "products": ["The 4-Hour Workweek"]
    }),
    json.dumps({
        "primary": "informative",
        "secondary": "conversational",
        "confidence": 0.9
    }),
    json.dumps({
        "type": "tutorial"
    }),
    json.dumps({
        "medical": False,
        "financial": False,
        "political": False
    })
)

cursor.execute(query, data)
conn.commit()

cursor.close()
conn.close()

print("Entity analysis inserted successfully")
