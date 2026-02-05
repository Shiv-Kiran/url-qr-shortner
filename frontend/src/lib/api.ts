import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
});

export type QRErrorCorrectionLevel = 'L' | 'M' | 'Q' | 'H';

export interface ShortenURLPayload {
  original_url: string;
  qr_error_correction?: QRErrorCorrectionLevel;
  qr_size?: number;
}

export interface ShortenURLResponse {
  short_code: string;
  original_url: string;
  short_url: string;
  qr_data_url: string;
  qr_error_correction: QRErrorCorrectionLevel;
  qr_size: number;
}

export interface QRCodeResponse {
  short_code: string;
  short_url: string;
  qr_data_url: string;
  qr_error_correction: QRErrorCorrectionLevel;
  qr_size: number;
}

export async function shortenURL(payload: ShortenURLPayload) {
  const response = await apiClient.post<ShortenURLResponse>('/urls', payload);
  return response.data;
}

export async function getURL(shortCode: string) {
  return apiClient.get(`/urls/${shortCode}`);
}

export async function getQRCode(shortCode: string, level?: QRErrorCorrectionLevel, size?: number) {
  const response = await apiClient.get<QRCodeResponse>(`/urls/${shortCode}/qr`, {
    params: {
      level,
      size,
    },
  });

  return response.data;
}
