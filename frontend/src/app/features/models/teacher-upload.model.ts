export interface DocumentUploadRequest {
  title: string;
  subject: string;
  grade_min: number;
  grade_max: number;
  content: string;
  product_id: string;
}

export interface DocumentUploadResponse {
  document_id: string;
  status: string;
  chunk_count: number;
  message: string;
}
