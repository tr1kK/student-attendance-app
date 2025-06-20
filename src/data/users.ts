// src/data/users.ts

export type Role = 'student' | 'teacher';

export interface User {
  id: string;
  password: string;
  name: string;
  role: Role;
}

export const users: User[] = [
  {
    id: 'stu001',
    password: '12345',
    name: 'Иванов И.И.',
    role: 'student',
  },
  {
    id: 'stu002',
    password: '23456',
    name: 'Петрова А.А.',
    role: 'student',
  },
  {
    id: 'teach001',
    password: 'admin1',
    name: 'Преподаватель С.С.',
    role: 'teacher',
  },
];
