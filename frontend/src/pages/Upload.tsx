import { useState } from "react";

interface AnalysisResult {
  analysis_id: string;
  transcript_id: string;
  creator_id: string;
  entities: {
    people: string[];
    tools: string[];
    brands: string[];
    products: string[];
    organizations: string[];
  };
  tone: {
    primary: string;
    secondary: string[];
    confidence: number;
  };
  style: {
    primary: string;
    confidence: number;
  };
  safety_flags: {
    sensitive_domains: string[];
    severity: string;
    requires_review: boolean;
  };
  created_at: string;
}

interface UploadResponse {
  status: string;
  result: AnalysisResult;
}

const Upload = () => {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<AnalysisResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0) return;
    setFile(e.target.files[0]);
    setError(null);
    setResult(null);
    setStatus("");
  };

  const handleUpload = async () => {
    if (!file) {
      setStatus("Please select a file");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      setLoading(true);
      setStatus("Uploading and processing... This may take a few minutes.");
      setError(null);
      setResult(null);

      console.log("Uploading file:", file.name);

      const res = await fetch("http://localhost:8080/api/entity-classification/upload", {
        method: "POST",
        body: formData,
        mode: "cors",
      });

      console.log("Response status:", res.status);

      if (!res.ok) {
        const errorText = await res.text();
        console.error("Error response:", errorText);
        throw new Error(`Upload failed (${res.status}): ${errorText}`);
      }

      const data: UploadResponse = await res.json();
      console.log("Success response:", data);
      
      if (data.status === "success" && data.result) {
        setResult(data.result);
        setStatus("Processing completed successfully!");
      } else {
        throw new Error("Invalid response from server");
      }
    } catch (err) {
      console.error("Upload error details:", err);
      let errorMessage = "Unknown error occurred";
      
      if (err instanceof TypeError && err.message.includes("fetch")) {
        errorMessage = "Cannot connect to server. Please ensure the backend is running on http://localhost:8080";
      } else if (err instanceof Error) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);
      setStatus("Error processing file");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page" style={{ padding: "20px", maxWidth: "900px", margin: "0 auto" }}>
      <h2>Upload Content for Analysis</h2>

      <div style={{ marginBottom: "20px" }}>
        <input
          type="file"
          accept=".mp4,.mov,.webm,.mp3,.wav,.txt"
          onChange={handleFileChange}
          disabled={loading}
          style={{ marginBottom: "10px" }}
        />

        {file && (
          <p style={{ color: "#666" }}>
            Selected file: <strong>{file.name}</strong> ({(file.size / 1024 / 1024).toFixed(2)} MB)
          </p>
        )}

        <button 
          onClick={handleUpload} 
          disabled={!file || loading}
          style={{ 
            marginTop: "12px",
            padding: "10px 20px",
            backgroundColor: loading ? "#ccc" : "#007bff",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: loading ? "not-allowed" : "pointer"
          }}
        >
          {loading ? "Processing..." : "Upload & Analyze"}
        </button>
      </div>

      {status && (
        <div style={{ 
          padding: "12px", 
          marginBottom: "16px",
          backgroundColor: error ? "#fee" : "#efe",
          border: `1px solid ${error ? "#fcc" : "#cfc"}`,
          borderRadius: "4px"
        }}>
          <p style={{ margin: 0, color: error ? "#c00" : "#060" }}>{status}</p>
        </div>
      )}

      {error && (
        <div style={{ 
          padding: "12px", 
          marginBottom: "16px",
          backgroundColor: "#fee",
          border: "1px solid #fcc",
          borderRadius: "4px"
        }}>
          <p style={{ margin: 0, color: "#c00", fontWeight: "bold" }}>Error:</p>
          <p style={{ margin: "8px 0 0 0", color: "#c00" }}>{error}</p>
          <details style={{ marginTop: "8px" }}>
            <summary style={{ cursor: "pointer", fontSize: "0.9em" }}>Troubleshooting tips</summary>
            <ul style={{ fontSize: "0.9em", marginTop: "8px" }}>
              <li>Check if the backend server is running on port 8080</li>
              <li>Check browser console (F12) for detailed error messages</li>
              <li>Verify CORS is enabled in the backend</li>
              <li>Ensure Python dependencies are installed</li>
            </ul>
          </details>
        </div>
      )}

      {result && (
        <div style={{ marginTop: "20px" }}>
          <h3>Analysis Results</h3>
          
          <div style={{ marginBottom: "16px" }}>
            <p><strong>Analysis ID:</strong> {result.analysis_id}</p>
            <p><strong>Transcript ID:</strong> {result.transcript_id}</p>
            <p><strong>Created At:</strong> {new Date(result.created_at).toLocaleString()}</p>
          </div>

          <div style={{ marginBottom: "16px" }}>
            <h4>Entities</h4>
            {Object.entries(result.entities).map(([key, values]) => (
              values.length > 0 && (
                <div key={key} style={{ marginBottom: "8px" }}>
                  <strong>{key.charAt(0).toUpperCase() + key.slice(1)}:</strong> {values.join(", ")}
                </div>
              )
            ))}
          </div>

          <div style={{ marginBottom: "16px" }}>
            <h4>Tone</h4>
            <p><strong>Primary:</strong> {result.tone.primary}</p>
            {result.tone.secondary.length > 0 && (
              <p><strong>Secondary:</strong> {result.tone.secondary.join(", ")}</p>
            )}
            <p><strong>Confidence:</strong> {(result.tone.confidence * 100).toFixed(1)}%</p>
          </div>

          <div style={{ marginBottom: "16px" }}>
            <h4>Style</h4>
            <p><strong>Primary:</strong> {result.style.primary}</p>
            <p><strong>Confidence:</strong> {(result.style.confidence * 100).toFixed(1)}%</p>
          </div>

          <div style={{ marginBottom: "16px" }}>
            <h4>Safety Flags</h4>
            <p><strong>Severity:</strong> {result.safety_flags.severity}</p>
            <p><strong>Requires Review:</strong> {result.safety_flags.requires_review ? "Yes" : "No"}</p>
            {result.safety_flags.sensitive_domains.length > 0 && (
              <p><strong>Sensitive Domains:</strong> {result.safety_flags.sensitive_domains.join(", ")}</p>
            )}
          </div>

          <details style={{ marginTop: "20px" }}>
            <summary style={{ cursor: "pointer", fontWeight: "bold" }}>View Raw JSON</summary>
            <pre style={{ 
              marginTop: "12px", 
              background: "#f5f5f5", 
              padding: "12px",
              borderRadius: "4px",
              overflow: "auto"
            }}>
              {JSON.stringify(result, null, 2)}
            </pre>
          </details>
        </div>
      )}
    </div>
  );
};

export default Upload;