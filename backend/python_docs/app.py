from fastapi import FastAPI, UploadFile, File, Header, HTTPException, Query
from dotenv import load_dotenv
import os

from backend.python_docs.quota import check_quota
from backend.python_docs.indexer import index_document
from backend.python_docs.search import search_documents

# Load environment variables
load_dotenv()

app = FastAPI()

# Environment config
MAX_UPLOAD_MB = int(os.getenv("MAX_UPLOAD_MB", "50"))
MAX_QUOTA_BYTES = MAX_UPLOAD_MB * 1024 * 1024

BASE_DIR = os.path.join(
    os.path.dirname(__file__),
    "storage",
    "users"
)

@app.post("/upload")
async def upload_document(
    file: UploadFile = File(...),
    user_id: str = Header(...)
):
    user_dir = os.path.join(BASE_DIR, user_id)
    os.makedirs(user_dir, exist_ok=True)

    content = await file.read()

    if not check_quota(user_dir, len(content), MAX_QUOTA_BYTES):
        raise HTTPException(
            status_code=400,
            detail=f"{MAX_UPLOAD_MB}MB quota exceeded"
        )

    path = os.path.join(user_dir, file.filename)
    with open(path, "wb") as f:
        f.write(content)

    index_document(user_dir, path)
    return {"status": "uploaded"}


@app.get("/search-docs")
def search_docs(
    q: str = Query(...),
    user_id: str = Header(...)
):
    user_dir = os.path.join(BASE_DIR, user_id)

    if not os.path.exists(user_dir):
        return []

    return search_documents(user_dir, q)
