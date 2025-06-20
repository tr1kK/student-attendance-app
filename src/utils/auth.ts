// src/utils/auth.ts
import { users, type User } from '../data/users';

export function authenticate(id: string, password: string): User | null {
  const user = users.find(u => u.id === id && u.password === password);
  if (!user) return null;

  localStorage.setItem('currentUser', JSON.stringify(user));
  return user;
}

export function getCurrentUser(): User | null {
  const stored = localStorage.getItem('currentUser');
  if (!stored) return null;
  return JSON.parse(stored);
}

export function logout() {
  localStorage.removeItem('currentUser');
}
