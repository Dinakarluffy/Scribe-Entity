import { useEffect, useState } from "react";
import { getResults } from "../api/api";

function Results() {
  const [data, setData] = useState<any[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchResults();
  }, []);

  const fetchResults = async () => {
    try {
      setLoading(true);
      setError("");

      const res = await getResults();

      console.log("GET RESULTS RESPONSE:", res.data); // remove after verification

      // Ensure response is an array
      const results = Array.isArray(res.data) ? res.data : [];

      setData(results);
    } catch (err: any) {
      console.error("Fetch results error:", err);

      if (err.response) {
        setError(`Server error: ${err.response.status}`);
      } else {
        setError("Network/CORS error – backend not reachable");
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <p>Loading results...</p>;
  }

  return (
    <div className="page">
      <h2>Analysis Results</h2>

      {error && <p className="error">{error}</p>}

      <table border={1} cellPadding={8} style={{ width: "100%" }}>
        <thead>
          <tr>
            <th>Analysis ID</th>
            <th>Transcript ID</th>
            <th>Creator ID</th>
            <th>People</th>
            <th>Tools</th>
            <th>Brands</th>
            <th>Products</th>
            <th>Tone</th>
            <th>Style</th>
            <th>Sensitive Domains</th>
            <th>Severity</th>
            <th>Requires Review</th>
            <th>Created At</th>
            <th>Updated At</th> {/* ✅ NEW COLUMN */}
          </tr>
        </thead>

        <tbody>
          {data.length === 0 && !error && (
            <tr>
              <td colSpan={14} style={{ textAlign: "center" }}>
                No results found
              </td>
            </tr>
          )}

          {data.map((row) => (
            <tr key={row.analysis_id}>
              <td>{row.analysis_id}</td>
              <td>{row.transcript_id}</td>
              <td>{row.creator_id}</td>
              <td>{row.entities?.people?.join(", ") || "-"}</td>
              <td>{row.entities?.tools?.join(", ") || "-"}</td>
              <td>{row.entities?.brands?.join(", ") || "-"}</td>
              <td>{row.entities?.products?.join(", ") || "-"}</td>
              <td>{row.tone?.primary || "-"}</td>
              <td>{row.style?.primary || "-"}</td>
              <td>
                {row.safety_flags?.sensitive_domains?.join(", ") || "None"}
              </td>
              <td>{row.safety_flags?.severity || "None"}</td>
              <td>
                {row.safety_flags?.requires_review ? "Yes" : "No"}
              </td>
              <td>
                {row.created_at
                  ? new Date(row.created_at).toLocaleString()
                  : "-"}
              </td>
              <td>
                {row.updated_at
                  ? new Date(row.updated_at).toLocaleString()
                  : "-"}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Results;
