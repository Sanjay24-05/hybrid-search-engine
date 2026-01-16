import json
import os
from backend.python_docs.extractor.txt import extract_text


def index_document(user_dir, file_path):
    text = extract_text(file_path)
    index_path = f"{user_dir}/index.json"

    index = []
    if os.path.exists(index_path):
        index = json.load(open(index_path))

    index.append({
        "file": file_path,
        "content": text
    })

    json.dump(index, open(index_path, "w"))
