// src/pages/LoginPage.tsx
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authenticate } from '../utils/auth';

const LoginPage = () => {
  const [id, setId] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = () => {
    const user = authenticate(id, password);
    if (!user) {
      setError('Неверный логин или пароль');
      return;
    }
    navigate(`/${user.role}`);
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <h2>Вход в систему</h2>
        <input
          placeholder="Номер (студенческий/служебный)"
          value={id}
          onChange={e => setId(e.target.value)}
        />
        <input
          placeholder="Пароль"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
        />
        <button onClick={handleLogin}>Войти</button>
        {error && <p className="error">{error}</p>}
      </div>
    </div>
  );
};

export default LoginPage;
