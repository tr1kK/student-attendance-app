import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import * as api from '../utils/api';
import type { User } from '../types';

const AdminPage = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [message, setMessage] = useState('');
  const navigate = useNavigate();
  const userName = localStorage.getItem('userName') || 'Admin';

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const data = await api.adminGetUsers();
        setUsers(data);
      } catch (error) {
        console.error('Failed to fetch users:', error);
        setMessage('Failed to load users.');
      }
    };
    fetchUsers();
  }, []);

  const handleDeleteUser = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await api.adminDeleteUser(id);
        setUsers(users.filter(user => user.id !== id));
        setMessage('User deleted successfully.');
      } catch (error) {
        console.error('Failed to delete user:', error);
        setMessage('Failed to delete user.');
      } finally {
        setTimeout(() => setMessage(''), 3000);
      }
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userRole');
    localStorage.removeItem('userName');
    navigate('/');
  };

  return (
    <div className="admin-page">
      <header className="page-header">
        <h1>Панель администратора</h1>
        <div className="user-info">
          <span>Добро пожаловать, {userName}</span>
          <button onClick={handleLogout} className="logout-btn">
            Logout
          </button>
        </div>
      </header>

      {message && <div className="message success">{message}</div>}

      <div className="page-content">
        <h2>Список пользователей</h2>
        <div className="attendance-table">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Identifier</th>
                <th>Name</th>
                <th>Email</th>
                <th>Role</th>
                <th>Group</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map(user => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>{user.identifier}</td>
                  <td>{user.name}</td>
                  <td>{user.email}</td>
                  <td>{user.role}</td>
                  <td>{user.group?.name || 'N/A'}</td>
                  <td>
                    <button className="btn-danger" onClick={() => handleDeleteUser(user.id)}>Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default AdminPage; 