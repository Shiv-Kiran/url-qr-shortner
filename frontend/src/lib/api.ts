import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
});

export async function shortenURL(originalURL: string) {
  return apiClient.post('/urls', { original_url: originalURL });
}

export async function getURL(shortCode: string) {
  return apiClient.get(`/urls/${shortCode}`);
}
