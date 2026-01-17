FROM python:3.9-slim

WORKDIR /app

# Install system dependencies (needed for some python packages)
RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and scripts for downloading
COPY requirements.txt download_model.py ./
RUN pip install --no-cache-dir -r requirements.txt
RUN python download_model.py

# Copy source code
COPY . .

# Expose port (Render sets PORT env var)
EXPOSE 5000

CMD ["python", "app.py"]
