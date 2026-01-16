import json
import os

def search_documents(user_dir: str, query: str):
    index_path = os.path.join(user_dir, "index.json")
    if not os.path.exists(index_path):
        return []

    with open(index_path, "r", encoding="utf-8") as f:
        index = json.load(f)

    results = []
    for doc in index:
        content = doc.get("content", "").lower()
        if query.lower() in content:
            results.append({
                "file": os.path.basename(doc["file"]),
                "snippet": content[:200]
            })

    return results
