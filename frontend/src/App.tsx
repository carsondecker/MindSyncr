import { Routes, Route, Navigate } from "react-router-dom"

import SignUpPage from "./pages/SignUpPage"
import LoginPage from "./pages/LoginPage"
import HomePage from "./pages/HomePage"
import RoomDetailsPage from "./pages/RoomDetailsPage"
import SessionsPage from "./pages/SessionPage"

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<SignUpPage />} />
      <Route path="/rooms/:id" element={<RoomDetailsPage />} />
      <Route path="/sessions/:id" element={<SessionsPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App