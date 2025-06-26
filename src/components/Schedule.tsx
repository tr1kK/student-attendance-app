import type { Lesson } from '../types';
import './Schedule.css';

interface ScheduleProps {
  lessons: Lesson[];
  onLessonClick: (lesson: Lesson) => void;
  userRole?: 'student' | 'teacher';
  attendedLessonIds?: string[];
}

const Schedule = ({ lessons, onLessonClick, userRole = 'student', attendedLessonIds = [] }: ScheduleProps) => {
  const daysOfWeek = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница'];

  const getLessonsByDay = (day: string) => {
    return lessons.filter(lesson => lesson.day === day);
  };

  return (
    <div className="schedule-container">
      <h2>Расписание на неделю</h2>
      <div className="schedule-grid">
        {daysOfWeek.map(day => (
          <div key={day} className="day-column">
            <h3>{day}</h3>
            {getLessonsByDay(day).length > 0 ? (
              getLessonsByDay(day).map(lesson => (
                <div key={lesson.id} className="lesson-card" onClick={() => onLessonClick(lesson)}>
                  <div className="lesson-info">
                    <strong>
                      {lesson.name}
                      {userRole === 'student' && attendedLessonIds.includes(String(lesson.id)) && <span className="checkmark"> ✔️</span>}
                    </strong>
                    <span>{lesson.time}</span>
                    <span>Ауд: {lesson.room}</span>
                    {userRole === 'teacher' && lesson.Group && <span>Группа: {lesson.Group.name}</span>}
                    <span>{lesson.teacher}</span>
                  </div>
                </div>
              ))
            ) : (
              <div className="lesson-card empty">
                <p>Нет занятий</p>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Schedule; 