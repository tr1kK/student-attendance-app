const API_BASE_URL = 'http://localhost:8080';

const getAuthToken = () => localStorage.getItem('authToken');

const apiFetch = async (url: string, options: RequestInit = {}) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  
  if (options.headers) {
    Object.assign(headers, options.headers);
  }

  const token = getAuthToken();
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'An unknown error occurred' }));
    throw new Error(errorData.error || 'Request failed');
  }

  return response.json();
};

export const login = (identifier: string, password: string) => {
  return apiFetch('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ identifier, password }),
  });
};

export const register = (userData: any) => {
  return apiFetch('/auth/register', {
    method: 'POST',
    body: JSON.stringify(userData),
  });
};

export const getGroups = () => {
  return apiFetch('/groups');
};

export const getLessons = () => {
  return apiFetch('/api/lessons');
};

export const generateCode = (lesson_id: number) => {
  return apiFetch(`/api/teacher/lessons/${lesson_id}/code`, {
    method: 'POST',
    body: JSON.stringify({ lesson_id }),
  });
};

export const deactivateCode = (lesson_id: number) => {
  return apiFetch(`/api/teacher/lessons/${lesson_id}/code`, {
    method: 'DELETE',
    body: JSON.stringify({ lesson_id }),
  });
};

export const getLessonAttendance = (lessonId: number) => {
  return apiFetch(`/api/teacher/attendance/${lessonId}`);
};

export const getStudentAttendance = () => {
  return apiFetch('/api/student/attendance');
};

export const submitAttendance = (lesson_id: number, code: string) => {
  return apiFetch('/api/student/attendance', {
    method: 'POST',
    body: JSON.stringify({ lesson_id, code }),
  });
};

// Admin API
export const adminGetUsers = () => apiFetch('/api/admin/users');
export const adminCreateUser = (user: any) => apiFetch('/api/admin/users', { method: 'POST', body: JSON.stringify(user) });
export const adminUpdateUser = (id: number, user: any) => apiFetch(`/api/admin/users/${id}`, { method: 'PUT', body: JSON.stringify(user) });
export const adminDeleteUser = (id: number) => apiFetch(`/api/admin/users/${id}`, { method: 'DELETE' });
export const adminGetGroups = () => apiFetch('/api/admin/groups'); 