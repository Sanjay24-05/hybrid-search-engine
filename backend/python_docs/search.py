import json
import os

def search_documents(user_dir, query):
    index_path = f"{user_dir}/index.json"
    if not os.path.exists(index_path):
        return []

    index = json.load(open(index_path))
    return [
        doc for doc in index
        if query.lower() in doc["content"].lower()
    ]
