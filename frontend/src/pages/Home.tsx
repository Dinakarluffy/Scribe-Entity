  import { useNavigate } from "react-router-dom";

  function Home() {
    const navigate = useNavigate();

    return (
      <div className="page">
        <h2>Entity & Classification Service</h2>

        <div className="home-buttons">
          <button onClick={() => navigate("/upload")}>Upload</button>
          <br />
          <button onClick={() => navigate("/results")}>Get Results</button>
          <br />
          <button onClick={() => navigate("/results-by-id")}>
            Get Result by ID
          </button>
        </div>
      </div>
    );
  }

  export default Home;
