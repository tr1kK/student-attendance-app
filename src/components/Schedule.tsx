import type { WeeklySchedule, Lesson } from '../types';
import { weeklySchedule } from '../data/schedule';

interface ScheduleProps {
  onLessonClick?: (lesson: Lesson) => void;
  showActions?: boolean;
  userRole?: 'student' | 'teacher';
  attendedLessonIds?: string[];
}

const Schedule = ({ onLessonClick, showActions = false, userRole, attendedLessonIds = [] }: ScheduleProps) => {
  const days = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница'];

  return (
    <div className="schedule">
      <h2>Расписание на неделю</h2>
      <div className="schedule-grid">
        {days.map(day => (
          <div key={day} className="day-column">
            <h3>{day}</h3>
            {weeklySchedule[day]?.map(lesson => (
              <div key={lesson.id} className="lesson-card">
                <div className="lesson-info">
                  <h4>
                    {lesson.name}
                    {attendedLessonIds.includes(lesson.id) && <span className="checkmark">✔️</span>}
                  </h4>
                  <p className="lesson-time">{lesson.time}</p>
                  <p className="lesson-room">Кабинет: {lesson.room}</p>
                </div>
                {showActions && onLessonClick && (
                  <button 
                    className="lesson-action-btn"
                    onClick={() => onLessonClick(lesson)}
                  >
                    {userRole === 'teacher' ? 'Начать' : 'Отметить посещение'}
                  </button>
                )}
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Schedule; 