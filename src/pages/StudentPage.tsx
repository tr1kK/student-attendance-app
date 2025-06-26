import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Schedule from '../components/Schedule';
import CodeEntryModal from '../components/CodeEntryModal';
import * as api from '../utils/api';
import type { Lesson, AttendanceRecord, User } from '../types';

const StudentPage = () => {
  const [user, setUser] = useState<User | null>(null);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [selectedLesson, setSelectedLesson] = useState<Lesson | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [attendanceRecords, setAttendanceRecords] = useState<AttendanceRecord[]>([]);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      const parsedUser = JSON.parse(storedUser);
      setUser(parsedUser);
    } else {
      // If no user, redirect to login
      navigate('/');
    }

    const fetchInitialData = async () => {
      try {
        const lessonsData = await api.getLessons();
        setLessons(lessonsData);
        
        const attendanceData = await api.getStudentAttendance();
        setAttendanceRecords(attendanceData);
      } catch (err) {
        setError('Не удалось загрузить данные.');
        console.error(err);
      }
    };
    fetchInitialData();
  }, [navigate]);

  const handleLessonClick = (lesson: Lesson) => {
    setSelectedLesson(lesson);
    setIsModalOpen(true);
  };

  const handleSubmitCode = async (code: string) => {
    if (!selectedLesson) return;
    try {
      setError('');
      setMessage('');
      const response = await api.submitAttendance(selectedLesson.id, code);
      setMessage(response.message || 'Посещение успешно отмечено!');
      
      // Refetch attendance records to update the UI
      const updatedAttendance = await api.getStudentAttendance();
      setAttendanceRecords(updatedAttendance);

      setIsModalOpen(false);
    } catch (err: any) {
      setError(err.message || 'Ошибка при отправке кода.');
    } finally {
      setTimeout(() => {
        setMessage('');
        setError('');
      }, 5000);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
    navigate('/');
  };

  const attendedLessonIds = attendanceRecords.map(record => String(record.lesson_id));

  return (
    <div className="container">
      <header>
        <h1>Портал студента</h1>
        <div className="user-info">
          <span>Добро пожаловать, {user?.name || 'Студент'}</span>
          <button onClick={handleLogout} className="logout-button">Выйти</button>
        </div>
      </header>
      <main>
        {message && <p className="message success">{message}</p>}
        {error && <p className="error">{error}</p>}
        <Schedule 
          lessons={lessons}
          onLessonClick={handleLessonClick}
          userRole="student"
          attendedLessonIds={attendedLessonIds}
        />

        <div className="attendance-history">
          <h2>История посещений</h2>
          {attendanceRecords.length === 0 ? (
            <p>У вас пока нет отметок о посещении</p>
          ) : (
            <div className="history-list">
              {attendanceRecords.map(record => (
                <div key={record.id} className="history-item">
                  <div className="history-info">
                    <h4>{record.lesson.name}</h4>
                    <p>{new Date(record.submitted_at).toLocaleDateString('ru-RU')}</p>
                    <p>{new Date(record.submitted_at).toLocaleTimeString('ru-RU')}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>

      {selectedLesson && (
        <CodeEntryModal
          lesson={selectedLesson}
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false);
            setSelectedLesson(null);
          }}
          onSubmit={handleSubmitCode}
        />
      )}
    </div>
  );
};

export default StudentPage; 