import { useState } from "react";

const Upload = () => {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState("");

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0) return;
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      setStatus("Please select a file");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      setStatus("Uploading...");
      const res = await fetch(
        "http://localhost:8080/api/transcript/upload",
        {
          method: "POST",
          body: formData,
        }
      );

      if (!res.ok) {
        throw new Error("Upload failed");
      }

      setStatus("File uploaded successfully");
    } catch {
      setStatus("Error uploading file");
    }
  };

  return (
    <div className="page">
      <h2>Upload Content</h2>

      <input
        type="file"
        accept=".mp4,.mov,.webm,.mp3,.wav,.txt"
        onChange={handleFileChange}
      />

      {file && (
        <p>
          Selected file: <strong>{file.name}</strong>
        </p>
      )}

      <button onClick={handleUpload} style={{ marginTop: "12px" }}>
        Upload
      </button>

      {status && <p>{status}</p>}
    </div>
  );
};

export default Upload;
