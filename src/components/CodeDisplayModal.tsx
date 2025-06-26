import type { Lesson } from '../types';
import './CodeDisplayModal.css';

interface CodeDisplayModalProps {
  lesson: Lesson;
  code: string;
  isOpen: boolean;
  onClose: () => void;
  onStop: () => void;
}

const CodeDisplayModal = ({ lesson, code, isOpen, onClose, onStop }: CodeDisplayModalProps) => {
  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <button onClick={onClose} className="modal-close-btn">&times;</button>
        <div className="modal-header">
          <h3>Код для занятия</h3>
          <p>{lesson.name} ({lesson.day}, {lesson.time})</p>
        </div>
        <div className="code-display-box">
          <div className="code">{code}</div>
          <p className="expiry-info">Этот код действителен в течение 5 минут.</p>
        </div>
        <div className="modal-actions">
          <button onClick={onStop} className="btn btn-danger">Остановить код</button>
          <button onClick={onClose} className="btn btn-secondary">Закрыть</button>
        </div>
      </div>
    </div>
  );
};

export default CodeDisplayModal; 