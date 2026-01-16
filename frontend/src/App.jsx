import { useState } from "react";
import { login, search, upload } from "./api/client";

export default function App() {
  const [token, setToken] = useState(() => localStorage.getItem("token"));
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [query, setQuery] = useState("");
  const [results, setResults] = useState({ documents: [], web: [] });
  const [file, setFile] = useState(null);

  async function handleLogin(e) {
    e.preventDefault();
    try {
      const data = await login(email, password);
      localStorage.setItem("token", data.token);
      setToken(data.token);
    } catch {
      alert("Login failed");
    }
  }

  async function handleSearch() {
    const data = await search(query, token, {
      docs: 1,
      duckduckgo: 1,
    });
    setResults(data);
  }

  async function handleUpload() {
    if (!file) return;
    await upload(file, token);
    alert("Uploaded");
  }

  if (!token) {
    return (
      <form onSubmit={handleLogin}>
        <h2>Login</h2>
        <input onChange={e => setEmail(e.target.value)} />
        <input type="password" onChange={e => setPassword(e.target.value)} />
        <button>Login</button>
      </form>
    );
  }

  return (
    <div>
      <h2>Hybrid Search Engine</h2>

      <input onChange={e => setQuery(e.target.value)} />
      <button onClick={handleSearch}>Search</button>

      <input type="file" onChange={e => setFile(e.target.files[0])} />
      <button onClick={handleUpload}>Upload</button>

      <div style={{ display: "flex" }}>
        <div style={{ width: "50%" }}>
          <h3>Documents</h3>
          {results.documents.map((d, i) => (
            <div key={i}>{d.snippet}</div>
          ))}
        </div>

        <div style={{ width: "50%" }}>
          <h3>Web</h3>
          {results.web.map((w, i) => (
            <div key={i}>{w.title}</div>
          ))}
        </div>
      </div>
    </div>
  );
}
