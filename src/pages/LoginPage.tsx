// src/pages/LoginPage.tsx
import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import * as api from '../utils/api';
import { jwtDecode } from 'jwt-decode';

const LoginPage = () => {
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async () => {
    try {
      setError('');
      const response = await api.login(identifier, password);
      if (response.token && response.user) {
        localStorage.setItem('authToken', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));

        switch (response.user.role) {
          case 'student':
            navigate('/student');
            break;
          case 'teacher':
            navigate('/teacher');
            break;
          case 'admin':
            navigate('/admin');
            break;
          default:
            navigate('/');
        }
      } else {
        setError('Ошибка входа: неверный ответ от сервера.');
      }
    } catch (err: any) {
      setError(err.message || 'Ошибка входа');
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <h2>Вход в систему</h2>
        {error && <p className="error">{error}</p>}
        <input
          placeholder="Логин"
          value={identifier}
          onChange={e => setIdentifier(e.target.value)}
        />
        <input
          placeholder="Пароль"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        <button onClick={handleLogin}>Войти</button>
        <p style={{ textAlign: 'center', marginTop: '1rem' }}>
          Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
