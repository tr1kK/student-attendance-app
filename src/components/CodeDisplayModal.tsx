import type { Lesson } from '../types';

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
      <div className="modal">
        <div className="modal-header">
          <h3>Код для посещения</h3>
          <button className="close-btn" onClick={onClose}>&times;</button>
        </div>
        <div className="modal-body">
          <p><strong>Предмет:</strong> {lesson.name}</p>
          <p><strong>Время:</strong> {lesson.time}</p>
          <p><strong>Кабинет:</strong> {lesson.room}</p>
          
          <div className="code-display">
            <h4>Код для студентов:</h4>
            <div className="code-box">
              <span className="code-text">{code}</span>
            </div>
            <p className="code-instruction">
              Покажите этот код студентам для отметки посещения
            </p>
          </div>
          
          <div className="modal-actions">
            <button onClick={onStop} className="btn-danger">
              Остановить код
            </button>
            <button onClick={onClose} className="btn-secondary">
              Закрыть
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CodeDisplayModal; 