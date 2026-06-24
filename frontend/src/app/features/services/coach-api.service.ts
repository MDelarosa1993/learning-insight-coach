import { HttpClient, HttpHeaders } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
  DocumentUploadRequest,
  DocumentUploadResponse,
} from '../models/teacher-upload.model';
import { Observable } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class CoachApiService {
  private baseUrl = 'http://localhost:8080/api/v1';
  private apiKey = 'dev-secret-change-me';
  private http = inject(HttpClient);

  private getHeaders(): HttpHeaders {
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'X-API-Key': this.apiKey,
    });
  }

  uploadDocument(data: DocumentUploadRequest): Observable<DocumentUploadResponse> {
    return this.http.post<DocumentUploadResponse>(
      `${this.baseUrl}/documents`,
      data,
      { headers: this.getHeaders() }
    );
  }
}
