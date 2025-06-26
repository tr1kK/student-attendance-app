import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Schedule from '../components/Schedule';
import CodeDisplayModal from '../components/CodeDisplayModal';
import AttendanceList from '../components/AttendanceList';
import * as api from '../utils/api';
import type { Lesson, AttendanceRecord, GeneratedCode, User } from '../types';

const TeacherPage = () => {
  const [user, setUser] = useState<User | null>(null);
  const [lessons, setLessons] = useState<Lesson[]>([]);
  const [selectedLesson, setSelectedLesson] = useState<Lesson | null>(null);
  const [isCodeModalOpen, setIsCodeModalOpen] = useState(false);
  const [currentCode, setCurrentCode] = useState<GeneratedCode | null>(null);
  const [attendanceRecords, setAttendanceRecords] = useState<any[]>([]);
  const [message, setMessage] = useState('');
  const [generatedCode, setGeneratedCode] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();
  
  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    } else {
      navigate('/');
    }

    const fetchLessons = async () => {
      try {
        const data = await api.getLessons();
        setLessons(data);
      } catch (err) {
        setError('Не удалось загрузить расписание.');
        console.error(err);
      }
    };

    fetchLessons();
  }, [navigate]);

  useEffect(() => {
    if (!selectedLesson) return;
    const intervalId = setInterval(async () => {
      try {
        const data = await api.getLessonAttendance(selectedLesson.id);
        setAttendanceRecords(data);
      } catch (error) {
        console.error('Failed to fetch attendance:', error);
      }
    }, 3000); // Poll every 3 seconds

    return () => clearInterval(intervalId);
  }, [selectedLesson]);

  const handleLessonClick = async (lesson: Lesson) => {
    setSelectedLesson(lesson);
    try {
      const code = await api.generateCode(lesson.id);
      setCurrentCode(code);
      setIsCodeModalOpen(true);
      const records = await api.getLessonAttendance(lesson.id);
      setAttendanceRecords(records);
    } catch (error: any) {
      console.error('Failed to generate code:', error);
      setMessage(error.message || 'Failed to start lesson.');
    }
  };

  const handleGenerateCode = async (lesson: Lesson) => {
    setSelectedLesson(lesson);
    setAttendanceRecords([]); // Clear old records
    try {
      setError('');
      setMessage('');
      const response = await api.generateCode(lesson.id);
      setGeneratedCode(response.code);
      setMessage(`Код для занятия '${lesson.name}' сгенерирован.`);
      setIsModalOpen(true);
      const records = await api.getLessonAttendance(lesson.id);
      setAttendanceRecords(records);
    } catch (error: any) {
      setError(error.message || 'Не удалось сгенерировать код.');
    }
  };

  const handleStopCode = async () => {
    if (!selectedLesson) return;
    try {
      await api.deactivateCode(selectedLesson.id);
      setMessage(`Код для занятия '${selectedLesson.name}' остановлен.`);
      setIsModalOpen(false);
      setGeneratedCode(null);
    } catch (err: any) {
      setError(err.message || 'Не удалось остановить код.');
    }
  };

  const handleRefreshAttendance = async () => {
    if (!selectedLesson) return;
    try {
      setMessage('Обновление списка...');
      const records = await api.getLessonAttendance(selectedLesson.id);
      setAttendanceRecords(records);
      setMessage('Список посещаемости обновлен.');
    } catch (err) {
      setError('Не удалось обновить список.');
    } finally {
      setTimeout(() => setMessage(''), 3000);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
    navigate('/');
  };

  return (
    <div className="container">
      <header>
        <h1>Портал преподавателя</h1>
        <div className="user-info">
          <span>Добро пожаловать, {user?.name || 'Преподаватель'}</span>
          <button onClick={handleLogout} className="logout-button">Выйти</button>
        </div>
      </header>
      <main>
        {message && <p className="message success">{message}</p>}
        {error && <p className="error">{error}</p>}
        <Schedule lessons={lessons} onLessonClick={handleGenerateCode} userRole="teacher" />
        {selectedLesson && (
          <div className="attendance-section">
            <div className="attendance-header">
              <h3>Посещаемость: {selectedLesson.name}</h3>
              <button onClick={handleRefreshAttendance} className="btn-secondary">Обновить список</button>
            </div>
            <AttendanceList attendanceRecords={attendanceRecords} />
          </div>
        )}
      </main>

      {selectedLesson && currentCode && (
        <CodeDisplayModal
          lesson={selectedLesson}
          code={currentCode.code}
          isOpen={isCodeModalOpen}
          onClose={() => setIsCodeModalOpen(false)}
          onStop={() => setIsCodeModalOpen(false)} // Stop functionality can be enhanced later
        />
      )}

      {isModalOpen && selectedLesson && generatedCode && (
        <CodeDisplayModal
          lesson={selectedLesson}
          code={generatedCode}
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onStop={handleStopCode}
        />
      )}
    </div>
  );
};

export default TeacherPage; 