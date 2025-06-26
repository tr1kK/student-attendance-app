import type { AttendanceRecord } from '../types';

interface AttendanceListProps {
  attendanceRecords: AttendanceRecord[];
}

const AttendanceList = ({ attendanceRecords }: AttendanceListProps) => {
  return (
    <div className="attendance-list">
      {attendanceRecords.length === 0 ? (
        <p>На данный момент отметившихся нет.</p>
      ) : (
        <ul>
          {attendanceRecords.map(record => (
            <li key={record.id}>
              <span>{record.student.name}</span>
              <span className="timestamp">{new Date(record.submitted_at).toLocaleTimeString('ru-RU')}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default AttendanceList; 