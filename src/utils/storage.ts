// src/utils/storage.ts
import type { AttendanceRecord, GeneratedCode } from '../types';

const STORAGE_KEYS = {
  ATTENDANCE_RECORDS: 'attendanceRecords',
  GENERATED_CODES: 'generatedCodes',
  CURRENT_USER: 'currentUser'
} as const;

// Attendance Records
export const getAttendanceRecords = (): AttendanceRecord[] => {
  const stored = localStorage.getItem(STORAGE_KEYS.ATTENDANCE_RECORDS);
  return stored ? JSON.parse(stored) : [];
};

export const saveAttendanceRecord = (record: AttendanceRecord): void => {
  const records = getAttendanceRecords();
  records.push(record);
  localStorage.setItem(STORAGE_KEYS.ATTENDANCE_RECORDS, JSON.stringify(records));
};

export const getStudentAttendanceRecords = (studentId: string): AttendanceRecord[] => {
  const records = getAttendanceRecords();
  return records.filter(record => record.studentId === studentId);
};

export const getLessonAttendanceRecords = (lessonId: string): AttendanceRecord[] => {
  const records = getAttendanceRecords();
  return records.filter(record => record.lessonId === lessonId);
};

// Generated Codes
export const getGeneratedCodes = (): GeneratedCode[] => {
  const stored = localStorage.getItem(STORAGE_KEYS.GENERATED_CODES);
  return stored ? JSON.parse(stored) : [];
};

export const saveGeneratedCode = (code: GeneratedCode): void => {
  const codes = getGeneratedCodes();
  // Prevent duplicates
  if (!codes.some(c => c.id === code.id)) {
    codes.push(code);
    localStorage.setItem(STORAGE_KEYS.GENERATED_CODES, JSON.stringify(codes));
  }
};

export const createNewCodeForLesson = (lessonId: string, lessonName: string, teacherId: string): GeneratedCode => {
  let codes = getGeneratedCodes();
  
  // Deactivate existing active codes for this lesson
  codes = codes.map(c => 
    (c.lessonId === lessonId && c.isActive) ? { ...c, isActive: false } : c
  );

  // Create new code
  const newCode: GeneratedCode = {
    id: Date.now().toString(),
    lessonId,
    lessonName,
    code: generateRandomCode(),
    teacherId,
    createdAt: new Date().toISOString(),
    expiresAt: new Date(Date.now() + 15 * 60 * 1000).toISOString(), // 15 minutes
    isActive: true
  };
  
  codes.push(newCode);
  
  setGeneratedCodes(codes);

  return newCode;
}

export const setGeneratedCodes = (codes: GeneratedCode[]): void => {
  localStorage.setItem(STORAGE_KEYS.GENERATED_CODES, JSON.stringify(codes));
};

export const getActiveCodeForLesson = (lessonId: string): GeneratedCode | null => {
  const codes = getGeneratedCodes();
  const now = new Date().toISOString();
  
  const activeCodes = codes
    .filter(code => 
      code.lessonId === lessonId && 
      code.isActive && 
      code.expiresAt > now
    )
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());

  return activeCodes[0] || null;
};

export const deactivateCode = (codeId: string): void => {
  const codes = getGeneratedCodes();
  const updatedCodes = codes.map(code => 
    code.id === codeId ? { ...code, isActive: false } : code
  );
  localStorage.setItem(STORAGE_KEYS.GENERATED_CODES, JSON.stringify(updatedCodes));
};

// Utility functions
export const generateRandomCode = (): string => {
  return Math.floor(10000 + Math.random() * 90000).toString();
};

export const isCodeValid = (code: string, lessonId: string): boolean => {
  const activeCode = getActiveCodeForLesson(lessonId);
  return activeCode?.code === code;
}; 