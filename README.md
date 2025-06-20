# Student Attendance Tracking System

A React + TypeScript + Vite application for tracking student attendance in classroom sessions. This system allows teachers to generate attendance codes and students to mark their attendance using these codes.

## Features

### 🔐 Authentication
- Student login with student ID and password
- Teacher login with staff ID and password
- Role-based routing to appropriate pages
- No registration - test users are hardcoded

### 🎓 Student Features
- View weekly class schedule
- Enter 5-digit attendance codes for active lessons
- View attendance history
- Real-time feedback on attendance marking

### 🧑‍🏫 Teacher Features
- View the same weekly schedule
- Generate 5-digit attendance codes for lessons
- View real-time attendance lists
- Stop/restart attendance codes
- Track which students have submitted codes

### 💾 Data Persistence
- All data stored in localStorage
- Attendance records
- Generated codes
- User sessions

## Test Users

### Students
- **ID:** `stu001` | **Password:** `12345` | **Name:** Иванов И.И.
- **ID:** `stu002` | **Password:** `23456` | **Name:** Петрова А.А.

### Teachers
- **ID:** `teach001` | **Password:** `admin1` | **Name:** Преподаватель С.С.

## Getting Started

### Prerequisites
- Node.js (version 16 or higher)
- npm or yarn

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd student-attendance-app
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

4. Open your browser and navigate to `http://localhost:5173`

## How to Use

### For Students

1. **Login**: Use your student ID and password
2. **View Schedule**: See your weekly class schedule
3. **Mark Attendance**: 
   - Click "Отметить посещение" on an active lesson
   - Enter the 5-digit code provided by your teacher
   - Submit to mark your attendance
4. **View History**: Check your attendance history at the bottom of the page

### For Teachers

1. **Login**: Use your staff ID and password
2. **View Schedule**: See the weekly class schedule
3. **Generate Codes**:
   - Click "Начать" on any lesson
   - A 5-digit code will be generated and displayed
   - Share this code with students
4. **Monitor Attendance**: 
   - View real-time attendance list below the schedule
   - See which students have submitted the code
5. **Stop Codes**: Click "Остановить код" to deactivate the current code

## Technical Details

### Project Structure
```
src/
├── components/          # Reusable UI components
│   ├── Schedule.tsx
│   ├── CodeEntryModal.tsx
│   ├── CodeDisplayModal.tsx
│   └── AttendanceList.tsx
├── pages/              # Main page components
│   ├── LoginPage.tsx
│   ├── StudentPage.tsx
│   └── TeacherPage.tsx
├── data/               # Static data and types
│   ├── users.ts
│   └── schedule.ts
├── utils/              # Utility functions
│   ├── auth.ts
│   └── storage.ts
├── types/              # TypeScript type definitions
│   └── index.ts
└── App.tsx             # Main application component
```

### Key Technologies
- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **localStorage** - Data persistence

### Data Models

#### User
```typescript
interface User {
  id: string;
  password: string;
  name: string;
  role: 'student' | 'teacher';
}
```

#### Lesson
```typescript
interface Lesson {
  id: string;
  name: string;
  day: string;
  time: string;
  teacher: string;
  room: string;
}
```

#### AttendanceRecord
```typescript
interface AttendanceRecord {
  id: string;
  lessonId: string;
  lessonName: string;
  studentId: string;
  studentName: string;
  timestamp: string;
  code: string;
}
```

#### GeneratedCode
```typescript
interface GeneratedCode {
  id: string;
  lessonId: string;
  lessonName: string;
  code: string;
  teacherId: string;
  createdAt: string;
  expiresAt: string;
  isActive: boolean;
}
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## Browser Support

This application works in all modern browsers that support:
- ES6+ features
- localStorage API
- CSS Grid and Flexbox

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
