from sentence_transformers import SentenceTransformer
import sys

def download():
    model_name = 'all-MiniLM-L6-v2'
    print(f"Downloading model: {model_name}")
    try:
        SentenceTransformer(model_name)
        print("Model downloaded successfully.")
    except Exception as e:
        print(f"Download failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    download()
