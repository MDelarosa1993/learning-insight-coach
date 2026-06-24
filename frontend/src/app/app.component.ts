import { Component } from '@angular/core';
import { TeacherUploadComponent } from './features/teacher-upload/teacher-upload.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [TeacherUploadComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  currentDocumentId = '';

  handleDocumentUploaded(documentId: string): void {
    this.currentDocumentId = documentId;
    console.log('Uploaded document ID:', documentId);
  }
}
