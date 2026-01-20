import axios from "axios";

const API_BASE_URL = "http://localhost:8080";

export interface AnalyzePayload {
  transcript_id: string;
  creator_id: string;
  transcript_text: string;
}

export const analyzeTranscript = async (payload: AnalyzePayload) => {
  return axios.post(
    `${API_BASE_URL}/api/entity-classification/analyze`,
    payload
  );
};

export const getResults = async () => {
  return axios.get(
    `${API_BASE_URL}/api/entity-classification/results`
  );
};

export const getResultById = async (analysisId: string) => {
  return axios.get(
    `${API_BASE_URL}/api/entity-classification/results/${analysisId}`
  );
};
