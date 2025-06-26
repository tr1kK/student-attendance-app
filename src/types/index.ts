// src/types/index.ts

export interface Group {
  id: number;
  name: string;
}

export interface User {
  id: number;
  identifier: string;
  name: string;
  email: string;
  role: 'student' | 'teacher' | 'admin';
  group_id?: number;
  group?: Group;
}

export interface Lesson {
  id: number;
  name: string;
  day: string;
  time: string;
  teacher: string;
  room: string;
  group_id?: number;
  group?: Group;
}

export interface AttendanceRecord {
  id: number;
  lesson_id: number;
  student_id: number;
  submitted_at: string;
  lesson: Lesson;
  student: User;
}

export interface GeneratedCode {
  id: number;
  lesson_id: number;
  code: string;
  expires_at: string;
  is_active: boolean;
}

export interface WeeklySchedule {
  [day: string]: Lesson[];
} 