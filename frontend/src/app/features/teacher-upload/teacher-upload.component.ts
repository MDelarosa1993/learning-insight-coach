import { Component, EventEmitter, inject, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CoachApiService } from '../services/coach-api.service';
import { DocumentUploadResponse } from '../models/teacher-upload.model';

@Component({
  selector: 'app-teacher-upload',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './teacher-upload.component.html',
  styleUrl: './teacher-upload.component.css',
})
export class TeacherUploadComponent {
  private coachApi: CoachApiService = inject(CoachApiService);
  @Output() documentUploaded = new EventEmitter<string>();

  title = '';
  subject = '';
  gradeMin: number | null = null;
  gradeMax: number | null = null;
  productId = '';
  content = '';

  documentId = '';
  message = '';
  errorMessage = '';
  isLoading = false;

  uploadDocument(): void {
    this.errorMessage = '';
    this.message = '';
    this.documentId = '';

    if (
      !this.title ||
      !this.subject ||
      !this.productId ||
      !this.content ||
      this.gradeMin === null ||
      this.gradeMax === null
    ) {
      this.errorMessage = 'Please fill out all fields before uploading.';
      return;
    }

    this.isLoading = true;

    this.coachApi
      .uploadDocument({
        title: this.title,
        subject: this.subject,
        grade_min: this.gradeMin,
        grade_max: this.gradeMax,
        content: this.content,
        product_id: this.productId,
      })
      .subscribe({
        next: (response) => {
          this.documentId = response.document_id;
          this.message = `Document uploaded successfully. Chunks created: ${response.chunk_count}`;
          this.documentUploaded.emit(response.document_id);
          this.isLoading = false;
        },
        error: (error) => {
          this.errorMessage =
            error?.error?.error || 'Failed to upload document.';
          this.isLoading = false;
        },
      });
  }
}
