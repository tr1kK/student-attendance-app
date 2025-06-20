import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Schedule from '../components/Schedule';
import CodeEntryModal from '../components/CodeEntryModal';
import { getCurrentUser } from '../utils/auth';
import { 
  getStudentAttendanceRecords, 
  saveAttendanceRecord, 
  isCodeValid,
  setGeneratedCodes
} from '../utils/storage';
import { addMessageListener, removeMessageListener } from '../utils/broadcast';
import type { Lesson, AttendanceRecord, GeneratedCode } from '../types';
import type { User } from '../data/users';

const StudentPage = () => {
  const [selectedLesson, setSelectedLesson] = useState<Lesson | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [attendanceRecords, setAttendanceRecords] = useState<AttendanceRecord[]>([]);
  const [message, setMessage] = useState('');
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const user = getCurrentUser();
    if (!user || user.role !== 'student') {
      navigate('/');
    } else {
      setCurrentUser(user);
    }
  }, [navigate]);

  useEffect(() => {
    if (currentUser) {
      const records = getStudentAttendanceRecords(currentUser.id);
      setAttendanceRecords(records);
    }
  }, [currentUser]);

  useEffect(() => {
    const handleSyncCodes = (event: MessageEvent) => {
      if (event.data?.type === 'SYNC_CODES') {
        const codes: GeneratedCode[] = event.data.payload;
        setGeneratedCodes(codes); 
      }
    };

    addMessageListener(handleSyncCodes);
    return () => removeMessageListener(handleSyncCodes);
  }, []);

  const handleLessonClick = (lesson: Lesson) => {
    setSelectedLesson(lesson);
    setIsModalOpen(true);
  };

  const handleCodeSubmit = (code: string) => {
    if (!selectedLesson || !currentUser) return;

    if (isCodeValid(code, selectedLesson.id)) {
      const newRecord: AttendanceRecord = {
        id: Date.now().toString(),
        lessonId: selectedLesson.id,
        lessonName: selectedLesson.name,
        studentId: currentUser.id,
        studentName: currentUser.name,
        timestamp: new Date().toISOString(),
        code: code
      };

      saveAttendanceRecord(newRecord);
      setAttendanceRecords(prev => [...prev, newRecord]);
      setMessage('Посещение успешно отмечено!');
      setIsModalOpen(false);
      setSelectedLesson(null);
      
      // Clear message after 3 seconds
      setTimeout(() => setMessage(''), 3000);
    } else {
      setMessage('Неверный код. Попробуйте еще раз.');
      setTimeout(() => setMessage(''), 3000);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('currentUser');
    navigate('/');
  };

  if (!currentUser) return null;

  const attendedLessonIds = attendanceRecords.map(record => record.lessonId);

  return (
    <div className="student-page">
      <header className="page-header">
        <h1>Страница студента</h1>
        <div className="user-info">
          <span>Добро пожаловать, {currentUser.name}</span>
          <button onClick={handleLogout} className="logout-btn">
            Выйти
          </button>
        </div>
      </header>

      {message && (
        <div className={`message ${message.includes('успешно') ? 'success' : 'error'}`}>
          {message}
        </div>
      )}

      <div className="page-content">
        <Schedule 
          onLessonClick={handleLessonClick}
          showActions={true}
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
                    <h4>{record.lessonName}</h4>
                    <p>{new Date(record.timestamp).toLocaleDateString('ru-RU')}</p>
                    <p>{new Date(record.timestamp).toLocaleTimeString('ru-RU')}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {selectedLesson && (
        <CodeEntryModal
          lesson={selectedLesson}
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false);
            setSelectedLesson(null);
          }}
          onSubmit={handleCodeSubmit}
        />
      )}
    </div>
  );
};

export default StudentPage; 