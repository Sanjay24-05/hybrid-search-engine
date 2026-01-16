from fastapi import FastAPI, UploadFile, File, Header, HTTPException
import os

from backend.python_docs.quota import check_quota
from backend.python_docs.indexer import index_document
from backend.python_docs.search import search_documents

app = FastAPI()

BASE_DIR = "storage/users"

@app.post("/upload")
async def upload_document(
    file: UploadFile = File(...),
    user_id: str = Header(...)
):
    user_dir = os.path.join(BASE_DIR, user_id)
    os.makedirs(user_dir, exist_ok=True)

    content = await file.read()

    if not check_quota(user_dir, len(content)):
        raise HTTPException(status_code=400, detail="50MB quota exceeded")

    path = os.path.join(user_dir, file.filename)
    with open(path, "wb") as f:
        f.write(content)

    index_document(user_dir, path)
    return {"status": "uploaded"}
