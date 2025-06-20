import type { Lesson, WeeklySchedule } from '../types';

export const lessons: Lesson[] = [
  {
    id: 'math-1',
    name: 'Математика',
    day: 'Понедельник',
    time: '09:00-10:30',
    teacher: 'Преподаватель С.С.',
    room: '101'
  },
  {
    id: 'physics-1',
    name: 'Физика',
    day: 'Понедельник',
    time: '10:45-12:15',
    teacher: 'Преподаватель С.С.',
    room: '102'
  },
  {
    id: 'chemistry-1',
    name: 'Химия',
    day: 'Вторник',
    time: '09:00-10:30',
    teacher: 'Преподаватель С.С.',
    room: '103'
  },
  {
    id: 'biology-1',
    name: 'Биология',
    day: 'Вторник',
    time: '10:45-12:15',
    teacher: 'Преподаватель С.С.',
    room: '104'
  },
  {
    id: 'history-1',
    name: 'История',
    day: 'Среда',
    time: '09:00-10:30',
    teacher: 'Преподаватель С.С.',
    room: '105'
  },
  {
    id: 'literature-1',
    name: 'Литература',
    day: 'Среда',
    time: '10:45-12:15',
    teacher: 'Преподаватель С.С.',
    room: '106'
  },
  {
    id: 'english-1',
    name: 'Английский язык',
    day: 'Четверг',
    time: '09:00-10:30',
    teacher: 'Преподаватель С.С.',
    room: '107'
  },
  {
    id: 'computer-science-1',
    name: 'Информатика',
    day: 'Четверг',
    time: '10:45-12:15',
    teacher: 'Преподаватель С.С.',
    room: '108'
  },
  {
    id: 'geography-1',
    name: 'География',
    day: 'Пятница',
    time: '09:00-10:30',
    teacher: 'Преподаватель С.С.',
    room: '109'
  },
  {
    id: 'art-1',
    name: 'Искусство',
    day: 'Пятница',
    time: '10:45-12:15',
    teacher: 'Преподаватель С.С.',
    room: '110'
  }
];

export const weeklySchedule: WeeklySchedule = {
  'Понедельник': lessons.filter(lesson => lesson.day === 'Понедельник'),
  'Вторник': lessons.filter(lesson => lesson.day === 'Вторник'),
  'Среда': lessons.filter(lesson => lesson.day === 'Среда'),
  'Четверг': lessons.filter(lesson => lesson.day === 'Четверг'),
  'Пятница': lessons.filter(lesson => lesson.day === 'Пятница'),
}; 