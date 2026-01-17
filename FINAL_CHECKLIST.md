# NexusSearch Deployment Success Checklist

If you see a 502 or 404, here is the final "Ground Truth" configuration:

### 1. Python Worker (Render)
- **Repo**: `Sanjay24-05/nexus-search`
- **Root Directory**: `worker`
- **Env Vars**: 
  - `MONGODB_URI`: (Your Atlas link)
- **Deployment**: Must use **"Clear Build Cache & Deploy"** to pre-download the model.

### 2. Go Gateway (Render)
- **Repo**: `Sanjay24-05/nexus-search`
- **Root Directory**: `gateway`
- **Env Vars**:
  - `WORKER_URL`: `https://nexus-search-pozi.onrender.com` (NO trailing slash)
  - `ALLOWED_ORIGINS`: `*` (Use asterisk for troubleshooting CORS)
  - `MONGODB_URI`: (Same as worker)
  - `JWT_SECRET`: (Any string)
  - `SERPAPI_KEY`: (Your key)

### 3. React Frontend (Vercel)
- **Repo**: `Sanjay24-05/nexus-search-2` (Latest linked repo)
- **Root Directory**: `frontend`
- **Output Directory**: `dist`
- **Env Vars**:
  - `VITE_API_URL`: `https://nexus-search-1.onrender.com` (NO trailing slash)

---
### Verification Steps
1. **Login/Register**: Works = Frontend -> Gateway -> MongoDB is OK.
2. **Search**: Works = Gateway -> SerpApi/DDG is OK.
3. **Upload**: Works = Gateway -> Worker -> MongoDB is OK.
4. **PKB Results**: Works = Gateway -> Worker -> MongoDB Vector Search is OK.
