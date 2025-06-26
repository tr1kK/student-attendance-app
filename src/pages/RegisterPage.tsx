import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import * as api from '../utils/api';
import type { Group } from '../types';

const RegisterPage = () => {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [groupId, setGroupId] = useState<string>('');
  const [groups, setGroups] = useState<Group[]>([]);
  const [error, setError] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const availableGroups = await api.getGroups();
        setGroups(availableGroups);
      } catch (err) {
        setError('Не удалось загрузить список групп. Пожалуйста, попробуйте позже.');
        console.error('Error fetching groups:', err);
      }
    };
    fetchGroups();
  }, []);

  const handleRegister = async () => {
    if (!groupId) {
      setError('Пожалуйста, выберите группу.');
      return;
    }
    try {
      setError('');
      setMessage('');
      const response = await api.register({
        identifier,
        password,
        name,
        email,
        group_id: Number(groupId),
      });
      setMessage(response.message || 'Регистрация прошла успешно!');
      setTimeout(() => navigate('/'), 2000);
    } catch (err: any) {
      setError(err.message || 'Ошибка регистрации');
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <h2>Регистрация студента</h2>
        {message && <p className="message success">{message}</p>}
        {error && <p className="error">{error}</p>}
        <input
          placeholder="Логин (ID студента)"
          value={identifier}
          onChange={e => setIdentifier(e.target.value)}
        />
        <input
          placeholder="Пароль"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        <input
          placeholder="ФИО"
          value={name}
          onChange={e => setName(e.target.value)}
        />
        <input
          placeholder="Email"
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
        />
        <select value={groupId} onChange={e => setGroupId(e.target.value)} required>
          <option value="" disabled>-- Выберите группу --</option>
          {groups.map(group => (
            <option key={group.id} value={group.id}>
              {group.name}
            </option>
          ))}
        </select>
        <button onClick={handleRegister}>Зарегистрироваться</button>
        <p style={{ textAlign: 'center', marginTop: '1rem' }}>
          Уже есть аккаунт? <Link to="/">Войти</Link>
        </p>
      </div>
    </div>
  );
};

export default RegisterPage; 