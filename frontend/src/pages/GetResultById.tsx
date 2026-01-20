import { useState } from "react";
import { getResultById } from "../api/api";

function GetResultById() {
  const [analysisId, setAnalysisId] = useState("");
  const [data, setData] = useState<any | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

 const fetchResult = async () => {
  if (!analysisId.trim()) {
    setError("Analysis ID is required");
    return;
  }

  try {
    setLoading(true);
    setError("");
    setData(null);

    const res = await getResultById(analysisId);

    console.log("API RESPONSE:", res.data); // 🔍 DEBUG (keep for now)

    setData(res.data);
  } catch (err: any) {
    console.error("Fetch error:", err);

    if (err.response) {
      setError(`Server error: ${err.response.status}`);
    } else {
      setError("Network/CORS error – check backend CORS");
    }
  } finally {
    setLoading(false);
  }
};


  return (
    <div className="page">
      <h2>Get Result by Analysis ID</h2>

      <input
        type="text"
        placeholder="Enter analysis_id"
        value={analysisId}
        onChange={(e) => setAnalysisId(e.target.value)}
        style={{ width: "420px" }}
      />

      <br /><br />

      <button onClick={fetchResult} disabled={loading}>
        {loading ? "Fetching..." : "Get Result"}
      </button>

      {error && <p className="error">{error}</p>}

      {data && (
        <table border={1} cellPadding={8} style={{ marginTop: 20 }}>
          <tbody>
            <tr>
              <th>Analysis ID</th>
              <td>{data.analysis_id}</td>
            </tr>
            <tr>
              <th>Transcript ID</th>
              <td>{data.transcript_id}</td>
            </tr>
            <tr>
              <th>Creator ID</th>
              <td>{data.creator_id}</td>
            </tr>
            <tr>
              <th>People</th>
              <td>{data.entities?.people?.join(", ") || "-"}</td>
            </tr>
            <tr>
              <th>Tools</th>
              <td>{data.entities?.tools?.join(", ") || "-"}</td>
            </tr>
            <tr>
              <th>Brands</th>
              <td>{data.entities?.brands?.join(", ") || "-"}</td>
            </tr>
            <tr>
              <th>Products</th>
              <td>{data.entities?.products?.join(", ") || "-"}</td>
            </tr>
            <tr>
              <th>Tone</th>
              <td>{data.tone?.primary || "-"}</td>
            </tr>
            <tr>
              <th>Style</th>
              <td>{data.style?.primary || "-"}</td>
            </tr>
            <tr>
              <th>Sensitive Domains</th>
              <td>
                {data.safety_flags?.sensitive_domains?.join(", ") || "None"}
              </td>
            </tr>
            <tr>
              <th>Severity</th>
              <td>{data.safety_flags?.severity || "None"}</td>
            </tr>
            <tr>
              <th>Requires Review</th>
              <td>{data.safety_flags?.requires_review ? "Yes" : "No"}</td>
            </tr>
            <tr>
              <th>Created At</th>
              <td>
                {data.created_at
                  ? new Date(data.created_at).toLocaleString()
                  : "-"}
              </td>
            </tr>
            <tr>
              <th>Updated At</th>
              <td>
                {data.updated_at
                  ? new Date(data.updated_at).toLocaleString()
                  : "-"}
              </td>
            </tr>
          </tbody>
        </table>
      )}
    </div>
  );
}

export default GetResultById;
