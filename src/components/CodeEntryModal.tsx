import { useState } from 'react';
import type { Lesson } from '../types';

interface CodeEntryModalProps {
  lesson: Lesson;
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (code: string) => void;
}

const CodeEntryModal = ({ lesson, isOpen, onClose, onSubmit }: CodeEntryModalProps) => {
  const [code, setCode] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (code.length !== 5) {
      setError('Код должен содержать 5 цифр');
      return;
    }
    if (!/^\d{5}$/.test(code)) {
      setError('Код должен содержать только цифры');
      return;
    }
    setError('');
    onSubmit(code);
    setCode('');
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h3>Отметить посещение</h3>
          <button className="close-btn" onClick={onClose}>&times;</button>
        </div>
        <div className="modal-body">
          <p><strong>Предмет:</strong> {lesson.name}</p>
          <p><strong>Время:</strong> {lesson.time}</p>
          <p><strong>Кабинет:</strong> {lesson.room}</p>
          
          <form onSubmit={handleSubmit}>
            <div className="code-input-group">
              <label htmlFor="code">Введите 5-значный код:</label>
              <input
                id="code"
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="00000"
                maxLength={5}
                className="code-input"
              />
            </div>
            {error && <p className="error-message">{error}</p>}
            <div className="modal-actions">
              <button type="button" onClick={onClose} className="btn-secondary">
                Отмена
              </button>
              <button type="submit" className="btn-primary">
                Отметить посещение
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CodeEntryModal; 