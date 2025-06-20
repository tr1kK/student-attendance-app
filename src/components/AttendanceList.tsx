import type { AttendanceRecord } from '../types';

interface AttendanceListProps {
  attendanceRecords: AttendanceRecord[];
  lessonName: string;
}

const AttendanceList = ({ attendanceRecords, lessonName }: AttendanceListProps) => {
  if (attendanceRecords.length === 0) {
    return (
      <div className="attendance-list">
        <h3>Посещаемость: {lessonName}</h3>
        <p className="no-attendance">Пока нет отметок о посещении</p>
      </div>
    );
  }

  return (
    <div className="attendance-list">
      <h3>Посещаемость: {lessonName}</h3>
      <p className="attendance-count">
        Отметились: {attendanceRecords.length} студентов
      </p>
      <div className="attendance-table">
        <table>
          <thead>
            <tr>
              <th>Студент</th>
              <th>Время отметки</th>
              <th>Код</th>
            </tr>
          </thead>
          <tbody>
            {attendanceRecords.map(record => (
              <tr key={record.id}>
                <td>{record.studentName}</td>
                <td>{new Date(record.timestamp).toLocaleString('ru-RU')}</td>
                <td>{record.code}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AttendanceList; 