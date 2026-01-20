import { Routes, Route, Navigate } from "react-router-dom"

import SignUpPage from "./pages/SignUpPage"
import LoginPage from "./pages/LoginPage"
import HomePage from "./pages/HomePage"
import RoomDetailsPage from "./pages/RoomDetailsPage"
import PresenterDashboardPage from "./pages/PresenterDashboardPage"
import ViewerDashboardPage from "./pages/ViewerDashboard"

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<SignUpPage />} />
      <Route path="/rooms/:room_id" element={<RoomDetailsPage />} />
      <Route path="/sessions/:session_id/presenter" element={<PresenterDashboardPage />} />
      <Route path="/sessions/:session_id/viewer" element={<ViewerDashboardPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App