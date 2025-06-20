// src/types/index.ts

export interface Lesson {
  id: string;
  name: string;
  day: string;
  time: string;
  teacher: string;
  room: string;
}

export interface AttendanceRecord {
  id: string;
  lessonId: string;
  lessonName: string;
  studentId: string;
  studentName: string;
  timestamp: string;
  code: string;
}

export interface GeneratedCode {
  id: string;
  lessonId: string;
  lessonName: string;
  code: string;
  teacherId: string;
  createdAt: string;
  expiresAt: string;
  isActive: boolean;
}

export interface WeeklySchedule {
  [day: string]: Lesson[];
} 