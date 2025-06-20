import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Schedule from '../components/Schedule';
import CodeDisplayModal from '../components/CodeDisplayModal';
import AttendanceList from '../components/AttendanceList';
import { getCurrentUser } from '../utils/auth';
import { 
  getLessonAttendanceRecords, 
  saveGeneratedCode, 
  getActiveCodeForLesson,
  deactivateCode,
  generateRandomCode,
  getGeneratedCodes,
  createNewCodeForLesson
} from '../utils/storage';
import { postMessage } from '../utils/broadcast';
import type { Lesson, AttendanceRecord, GeneratedCode } from '../types';

const TeacherPage = () => {
  const [selectedLesson, setSelectedLesson] = useState<Lesson | null>(null);
  const [isCodeModalOpen, setIsCodeModalOpen] = useState(false);
  const [currentCode, setCurrentCode] = useState<string>('');
  const [attendanceRecords, setAttendanceRecords] = useState<AttendanceRecord[]>([]);
  const [activeCode, setActiveCode] = useState<GeneratedCode | null>(null);
  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  const currentUser = getCurrentUser();

  useEffect(() => {
    if (!currentUser || currentUser.role !== 'teacher') {
      navigate('/');
      return;
    }
  }, [currentUser, navigate]);

  useEffect(() => {
    if (!selectedLesson) return;

    const intervalId = setInterval(() => {
      const records = getLessonAttendanceRecords(selectedLesson.id);
      if (JSON.stringify(records) !== JSON.stringify(attendanceRecords)) {
        setAttendanceRecords(records);
      }
    }, 2000); 

    return () => clearInterval(intervalId);
  }, [selectedLesson, attendanceRecords]);

  const handleLessonClick = (lesson: Lesson) => {
    setSelectedLesson(lesson);
    
    // Check for an existing active code, otherwise create a new one.
    let codeToDisplay = getActiveCodeForLesson(lesson.id);
    if (!codeToDisplay) {
      codeToDisplay = createNewCodeForLesson(lesson.id, lesson.name, currentUser!.id);
    }
    
    setActiveCode(codeToDisplay);
    setCurrentCode(codeToDisplay.code);
    setIsCodeModalOpen(true);

    // Broadcast the updated state of all codes
    const allCodes = getGeneratedCodes();
    postMessage({ type: 'SYNC_CODES', payload: allCodes });

    // Load attendance records for this lesson
    const records = getLessonAttendanceRecords(lesson.id);
    setAttendanceRecords(records);
  };

  const handleStopCode = () => {
    if (activeCode) {
      deactivateCode(activeCode.id);

      // Broadcast the entire state of codes
      const allCodes = getGeneratedCodes();
      postMessage({ type: 'SYNC_CODES', payload: allCodes });
      
      setActiveCode(null);
      setCurrentCode('');
      setIsCodeModalOpen(false);
      setMessage('Код деактивирован');
      setTimeout(() => setMessage(''), 3000);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('currentUser');
    navigate('/');
  };

  if (!currentUser) return null;

  return (
    <div className="teacher-page">
      <header className="page-header">
        <h1>Страница преподавателя</h1>
        <div className="user-info">
          <span>Добро пожаловать, {currentUser.name}</span>
          <button onClick={handleLogout} className="logout-btn">
            Выйти
          </button>
        </div>
      </header>

      {message && (
        <div className={`message ${message.includes('деактивирован') ? 'success' : 'error'}`}>
          {message}
        </div>
      )}

      <div className="page-content">
        <Schedule 
          onLessonClick={handleLessonClick}
          showActions={true}
          userRole="teacher"
        />

        {selectedLesson && (
          <AttendanceList 
            attendanceRecords={attendanceRecords}
            lessonName={selectedLesson.name}
          />
        )}
      </div>

      {selectedLesson && currentCode && (
        <CodeDisplayModal
          lesson={selectedLesson}
          code={currentCode}
          isOpen={isCodeModalOpen}
          onClose={() => setIsCodeModalOpen(false)}
          onStop={handleStopCode}
        />
      )}
    </div>
  );
};

export default TeacherPage; 